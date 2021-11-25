package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type ArtWork struct {
}

type User struct {
	ID      string         `json:"ID"`
	Ownlist map[string]int `json:"Ownlist"`
}

type Works struct {
	UCIcode       string         `json:"UCIcode"`
	Title         string         `json:"Title"`
	Artist        string         `json:"Artist"`
	InitialPrice  int            `json:"InitialPrice"`
	TotalProperty int            `json:"TotalProperty"`
	Owners        map[string]int `json:"Owners"`
	Status        string         `json:"Status"`
	FinalPrice    int            `json:"FinalPrice"`
}

const COMPANY_ID = "myCompany"

func (a *ArtWork) Init(APIstub shim.ChaincodeStubInterface) peer.Response {

	data := User{ID: COMPANY_ID, Ownlist: map[string]int{}}
	dataAsBytes, _ := json.Marshal(data)

	APIstub.PutState(COMPANY_ID, dataAsBytes)

	return shim.Success(nil)
}

func (a *ArtWork) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {

	fn, args := APIstub.GetFunctionAndParameters()

	if fn == "addUser" {
		return a.addUser(APIstub, args)
	} else if fn == "addWork" {
		return a.addWork(APIstub, args)
	} else if fn == "getInfos" {
		return a.getInfos(APIstub, args)
	} else if fn == "tradeProps" {
		return a.tradeProps(APIstub, args)
	} else if fn == "endTradeProps" {
		return a.endTradeProps(APIstub, args)
	} else if fn == "getHistory" {
		return a.getHistory(APIstub, args)
	}

	return shim.Error("Not supported chaincode function")
}

func (a *ArtWork) addUser(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments")
	}

	userID := args[0]

	// key: ID (args[0]), value: ID+OwnList (args[0]+map[string]int{}) -> marshal (make data to byte)
	data := User{ID: userID, Ownlist: map[string]int{}}
	dataAsBytes, _ := json.Marshal(data)

	APIstub.PutState(userID, dataAsBytes)

	return shim.Success(nil)
}

func (a *ArtWork) addWork(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect arguments")
	}

	// key: UCIcode (args[0]), value: UCIcode+Title+Artist+InitialPrice+TotalProperty+Status+Owners+FinalPrice  (args[0]+args[1]+args[2]+args[3]+args[4~]+1+[])
	UCIcode := args[0]
	title := args[1]
	artist := args[2]
	initialPrice, _ := strconv.Atoi(args[3])
	totalProperty, _ := strconv.Atoi(args[4])
	owners := map[string]int{COMPANY_ID: totalProperty}

	workData := Works{UCIcode: UCIcode, Title: title, Artist: artist, InitialPrice: initialPrice, TotalProperty: totalProperty, Owners: owners, Status: "ENROLL", FinalPrice: -1}
	workDataAsBytes, _ := json.Marshal(workData)

	APIstub.PutState(UCIcode, workDataAsBytes)

	companyAsBytes, _ := APIstub.GetState(COMPANY_ID)
	company := User{}
	json.Unmarshal(companyAsBytes, &company)

	company.Ownlist[UCIcode] = totalProperty
	companyAsBytes, _ = json.Marshal(company)
	APIstub.PutState(COMPANY_ID, companyAsBytes)

	return shim.Success(nil)
}

func (a *ArtWork) getInfos(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments")
	}

	pID := args[0]

	idAsBytes, _ := APIstub.GetState(pID)
	if idAsBytes == nil {
		shim.Error("Invaild ID")
	}

	return shim.Success(idAsBytes)
}

func (a *ArtWork) tradeProps(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect arguments")
	}

	sellerID := args[0]
	buyerID := args[1]
	workID := args[2]
	propsNum := args[3]

	sellerAsBytes, _ := APIstub.GetState(sellerID)
	buyerAsBytes, _ := APIstub.GetState(buyerID)
	workAsBytes, _ := APIstub.GetState(workID)
	transferNum, _ := strconv.Atoi(propsNum)

	seller := User{}
	buyer := User{}
	work := Works{}

	json.Unmarshal(sellerAsBytes, &seller)
	json.Unmarshal(buyerAsBytes, &buyer)
	json.Unmarshal(workAsBytes, &work)

	// chake work's status
	if work.Status == "END" {
		return shim.Error("This work can't trade anymore, please check again")
	}

	// check seller have work
	_, haveWork := seller.Ownlist[workID]
	_, isOwnWork := work.Owners[sellerID]
	if !(haveWork) || !(isOwnWork) {
		return shim.Error("Seller doesn't have work, please check again")
	}

	// check seller have enough ownership
	if seller.Ownlist[workID] < transferNum {
		return shim.Error("Seller doesn't have enough ownership to trade")
	}

	seller.Ownlist[workID] -= transferNum
	work.Owners[sellerID] -= transferNum
	if _, ishaveWork := buyer.Ownlist[args[2]]; ishaveWork { // user1 and user2 both have work

		buyer.Ownlist[workID] += transferNum
		work.Owners[buyerID] += transferNum
	} else if seller.Ownlist[workID] > transferNum { // user1 only have work and give ownership partially

		buyer.Ownlist[workID] = transferNum
		work.Owners[buyerID] = transferNum
	} else { // user1 only have work and give all ownership

		delete(work.Owners, sellerID)
		delete(seller.Ownlist, workID)
		work.Owners[buyerID] = transferNum
		buyer.Ownlist[workID] = transferNum
	}

	work.Status = "TRADING"

	sellerAsBytes, _ = json.Marshal(seller)
	buyerAsBytes, _ = json.Marshal(buyer)
	workAsBytes, _ = json.Marshal(work)

	APIstub.PutState(sellerID, sellerAsBytes)
	APIstub.PutState(buyerID, buyerAsBytes)
	APIstub.PutState(workID, workAsBytes)

	return shim.Success(nil)
}

func (a *ArtWork) endTradeProps(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect arguments")
	}

	workID := args[0]
	finalPrice, _ := strconv.Atoi(args[1])
	if finalPrice < 0 {
		return shim.Error("Please, check the price")
	}

	workAsBytes, _ := APIstub.GetState(workID)
	work := Works{}
	json.Unmarshal(workAsBytes, &work)

	work.Status = "END"
	work.FinalPrice = finalPrice

	workAsBytes, _ = json.Marshal(work)
	APIstub.PutState(workID, workAsBytes)

	return shim.Success(nil)
}

func (a *ArtWork) getHistory(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments")
	}

	ID := args[0]

	fmt.Printf("- start getHistoryForKey: %s\n", ID)

	resultsIterator, err := APIstub.GetHistoryForKey(ID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the work
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {

		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Values\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForKey returning: \n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(ArtWork))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

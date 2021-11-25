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
	Status        string         `json:"Status"`
	Owners        map[string]int `json:"Owners"`
	FinalPrice    int            `json:"FinalPrice"`
}

type ArtWork struct {
}

func (a *ArtWork) Init(APIstub shim.ChaincodeStubInterface) peer.Response {

	var ownList = make(map[string]int)
	data := User{ID: "myCompany", Ownlist: ownList}
	dataAsBytes, _ := json.Marshal(data)

	err := APIstub.PutState("myCompany", dataAsBytes)
	if err != nil {
		return shim.Error("Failed to Init")
	}

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

	var ownList = make(map[string]int)
	// key: ID (args[0]), value: ID+OwnList (args[0]+[]) -> marshal (make data to byte)
	data := User{ID: args[0], Ownlist: ownList}
	dataAsBytes, _ := json.Marshal(data)

	err := APIstub.PutState(args[0], dataAsBytes)
	if err != nil {
		return shim.Error("Failed to add user")
	}

	return shim.Success(nil)
}

func (a *ArtWork) addWork(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println(len(args))
	if len(args) != 5 {
		return shim.Error("Incorrect arguments")
	}

	// key: UCIcode (args[0]), value: UCIcode+Title+Artist+InitialPrice+TotalProperty+Status+Owners+FinalPrice  (args[0]+args[1]+args[2]+args[3]+args[4~]+1+[])
	price, _ := strconv.Atoi(args[3])
	property, _ := strconv.Atoi(args[4])
	owners := map[string]int{"myCompany": property}

	data := Works{UCIcode: args[0], Title: args[1], Artist: args[2], InitialPrice: price, TotalProperty: property, Status: "ENROLL", Owners: owners, FinalPrice: -1}
	dataAsBytes, _ := json.Marshal(data)

	err := APIstub.PutState(args[0], dataAsBytes)
	if err != nil {
		return shim.Error("Failed to add work")
	}

	companyAsBytes, _ := APIstub.GetState("myCompany")
	company := User{}
	json.Unmarshal(companyAsBytes, &company)

	company.Ownlist[args[0]] = property
	err = APIstub.PutState("myCompany", companyAsBytes)
	if err != nil {
		return shim.Error("Failed to add work")
	}

	return shim.Success(nil)
}

func (a *ArtWork) getInfos(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect arguments")
	}

	idAsBytes, _ := APIstub.GetState(args[0])
	if idAsBytes == nil {
		shim.Error("Invaild ID")
	}

	return shim.Success(idAsBytes)
}

func (a *ArtWork) tradeProps(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect arguments")
	}

	userAsBytes1, _ := APIstub.GetState(args[0])
	userAsBytes2, _ := APIstub.GetState(args[1])
	workAsBytes, _ := APIstub.GetState(args[2])
	transferNum, _ := strconv.Atoi(args[3])

	user1 := User{}
	user2 := User{}
	work := Works{}

	json.Unmarshal(userAsBytes1, &user1)
	json.Unmarshal(userAsBytes2, &user2)
	json.Unmarshal(workAsBytes, &work)

	// chake work's status
	if work.Status == "END" {
		return shim.Error("This work can't trade anymore, please check again")
	}

	// check seller have work
	_, haveWork := user1.Ownlist[args[2]]
	_, isOwnWork := work.Owners[args[0]]
	if !(haveWork) || !(isOwnWork) {
		return shim.Error("Seller doesn't have work, please check again")
	}

	// check seller have enough ownership
	if user1.Ownlist[args[2]] < transferNum {
		return shim.Error("Seller doesn't have enough ownership to trade")
	}

	user1.Ownlist[args[2]] -= transferNum
	work.Owners[args[0]] -= transferNum
	if _, ishaveWork := user2.Ownlist[args[2]]; ishaveWork { // user1 and user2 both have work

		user2.Ownlist[args[2]] += transferNum
		work.Owners[args[1]] += transferNum
	} else if user1.Ownlist[args[2]] > transferNum { // user1 only have work and give ownership partially

		work.Owners[args[1]] = transferNum
		user2.Ownlist[args[2]] = transferNum
	} else { // user1 only have work and give all ownership

		delete(work.Owners, args[0])
		delete(user1.Ownlist, args[2])
		work.Owners[args[1]] = transferNum
		user2.Ownlist[args[2]] = transferNum
	}

	work.Status = "TRADING"

	userAsBytes1, _ = json.Marshal(user1)
	userAsBytes2, _ = json.Marshal(user2)
	workAsBytes, _ = json.Marshal(work)

	APIstub.PutState(args[0], userAsBytes1)
	APIstub.PutState(args[1], userAsBytes2)
	APIstub.PutState(args[2], workAsBytes)

	return shim.Success(nil)
}

func (a *ArtWork) endTradeProps(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect arguments")
	}

	workAsBytes, _ := APIstub.GetState(args[0])
	work := Works{}
	json.Unmarshal(workAsBytes, &work)

	work.Status = "END"

	finalPrice, _ := strconv.Atoi(args[1])
	if finalPrice < 0 {
		return shim.Error("Please, check the price")
	}
	work.FinalPrice = finalPrice

	workAsBytes, _ = json.Marshal(work)
	APIstub.PutState(args[0], workAsBytes)

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

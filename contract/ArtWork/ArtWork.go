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
	InitalPrice   int            `json:"InitalPrice"`
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

	var result string
	var err error
	if fn == "addUser" {
		result, err = a.addUser(APIstub, args)
	} else if fn == "addWork" {
		result, err = a.addWork(APIstub, args)
	} else if fn == "tradeProps" {
		result, err = a.tradeProps(APIstub, args)
	} else if fn == "endTradeProps" {
		result, err = a.endTradeProps(APIstub, args)
	} else if fn == "getHistory" {
		result, err = a.getHistory(APIstub, args)
	} else {
		return shim.Error("Not supported chaincode function")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func (a *ArtWork) addUser(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	var ownList = make(map[string]int)
	// key: ID (args[0]), value: ID+OwnList (args[0]+[]) -> marshal (make data to byte)
	data := User{ID: args[0], Ownlist: ownList}
	dataAsBytes, _ := json.Marshal(data)

	err := APIstub.PutState(args[0], dataAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to add user : %s", err)
	}

	return string(dataAsBytes), nil
}

func (a *ArtWork) addWork(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println(len(args))
	if len(args) != 5 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	// key: UCIcode (args[0]), value: UCIcode+Title+Artist+InitalPrice+TotalProperty+Status+Owners+FinalPrice  (args[0]+args[1]+args[2]+args[3]+args[4~]+1+[])
	price, _ := strconv.Atoi(args[3])
	property, _ := strconv.Atoi(args[4])
	owners := map[string]int{"myCompany": property}

	data := Works{UCIcode: args[0], Title: args[1], Artist: args[2], InitalPrice: price, TotalProperty: property, Status: "ENROLL", Owners: owners, FinalPrice: -1}
	dataAsBytes, _ := json.Marshal(data)

	err := APIstub.PutState(args[0], dataAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to add work : %s", err)
	}

	companyAsBytes, _ := APIstub.GetState("myCompany")
	company := User{}
	json.Unmarshal(companyAsBytes, &company)

	company.Ownlist[args[0]] = property
	err = APIstub.PutState("myCompany", companyAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to add work : %s", err)
	}

	return string(dataAsBytes), nil
}

func (a *ArtWork) tradeProps(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect arguments")
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
		return "", fmt.Errorf("This work can't trade anymore, please check again")
	}

	// check seller have work
	_, haveWork := user1.Ownlist[args[2]]
	_, isOwnWork := work.Owners[args[0]]
	if !(haveWork) || !(isOwnWork) {
		return "", fmt.Errorf("Seller doesn't have work, please check again")
	}

	// check seller have enough ownership
	if user1.Ownlist[args[2]] < transferNum {
		return "", fmt.Errorf("Seller doesn't have enough ownership to trade")
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

	return string(userAsBytes1) + string(userAsBytes2) + string(workAsBytes), nil
}

func (a *ArtWork) endTradeProps(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	workAsBytes, _ := APIstub.GetState(args[0])
	work := Works{}
	json.Unmarshal(workAsBytes, &work)

	work.Status = "END"

	finalPrice, _ := strconv.Atoi(args[1])
	if finalPrice < 0 {
		return "", fmt.Errorf("Please, check the price")
	}
	work.FinalPrice = finalPrice

	workAsBytes, _ = json.Marshal(work)
	APIstub.PutState(args[0], workAsBytes)

	return string(workAsBytes), nil
}

func (a *ArtWork) getHistory(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}
	ID := args[0]

	resultsIterator, err := APIstub.GetHistoryForKey(ID)
	if err != nil {
		return "", fmt.Errorf("Failed to get history : %s with error: %s", args[0], err)
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the work
	var buffer bytes.Buffer
	buffer.WriteString("start getHistory: " + ID + "-")
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {

		response, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf("Failed to get history of work : %s with error: %s", args[0], err)
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
	buffer.WriteString("] - getHistory returning")

	return buffer.String(), nil
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(ArtWork))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type User struct {
	ID      string         `json:"ID"`
	Ownlist map[string]int `json:"Ownlist"`
}

type Works struct {
	UCIcode        string         `json:"UCIcode"`
	WorkName       string         `json:"WorkName"`
	Artist         string         `json:"Artist"`
	InitalPrice    int            `json:"InitalPrice"`
	WorkInfo       []string       `json:"WorkInfo"`
	TotalOwnerShip int            `json:"TotalOwnership"`
	Owners         map[string]int `json:"Owners"`
	FinalPrice     int            `json:"FinalPrice"`
}

type ArtWork struct {
}

func (a *ArtWork) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
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
	} else if fn == "getUser" {
		result, err = a.getUser(APIstub, args)
	} else if fn == "getWork" {
		result, err = a.getWork(APIstub, args)
	} else if fn == "delUser" {
		result, err = a.delUser(APIstub, args)
	} else if fn == "delWork" {
		result, err = a.delWork(APIstub, args)
	} else if fn == "queryAllWorks" {
		result, err = a.queryAllWorks(APIstub, args)
	} else if fn == "searchWorksBasedOnWorkName" {
		result, err = a.searchWorksBasedOnWorkName(APIstub, args)
	} else if fn == "searchWorksBasedOnArtist" {
		result, err = a.searchWorksBasedOnArtist(APIstub, args)
	} else if fn == "transferOwnerShip" {
		result, err = a.transferOwnerShip(APIstub, args)
	} else if fn == "splitOwnerShip" {
		result, err = a.splitOwnerShip(APIstub, args)
	} else if fn == "recordPriceOfSoldWork" {
		result, err = a.recordPriceOfSoldWork(APIstub, args)
	} else if fn == "getHistoryForWork" {
		result, err = a.getHistoryForWork(APIstub, args)
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
	if len(args) < 4 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	// key: UCIcode (args[0]), value: UCIcode+WorkName+Artist+InitalPrice+WorkInfo+TotalOwnerShip+Owners  (args[0]+args[1]+args[2]+args[3]+args[4~]+1+[])
	price, _ := strconv.Atoi(args[3])
	var infos []string
	for i := 4; i < len(args); i++ {
		fmt.Println(args[i])
		infos = append(infos, args[i])
	}
	owners := map[string]int{"myCompany": 1}

	data := Works{UCIcode: args[0], WorkName: args[1], Artist: args[2], InitalPrice: price, WorkInfo: infos, TotalOwnerShip: 1, Owners: owners, FinalPrice: -1}
	dataAsBytes, _ := json.Marshal(data)

	err := APIstub.PutState(args[0], dataAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to add work : %s", err)
	}

	companyAsBytes, _ := APIstub.GetState("myCompany")
	company := User{}
	json.Unmarshal(companyAsBytes, &company)

	company.Ownlist[args[0]] = 1
	err = APIstub.PutState("myCompany", companyAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to add work : %s", err)
	}

	//  ==== Index the works to enable workName-based range queries ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~WorkName~UCIcode.
	//  This will enable very efficient state range queries based on composite keys matching indexName~workName~*
	indexName := "workName~UCIcode"
	colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{data.WorkName, data.UCIcode})
	if err != nil {
		return "", fmt.Errorf("Failed to add index : %s", err)
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value := []byte{0x00}
	APIstub.PutState(colorNameIndexKey, value)

	indexName = "artist~UCIcode"
	colorNameIndexKey, err = APIstub.CreateCompositeKey(indexName, []string{data.Artist, data.UCIcode})
	if err != nil {
		return "", fmt.Errorf("Failed to add index : %s", err)
	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	value = []byte{0x00}
	APIstub.PutState(colorNameIndexKey, value)

	return string(dataAsBytes), nil
}

func (a *ArtWork) getUser(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	value, err := APIstub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get user : %s with error: %s", args[0], err)
	}

	if value == nil {
		return "", fmt.Errorf("Asset not found : %s", args[0])
	}

	return string(value), nil
}

func (a *ArtWork) getWork(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	value, err := APIstub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset : %s with error: %s", args[0], err)
	}

	if value == nil {
		return "", fmt.Errorf("Asset not found : %s", args[0])
	}

	return string(value), nil
}

func (a *ArtWork) delUser(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	userAsBytes, _ := APIstub.GetState(args[0])
	user := User{}
	json.Unmarshal(userAsBytes, &user)

	if len(user.Ownlist) != 0 {
		return "", fmt.Errorf("User have works")
	}

	err := APIstub.DelState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete user : %s with error: %s", args[0], err)
	}

	return "User is deleted", nil
}

func (a *ArtWork) delWork(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	workAsBytes, _ := APIstub.GetState(args[0])
	work := Works{}
	json.Unmarshal(workAsBytes, &work)

	if work.FinalPrice == -1 {
		return "", fmt.Errorf("Work doesn't sell yet")
	}

	err := APIstub.DelState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete work : %s with error: %s", args[0], err)
	}

	return "Work is deleted", nil
}

func (a *ArtWork) queryAllWorks(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return "", fmt.Errorf("Failed to get works")
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf("Failed to get works")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("] - queryAllWorks")

	return buffer.String(), nil
}

func (a *ArtWork) searchWorksBasedOnWorkName(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	workName := args[0]

	// Query the workName~UCIcode index by workName
	// This will execute a key range query on all keys starting with 'workName'
	coloredMarbleResultsIterator, err := APIstub.GetStateByPartialCompositeKey("workName~UCIcode", []string{workName})
	if err != nil {
		return "", fmt.Errorf("Failed to get works based on" + args[0])
	}
	defer coloredMarbleResultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	// Iterate through result set and for each marble found, transfer to newOwner
	var i int
	for i = 0; coloredMarbleResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble name from the composite key
		responseRange, err := coloredMarbleResultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf("Failed to get works based on" + args[0])
		}

		// get the workName and UCIcode from workName~UCIcode composite key
		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return "", fmt.Errorf("Failed to get work based on" + args[0])
		}
		returnedUCIcode := compositeKeyParts[1]

		workAsBytes, _ := APIstub.GetState(returnedUCIcode)
		work := Works{}
		json.Unmarshal(workAsBytes, &work)

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedUCIcode)
		buffer.WriteString("\"")
		buffer.WriteString(", \"WorkName\":")
		buffer.WriteString("\"")
		buffer.WriteString(work.WorkName)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Artist\":")
		buffer.WriteString("\"")
		buffer.WriteString(work.Artist)
		buffer.WriteString("\"")
		buffer.WriteString(", \"InitalPrice\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(work.InitalPrice))
		buffer.WriteString("\"")
		buffer.WriteString(", \"WorkInfo\":")
		buffer.WriteString("\"")
		buffer.WriteString(strings.Join(work.WorkInfo, ", "))
		buffer.WriteString("\"")
		buffer.WriteString(", \"TotalOwnerShip\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(work.TotalOwnerShip))
		buffer.WriteString("\"")
		buffer.WriteString(", \"Owners\":")
		buffer.WriteString("\"")
		owners, _ := json.Marshal(work.Owners)
		buffer.WriteString(string(owners))
		buffer.WriteString("\"")
		buffer.WriteString(", \"FinalPrice\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(work.FinalPrice))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("] - search works based on work's name :" + args[0])

	return buffer.String(), nil
}

func (a *ArtWork) searchWorksBasedOnArtist(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	artist := args[0]

	// Query the artist~UCIcode index by workName
	// This will execute a key range query on all keys starting with 'artist'
	coloredMarbleResultsIterator, err := APIstub.GetStateByPartialCompositeKey("artist~UCIcode", []string{artist})
	if err != nil {
		return "", fmt.Errorf("Failed to get works based on" + args[0])
	}
	defer coloredMarbleResultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	// Iterate through result set and for each marble found, transfer to newOwner
	var i int
	for i = 0; coloredMarbleResultsIterator.HasNext(); i++ {
		responseRange, err := coloredMarbleResultsIterator.Next()
		if err != nil {
			return "", fmt.Errorf("Failed to get works based on" + args[0])
		}

		// get the artist and UCIcode from artist~UCIcode composite key
		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return "", fmt.Errorf("Failed to get work based on" + args[0])
		}
		returnedArtist := compositeKeyParts[1]

		workAsBytes, _ := APIstub.GetState(returnedArtist)
		work := Works{}
		json.Unmarshal(workAsBytes, &work)

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(returnedArtist)
		buffer.WriteString("\"")
		buffer.WriteString(", \"WorkName\":")
		buffer.WriteString("\"")
		buffer.WriteString(work.WorkName)
		buffer.WriteString("\"")
		buffer.WriteString(", \"Artist\":")
		buffer.WriteString("\"")
		buffer.WriteString(work.Artist)
		buffer.WriteString("\"")
		buffer.WriteString(", \"InitalPrice\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(work.InitalPrice))
		buffer.WriteString("\"")
		buffer.WriteString(", \"WorkInfo\":")
		buffer.WriteString("\"")
		buffer.WriteString(strings.Join(work.WorkInfo, ", "))
		buffer.WriteString("\"")
		buffer.WriteString(", \"TotalOwnerShip\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(work.TotalOwnerShip))
		buffer.WriteString("\"")
		buffer.WriteString(", \"Owners\":")
		buffer.WriteString("\"")
		owners, _ := json.Marshal(work.Owners)
		buffer.WriteString(string(owners))
		buffer.WriteString("\"")
		buffer.WriteString(", \"FinalPrice\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(work.FinalPrice))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("] - search works based on work's name :" + args[0])

	return buffer.String(), nil
}

func (a *ArtWork) transferOwnerShip(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

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

	userAsBytes1, _ = json.Marshal(user1)
	userAsBytes2, _ = json.Marshal(user2)
	workAsBytes, _ = json.Marshal(work)

	APIstub.PutState(args[0], userAsBytes1)
	APIstub.PutState(args[1], userAsBytes2)
	APIstub.PutState(args[2], workAsBytes)

	return string(userAsBytes1) + string(userAsBytes2) + string(workAsBytes), nil
}

func (a *ArtWork) splitOwnerShip(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	workAsBytes, _ := APIstub.GetState(args[0])
	work := Works{}

	newOwnerShipNum, _ := strconv.Atoi(args[1])
	json.Unmarshal(workAsBytes, &work)

	work.TotalOwnerShip = newOwnerShipNum
	work.Owners["myCompany"] = newOwnerShipNum

	workAsBytes, _ = json.Marshal(work)
	err := APIstub.PutState(args[0], workAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to update ownership number : %s", err)
	}

	companyAsBytes, _ := APIstub.GetState("myCompany")
	company := User{}
	json.Unmarshal(companyAsBytes, &company)

	company.Ownlist[args[0]] = newOwnerShipNum
	APIstub.PutState("myCompany", companyAsBytes)

	return string(workAsBytes), nil
}

func (a *ArtWork) recordPriceOfSoldWork(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	workAsBytes, _ := APIstub.GetState(args[0])
	work := Works{}
	json.Unmarshal(workAsBytes, &work)

	finalPrice, _ := strconv.Atoi(args[1])
	if finalPrice < 0 {
		return "", fmt.Errorf("Please, check the price")
	}
	work.FinalPrice = finalPrice

	workAsBytes, _ = json.Marshal(work)
	APIstub.PutState(args[0], workAsBytes)

	return string(workAsBytes), nil
}

func (a *ArtWork) getHistoryForWork(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments")
	}
	workCode := args[0]

	resultsIterator, err := APIstub.GetHistoryForKey(workCode)
	if err != nil {
		return "", fmt.Errorf("Failed to get history of work : %s with error: %s", args[0], err)
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the work
	var buffer bytes.Buffer
	buffer.WriteString("start getHistoryForMarble: " + workCode + "-")
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

		buffer.WriteString(", \"WorkName\":")
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
	buffer.WriteString("] - getHistoryForWork returning")

	return buffer.String(), nil
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(ArtWork))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

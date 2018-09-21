package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct{
}

type Task struct {
	Id string `json:"id"`
	Timecoin string `json:"timecoin"`
	Publisher string `json:"publisher"`
	Tasktype string `json:"tasktype"`
	Title string `json:"title"`
	Accepted string `json:"accepted"`
	Completed string `json:"completed"`
	Owner string `json:"owner"`
}

//{"Args":["init","zpl"]}
//original value: 0 asset
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Initing.......==============================================================")
	_, args := APIstub.GetFunctionAndParameters()
	var person string 
	var Aval int 
	var err error 

	Aval = 100
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	person = args[0]
	fmt.Printf("Asset = %d", Aval)

	//write the state to the ledger
	err = APIstub.PutState(person,[]byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryTask" {
		return s.queryTask(APIstub, args)
	}else if function == "queryPeople" {
		return s.queryPeople(APIstub,args)
	} else if function == "createTask" {
		return s.createTask(APIstub, args)
	} else if function == "createPeople" {
		return s.createPeople(APIstub)
	}else if function == "queryAllTasks" {
		return s.queryAllTasks(APIstub)
	} else if function == "changeTaskOwner" {
		return s.changeTaskOwner(APIstub, args)
	} else if function == "changeTaskState" { //completed or not
		return s.changeTaskState(APIstub,args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

//{"Args":["queryTask","task0"]}
func (s *SmartContract) queryTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	taskAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(taskAsBytes)
}

//{"Args":["queryPeople","people1"]}
func (s *SmartContract) queryPeople(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	peopleAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(peopleAsBytes)
}

//"institution"
//{"Args":["createTask","task0","100","zpl","person","help cleaning the window","not","not","None"]}
func (s *SmartContract) createTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var PersonName string
	var err error
	var TimeCoin,Aval int
	var res []byte
	TimeCoin, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	if args[3] == "person" {

	}
	if args[3] == "institution" {
		var task = Task{Id: args[0], Timecoin: args[1], Publisher: args[2], Tasktype: args[3],Title:args[4],
			Accepted:args[5],Completed: args[6], Owner: args[7]}
		taskAsBytes, _ := json.Marshal(task)
		res = taskAsBytes
		APIstub.PutState(args[0], taskAsBytes)
	}else if args[3] == "person"{
		PersonName = args[2]
		Avalbytes,err := APIstub.GetState(PersonName)
		if err != nil {
			return shim.Error("Failed to get state")
		}
		if Avalbytes == nil {
			return shim.Error("Entity not found")
		}
		Aval, _ = strconv.Atoi(string(Avalbytes))
		Aval -= TimeCoin  //decide whether Aval >= 0
		fmt.Printf("Person Aval = %d\n", Aval) 
		if Aval >= 0 {
			err = APIstub.PutState(PersonName, []byte(strconv.Itoa(Aval)))
			if err != nil {
				return shim.Error(err.Error())
			}
			var task = Task{Id: args[0], Timecoin: args[1], Publisher: args[2], Tasktype: args[3],Title:args[4],
			Accepted:args[5],Completed: args[6], Owner: args[7]}
			taskAsBytes, _ := json.Marshal(task)
			res = taskAsBytes
			APIstub.PutState(args[0], taskAsBytes)
		}else {
			return shim.Error("Insufficient account balance!")
		}
	}
	return shim.Success(res)
}

//{"Args":["createPeople","zpl"]} 
//0 asset
func (s *SmartContract) createPeople(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Person Init")
	_, args := APIstub.GetFunctionAndParameters()
	var person string 
	var Aval int 
	var err error 

	Aval = 0
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	person = args[0]
	fmt.Printf("Asset = %d", Aval)

	//write the state to the ledger
	err = APIstub.PutState(person,[]byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)	
}


func (s *SmartContract) queryAllTasks(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Task0"
	endKey := "Task999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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
	buffer.WriteString("]")

	fmt.Printf("- queryAllTasks:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//if a task is accepted, then we change its owner
//{"Args":["changeTaskOwner","task0","zpl"]}
func (s *SmartContract) changeTaskOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	taskAsBytes, _ := APIstub.GetState(args[0])
	task := Task{}

	json.Unmarshal(taskAsBytes, &task)
	task.Owner = args[1]

	taskAsBytes, _ = json.Marshal(task)
	APIstub.PutState(args[0], taskAsBytes)

	return shim.Success(nil)
}

//if a task is completed, then we change its the ledger and its state
//{"Args":["changeTaskState","task0","zpl"]}
func (s *SmartContract) changeTaskState(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var PersonName string
	var Timecoin int //Transaction value
	var err error
	var Aval int 

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	taskAsBytes, _ := APIstub.GetState(args[0])
	task := Task{}

	json.Unmarshal(taskAsBytes, &task)
	task.Completed = "Yes"

	Timecoin, err = strconv.Atoi(task.Timecoin)
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

	taskAsBytes, _ = json.Marshal(task)
	APIstub.PutState(args[0], taskAsBytes)

	
	PersonName = args[1]
	Avalbytes,err := APIstub.GetState(PersonName)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))
	Aval = Aval + Timecoin

	err = APIstub.PutState(PersonName, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
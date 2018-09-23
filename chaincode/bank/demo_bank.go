package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
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
	PublishTime string `json:"publishtime"`
	AcceptedTime string `json:"acceptedtime"`
	CompletedTime string `json:"completedtime"`
}

type People struct {
	Name string `json:"name"`
	Asset string `json:"asset"`
	PublishedTask string `json:"publishedtask"`
	AcceptedTask string `json:"acceptedtask"`
	CompletedTask string `json:"completedtask"`
}

//{"Args":["init","zpl"]}
//original value: 0 asset
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Person Init")
	_, args := APIstub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var person = People{Name: args[0], Asset: "0", PublishedTask: "None", AcceptedTask: "None", CompletedTask: "None"}
	
	personAsBytes, _ := json.Marshal(person)
	//write the state to the ledger
	err := APIstub.PutState(args[0], personAsBytes)	
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
	PersonName = args[2]
	TimeCoin, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}
	
	current_time := time.Now()
	timestamp := current_time.Unix()
    t_unix := time.Unix(timestamp, 0)
    t := t_unix.Format("2006-01-02 03:04:05 PM")

	if args[3] == "institution" {
		var task = Task{Id: args[0], Timecoin: args[1], Publisher: args[2], Tasktype: args[3],Title:args[4],
			Accepted:args[5],Completed: args[6], Owner: args[7], PublishTime: t, AcceptedTime:"None", CompletedTime: "None"}
		taskAsBytes, _ := json.Marshal(task)
		APIstub.PutState(args[0], taskAsBytes)
	}else if args[3] == "person" {
		peopleAsBytes, err := APIstub.GetState(PersonName)
		if err != nil {
			return shim.Error("Failed to get state")
		}
		if peopleAsBytes == nil {
			return shim.Error("Entity not found")
		}
		people := People{}
		json.Unmarshal(peopleAsBytes, &people)
		Aval, _ = strconv.Atoi(string(people.Asset))
		Aval = Aval - TimeCoin
		if Aval >= 0 {
			people.Asset = strconv.Itoa(Aval)
			if people.PublishedTask  == "None" {
				people.PublishedTask = args[0]
			}else{
				people.PublishedTask = people.PublishedTask + "," + args[0]
			}
			peopleAsBytes, _ = json.Marshal(people)
			APIstub.PutState(PersonName, peopleAsBytes)

			var task = Task{Id: args[0], Timecoin: args[1], Publisher: args[2], Tasktype: args[3],Title:args[4],
			Accepted:args[5],Completed: args[6], Owner: args[7], PublishTime: t, AcceptedTime:"None", CompletedTime: "None"}
			taskAsBytes, _ := json.Marshal(task)
			APIstub.PutState(args[0], taskAsBytes)
		}else {
			return shim.Error("Insufficient account balance!")
		}
	}
	return shim.Success(nil)
}

//{"Args":["createPeople","zpl"]} 
//0 asset
func (s *SmartContract) createPeople(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Person Creating")
	_, args := APIstub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var person = People{Name: args[0], Asset: "0", PublishedTask: "None", AcceptedTask: "None", CompletedTask: "None"}
	
	personAsBytes, _ := json.Marshal(person)
	//write the state to the ledger
	err := APIstub.PutState(args[0], personAsBytes)	
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (s *SmartContract) queryAllTasks(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "task0"
	endKey := "task999"

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
//when a task is accepted, this function must be called!!!
func (s *SmartContract) changeTaskOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	taskAsBytes, _ := APIstub.GetState(args[0])
	task := Task{}

	current_time := time.Now()
	timestamp := current_time.Unix()
    t_unix := time.Unix(timestamp, 0)
    t := t_unix.Format("2006-01-02 03:04:05 PM")


	json.Unmarshal(taskAsBytes, &task)
	task.Owner = args[1]
	task.Accepted = "yes"
	task.AcceptedTime = t
	taskAsBytes, _ = json.Marshal(task)
	APIstub.PutState(args[0], taskAsBytes)

	//update person's accepted task
	peopleAsBytes, _ := APIstub.GetState(args[1])
	people := People{}
	json.Unmarshal(peopleAsBytes, &people)
	if people.AcceptedTask  == "None" {
		people.AcceptedTask = args[0]
	}else{
		people.AcceptedTask = people.AcceptedTask + "," + args[0]
	}
	peopleAsBytes, _ = json.Marshal(people)
	APIstub.PutState(args[1], peopleAsBytes)
	return shim.Success(nil)
}

//if a task is completed, then we change its the ledger and its state
//{"Args":["changeTaskState","task0","zpl"]}
func (s *SmartContract) changeTaskState(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var Timecoin int //Transaction value
	var err error
	var Aval int 

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	//update task information
	current_time := time.Now()
	timestamp := current_time.Unix()
    t_unix := time.Unix(timestamp, 0)
    t := t_unix.Format("2006-01-02 03:04:05 PM")



	taskAsBytes, _ := APIstub.GetState(args[0])
	task := Task{}

	json.Unmarshal(taskAsBytes, &task)
	task.Completed = "Yes"
	task.CompletedTime = t
	Timecoin, err = strconv.Atoi(task.Timecoin)
	taskAsBytes, _ = json.Marshal(task)
	APIstub.PutState(args[0], taskAsBytes)

	peopleAsBytes, err := APIstub.GetState(args[1])
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if peopleAsBytes == nil {
		return shim.Error("Entity not found")
	}
	people := People{}
	json.Unmarshal(peopleAsBytes, &people)
	Aval, _ = strconv.Atoi(string(people.Asset))
	Aval = Aval + Timecoin
	people.Asset = strconv.Itoa(Aval)

	if people.CompletedTask  == "None" {
		people.CompletedTask = args[0]
	}else{
		people.CompletedTask = people.CompletedTask + "," + args[0]
	}
	peopleAsBytes, _ = json.Marshal(people)
	APIstub.PutState(args[1], peopleAsBytes)
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
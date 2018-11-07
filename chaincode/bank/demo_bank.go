package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"regexp"
)

type SmartContract struct{
}

type Task struct {  //single task  only one man can accept the task
	Id string `json:"id"`
	Timecoin string `json:"timecoin"`
	Publisher string `json:"publisherid"`
	PublisherName string `json:"publishername"`
	Tasktype string `json:"tasktype"` //ins_created / person_created
	Title string `json:"title"`
	Accepted string `json:"accepted"` // 0 not accepted 1 accepted
	Completed string `json:"completed"` // 0 not completed 1 completed
	Owner string `json:"owner"`
	OwnerName string `json:"ownername"`
	PublishTime string `json:"publishtime"`
	AcceptedTime string `json:"acceptedtime"`
	CompletedTime string `json:"completedtime"`
	Canceled string `json:"canceled"` //whether it is canceled
}

type GroupTask struct { //group task --- more than one man can accept
	Id string `json:"id"`
	Timecoin string `json:"timecoin"`
	Publisher string `json:"publisher"`
	PublisherName string `json:"publishername"`
	Tasktype string `json:"tasktype"` //ins_created/person_created
	Title string `json:"title"`
	Accepted string `json:"accepted"` // 0 not accepted 1 accepted
	Completed string `json:"completed"` // 0 not completed 1 completed
	Owner string `json:"owner"`
	OwnerName string `json:"ownername"`
	PublishTime string `json:"publishtime"`
	AcceptedTime string `json:"acceptedtime"`
 	CompletedTime string `json:"completedtime"`
 	StartTime string `json:"starttime"` //all people needed is found
 	AvailableNumber string `json:"availablenumber"` //still need x people
	Needpeople string `json:"Needpeople"` // number of needed people 
	Canceled string `json:"canceled"` //whether it is canceled
}

type People struct {
	Name string `json:"name"`
	Asset string `json:"asset"`
	PublishedTask string `json:"publishedtask"`
	AcceptedTask string `json:"acceptedtask"`
	CompletedTask string `json:"completedtask"`
	Creditscore string `json:"creditscore"`
	Disputedtask string `json:"disputedtask"` // zhengyi task ID
	IndexofPeople string `json:"indexofpeople"` //unique index of a people, start: person0
	Blacklist string `json:"blacklist"` // 1 = blacklist 0 = not in blacklist
}

type Institution struct {
	Name string `json:"name"`
	Asset string `json:"asset"`
	Authority string `json:"authority"` // 0 = cannot issue a task 1 = able to issue a task
	PublishedTask string `json:"publishedtask"` // TODO: modify to publishedtask
	Disputedtask string `json:"disputedtask"`
	Creditscore string `json:"creditscore"`
	IndexofIns string `json:"index"` // unique index of a institution, start: ins0
	Blacklist string `json:"blacklist"`
}

var ins_cnt int
var task_cnt int
var per_cnt int

//{"Args":["init"]}
//initializing the global variable
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Initing the environment......")
	ins_cnt = 0
	task_cnt = 0
	per_cnt = 0
	fmt.Println("Initialization finished !")
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "query" {
		return s.query(APIstub, args)
	} else if function == "createGroupTask" {
		return s.createGroupTask(APIstub,args)
	} else if function == "createTask" {
		return s.createTask(APIstub, args)
	} else if function == "createPeople" {
		return s.createPeople(APIstub)
	} else if function == "createInstitution" { //without authority 
		return s.createInstitution(APIstub) 
	} else if function == "registerInstitution" { //have authority
		return s.registerInstitution(APIstub)
	} else if function == "giveInstitutionCoin" { //monthly clear and monthly give
		return s.giveInstitutionCoin(APIstub)
	} else if function == "queryAllTasks" {
		return s.queryAllTasks(APIstub)
	} else if function == "acceptSingleTask" {
		return s.acceptSingleTask(APIstub, args)
	} else if function == "completeSingleTask" { //completed or not 
		return s.completeSingleTask(APIstub,args)
	} else if function == "acceptGroupTask" { //add group task owner
 		return s.acceptGroupTask(APIstub,args)
	} else if function == "completeGroupTask" {
		return s.completeGroupTask(APIstub,args)
	} else if function == "recordDisputedGroupTask" { //record group task
		return s.recordDisputedGroupTask(APIstub,args)
	} else if function == "recordDisputedTask" { // for single task
		return s.recordDisputedTask(APIstub,args)
	} else if function == "cancelTask" { 
		return s.cancelTask(APIstub,args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

//{"Args":["query","task0"]}
//{"Args":["query","person0"]}
//{"Args":["query","ins0"]}
func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	taskAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(taskAsBytes)
}

//"institution"
//{"Args":["createTask","task0","100","zpl","person","help cleaning the window","person0"]}  the index of zpl is person0 
//{"Args":["createTask","task0","100","inpluslab","institution","help cleaning the window","ins0"]} the index of inpluslab is ins0
func (s *SmartContract) createTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var Id string
	var err error
	var TimeCoin,Aval int
	var Name string
	Id = args[5]
	Name = args[2]
	TimeCoin, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	
	current_time := time.Now()
	timestamp := current_time.Unix()
    t_unix := time.Unix(timestamp, 0)
    t := t_unix.Format("2006-01-02 03:04:05 PM")

	if args[3] == "institution" {
		var task = Task{Id: args[0], Timecoin: args[1], Publisher: Id, PublisherName: Name ,Tasktype: args[3],Title:args[4],
			Accepted: "0",Completed: "0", Owner: "None", OwnerName: "None", PublishTime: t, AcceptedTime:"None", CompletedTime: "None",
			Canceled: "0"}
		taskAsBytes, _ := json.Marshal(task)
		
		var balance int
		institutionAsBytes,_:= APIstub.GetState(Id)
		institution := Institution{}
		json.Unmarshal(institutionAsBytes,&institution)
		balance,_ = strconv.Atoi(institution.Asset)
		if institution.Authority == "0" || institution.Blacklist == "1" {
			return shim.Error("Institution have no access to issue a task")
		}else {
			if balance < TimeCoin {
				return shim.Error("Insufficient account balance!")
			}else {
				balance -= TimeCoin
			}
			institution.Asset = strconv.Itoa(balance)
			if institution.PublishedTask == "None" {
				institution.PublishedTask = args[0]
			}else {
				institution.PublishedTask = institution.PublishedTask + "," + args[0]
			}
		}
		APIstub.PutState(args[0], taskAsBytes)
		institutionAsBytes,_ = json.Marshal(institution)
		APIstub.PutState(Id,institutionAsBytes)
	}else if args[3] == "person" {
		peopleAsBytes, err := APIstub.GetState(Id)
		if err != nil {
			return shim.Error("Failed to get state")
		}
		if peopleAsBytes == nil {
			return shim.Error("Entity not found")
		}
		people := People{}
		json.Unmarshal(peopleAsBytes, &people)
		if people.Blacklist == "1" {
			return shim.Error("Blacklist users cannot issue a task")
		}

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
			APIstub.PutState(Id, peopleAsBytes)

			var task = Task{Id: args[0], Timecoin: args[1], Publisher: Id, PublisherName: Name, Tasktype: args[3],Title:args[4],
			Accepted:"0",Completed: "0", Owner: "None",OwnerName: "None", PublishTime: t, AcceptedTime:"None", CompletedTime: "None",
			Canceled:"0"}
			taskAsBytes, _ := json.Marshal(task)
			APIstub.PutState(args[0], taskAsBytes)
		}else {
			return shim.Error("Insufficient account balance!")
		}
	}
	task_cnt += 1

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}


//need10people
//{"Args":["createGroupTask","task0","100","inpluslab","institution","help cleaning the window","ins0","10"]} 
func (s *SmartContract) createGroupTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var PersonName string
	var err error
	var TimeCoin,Aval int
	PersonName = args[5]
	TimeCoin, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}
	
	current_time := time.Now()
	timestamp := current_time.Unix()
    t_unix := time.Unix(timestamp, 0)
    t := t_unix.Format("2006-01-02 03:04:05 PM")

	if args[3] == "institution" {
		var task = GroupTask{Id: args[0], Timecoin: args[1], Publisher: args[5], PublisherName: args[2], Tasktype: args[3],Title:args[4],
			Accepted: "0",Completed: "0", Owner: "None", PublishTime: t, AcceptedTime:"None", CompletedTime: "None",
			StartTime: "None", AvailableNumber: args[6] ,Needpeople: args[6], Canceled: "0"}
		
		taskAsBytes, _ := json.Marshal(task)

		institutionAsBytes,_:= APIstub.GetState(PersonName)
		institution := Institution{}
		json.Unmarshal(institutionAsBytes,&institution)
		if institution.Authority == "0" || institution.Blacklist == "1" {
			return shim.Error("Institution have no access to issue a task")
		}else {
			var balance int
			balance,_ = strconv.Atoi(institution.Asset)
			if balance < TimeCoin {
				return shim.Error("Insufficient balance !")
			}else {
				balance -= TimeCoin
			}
			institution.Asset = strconv.Itoa(balance)
			if institution.PublishedTask == "None" {
				institution.PublishedTask = args[0]
			}else {
				institution.PublishedTask = institution.PublishedTask + "," + args[0]
			}
		}
		APIstub.PutState(args[0], taskAsBytes)
		institutionAsBytes,_ = json.Marshal(institution)
		APIstub.PutState(PersonName,institutionAsBytes)
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
		if people.Blacklist == "1" {
			return shim.Error("Blacklist users cannot issue a task")
		}

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

			var task = GroupTask{Id: args[0], Timecoin: args[1], Publisher: args[5], PublisherName: args[2], Tasktype: args[3],Title:args[4],
			Accepted:"0",Completed: "0", Owner: "None",OwnerName:"None", PublishTime: t, AcceptedTime:"None", CompletedTime: "None",
			StartTime: "None", AvailableNumber: args[6], Needpeople: args[6], Canceled: "0"}
			taskAsBytes, _ := json.Marshal(task)
			APIstub.PutState(args[0], taskAsBytes)
		}else {
			return shim.Error("Insufficient account balance!")
		}
	}
	task_cnt += 1
	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}


//{"Args":["createPeople","zpl"]} 
//0 asset
//updated
func (s *SmartContract) createPeople(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Person Creating")
	_, args := APIstub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	index := strconv.Itoa(per_cnt)
	per_index := "person" + index

	var person = People{Name: args[0], Asset: "0", PublishedTask: "None", AcceptedTask: "None", CompletedTask: "None",
	 Creditscore: "100",Disputedtask: "None", IndexofPeople: per_index, Blacklist: "0"}
	personAsBytes, _ := json.Marshal(person)
	//write the state to the ledger
	err := APIstub.PutState(per_index, personAsBytes)	
	if err != nil {
		return shim.Error(err.Error())
	}
	per_cnt += 1

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}


//{"Args":["createInstitution","inpluslab"]} 
//0 asset
//updated
func (s *SmartContract) createInstitution(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Institution Creating")
	_, args := APIstub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	index := strconv.Itoa(ins_cnt)
	ins_index := "ins" + index

	var institution = Institution{Name: args[0], Asset: "0", Authority: "0", PublishedTask: "None", Disputedtask: "None",
	 Creditscore: "100", IndexofIns: ins_index ,Blacklist: "0"}

	institutionAsBytes, _ := json.Marshal(institution)

	//write the state to the ledger
	err := APIstub.PutState(ins_index, institutionAsBytes)	
	if err != nil {
		return shim.Error(err.Error())
	}
	ins_cnt += 1

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}

// inpluslab ----- index
//{"Args":["registerInstitution","ins0"]} 
//change authority of an institution
//updated
func (s *SmartContract) registerInstitution(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Institution Registering")
	_, args := APIstub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	institutionAsBytes,_ := APIstub.GetState(args[0])
	institution := Institution{}
	json.Unmarshal(institutionAsBytes, &institution)
	institution.Authority = "1"
	institution.Asset = strconv.Itoa(1000)


	institutionAsBytes,_ = json.Marshal(institution)
	APIstub.PutState(args[0],institutionAsBytes)

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}

//{"Args":["giveInstitutionCoin"]}
//monthly called function !!
func (s *SmartContract) giveInstitutionCoin(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Send TimeCoin to every Institution")
	for i:= 0; i < ins_cnt; i ++ {
		Ins_id := "ins"
		Ins_id += strconv.Itoa(i)
		institutionAsBytes,_ := APIstub.GetState(Ins_id)
		institution := Institution{}
		json.Unmarshal(institutionAsBytes, &institution)
		if institution.Authority == "0" {
			continue
		}
		var creditscore int
		var timecoin int
		creditscore,_ = strconv.Atoi(institution.Creditscore)
		timecoin = 10*creditscore
		institution.Asset = strconv.Itoa(timecoin)
		institutionAsBytes,_ = json.Marshal(institution)
		APIstub.PutState(Ins_id,institutionAsBytes)
	}

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}


func (s *SmartContract) queryAllTasks(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "task0"
	endKey := "task" + strconv.Itoa(task_cnt)
	//endKey := "task999"

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
//{"Args":["acceptSingleTask","task0","person0","zpl"]} unique-id of a person, person_name
//when a single task is accepted, this function must be called!!!
func (s *SmartContract) acceptSingleTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	taskAsBytes, _ := APIstub.GetState(args[0])
	task := Task{}

	current_time := time.Now()
	timestamp := current_time.Unix()
    t_unix := time.Unix(timestamp, 0)
    t := t_unix.Format("2006-01-02 03:04:05 PM")


	json.Unmarshal(taskAsBytes, &task)
	task.Owner = args[1]
	task.OwnerName = args[2]
	task.Accepted = "1"
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

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}

//if a task is completed, then we change the ledger and its state
//{"Args":["completeSingleTask","task0","person0"]} unique-id of a person
func (s *SmartContract) completeSingleTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
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
	task.Completed = "1"
	task.CompletedTime = t
	Timecoin, err = strconv.Atoi(task.Timecoin)
	taskAsBytes, _ = json.Marshal(task)
	APIstub.PutState(args[0], taskAsBytes)
	
	if task.Tasktype == "person" { // add the publisher's creditscore
		publisher := task.Publisher
		publisherAsBytes,_:= APIstub.GetState(publisher)
		publisher_people := People{}
		json.Unmarshal(publisherAsBytes,&publisher_people)
		creditscore,_ := strconv.Atoi(string(publisher_people.Creditscore))
		if creditscore >= 100 && creditscore < 110 {
			creditscore = creditscore + 1
		}else if creditscore < 100 {
			creditscore = creditscore + 2
		}
		publisher_people.Creditscore = strconv.Itoa(creditscore)
		publisherAsBytes,_ = json.Marshal(publisher_people)
		APIstub.PutState(publisher,publisherAsBytes)
	} else {
		publisher := task.Publisher
		publisherAsBytes,_:= APIstub.GetState(publisher)
		publisher_people := Institution{}
		json.Unmarshal(publisherAsBytes,&publisher_people)
		creditscore,_ := strconv.Atoi(string(publisher_people.Creditscore))
		if creditscore < 110 {
			creditscore = creditscore + 1
		}
		publisher_people.Creditscore = strconv.Itoa(creditscore)
		publisherAsBytes,_ = json.Marshal(publisher_people)
		APIstub.PutState(publisher,publisherAsBytes)
	}

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

	Creditscore, _ := strconv.Atoi(string(people.Creditscore))

	Timecoin = Timecoin * Creditscore / 100
	Aval = Aval + Timecoin
	Aval = int(Aval)

	if Creditscore >= 100 && Creditscore < 110{
		Creditscore = Creditscore + 1
	}else if Creditscore < 100{
		Creditscore = Creditscore + 2
	}
	
	people.Asset = strconv.Itoa(Aval)
	people.Creditscore = strconv.Itoa(Creditscore)
	
	if people.CompletedTask  == "None" {
		people.CompletedTask = args[0]
	}else{
		people.CompletedTask = people.CompletedTask + "," + args[0]
	}
	peopleAsBytes, _ = json.Marshal(people)
	APIstub.PutState(args[1], peopleAsBytes)

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}

// {"Args":["recordDisputedTask","task0"]}
func (s *SmartContract) recordDisputedTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	taskAsBytes, _ := APIstub.GetState(args[0])
	task := Task{}
	json.Unmarshal(taskAsBytes, &task)
	publisher := task.Publisher
	owner := task.Owner
	Timecoin,_:= strconv.Atoi(task.Timecoin)
	penalty := Timecoin / 2
	penalty = int(penalty)

	publisherAsBytes,_:= APIstub.GetState(publisher)
	ownerAsBytes,_:= APIstub.GetState(owner)
	
	if task.Tasktype == "institution" {
		publisher_institution := Institution{}
		json.Unmarshal(publisherAsBytes, &publisher_institution)
		if publisher_institution.Disputedtask == "None" {
			publisher_institution.Disputedtask = task.Id
		}else {
			publisher_institution.Disputedtask = publisher_institution.Disputedtask + "," + task.Id		
		}
		publisher_creditscore,_ := strconv.Atoi(string(publisher_institution.Creditscore))
		if (publisher_creditscore - 5) <= 0 {
			publisher_creditscore = 0
			publisher_institution.Blacklist = "1"
		}else {
			publisher_creditscore = publisher_creditscore - 5
		}
		publisher_institution.Creditscore = strconv.Itoa(publisher_creditscore)
		publisherAsBytes,_ = json.Marshal(publisher_institution)
		APIstub.PutState(task.Publisher,publisherAsBytes)
	}else{
		publisher_people := People{}
		json.Unmarshal(publisherAsBytes, &publisher_people)

		if publisher_people.Disputedtask == "None" {
			publisher_people.Disputedtask = task.Id
		}else {
			publisher_people.Disputedtask = publisher_people.Disputedtask + "," + task.Id		
		}
		
		publisher_creditscore,_ := strconv.Atoi(string(publisher_people.Creditscore))

		if (publisher_creditscore - 5) <= 0 {
			publisher_creditscore = 0
			publisher_people.Blacklist = "1"
		}else {
			publisher_creditscore = publisher_creditscore - 5
		}
		publisher_people.Creditscore = strconv.Itoa(publisher_creditscore)
		publisherAsBytes,_ = json.Marshal(publisher_people)
		APIstub.PutState(task.Publisher,publisherAsBytes)
	}

	owner_people := People{}
	json.Unmarshal(ownerAsBytes, &owner_people)
	if owner_people.Disputedtask == "None" {
		owner_people.Disputedtask = task.Id
	}else {
		owner_people.Disputedtask = owner_people.Disputedtask + "," + task.Id
	}
	owner_creditscore,_ := strconv.Atoi(string(owner_people.Creditscore))
	if (owner_creditscore - 5) <= 0 {
		owner_creditscore = 0
		owner_people.Blacklist = "1"
	}else {
		owner_creditscore = owner_creditscore - 5
	}
	owner_people.Creditscore = strconv.Itoa(owner_creditscore)
	//write to the ledger	
	ownerAsBytes,_ = json.Marshal(owner_people)
	APIstub.PutState(task.Owner,ownerAsBytes)

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}

// {"Args":["acceptGroupTask","task0","person0","zpl"]}
func (s *SmartContract) acceptGroupTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	taskAsBytes, _ := APIstub.GetState(args[0])
	task := GroupTask{}

	current_time := time.Now()
	timestamp := current_time.Unix()
    t_unix := time.Unix(timestamp, 0)
    t := t_unix.Format("2006-01-02 03:04:05 PM")

	json.Unmarshal(taskAsBytes, &task)
	var available int
	available,_ = strconv.Atoi(task.AvailableNumber)
	if available < 1{
		return shim.Error("Enough People!")
	}else {
		available -= 1
	}
	task.AvailableNumber = strconv.Itoa(available)

	if task.Owner == "None" {
		task.Owner = args[1]
		task.OwnerName = args[2]
		task.AcceptedTime = t
	}else {
		task.Owner = task.Owner + " " + args[1]
		task.OwnerName = task.OwnerName + "," + args[2]
	}
	
	if available == 0 {
		task.StartTime = t
	}

	task.Accepted = "1"
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

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}

// {"Args":["completeGroupTask","task0"]}
func (s *SmartContract) completeGroupTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	taskAsBytes, _ := APIstub.GetState(args[0])
	task := GroupTask{}

	current_time := time.Now()
	timestamp := current_time.Unix()
    t_unix := time.Unix(timestamp, 0)
    t := t_unix.Format("2006-01-02 03:04:05 PM")

	json.Unmarshal(taskAsBytes, &task)

	task.Completed = "1"
	task.CompletedTime = t
	Timecoin, _ := strconv.Atoi(task.Timecoin)
	taskAsBytes, _ = json.Marshal(task)
	APIstub.PutState(args[0], taskAsBytes)
	
	if task.Tasktype == "person" { // add the publisher's creditscore
		publisher := task.Publisher
		publisherAsBytes,_:= APIstub.GetState(publisher)
		publisher_people := People{}
		json.Unmarshal(publisherAsBytes,&publisher_people)
		creditscore,_ := strconv.Atoi(string(publisher_people.Creditscore))
		if creditscore >= 100 && creditscore < 110 {
			creditscore = creditscore + 1
		}else if creditscore < 100 {
			creditscore = creditscore + 2
		}
		publisher_people.Creditscore = strconv.Itoa(creditscore)
		publisherAsBytes,_ = json.Marshal(publisher_people)
		APIstub.PutState(publisher,publisherAsBytes)
	} else {
		publisher := task.Publisher
		publisherAsBytes,_:= APIstub.GetState(publisher)
		publisher_people := Institution{}
		json.Unmarshal(publisherAsBytes,&publisher_people)
		creditscore,_ := strconv.Atoi(string(publisher_people.Creditscore))
		if creditscore < 110 {
			creditscore = creditscore + 1
		}
		publisher_people.Creditscore = strconv.Itoa(creditscore)
		publisherAsBytes,_ = json.Marshal(publisher_people)
		APIstub.PutState(publisher,publisherAsBytes)
	}

	needed_people,_ := strconv.Atoi(task.Needpeople)
	reg := regexp.MustCompile(`person\d{1,}`)
	all_id := reg.FindAllString(task.Owner, -1)

	for i:= 0; i < needed_people; i++ {
		peopleAsBytes, err := APIstub.GetState(all_id[i])
		if err != nil {
			return shim.Error("Failed to get state")
		}
		if peopleAsBytes == nil {
			return shim.Error("Entity not found")
		}
		people := People{}
		json.Unmarshal(peopleAsBytes, &people)
		Aval, _ := strconv.Atoi(string(people.Asset))
		Creditscore, _ := strconv.Atoi(string(people.Creditscore))

		TimeValue := (Timecoin * Creditscore / (100*needed_people))
		Aval = Aval + TimeValue
		Aval = int(Aval)

		if Creditscore >= 100 && Creditscore < 110{
			Creditscore = Creditscore + 1
		}else if Creditscore < 100{
			Creditscore = Creditscore + 2
		}
		
		people.Asset = strconv.Itoa(Aval)
		people.Creditscore = strconv.Itoa(Creditscore)
		
		if people.CompletedTask  == "None" {
			people.CompletedTask = args[0]
		}else{
			people.CompletedTask = people.CompletedTask + "," + args[0]
		}
		peopleAsBytes, _ = json.Marshal(people)
		APIstub.PutState(all_id[i], peopleAsBytes)		
	}

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}

//{"Args":["recordDisputedGroupTask","task0"]}
//
func (s *SmartContract) recordDisputedGroupTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	taskAsBytes, _ := APIstub.GetState(args[0])
	task := GroupTask{}
	json.Unmarshal(taskAsBytes, &task)
	publisher := task.Publisher

	Timecoin,_:= strconv.Atoi(task.Timecoin)
	penalty := Timecoin / 2
	penalty = int(penalty)

	publisherAsBytes,_:= APIstub.GetState(publisher)
	
	if task.Tasktype == "institution" {
		publisher_institution := Institution{}
		json.Unmarshal(publisherAsBytes, &publisher_institution)
		if publisher_institution.Disputedtask == "None" {
			publisher_institution.Disputedtask = task.Id
		}else {
			publisher_institution.Disputedtask = publisher_institution.Disputedtask + "," + task.Id		
		}
		publisher_creditscore,_ := strconv.Atoi(string(publisher_institution.Creditscore))
		if (publisher_creditscore - 5) <= 0 {
			publisher_creditscore = 0
			publisher_institution.Blacklist = "1"
		}else {
			publisher_creditscore = publisher_creditscore - 5
		}
		publisher_institution.Creditscore = strconv.Itoa(publisher_creditscore)
		publisherAsBytes,_ = json.Marshal(publisher_institution)
		APIstub.PutState(task.Publisher,publisherAsBytes)
	}else{
		publisher_people := People{}
		json.Unmarshal(publisherAsBytes, &publisher_people)

		if publisher_people.Disputedtask == "None" {
			publisher_people.Disputedtask = task.Id
		}else {
			publisher_people.Disputedtask = publisher_people.Disputedtask + "," + task.Id		
		}

		publisher_creditscore,_ := strconv.Atoi(string(publisher_people.Creditscore))

		if (publisher_creditscore - 5) <= 0 {
			publisher_creditscore = 0
			publisher_people.Blacklist = "1"
		}else {
			publisher_creditscore = publisher_creditscore - 5
		}
		publisher_people.Creditscore = strconv.Itoa(publisher_creditscore)
		publisherAsBytes,_ = json.Marshal(publisher_people)
		APIstub.PutState(task.Publisher,publisherAsBytes)
	}

	reg := regexp.MustCompile(`person\d{1,}`)
	all_id := reg.FindAllString(task.Owner, -1)
	for i:= 0; i < len(all_id); i ++ {
		ownerAsBytes,_:= APIstub.GetState(all_id[i])
		owner_people := People{}
		json.Unmarshal(ownerAsBytes, &owner_people)
		if owner_people.Disputedtask == "None" {
			owner_people.Disputedtask = task.Id
		}else {
			owner_people.Disputedtask = owner_people.Disputedtask + "," + task.Id
		}
		owner_creditscore,_ := strconv.Atoi(string(owner_people.Creditscore))
		if (owner_creditscore - 5) <= 0 {
			owner_creditscore = 0
			owner_people.Blacklist = "1"
		}else {
			owner_creditscore = owner_creditscore - 5
		}
		owner_people.Creditscore = strconv.Itoa(owner_creditscore)
		//write to the ledger	
		ownerAsBytes,_ = json.Marshal(owner_people)
		APIstub.PutState(all_id[i],ownerAsBytes)		
	}

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}


//{"Args":["cancelTask","task0","group","institution"]} group/single  publisher_type
func (s *SmartContract) cancelTask(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	taskAsBytes, _ := APIstub.GetState(args[0])
	var TimeCoin int
	var publisher string

	if args[1] == "group" {
		task := GroupTask{}
		json.Unmarshal(taskAsBytes, &task)
		publisher = task.Publisher
		TimeCoin,_ = strconv.Atoi(task.Timecoin)
		task.Canceled = "1"
		taskAsBytes,_ = json.Marshal(task)
		APIstub.PutState(args[0],taskAsBytes)

	}else {
		task := Task{}
		json.Unmarshal(taskAsBytes, &task)
		publisher = task.Publisher
		TimeCoin,_ = strconv.Atoi(task.Timecoin)
		task.Canceled = "1"
		taskAsBytes,_ = json.Marshal(task)
		APIstub.PutState(args[0],taskAsBytes)
	}	
	
	publisherAsBytes, _ := APIstub.GetState(publisher)
	

	if args[2] == "institution" {
		publisher_people := Institution{}
		json.Unmarshal(publisherAsBytes, &publisher_people)
		Asset,_ := strconv.Atoi(publisher_people.Asset)
		Asset += TimeCoin
		publisher_people.Asset = strconv.Itoa(Asset)

		publisherAsBytes,_ = json.Marshal(publisher_people)
		APIstub.PutState(publisher,publisherAsBytes)
	}else {
		publisher_people := People{}
		json.Unmarshal(publisherAsBytes, &publisher_people)
		Asset,_ := strconv.Atoi(publisher_people.Asset)
		Asset += TimeCoin
		publisher_people.Asset = strconv.Itoa(Asset)

		publisherAsBytes,_ = json.Marshal(publisher_people)
		APIstub.PutState(publisher,publisherAsBytes)
	}

	transaction := APIstub.GetTxID()
	transactionAsBytes, _ := json.Marshal(transaction)
	return shim.Success(transactionAsBytes)
}


// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
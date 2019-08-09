package main



// Step0.引用套件
import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"

)

type User struct {
	UserID string `json:"userid"`
	Userval int `json:"userval"`
	Times int `json:"times"`
	Moviename string `json:"moviename"`
}

type MCompany struct {
	Cname string `json:"cname"`
	MCompanyval int `json:"mcompanyval"`
	Movie string `json:"moviename"`
}

// Step1.定義資產
type SampleAsset struct{


}

// Step2.定義合約
type SampleChaincode struct {


}
var movielist []string

// Step3.定義合約可以操作的方法

// Function-1. 插入新User進入區塊鏈
// INPUT PATTERN(userid,userval)
func (t *SampleChaincode) adduser(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	a,_ := strconv.Atoi(args[1])
	var Userid= User{UserID: args[0], Userval: a, Times: 0}
	Userbytes, _ := json.Marshal(Userid)
	err := stub.PutState(args[0], Userbytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to add a user: %s", args[0]))
	}
	return shim.Success(nil)
}



// Function-2 Transaction makes payment of X units from A to B
// Input pattern (User,Mcompany,Moviename,Mprice)
func (t *SampleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var Userid, MCompanyid, Moviename string
	var Mprice int                    // Transaction value
	var err error


	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}


	var user = User{}
	var company = MCompany{}
	Userid = args[0]
	MCompanyid = args[1]
	Moviename = args[2]
	Mprice,_ = strconv.Atoi(args[3])
	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Userbytes, err := stub.GetState(Userid)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Userbytes == nil {
		return shim.Error("Entity not found")
	}
	err = json.Unmarshal(Userbytes, &user)
	if err != nil {
		fmt.Println("error:", err)
	}

	MCompanybytes, err := stub.GetState(MCompanyid)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if MCompanybytes == nil {
		return shim.Error("Entity not found")
	}

	err = json.Unmarshal(MCompanybytes, &company)
	if err != nil {
		fmt.Println("error:", err)
	}
	// Perform the execution
	Mprice, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	user.Userval = user.Userval  - Mprice
	user.Moviename = args[2]
	company.MCompanyval = company.MCompanyval + Mprice
	company.Movie = args[2]
	fmt.Printf("Userval = %d, MCompanyval = %d\n", user.Userval, company.MCompanyval)
	movielist = append(movielist,Moviename)
	user.Times = user.Times + 1
	//Giving $25 for every 3 times
	if user.Times % 3 == 0{
		user.Userval = user.Userval + 25
	}

	Userbytes1, _ := json.Marshal(user)
	err = stub.PutState(args[0], Userbytes1)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to update user: %s", args[0]))
	}

	MCompanybytes1, _ := json.Marshal(company)
	err = stub.PutState(args[1], MCompanybytes1)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to update company: %s", args[0]))
	}

	return shim.Success(nil)

}

// Function 3 Query user info
//Input pattern (userid)
func (t *SampleChaincode) quser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	Userbytes, _ := stub.GetState(args[0])

	if Userbytes == nil {
		return shim.Error("User not found")
	}
	return shim.Success(Userbytes)
}
// Function 4 Query movielist_history
//Input pattern moviecompanyid
func (t *SampleChaincode) qmoviellist(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	var MCompanyid string
	var company = MCompany{}
	MCompanyid = args[0]
	MCompanybytes, err := stub.GetState(MCompanyid)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if MCompanybytes == nil {
		return shim.Error("Entity not found")
	}

	err = json.Unmarshal(MCompanybytes, &company)
	if err != nil {
		fmt.Println("error:", err)
	}


	return shim.Success(nil)
}



// Step4.定義Init方法 INPUT PATTERN(Cname,MCompanyval)
func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Initializing")
	_, args := stub.GetFunctionAndParameters()
	var err error


	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	// Initialize the chaincode
	a,_ := strconv.Atoi(args[1])

	var MCompany= MCompany{Cname: args[0], MCompanyval: a, Movie:"" }

	MCompanybytes, _ := json.Marshal(MCompany)
	err = stub.PutState(args[0], MCompanybytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to init a company: %s", args[0]))
	}
	return shim.Success(nil)
}


// Step5.定義Invoke方法
func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	function, args := stub.GetFunctionAndParameters()

	if function == "invoke" {
		return t.invoke(stub, args)
	} else if function == "adduser" {
		return t.adduser(stub, args)
	}	else if function == "quser" {
		return t.quser(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}



// Step6.定義main方法
func main() {
	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
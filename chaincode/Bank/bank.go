package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	/* "encoding/gob"
	"bytes" */
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type bankInfo struct {
	BankName              string `json:"BankName"`
	BankBranch            string `json:"BankBranch"`
	Bankcode              string `json:"BankCode"`
	BankWalletID          string `json:"MainWallet"`      //will take the values for the respective wallet from the user
	BankAssetWalletID     string `json:"AssetWallet"`     //will take the values for the respective wallet from the user
	BankChargesWalletID   string `json:"ChargesWallet"`   //will take the values for the respective wallet from the user
	BankLiabilityWalletID string `json:"LiabilityWallet"` //will take the values for the respective wallet from the user
	TDSreceivableWalletID string `json:"TDSWallet"`       //will take the values for the respective wallet from the user
}
type bankInfoVal struct {
	BankName              string `json:"BankName"`
	BankWalletID          string `json:"MainWallet"`      //will take the values for the respective wallet from the user
	BankWalletIDBal       string `json:"MainWalletBal"`      //will take the values for the respective wallet from the user
	BankAssetWalletID     string `json:"AssetWallet"`     //will take the values for the respective wallet from the user
	BankAssetWalletIDBal  string `json:"AssetWalletBal"`     //will take the values for the respective wallet from the user
	BankChargesWalletID   string `json:"ChargesWallet"`   //will take the values for the respective wallet from the user
	BankChargesWalletIDBal string `json:"ChargesWalletBal"`   //will take the values for the respective wallet from the user
	BankLiabilityWalletID string `json:"LiabilityWallet"` //will take the values for the respective wallet from the user
	BankLiabilityWalletIDBal string `json:"LiabilityWalletBal"` //will take the values for the respective wallet from the user
	TDSreceivableWalletID string `json:"TDSWallet"`       //will take the values for the respective wallet from the user
	TDSreceivableWalletIDBal string `json:"TDSWalletBal"`       //will take the values for the respective wallet from the user
}

var jsonresp string
var err1 error
var jsonresp1 []string
var err2 error
// toChaincodeArgs returns byte arrau of string of arguments, so it can be passed to other chaincodes
func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	bank := bankInfo{}
	indexName := "Bankcode~BankBranch"
	codeBranchKey, err := stub.CreateCompositeKey(indexName, []string{bank.Bankcode, bank.BankBranch})
	if err != nil {
		return shim.Error("Unable to create composite key Bankcode~BankBranch in bankcc")
	}
	value := []byte{0x00}
	stub.PutState(codeBranchKey, value)
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "writeBankInfo" {
		//Creates a new Bank Information
		return writeBankInfo(stub, args)
	} else if function == "getBankInfo" {
		//Retrieves the Bank information
		return getBankInfo(stub, args)
	}  else if function == "getWalletID" {
		//Returns the walletID for the required wallet type
		return getWalletID(stub, args)
	} else if function == "bankIDexists" {
		//To check the BankId existence
		return bankIDexists(stub, args[0])  
	} else if function == "getWalletsofBank" {
		//To check the BankId existence
		//return getWalletsofBank(stub, []string{"main","asset","charges"})
		return getWalletsofBank(stub, args)

	} else if function == "getWBank" {
		//Returns the walletID for the required bank wallet type
		//return getWBank(stub, args)
		jsonresp,err1 := getWBank(stub, args)
		if err1 != nil {
			return shim.Error(err1.Error())
		} else {
			//return shim.Success([]byte(strings.Join(jsonresp,"")))
			return shim.Success([]byte(jsonresp))
		}
	} else if function == "getBankInfoBalTemp" {
		//Returns the walletID for the required bank wallet type
		//return getWBank(stub, args)
		jsonrespasstruct,err1 := getBankInfoBalTemp(stub, args)
		if err1 != nil {
			return shim.Error(err1.Error())
		} else {
			//return shim.Success([]byte(strings.Join(jsonresp,"")))
			return shim.Success([]byte(jsonrespasstruct))
		}
	}
	return shim.Error("No function named " + function + " in Banksssssssss")

}

// To write the bank info into the ledger
func writeBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("********************While Writing Bank*************************")
	//Checking argument length
	if len(args) != 9 {
		xLenStr := strconv.Itoa(len(args)) //needed?!
		return shim.Error("Invalid number of arguments in writeBankInfo (required:9) given:" + xLenStr)
	}
	fmt.Println("args[0]", args[0])
	fmt.Println("args[1]", args[1])
	fmt.Println("args[2]", args[2])
	fmt.Println("args[3]", args[3])
	fmt.Println("args[4]", args[5])
	fmt.Println("args[5]", args[5])
	fmt.Println("args[6]", args[6])
	fmt.Println("args[7]" ,args[7])
	fmt.Println("args[8] ",args[8])


	//Checking Bank ID existence
	response := bankIDexists(stub, args[0])
	if response.Status != shim.OK {
		return shim.Error("bankcc: " + response.Message)
	}

	//Checking existence of Bank code (!*!)
	codeBranchIterator, err := stub.GetStateByPartialCompositeKey(" Bankcode~BankBranch ", []string{args[3]})
	codeBranchData, err := codeBranchIterator.Next()
	if codeBranchData != nil {
		return shim.Error("bankcc : " + "Bank code already exist: " + args[3])
	}
	defer codeBranchIterator.Close()
	indexName := "Bankcode~BankBranch"
	codeBranchKey, err := stub.CreateCompositeKey(indexName, []string{args[3],args[2]})
	if err != nil {
		return shim.Error("Unable to create composite key Bankcode~BankBranch in bankcc")
	}
	value := []byte{0x00}
	stub.PutState(codeBranchKey, value)
	hash := sha256.New()

	//TODO:	check for collisions

	// Hashing bankWalletId
	BankWalletStr := args[3] + "BankWallet"
	hash.Write([]byte(BankWalletStr))
	md := hash.Sum(nil)
	BankWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BankWalletIDsha ",BankWalletIDsha)
	createWallet(stub, BankWalletIDsha, args[4])
	

	// Hashing bankAssetWalletId
	BankAssetWalletStr := args[3] + "BankAssetWallet"
	hash.Write([]byte(BankAssetWalletStr))
	md = hash.Sum(nil)
	BankAssetWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BankAssetWalletIDsha ",BankAssetWalletIDsha)
	createWallet(stub, BankAssetWalletIDsha, args[5])

	// Hashing BankChargesWalletID
	BankChargesWalletStr := args[3] + "BankChargesWallet"
	hash.Write([]byte(BankChargesWalletStr))
	md = hash.Sum(nil)
	BankChargesWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BankChargesWalletIDsha ",BankChargesWalletIDsha)
	createWallet(stub, BankChargesWalletIDsha, args[6])

	// Hashing BankLiabilityWalletID
	BankLiabilityWalletStr := args[3] + "BankLiabilityWallet"
	hash.Write([]byte(BankLiabilityWalletStr))
	md = hash.Sum(nil)
	BankLiabilityWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BankLiabilityWalletIDsha ",BankLiabilityWalletIDsha)
	createWallet(stub, BankLiabilityWalletIDsha, args[7])

	// Hashing TDSreceivableWalletID
	TDSreceivableWalletStr := args[3] + "TDSreceivableWallet"
	hash.Write([]byte(TDSreceivableWalletStr))
	md = hash.Sum(nil)
	TDSreceivableWalletIDsha := hex.EncodeToString(md)
	fmt.Print("TDSreceivableWalletIDsha ",TDSreceivableWalletIDsha)
	createWallet(stub, TDSreceivableWalletIDsha,  args[8])

	//args[0] -> bankID | creating a bank struct obj and writing it to the ledger
	bank := bankInfo{args[1], args[2], args[3], BankWalletIDsha, BankAssetWalletIDsha, BankChargesWalletIDsha, BankLiabilityWalletIDsha, TDSreceivableWalletIDsha}
	bankBytes, err := json.Marshal(bank)
	if err != nil {
		return shim.Error("Unable to Marshal the json file " + err.Error())
	}

	err = stub.PutState(args[0], bankBytes)
	if err != nil {
		return shim.Error("bankcc : " + err.Error())
	}
	fmt.Println("********************End Writing *************************")
	return shim.Success([]byte("Succefully written into the ledger"))
}

func bankIDexists(stub shim.ChaincodeStubInterface, bankID string) pb.Response {
	ifExists, _ := stub.GetState(bankID)
	if ifExists != nil {
		fmt.Println(ifExists) //needed!?
		return shim.Error("bankcc : " + "BankId " + bankID + " exits. Cannot create new ID")
	}
	return shim.Success(nil)
}

func createWallet(stub shim.ChaincodeStubInterface, walletID string, amt string) pb.Response {
	//Calling wallet Chaincode to create new wallet
	chaincodeArgs := toChaincodeArgs("newWallet", walletID, amt)
	response := stub.InvokeChaincode("walletcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("bankcc : " + "Unable to create new wallet from bank")
	}
	return shim.Success([]byte("created new wallet from bank"))
}

func getBankInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("bankcc : " + "Invalid number of arguments in getBankInfo (required:1) given:" + xLenStr)
	}

	// getting bank info from the ledger, args[0] -> bankID
	bankInfoBytes, err := stub.GetState(args[0])

	if err != nil {
		return shim.Error("bankcc : " + "Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		return shim.Error("bankcc : " + "Data does not exist for " + args[0])
	}

	// unmarshalling bankInfoBytes into bank structure
	bank := bankInfo{}
	
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		return shim.Error("bankcc : " + "Uable to paser into the json format")
	}
	x := fmt.Sprintf("%+v", bank)
	fmt.Println("BankInfo : %s\n", x)
	//fmt.Println("")
	return shim.Success([]byte(x))
	
}
/* 
func getBankInfoTemp(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("bankcc : " + "Invalid number of arguments in getBankInfo (required:1) given:" + xLenStr)
	}

	// getting bank info from the ledger, args[0] -> bankID
	bankInfoBytes, err := stub.GetState(args[0])

	if err != nil {
		return shim.Error("bankcc : " + "Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		return shim.Error("bankcc : " + "Data does not exist for " + args[0])
	}

	// unmarshalling bankInfoBytes into bank structure
	bank := bankInfoVal{}
	
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		return shim.Error("bankcc : " + "Uable to paser into the json format")
	}
	x := fmt.Sprintf("%+v", bank)
	fmt.Println("BankInfo : %s\n", x)
	//fmt.Println("")
	return shim.Success([]byte(x))
	
} */
//
func getBankInfoBalTemp(stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {
	fmt.Println("********************Start getBankInfoBalTemp *************************")
		bankvalasstruct := bankInfoVal{}
		bank, err := stub.GetState(args[0])
		if err != nil {
				return nil,fmt.Errorf("bankcc : " + "Unable to fetch the state" + err.Error())
		}
		if bank == nil {
				return nil,fmt.Errorf("bankcc : " + "Data does not exist for " + args[0])
		}	
		err = json.Unmarshal(bank, &bankvalasstruct)
		if err != nil {
				return nil,fmt.Errorf("bankcc : " + "Uable to paser into the json format")
		}	
		data, err := json.Marshal(bankvalasstruct)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Println("Returning value of getBankInfoBalTemp  %s",data)
		fmt.Println("********************Start getBankInfoBalTemp *************************")
	 return data,nil 
	}

/* func getBankInfoBal(stub shim.ChaincodeStubInterface, args []string) (string,error) {
/* 
	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		//return shim.Error("bankcc : " + "Invalid number of arguments in getBankInfo (required:1) given:" + xLenStr)
		return "",fmt.Errorf("bankcc : " + "Invalid number of arguments in getBankInfo (required:1) given:" + xLenStr)
	}
 
	// getting bank info from the ledger, args[0] -> bankID
	bankInfoBytes, err := stub.GetState(args[0])

	if err != nil {
		//return shim.Error("bankcc : " + "Unable to fetch the state" + err.Error())
		return "",fmt.Errorf("bankcc : " + "Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		//return shim.Error("bankcc : " + "Data does not exist for " + args[0])
		return "",fmt.Errorf("bankcc : " + "Data does not exist for " + args[0])
	}

	// unmarshalling bankInfoBytes into bank structure
	bank := bankInfoVal{}
	//
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		//return shim.Error("bankcc : " + "Uable to paser into the json format")
		return "",fmt.Errorf("bankcc : " + "Uable to paser into the json format")
	}
	
	x := fmt.Sprintf("%+v", bank)

	fmt.Println("BankInfoBankInfo : %s\n", x)
	

	return string(bankInfoBytes),nil
	
} */
//func getWalletsofBank(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	func getWalletsofBank(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	//var i int;
	//var resp []string;
	//var allResponse []byte;
	fmt.Println("********************While Reading Bank getWalletsofBank *************************")
	fmt.Println("args []",args [0], " ",args [1], " ",args [2], " ",args [3], " ",args [4], " ",args [5])
	var response1 pb.Response  ;
	var response2 pb.Response  ;
	var response3 pb.Response  ;
	var response4 pb.Response  ;
	var response5 pb.Response  ;

	var intarg1 []string;
	var intarg2 []string;
	var intarg3 []string;
	var intarg4 []string;
	var intarg5 []string;

	intarg1=append(intarg1,args[0])
	intarg1=append(intarg1,args[1])

	intarg2=append(intarg2,args[0])
	intarg2=append(intarg2,args[2])

	intarg3=append(intarg3,args[0])
	intarg3=append(intarg3,args[3])

	intarg4=append(intarg4,args[0])
	intarg4=append(intarg4,args[4])

	intarg5=append(intarg5,args[0])
	intarg5=append(intarg5,args[5])

	
	var comres []string;  
	var bankValues []byte;
	
	mainWalletId,err1 := getWBank(stub ,intarg1)
	assetWalletId,err1 := getWBank(stub ,intarg2)
	chargesWalletId,err1 := getWBank(stub ,intarg3)
	liabilityWalletId,err1 := getWBank(stub ,intarg4)
	tdsWalletId,err1 := getWBank(stub ,intarg5)
	if err1 != nil {
		//return shim.Error("Unable to get Bank Info ")
	//	return []string{""},fmt.Errorf("Unable to get Wallet Id")
	//return "",fmt.Errorf("Unable to get Wallet Id")
	}
	bankvalasstruct := bankInfoVal{}
	bankValues,err56 := getBankInfoBalTemp(stub, args)
	if err56 != nil {
		return shim.Error("Cant't get value for Bank"+ err56.Error())
    }  
	 err := json.Unmarshal(bankValues, &bankvalasstruct)
	if err != nil {
		return shim.Error("error while unmarshalling bankValues"+ err.Error())
    }  
	data, err1 := json.Marshal(bankvalasstruct)
	if err1 != nil {
		return shim.Error("error while marshaling  bankvalasstruct" + err1.Error() )
	} 
	fmt.Println("Before appending balance to json %s",data)


	chaincodeArgs1 := toChaincodeArgs("getWallet", mainWalletId)
	response1 = stub.InvokeChaincode("walletcc", chaincodeArgs1, "myc")
	comres=append(comres,string(response1.Payload))
	bankvalasstruct.BankWalletIDBal = string(response1.Payload);
	fmt.Println("mainWalletId balance %s",string(response1.Payload))
		
	chaincodeArgs2 := toChaincodeArgs("getWallet", assetWalletId)
	response2 = stub.InvokeChaincode("walletcc", chaincodeArgs2, "myc")
	comres=append(comres,string(response2.Payload))
	bankvalasstruct.BankAssetWalletIDBal=string(response2.Payload);
	fmt.Println("assetwallet balance %s",string(response2.Payload))
	/* slice1 := make([]byte, 0, len(comres)+1)
		comresAsByteArray1 := response2.Payload
		comres = append(slice1, byte(len(comresAsByteArray1))); */
	//bankValues.BankAssetWalletIDBal=string(response2.Payload);
	/* slice1 := make([]byte, 0, len(comres)+1)
		comresAsByteArray1 := response2.Payload
		comres = append(slice1, byte(len(comresAsByteArray1))); */
		
	chaincodeArgs3 := toChaincodeArgs("getWallet", chargesWalletId)
	response3 = stub.InvokeChaincode("walletcc", chaincodeArgs3, "myc")
	fmt.Println("chargesWalletId balance %s",string(response3.Payload))
	comres=append(comres,string(response3.Payload))
	bankvalasstruct.BankChargesWalletIDBal = string(response3.Payload);
	
	//bankValues.BankChargesWalletIDBal=string(response3.Payload);
	/* slice2 := make([]byte, 0, len(comres)+1)
		comresAsByteArray2 := response3.Payload
		comres = append(slice2, byte(len(comresAsByteArray2))); */
		
	chaincodeArgs4 := toChaincodeArgs("getWallet", liabilityWalletId)
	response4 = stub.InvokeChaincode("walletcc", chaincodeArgs4, "myc")
	fmt.Println("liabilityWallet balance %s",string(response4.Payload))
	comres=append(comres,string(response4.Payload))
	bankvalasstruct.BankLiabilityWalletIDBal = string(response4.Payload);
	//bankValues.BankLiabilityWalletIDBal=string(response4.Payload);
	/* slice3 := make([]byte, 0, len(comres)+1)
		comresAsByteArray3 := response4.Payload
		comres = append(slice3, byte(len(comresAsByteArray3))); */
			
	chaincodeArgs5 := toChaincodeArgs("getWallet", tdsWalletId)
	response5 = stub.InvokeChaincode("walletcc", chaincodeArgs5, "myc")
	comres=append(comres,string(response5.Payload))
	bankvalasstruct.TDSreceivableWalletIDBal = string(response5.Payload);
	//bankValues.TDSreceivableWalletIDBal=string(response5.Payload);
	/* slice4 := make([]byte, 0, len(comres)+1)
	comresAsByteArray4 := response5.Payload */
	
	fmt.Println("tdsWalletId balance %s",string(response5.Payload))

//	bankvalasstruct.BankChargesWalletIDBal = string(response3.Payload);
	data, err1 = json.Marshal(bankvalasstruct)
	if err1 != nil {
		return shim.Error("error while marshaling  bankvalasstruct" + err1.Error() )
	} 
	fmt.Println("comres  %s",comres);
	fmt.Println("Response returning from getWalletsofBank %s",data)
		if response1.Status  != shim.OK {
			return shim.Error("bankcc : " + "Unable to get wallet of bank")
	}
	fmt.Println("********************End Reading *************************")
	return shim.Success(data)
}

func getWalletID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return  shim.Error("bankcc : " + "Invalid number of arguments in getWalletID(bank) (required:2) given:" + xLenStr)
	}
	bankInfoBytes, err := stub.GetState(args[0])
	if err != nil {
		return  shim.Error("bankcc : " + "Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		return shim.Error("bankcc : " + "Data does not exist for " + args[0])
	}
	bank := bankInfo{}
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		return  shim.Error("bankcc : " + "Uable to paser into the json format ")
	}
	jsonString := fmt.Sprintf("%s", bank)
	fmt.Println("jsonString %s",jsonString)
	walletID := ""

	switch args[1] {
	case "main":
		walletID = bank.BankWalletID
	case "asset":
		walletID = bank.BankAssetWalletID
	case "charges":
		walletID = bank.BankChargesWalletID
	case "liability":
		walletID = bank.BankLiabilityWalletID
	case "tds":
		walletID = bank.TDSreceivableWalletID
	}
	return shim.Success([]byte(walletID))
}

func getWBank(stub shim.ChaincodeStubInterface, args []string) (string, error) {

fmt.Println("********************Start getWBank *************************")
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return "", fmt.Errorf("bankcc : " + "Invalid number of arguments in getWalletID(bank) (required:2) given:" + xLenStr)
	}
	bankInfoBytes, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("bankcc : " + "Unable to fetch the state" + err.Error())
	}
	if bankInfoBytes == nil {
		return "" ,fmt.Errorf("bankcc : " + "Data does not exist for " + args[0])
	}
	bank := bankInfo{}
	err = json.Unmarshal(bankInfoBytes, &bank)
	if err != nil {
		return "", fmt.Errorf("bankcc : " + "Uable to paser into the json format")
	}
	walletID := ""
	switch args[1] {
	case "main":
		walletID = bank.BankWalletID
	case "asset":
		walletID = bank.BankAssetWalletID
	case "charges":
		walletID = bank.BankChargesWalletID
	case "liability":
		walletID = bank.BankLiabilityWalletID
	case "tds":
		walletID = bank.TDSreceivableWalletID
	}
	
	fmt.Println("walletID based on case  %s " ,walletID)
	fmt.Println("********************End getWBank *************************")
	return string(walletID),nil;
}
func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("bankcc : "+"Error starting Bank chaincode: %s\n", err)
	}

}

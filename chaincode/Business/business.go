package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type businessInfo struct {
	BusinessName                         string `json:"Name"`
	BusinessAcNo                         string `json:"AcNo"`
	BusinessLimit                        int64  `json:"Limit"`
	BusinessWalletID                     string `json:"MainWallet"`      //will take the values for the respective wallet from the user
	BusinessLoanWalletID                 string `json:"LoanWallet"`      //will take the values for the respective wallet from the user
	BusinessLiabilityWalletID            string `json:"LiabilityWallet"` //will take the values for the respective wallet from the user
	MaxROI                               int64  `json:"MaxROI"`
	MinROI                               int64  `json:"MinROI"`
	BusinessPrincipalOutstandingWalletID string `json:"POsWallet"` //will take the values for the respective wallet from the user
	BusinessChargesOutstandingWalletID   string `json:"COsWallet"` //will take the values for the respective wallet from the user
	
}

type businessInfoVal struct {
	BusinessName                         string `json:"Name"`
	BusinessWalletID                     string `json:"MainWallet"`
	BusinessWalletBalance                string `json:"MainWalletBalance"`      
	BusinessLoanWalletID                 string `json:"LoanWallet"`     
	BusinessLoanWalletBalance            string `json:"LoanWalletBalance"` 
	BusinessLiabilityWalletID            string `json:"LiabilityWallet"` 
	BusinessLiabilityWalletBalance       string `json:"LiabilityWalletBalance"`
	MaxROI                               int64  `json:"MaxROI"`
	MinROI                               int64  `json:"MinROI"`
	BusinessPrincipalOutstandingWalletID string `json:"POsWallet"` 
	BusinessPrincipalOutstandingWalletBalance string `json:"POsWalletBalance"`
	BusinessChargesOutstandingWalletID   string `json:"COsWallet"` 
	BusinessChargesOutstandingWalletBalance   string `json:"COsWalletBalance"` 
}
// toChaincodeArgs returns byte array of string of arguments, so it can be passed to other chaincodes
func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	bis := businessInfo{}
	indexName := "BusinessAcNo~BusinessName"
	acntNoNameKey, err := stub.CreateCompositeKey(indexName, []string{bis.BusinessAcNo, bis.BusinessName})
	if err != nil {
		return shim.Error("businesscc: " + "Unable to create composite key BusinessAcNo~BusinessName in businesscc")
	}
	value := []byte{0x00}
	stub.PutState(acntNoNameKey, value)
	return shim.Success(nil)
}
func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "putNewBusinessInfo" {
		//Creates a new Business Information
		return putNewBusinessInfo(stub, args)
	} else if function == "getBusinessInfo" {
		//Retrieves the Business information
		return getBusinessInfo(stub, args)
	} else if function == "getWalletID" {
		//Returns the walletID for the required wallet type
		return getWalletID(stub, args)
	} else if function == "bisIDexists" {
		//To check the BusinessId existence
		return bisIDexists(stub, args[0])
	} else if function == "bAccNoexists" {
		//To check the Acc No existence
		return bAccNoexists(stub, args[0])
	} else if function == "updateBusinessInfo" {
		//Updates Business Limit / MAX ROI / MAX ROI if required
		return updateBusinessInfo(stub, args)
	} else if function == "getWalletsofBusiness" {
		return getWalletsofBusiness(stub, args)
	} else if function == "getBusinessWallet" {
		jsonresp,err1 := getBusinessWallet(stub, args)
		if err1 != nil {
			return shim.Error(err1.Error())
		} else {
			return shim.Success([]byte(jsonresp))
		}		
	} else if function == "getBusinessInfoVal" {
		//Returns the walletID for the required bank wallet type
		//return getWBank(stub, args)
		jsonrespasstruct,err1 := getBusinessInfoVal(stub, args)
		if err1 != nil {
			return shim.Error(err1.Error())
		} else {
			//return shim.Success([]byte(strings.Join(jsonresp,"")))
			return shim.Success([]byte(jsonrespasstruct))
		}
	}
	return shim.Error("businesscc: " + "No function named " + function + " in Businessssssss")
}

func putNewBusinessInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************While Writing Business *************************")
	if len(args) != 11 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("businesscc: " + "Invalid number of arguments in putNewBusinessInfo (required:11) given:" + xLenStr)

	}
	fmt.Print("args[0] &[1]  &[2]", args[0] ," -----> ",args[1]," -----> ",args[2])
	fmt.Print("args[3] &[4]args[5] ", args[3] ," -----> ",args[4]," -----> ",args[5])
	fmt.Print("args[6] &[7]", args[6] ," -----> ",args[7]," -----> ",args[8]," -----> ",args[9])
	response := bisIDexists(stub, args[0])
	if response.Status != shim.OK {
		return shim.Error("businesscc: " + response.Message)
	}
	response1 := bAccNoexists(stub, args[0])
	if response1.Status != shim.OK {
		return shim.Error("businesscc: " + response1.Message)
	}
	businessLimitConv, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error("businesscc: " + err.Error())
	}
	if businessLimitConv <= 0 {
		return shim.Error("businesscc: " + "Invalid Business Limit value: " + args[3])
	}
	hash := sha256.New()
	// Hashing BusinessWalletID
	BusinessWalletStr := args[2] + "BusinessWallet"
	hash.Write([]byte(BusinessWalletStr))
	md := hash.Sum(nil)
	BusinessWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BusinessWalletIDsha ",BusinessWalletIDsha)
	createWallet(stub, BusinessWalletIDsha, args[6])
	
	// Hashing BusinessLoanWalletID
	BusinessLoanWalletStr := args[2] + "BusinessLoanWallet"
	hash.Write([]byte(BusinessLoanWalletStr))
	md = hash.Sum(nil)
		
	BusinessLoanWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BusinessLoanWalletIDsha ",BusinessLoanWalletIDsha)
	createWallet(stub, BusinessLoanWalletIDsha, args[7])
	// Hashing BusinessLiabilityWalletID
	BusinessLiabilityWalletStr := args[2] + "BusinessLiabilityWallet"
	hash.Write([]byte(BusinessLiabilityWalletStr))
	md = hash.Sum(nil)
	BusinessLiabilityWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BusinessLiabilityWalletIDsha ",BusinessLiabilityWalletIDsha)
	createWallet(stub, BusinessLiabilityWalletIDsha, args[8])

	maxROIconvertion, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		fmt.Println("Invalid Maximum ROI: %s\n", args[5])	
		return shim.Error("Businesscc: " + err.Error())
	}
	if maxROIconvertion <= 0 {
		return shim.Error(" Businesscc: " + " Invalid Max ROI value: " + args[5])
	}
	minROIconvertion, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		fmt.Println("Invalid Minimum ROI: %s\n", args[4])
		return shim.Error("Businesscc: " + err.Error())
	}
	if minROIconvertion <= 0 {
		return shim.Error("Businesscc: " + "Invalid Min ROI value: " + args[4])
	}
	// Hashing BusinessPrincipalOutstandingWalletID
	BusinessPrincipalOutstandingWalletStr := args[2] + "BusinessPrincipalOutstandingWallet"
	hash.Write([]byte(BusinessPrincipalOutstandingWalletStr))
	md = hash.Sum(nil)
	BusinessPrincipalOutstandingWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BusinessPrincipalOutstandingWalletIDsha ",BusinessPrincipalOutstandingWalletIDsha)
	createWallet(stub, BusinessPrincipalOutstandingWalletIDsha, args[9])
	// Hashing BusinessChargesOutstandingWalletID
	BusinessInterestOutstandingWalletStr := args[2] + "BusinessInterestOutstandingWallet"
	hash.Write([]byte(BusinessInterestOutstandingWalletStr))
	md = hash.Sum(nil)
	BusinessChargesOutstandingWalletIDsha := hex.EncodeToString(md)
	fmt.Print("BusinessChargesOutstandingWalletIDsha ",BusinessChargesOutstandingWalletIDsha)
	createWallet(stub, BusinessChargesOutstandingWalletIDsha, args[10])
	newInfo := &businessInfo{args[1], args[2], businessLimitConv, BusinessWalletIDsha, BusinessLoanWalletIDsha, BusinessLiabilityWalletIDsha, maxROIconvertion, minROIconvertion, BusinessPrincipalOutstandingWalletIDsha, BusinessChargesOutstandingWalletIDsha}
	newInfoBytes, _ := json.Marshal(newInfo)
	err = stub.PutState(args[0], newInfoBytes) // businessID = args[0]
	if err != nil {
		return shim.Error("businesscc: " + err.Error())
	}
	fmt.Println("Successfully added buissness " + args[1] + " to the ledger")
	fmt.Println("******************** End Writing Business *************************")
	return shim.Success([]byte("Successfully added buissness " + args[1] + " to the ledger"))
}

func createWallet(stub shim.ChaincodeStubInterface, walletID string, amt string) pb.Response {
	//Calling the wallet Chaincode to create new wallet
	chaincodeArgs := toChaincodeArgs("newWallet", walletID, amt)
	response := stub.InvokeChaincode("walletcc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("businesscc: " + "Unable to create new wallet from business")
	}
	return shim.Success([]byte("businesscc: " + "created new wallet from business"))
}

func getBusinessInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("businesscc: " + "Invalid number of arguments in getBusinessInfo (required:1) given:" + xLenStr)
	}
	parsedBusinessInfo := businessInfo{}
	businessIDvalue, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("businesscc: " + "Failed to get the business information: " + err.Error())
	} else if businessIDvalue == nil {
		return shim.Error("businesscc: " + "No information is avalilable on this businessID " + args[0])
	}
	err = json.Unmarshal(businessIDvalue, &parsedBusinessInfo)
	if err != nil {
		return shim.Error("businesscc: " + "Unable to parse businessInfo into the structure " + err.Error())
	}
	return shim.Success(nil)
}

func bisIDexists(stub shim.ChaincodeStubInterface, bisID string) pb.Response {
	ifExists, _ := stub.GetState(bisID)
	if ifExists != nil {
		fmt.Println(ifExists)
		return shim.Error("businesscc: " + "BusinessId " + bisID + " exits. Cannot create new ID")
	}
	return shim.Success(nil)
}
func bAccNoexists(stub shim.ChaincodeStubInterface, bacc string) pb.Response {
	ifExists, _ := stub.GetState(bacc)
	if ifExists != nil {
		fmt.Println(ifExists)
		return shim.Error("businesscc: " + "Business Account No  " + bacc + " exits. Cannot create new ID")
	}
	return shim.Success(nil)
}
func updateBusinessInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/*
		args[0] -> BusinessId
		args[1] -> Business Limit / MAX ROI / MAX ROI
		args[2] -> values
	*/
	if len(args) != 3 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("businesscc: " + "Invalid number of arguments in updateBusinessInfo(business) (required:3) given:" + xLenStr)
	}

	parsedBusinessInfo := businessInfo{}
	businessIDvalue, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("businesscc: " + "Failed to get the business information(updateBusinessInfo): " + err.Error())
	} else if businessIDvalue == nil {
		return shim.Error("businesscc: " + "No information is avalilable on this (updateBusinessInfo) businessID " + args[0])
	}

	err = json.Unmarshal(businessIDvalue, &parsedBusinessInfo)
	if err != nil {
		return shim.Error("businesscc: " + "Unable to parse businessInfo into the structure(updateBusinessInfo) " + err.Error())
	}

	lowerStr := strings.ToLower(args[1])

	value, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return shim.Error("businesscc: " + "value (updateBusinessInfo):" + err.Error())
	}

	if lowerStr == "business limit" {
		parsedBusinessInfo.BusinessLimit = value
	} else if lowerStr == "max roi" {
		parsedBusinessInfo.MaxROI = value
	} else if lowerStr == "min roi" {
		parsedBusinessInfo.MinROI = value
	}

	parsedBusinessInfoBytes, _ := json.Marshal(parsedBusinessInfo)
	err = stub.PutState(args[0], parsedBusinessInfoBytes)
	if err != nil {
		return shim.Error("businesscc: " + "Error in updating business: " + err.Error())
	}

	return shim.Success([]byte("businesscc: " + "Successfully updated Business " + args[0]))

}
func getWalletsofBusiness(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	fmt.Println("********************While Reading Business *************************")
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
	var businessValues []byte;

	var err4 error;
	mainWalletId,err1 := getBusinessWallet(stub ,intarg1)
	loanWalletId,err1 := getBusinessWallet(stub ,intarg2)
	chargesWalletId,err1 := getBusinessWallet(stub ,intarg3)
	liabilityWalletId,err1 := getBusinessWallet(stub ,intarg4)
	posWalletId,err1 := getBusinessWallet(stub ,intarg5)
	if err1 != nil {

	}
	fmt.Print("loanWalletId ",loanWalletId)
	fmt.Print("posWalletId ",posWalletId)
	businessvalasstruct := businessInfoVal{}
	businessValues,err4 = getBusinessInfoVal(stub, args)
	
	 err := json.Unmarshal(businessValues, &businessvalasstruct)
	if err != nil {
		return shim.Error("error while unmarshalling businessvalasstruct"+ err.Error())
    }  
	data, err1 := json.Marshal(businessvalasstruct)
	if err1 != nil {
		return shim.Error("error while marshaling  businessvalasstruct" + err1.Error() )
	} 
	//fmt.Println("Before appending balance %s",data)

	if err4 != nil {
		return shim.Error("Unable to get Business Info>>>> "+err4.Error())
	} else {
		//	fmt.Print("businessValues " ,businessValues)
	}
	
	chaincodeArgs1 := toChaincodeArgs("getWallet", mainWalletId)
	response1 = stub.InvokeChaincode("walletcc", chaincodeArgs1, "myc")
	comres=append(comres,string(response1.Payload))

	businessvalasstruct.BusinessWalletBalance = string(response1.Payload);
	fmt.Println("mainwallet balance %s",string(response1.Payload))

	chaincodeArgs2 := toChaincodeArgs("getWallet", loanWalletId)
	response2 = stub.InvokeChaincode("walletcc", chaincodeArgs2, "myc")
	comres=append(comres,string(response2.Payload))
	businessvalasstruct.BusinessLoanWalletBalance=string(response2.Payload);	
	fmt.Println("loanwallet balance %s",string(response2.Payload))

	chaincodeArgs3 := toChaincodeArgs("getWallet", chargesWalletId)
	response3 = stub.InvokeChaincode("walletcc", chaincodeArgs3, "myc")
	comres=append(comres,string(response3.Payload))
	businessvalasstruct.BusinessChargesOutstandingWalletBalance=string(response3.Payload);	
	fmt.Println("chargesWalletId balance %s",string(response3.Payload))

	chaincodeArgs4 := toChaincodeArgs("getWallet", liabilityWalletId)
	response4 = stub.InvokeChaincode("walletcc", chaincodeArgs4, "myc")
	fmt.Println("liabilityWallet balance %s",string(response4.Payload))
	comres=append(comres,string(response4.Payload))
	businessvalasstruct.BusinessLiabilityWalletBalance = string(response4.Payload);
				
	chaincodeArgs5 := toChaincodeArgs("getWallet", posWalletId)
	response5 = stub.InvokeChaincode("walletcc", chaincodeArgs5, "myc")
	comres=append(comres,string(response5.Payload))
	businessvalasstruct.BusinessPrincipalOutstandingWalletBalance = string(response5.Payload);
	
	fmt.Println("BusinessPrincipalOutstandingWalletBalance  %s",string(response5.Payload))

	//bankvalasstruct.BankChargesWalletIDBal = string(response3.Payload);
	data, err1 = json.Marshal(businessvalasstruct)
	if err1 != nil {
		return shim.Error("after adding error while marshaling  bankvalasstruct" + err1.Error() )
	} 	
	fmt.Println("comres as string %s",comres);
//	fmt.Println("Returning value getWalletsofBusiness %s",data)

	fmt.Println("********************End Reading Business *************************")
	return shim.Success(data)
}
func getBusinessWallet(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	fmt.Println("********************Start getBusinessWallet *************************")
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return "", fmt.Errorf("business : " + "Invalid number of arguments in getWalletID(business) (required:2) given:" + xLenStr)
	}
	businessInfoBytes, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("businesscc : " + "Unable to fetch the state" + err.Error())
	}
	if businessInfoBytes == nil {
		return "" ,fmt.Errorf("businesscc : " + "Data does not exist for " + args[0])
	}
	business := businessInfo{}
	err = json.Unmarshal(businessInfoBytes, &business)
	if err != nil {
		
		return "", fmt.Errorf("bankcc : " + "Uable to paser into the json format")
	}
	walletID := ""
	switch args[1] {
	case "main":
		walletID = business.BusinessWalletID
	case "loan":
		walletID = business.BusinessLoanWalletID
	case "charges":
		walletID = business.BusinessChargesOutstandingWalletID
	case "liability":
		walletID = business.BusinessLiabilityWalletID
	case "POS":
		walletID = business.BusinessPrincipalOutstandingWalletID
	}	
	fmt.Println("walletID based on case %s " ,walletID)
	fmt.Println("********************End getBusinessWallet *************************")
	return string(walletID),nil;
}

func getBusinessInfoVal(stub shim.ChaincodeStubInterface, args []string) ([]byte,error) {

	fmt.Println("******************** Start getBusinessInfoVal *************************")
	fmt.Println("args ",args[0])
	businessvalasstruct := businessInfoVal{}
	business, err := stub.GetState(args[0])
	if err != nil {
		return nil,fmt.Errorf("businesscc : " + "Unable to fetch the state" + err.Error())
	}
	if business == nil {
			return nil,fmt.Errorf("businesscc : " + "Data does not exist for " + args[0])
	}
	err = json.Unmarshal(business, &businessvalasstruct)
	if err != nil {
		return nil,fmt.Errorf("businesscc : " + "Uable to paser into the json format")
	}	
	data, err := json.Marshal(businessvalasstruct)
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Println("Returning from getBusinessInfoVal %s",data)
return data,nil 
}

func getWalletID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("businesscc: " + "Invalid number of arguments in getWalletId(business) (required:2) given:" + xLenStr)
	}
	parsedBusinessInfo := businessInfo{}
	businessIDvalue, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("businesscc: " + "Failed to get the business information: " + err.Error())
	} else if businessIDvalue == nil {
		return shim.Error("businesscc: " + "No information is avalilable on this businessID " + args[0])
	}
	err = json.Unmarshal(businessIDvalue, &parsedBusinessInfo)
	if err != nil {
		return shim.Error("businesscc: " + "Unable to parse into the structure " + err.Error())
	}
	walletID := ""
	switch args[1] {
	case "main":
		walletID = parsedBusinessInfo.BusinessWalletID
	case "loan":
		walletID = parsedBusinessInfo.BusinessLoanWalletID
	case "liability":
		walletID = parsedBusinessInfo.BusinessLiabilityWalletID
	case "principalOut":
		walletID = parsedBusinessInfo.BusinessPrincipalOutstandingWalletID
	case "chargesOut":
		walletID = parsedBusinessInfo.BusinessChargesOutstandingWalletID
	default:
		return shim.Error("businesscc: " + "There is no wallet of this type in Business :" + args[1])
	}
	fmt.Println("Wallets based on case ",walletID)

	return shim.Success([]byte(walletID))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Println("businesscc: "+"Error starting Business chaincode: %s\n", err)
	}

}

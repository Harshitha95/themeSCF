package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type pprInfo struct {
	ProgramID                         string  `json:"ProgramID"`                         //[1]
	BusinessID                        string  `json:"BusinessID"`                        //[2]
	Relationship                      string  `json:"Relationship"`                      //[3]
	ProgramBusinessLimit              int64   `json:"ProgramBusinessLimit"`              //[4]
	ProgramBusinessROI                float64 `json:"ProgramBusinessROI"`                //[5]
	ProgramBusinessDiscountPeriod     int     `json:"ProgramBusinessDiscountPeriod"`     //[6]
	ProgramBusinessDiscountPercentage string  `json:"ProgramBusinessDiscountPercentage"` //[7]//use float64 for parsing
	StaleDays                         int     `json:"StaleDays"`                         //[8]
	RepaymentAcNo                     string  `json:"RepaymentAcNo"`                     //[9]
	RepaymentWalletID                 string  `json:"RepaymentWalletID"`                 //will be taken from business Id
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	indexName := "ProgramID~BusinessID~DiscountPercentage"
	ppr := pprInfo{}
	prgrmBusPercentageKey, err := stub.CreateCompositeKey(indexName, []string{ppr.ProgramID, ppr.BusinessID, ppr.ProgramBusinessDiscountPercentage})
	if err != nil {
		return shim.Error("pprcc: " + "Unableto create composite key ProgramID~BusinessID~DiscountPercentage :" + err.Error())

	}
	value := []byte{0x00}
	stub.PutState(prgrmBusPercentageKey, value)
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "createPPR" {
		//Creates a new PPR Information
		return createPPR(stub, args)
	} else if function == "seePPR" {
		//Retrieves the PPR Information
		return seePPR(stub, args)
	} else if function == "pprIDexists" {
		//Checks the existence of PprID
		return pprIDexists(stub, args[0])
	} else if function == "getDiscountPercentage" {
		//Returns DiscountPercentage
		return discountPercentage(stub, args)
	} else if function == "updatePPR" {
		/*
			Parameters for Value Calculation
			    a. Program Business Limit
			    b. Program Business ROI
			    c. Program Business Discount Percentage
				d. Program Business Discount Period
			are updated
		*/
		return updatePPR(stub, args)
	}
	return shim.Error("pprcc: " + "No function named " + function + " in PPRsssssss")
}

func createPPR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("******************** Start createPPR *************************")
	fmt.Println("args[0]", args[0])
	fmt.Println("args[1]", args[1])
	fmt.Println("args[2]", args[2])
	fmt.Println("args[3]", args[3])
	fmt.Println("args[4]", args[4])
	fmt.Println("args[5]", args[5])
	fmt.Println("args[6]", args[6])
	fmt.Println("args[7]" ,args[7])
	fmt.Println("args[8] ",args[8])
	if len(args) != 10 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("pprcc: " + "Invalid number of arguments in createPPR (required:10) given:" + xLenStr)
	}

	//Checking existence of PprID
	response := pprIDexists(stub, args[0])
	if response.Status != shim.OK {
		return shim.Error("pprcc: " + response.Message)
	}

	//Checking existence of businessID
	chaincodeArgs := toChaincodeArgs("bisIDexists", args[2])
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("pprcc: " + "BusinessId " + args[2] + " does not exits")
	}

	//Checking existence of ProgramID
	chaincodeArgs = toChaincodeArgs("programIDexists", args[1])
	response = stub.InvokeChaincode("programcc", chaincodeArgs, "myc")
	if response.Status == shim.OK {
		return shim.Error("pprcc: " + "ProgramId " + args[1] + " does not exits")
	}

	relationship := map[string]bool{
		"seller": true,
		"vendor": true,
		"buyer":  true,
		"dealer": true,
	}

	relationshipLower := strings.ToLower(args[3])

	if !relationship[relationshipLower] {
		return shim.Error("pprcc: " + "Invalid relationship " + relationshipLower)
	}

	// ProgramBusinessLimit -> PBLimit
	PBLimit, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error("pprcc: " + err.Error())
	}

	//ProgramBusinessROI -> PBroi
	PBroi, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return shim.Error("pprcc: " + err.Error())
	}

	//ProgramBusinessDiscountPeriod -> PBDperiod
	PBDperiod, err := strconv.Atoi(args[6])
	if err != nil {
		return shim.Error("pprcc: " + err.Error())
	}

	//ProgramBusinessDiscountPercentage -> PBDpercentange
	_, err = strconv.ParseFloat(args[7], 64)
	if err != nil {
		return shim.Error("pprcc: " + err.Error())
	}

	//StaleDays -> sDays
	sDays, err := strconv.Atoi(args[8])
	if err != nil {
		return shim.Error("pprcc: " + err.Error())
	}

	//Wallet ID for repayment
	chaincodeArgs = toChaincodeArgs("getWalletID", args[2], "main")
	response = stub.InvokeChaincode("businesscc", chaincodeArgs, "myc")
	if response.Status != shim.OK {
		return shim.Error("pprcc: " + response.Message)
	}
	repayWalletID := string(response.GetPayload())

	ppr := pprInfo{args[1], args[2], relationshipLower, PBLimit, PBroi, PBDperiod, args[7], sDays, args[9], repayWalletID}
	pprBytes, err := json.Marshal(ppr)
	err = stub.PutState(args[0], pprBytes)
	fmt.Println("******************** end  createPPR *************************")
	return shim.Success([]byte("Successfully added PPR ", args[0]))
	
}

func pprIDexists(stub shim.ChaincodeStubInterface, pprID string) pb.Response {
	fmt.Println("********** Inside pprIDexists with arg ", pprID);
	ifExists, _ := stub.GetState(pprID)
	if ifExists != nil {
		fmt.Println(ifExists)
		return shim.Error("pprcc: " + "PprId " + pprID + " exits. Cannot create new ID")
	}
	fmt.Println("********** End pprIDexists")
	return shim.Success(nil)
}

func updatePPR(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********** Inside updateProgramInfo with args ", args[0] ," ", args[1] ," ", args[2]);
	if len(args) != 3 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("pprcc: " + "Invalid number of arguments in updatePPR(PPR) (required:3) given:" + xLenStr)
	}
	/*
		args[0] -> pprID
		args[1] -> Program Business Limit, Program Business ROI,
				   Program Business Discount Percentage, Program Business Discount Period
		args[2] -> values
	*/
	pprObject := pprInfo{}
	pprBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("pprcc: " + "updatePPR(PPR)" + err.Error())
	}
	err = json.Unmarshal(pprBytes, &pprObject)
	lowerStr := strings.ToLower(args[1])

	if lowerStr == "program business limit" {
		//Changing Program Business Limit
		PBLimit, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return shim.Error("pprcc: " + "updatePPR(PPR) Program Business Limit" + err.Error())
		}
		pprObject.ProgramBusinessLimit = PBLimit
	} else if lowerStr == "program business roi" {
		//Changing Program Business ROI
		PBroi, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return shim.Error("pprcc: " + "updatePPR(PPR) Program Business ROI" + err.Error())
		}
		pprObject.ProgramBusinessROI = PBroi
	} else if lowerStr == "program business discount percentage" {
		//Changing Program Business Discount Percentage
		_, err = strconv.ParseFloat(args[2], 64)
		if err != nil {
			return shim.Error("pprcc: " + "updatePPR(PPR) Program Business Discount Percentage" + err.Error())
		}
		pprObject.ProgramBusinessDiscountPercentage = args[2]
	} else if lowerStr == "program business discount period" {
		//Chainging Program Business Discount Period
		PBDperiod, err := strconv.Atoi(args[2])
		if err != nil {
			return shim.Error("pprcc: " + "updatePPR(PPR) Program Business Discount Period" + err.Error())
		}
		pprObject.ProgramBusinessDiscountPeriod = PBDperiod
	}
	fmt.Println("********** End updatePPR")
	return shim.Success(nil)

}

func discountPercentage(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	prgrmBusPercentageIte, err := stub.GetStateByPartialCompositeKey("ProgramID~BusinessID~DiscountPercentage", []string{args[0], args[1]})
	prgrmBusPercentageData, _ := prgrmBusPercentageIte.Next()
	_, data, err := stub.SplitCompositeKey(prgrmBusPercentageData.Key)
	if err != nil {
		return shim.Error("pprcc: " + "Error spliting composite key ProgramID~BusinessID~DiscountPercentage (ppr):" + err.Error())
	}
	defer prgrmBusPercentageIte.Close()
	return shim.Success([]byte(data[2]))
}

func seePPR(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********** Inside updateProgramInfo with args ", args[0])
	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("pprcc: " + "Invalid number of arguments in seePPR (required:1) given:" + xLenStr)
	}
	pprObject := pprInfo{}
	pprArray, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("pprcc: " + err.Error())
	}
	err = json.Unmarshal(pprArray, &pprObject)
	pprString := fmt.Sprintf("%+v", pprObject)
	fmt.Println("pprString ",pprString)
	fmt.Println("********** End seePPR")
	return shim.Success([]byte(pprString))
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("pprcc: "+"Error starting PPR chaincode: %s\n", err)
	}
}

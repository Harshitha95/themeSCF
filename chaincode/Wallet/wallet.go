package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type chainCode struct {
}

type walletsInfo struct {
	Balance float64 `json:"balance"`
}

func (c *chainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (c *chainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "newWallet" {
		return newWallet(stub, args)
	} else if function == "getWallet" {
		return getWallet(stub, args)
	} else if function == "updateWallet" {
		return updateWallet(stub, args)
	}
	return shim.Error("walletcc: " + "No function named " + function + " in Wallet")

}

//Creating new Wallet

func newWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("********************While Creating Wallet*************************")
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("walletcc: " + "Invalid number of arguments in newWallet (required:2) given:" + xLenStr)
	}
	fmt.Println("args ",args[0] ," ",args[1]);
	bal64, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("walletcc: " + err.Error())
	}
	//TODO Check for duplicate WalletId before inserting 
	ifExists, err := stub.GetState(args[0])
	if ifExists != nil {
		return shim.Error("walletcc: " + "WalletId " + args[0] + " exits. Cannot create new ID")
	}

	bal := walletsInfo{bal64}
	balBytes, _ := json.Marshal(bal)
	err = stub.PutState(args[0], balBytes)
	fmt.Println("********************End Creating Wallet*************************")
	return shim.Success(nil)
}

func getWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("******************** Start get Wallet*************************")
	if len(args) != 1 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("walletcc: " + "Invalid number of arguments in getWallet (required:1) given: " + xLenStr)
	}
	balBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("walletcc: " + err.Error())
	} else if balBytes == nil {
		return shim.Error("walletcc: " + "No data exists on this WalletId: " + args[0])
	}
	bal := walletsInfo{}
	err = json.Unmarshal(balBytes, &bal)
	if err != nil {
		return shim.Error("walletcc: " + err.Error())
	}
	balInt := int64(bal.Balance)
	balStr := strconv.FormatInt(balInt, 10)
	fmt.Println(" getWallet WalletId ",args[0])
	fmt.Println(" getWallet Converted to String ",balStr)
	fmt.Println("******************** End get Wallet*************************")
	return shim.Success([]byte(balStr))
}

func updateWallet(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	/*
	*args[0] -> WalletID
	*args[1] -> Wallet Ballance
	 */
	 fmt.Println("******************** Start Update Wallet*************************")
	 fmt.Println("args ",args[0]," ",args[1])
	if len(args) != 2 {
		xLenStr := strconv.Itoa(len(args))
		return shim.Error("walletcc: " + "Invalid number of arguments in Wallet Updation (required:2) given: " + xLenStr)
	}
	balBytes, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error("walletcc: " + err.Error())
	} else if balBytes == nil {
		return shim.Error("walletcc: " + "No data exists on this WalletId: " + args[0])
	}
	bal := walletsInfo{}
	err = json.Unmarshal(balBytes, &bal)
	if err != nil {
		return shim.Error("walletcc: " + err.Error())
	}
	bal.Balance, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("walletcc: " + "Error in Wallet updation parse int" + err.Error())
	}
	balBytes, _ = json.Marshal(bal)
	err = stub.PutState(args[0], balBytes)
	if err != nil {
		return shim.Error("walletcc: " + "Error in Wallet updation " + err.Error())
	}
	fmt.Printf("Balance for  %s : %f\n", args[0], bal.Balance)
	fmt.Println("******************** End Update Wallet*************************")
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(chainCode))
	if err != nil {
		fmt.Printf("walletcc: "+"Error starting Wallet chaincode: %s\n", err)
	}
}

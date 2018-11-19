/*
 * it should return 2 values : walletID, openBal, txnBalance
 * it takes arguments : participantID, walletType, ccName, cAmt, dAmt



 */

//#######################################################################################################################
	//Calling for updating Bank Main_Wallet
	//####################################################################################################################

	cAmtString := "0"
	dAmtString := args[5]

	walletID, openBalString, txnBalString, err := getWalletInfo(args[6], "main", "bankcc", cAmt, dAmt)
	if err != nil {
		return shim.Error(err.Error())
	}



	func getWalletInfo(participantID string, walletType string, ccName string, cAmtStr int64, dAmtStr int64) string, string, string, error {

		// STEP-1
		// using FromID, get a walletID from bank structure
		// bankID = bankID
		
		chaincodeArgs := util.ToChaincodeArgs("getWalletID", participantID, walletType)
		response := stub.InvokeChaincode("bankcc", chaincodeArgs, "myc")
		if response.Status != shim.OK {
			return "", "", "", error.new(response.Message)
		}
		walletID := string(response.GetPayload())

		// STEP-2
		// getting Balance from walletID
		// walletFcn := "getWallet"
		walletArgs := util.ToChaincodeArgs("getWallet", walletID)
		walletResponse := stub.InvokeChaincode("walletcc", walletArgs, "myc")
		if walletResponse.Status != shim.OK {
			return "", "", "", error.new(walletResponse.Message)
		}
		openBalString := string(walletResponse.Payload)

		openBal, err := strconv.ParseInt(openBalString, 10, 64)
		if err != nil {
			return "", "", "", error.new("Error in converting the openBalance")
		}
		cAmt, err := strconv.ParseInt(cAmtStr, 10, 64)
		if err != nil {
			return "", "", "", error.new("Error in converting the cAmt")
		}
		dAmt, err := strconv.ParseInt(dAmtStr, 10, 64)
		if err != nil {
			return "", "", "", error.new("Error in converting the dAmt")
		}

		txnBal := openBal - dAmt + cAmt
		txnBalString := strconv.FormatInt(txnBal, 10)

		// STEP-3
		// update wallet of ID walletID here, and write it to the wallet_ledger
		// walletFcn := "updateWallet"
		walletArgs = util.ToChaincodeArgs("updateWallet", walletID, string(txnBal))
		walletResponse = stub.InvokeChaincode("walletcc", walletArgs, "myc")
		if walletResponse.Status != shim.OK {
			return "", "", "", error.new(walletResponse.Message)
		}
		
		return walletID, openBalString, txnBalString, nil
	}
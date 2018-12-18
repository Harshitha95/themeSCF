var express = require('express')

var app = express();

var router = express.Router()
var url = require('url');

const util = require('util')
//proposalresponse payload(shim.sucess)
//const protoLoader = require('@grpc/proto-loader');
var path = require('path');
var grpc = require('grpc');
//var utils = require('../node_modules/fabric-client/lib/utils.js');
//var blockdecoder = utils.getLogger('BlockDecoder.js');
//var _responseProto = protoLoader.load('/node_modules/fabric-client/lib/protos/peer/proposal_response.proto');
//var _responseProto = grpc.load(__dirname + '/../node_modules/fabric-client/lib/protos/peer/proposal_response.proto').protos;
//var _responseProto = path.join(__dirname + './node_modules/fabric-client/lib/protos/peer/proposal_response.proto').protos;
//var proto_proposal_response_payload = _responseProto.ProposalResponsePayload.decode(proposal_response_payload_bytes);
// define the home page route
router.get('/', function (req, res) {
	'use strict';
	/*
* Copyright IBM Corp All Rights Reserved
*
* SPDX-License-Identifier: Apache-2.0
*/	/*
	 * Chaincode Invoke
	 */
	//var url_parts = url.parse(req.url, true);
    //var query = url_parts.query;
    //var bodyParser = require('body-parser');
    //app.use(bodyParser.json());
    var data=req.body;
    console.log("post data as JSON"+JSON.stringify(req.body));
 

    console.log("req.body.chaincodeId "+req.body.chaincodeId);
    console.log("req.body.methodToBeCalled "+req.body.methodToBeCalled);
    console.log("req.body.dataArguments "+req.body.dataArguments);
  //  alert("req.body.dataArguments", req.body.dataArguments);
//console.log("Length "+ req.body.count);
console.log("Length of Object  "+ Object.keys(req.body).length+" Data Arguments "+ Object.keys(req.body.dataArguments).length);
 var  functionArgs = [];
//iterating the 3rd 
for (var i in req.body.dataArguments) {
       console.log("req.body.dataArguments "+req.body.dataArguments[i])
       functionArgs.push(req.body.dataArguments[i])
}
console.log("functionArgs "+functionArgs); 
/**	var myArgs = JSON.parse(query.arguments);
var  myArgs = [];
myArgs.push("bankcc");
myArgs.push("writeBankInfo");
myArgs.push(req.body.bankId);//
myArgs.push(req.body.bankName);
myArgs.push(req.body.bankBranch);
myArgs.push(req.body.bankCode);
myArgs.push(req.body.bWB);
myArgs.push("1000");
myArgs.push("1000");
myArgs.push("1000");
myArgs.push("1000");
var count=0;
data.forEach(function (myArgs) {
   
    console.log("Count "+count++);
});
var fcn_args = myArgs.slice(2);**/


//console.log("Query Params" +query.arguments);
	
	var Fabric_Client = require('fabric-client');
	var path = require('path');
	var util = require('util');
	var os = require('os');

	var proposalResponses ;
		var proposal ;
	//
	var fabric_client = new Fabric_Client();

	// setup the fabric network
	var channel = fabric_client.newChannel('myc');
	var peer = fabric_client.newPeer('grpc://localhost:7051');
	channel.addPeer(peer);
	var order = fabric_client.newOrderer('grpc://localhost:7050')
	channel.addOrderer(order);

	//
	var member_user = null;
	var store_path = path.join(__dirname, '/../hfc-key-store');
	console.log('Store path:' + store_path);
	var tx_id = null;

	// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
	Fabric_Client.newDefaultKeyValueStore({
		path: store_path
	}).then((state_store) => {
		// assign the store to the fabric client
		fabric_client.setStateStore(state_store);
		var crypto_suite = Fabric_Client.newCryptoSuite();
		// use the same location for the state store (where the users' certificate are kept)
		// and the crypto store (where the users' keys are kept)
		var crypto_store = Fabric_Client.newCryptoKeyStore({ path: store_path });
		crypto_suite.setCryptoKeyStore(crypto_store);
		fabric_client.setCryptoSuite(crypto_suite);

		// get the enrolled user from persistence, this user will sign all requests
		return fabric_client.getUserContext('user1', true);
	}).then((user_from_store) => {
		if (user_from_store && user_from_store.isEnrolled()) {
			console.log('Successfully loaded user1 from persistence');
			member_user = user_from_store;
		} else {
			throw new Error('Failed to get user1.... run registerUser.js');
		}

		// get a transaction id object based on the current user assigned to fabric client
		tx_id = fabric_client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);

		// createCar chaincode function - requires 5 args, ex: args: ['CAR12', 'Honda', 'Accord', 'Black', 'Tom'],
		// changeCarOwner chaincode function - requires 2 args , ex: args: ['CAR10', 'Dave'],
		// must send the proposal to endorsing peers
        
        var request = {
			//targets: let default to the peer assigned to the client
			chaincodeId: req.body.chaincodeId,
			fcn: req.body.methodToBeCalled,
			args: functionArgs,
			chainId: 'myc',
			txId: tx_id
		};
        // send the transaction proposal to the peers
        
		return channel.sendTransactionProposal(request);
	}).then((results) => {
		// Declared Globaly 
		 proposalResponses = results[0];
		 proposal = results[1];
		let isProposalGood = false;
		if (proposalResponses && proposalResponses[0].response &&
			proposalResponses[0].response.status === 200) {
			isProposalGood = true;
            console.log('Transaction proposal was good'+ proposalResponses[0].response.message+ "Hello "+ proposalResponses[0].response.payload + proposalResponses[0].response.status + "End ");
          //  res.send()
		} else {
			console.error('Transaction proposal was bad '+ proposalResponses[0].response.message+ proposalResponses[0].response.status +'No dummy' );
		}
		if (isProposalGood) {
			//console.log("resp payload  "+ blockdecoder.decodeProposalResponsePayload(proposalResponses[0].payload));
			console.log(util.format(
				'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s",payload - %s',
				proposalResponses[0].response.status, proposalResponses[0].response.message,proposalResponses[0].payload));
				//console.log("payload response"+proposalResponses[0].payload);
			// build up the request for the orderer to have the transaction committed
			console.log("Length of result "+proposalResponses[0].response.payload.length)
			var request = {
				proposalResponses: proposalResponses,
				proposal: proposal
			};
			// set the transaction listener and set a timeout of 30 sec
			// if the transaction did not get committed within the timeout period,
			// report a TIMEOUT status
			var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
			var promises = [];

			var sendPromise = channel.sendTransaction(request);
			promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

			// get an eventhub once the fabric client has a user assigned. The user
			// is required bacause the event registration must be signed
			let event_hub = fabric_client.newEventHub();
			event_hub.setPeerAddr('grpc://localhost:7053');

			// using resolve the promise so that result status may be processed
			// under the then clause rather than having the catch clause process
			// the status
			let txPromise = new Promise((resolve, reject) => {
				let handle = setTimeout(() => {
					event_hub.disconnect();
					resolve({ event_status: 'TIMEOUT' }); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
				}, 30000);
				event_hub.connect();
				event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
					// this is the callback for transaction event status
					// first some clean up of event listener
					clearTimeout(handle);
					event_hub.unregisterTxEvent(transaction_id_string);
					event_hub.disconnect();

					// now let the application know what happened
					var return_status = { event_status: code, tx_id: transaction_id_string };
					if (code !== 'VALID') {
						console.error('The transaction was invalid, code = ' + code);
						resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
					} else {
						console.log('The transaction has been committed on peer ' + event_hub._ep._endpoint.addr);
						resolve(return_status);
					}
				}, (err) => {
					//this is the callback if something goes wrong with the event registration or processing
					reject(new Error('There was a problem with the eventhub ::' + err));
				});
			});
			promises.push(txPromise);

			return Promise.all(promises);
		} else {
			console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
			throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
		}
	}).then((results) => {
		console.log('Send transaction promise and event listener promise have completed');
		// check the results in the order the promises were added to the promise all list
		if (results && results[0] && results[0].status === 'SUCCESS') {
			console.log(util.format('Successfully sent transaction to the orderer."%s" ' ,proposalResponses[0].response.message));
			console.log(util.format('Successfully sent transaction to the orderer wesdfjkgh  as j"%j" ' ,proposalResponses[0].response.payload));
            //res.send('Successfully sent transaction to the orderer.');
		} else {
            console.error('Failed to order the transaction. Error code: ' + response.status);
            //res.send('Failed to order the transaction.');
		}
		if (results && results[1] || results[1].event_status === 'VALID') {
			console.log(util.format("Successfully committed the change to the ledger by the peer %s",proposalResponses[0].response.payload.length));

		} else {
			console.log('Transaction failed to be committed to the ledger due to ::' + results[1].event_status);
        }
        /**if(results[0].status === 'SUCCESS' || results[1].event_status === 'VALID'){
         res.send("Transaction Success" +proposalResponses[0].response.message);
        } else {
            res.send("Transaction Failed " + response.status + " and "+ results[1].event_status);
        }**/
       res.send(proposalResponses[0].response.payload);
	   // res.send('Transaction ');
	}).catch((err) => {
        console.error('Failed to invoke successfully : ' + err);
        res.send("Failed to invoke " +err);
	});
	//res.send('Successfully posted')
})
module.exports = router;

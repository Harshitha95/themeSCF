# ThemechainSCF0.1V
Golang

# To work on CLI, follow the steps.
Dev Mode Hyperledger Fabric

This document help you to test chaincode from cli 
STEP 1	
	- go to first-network inside this directory
	- run
		./byfn.sh -m  down
STEP 2
- in the same directory run
./byfn.sh -m generate -c myc
-c myc – chanel name
STEP 3
	- after that run
docker-compose -f docker-compose-cli.yaml up -d
Open another terminal
docker exec -it cli bash    - now you're inside the cli container
Note: moving to specific peer 
docker exec -it peer0.org1.example.com  bash
docker exec -it peer1.org1.example.com bash
	- type 'ls' and you'll see a directory called scripts
	- in that 'scripts/' there are two bash scripts :
		1) chan.sh : this is used to create channel, join peers to channel and specify anchor peer etc...
(or) follow
peer channel create -o orderer.example.com:7050 -c myc -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
peer channel join -b myc.block
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer channel join -b myc.block

peer channel update -o orderer.example.com:7050 -c myc -f ./channel-artifacts/Org1MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer channel update -o orderer.example.com:7050 -c myc -f ./channel-artifacts/Org2MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
2)  Run cmd  go get  github.com/google/uuid

3)  installCC.sh : contains the code to install and instantiate the business code.
Run the install and instantiate code manualy. 
Follow cciInstall.sh
peer chaincode install -n bankcc -v 1.0 -p github.com/chaincode/Bank/
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n bankcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
Change the logic code name to next appropriate chaincode. (business, bank, loan, instrument, program, ppr etc.., )



3) invokeCC.sh : this'll invoke all necessary ledgers.
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n bankcc -c '{"Args":["writeBankInfo","1bank","kvb","chennai","40A","2333s673sxx78","sdr3cfgtdui3","23rfs6vhj148b","897vhessety","86zs0lhtd"]}' -C myc
Change the chaincode name and arguments accordingly.

Note:
We can install any number of versions of chaincode, but instantiate only one version of code.
The instantiated code is saved in docker container as an image.
If u need to instaintiate another version of code , u need to remove the docker image and then re-instantiate.




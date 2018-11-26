# ThemechainSCF0.1V
Golang


#Pre-requisite

https://hyperledger.github.io/composer/latest/prereqs-ubuntu.sh
1.Ubuntu 16.04 LTS
2.Install Ubuntu 18.04 LTS
3.Install prereqs to install Hyperledger fabric V1.2
https://hyperledger-fabric.readthedocs.io/en/release-1.2/build_network.html

4. for binaries
https://github.com/hyperledger/fabric/blob/release-1.2/scripts/bootstrap.sh
curl the file ..and save it(with full rights chmod U+X) and run it.

5.get the code from git 
git clone https://github.com/uchihamalolan/EncoreBlockchain

6.go tom .bashrc file
vi ~/.bashrc add Path to hyperledger first network
export PATH="$HOME/hyperledger/fabric-samples/bin:$PATH"
OR similiar variable for bin to be used
 
7. reload the bash by 
source ~/.bashrc

8.CD fabric-samples/Firstnetwork/ run ./byfn.sh generate ( optionally use channel name by giving -c "channelname"
whih will generate the certificate ,keys...

9.CD EncoreBlockchain/Node/  run ./startfabric.sh

10. npm install

11.node enrollAdmin.js
12.node registerUser.js

13 ./installCC.sh
 install chaincode membership certificate location and ,version(-v) , chaincodename , and chaincode(.go file) location
 instansiate chaincode membership certificate location and , chaincodename(-n) , and chaincode(.go file) location ,channelname(-C) ,
arguments(-c) , validation (-P)



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




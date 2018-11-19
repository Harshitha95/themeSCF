docker-compose -f docker-compose-cli.yaml up -d

docker exec -it cli bash

peer channel create -o orderer.example.com:7050 -c myc -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

peer channel join -b myc.block

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer channel join -b myc.block

peer channel update -o orderer.example.com:7050 -c myc -f ./channel-artifacts/Org1MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer channel update -o orderer.example.com:7050 -c myc -f ./channel-artifacts/Org2MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem






peer chaincode install -n bankcc -v 1.0 -p github.com/chaincode/Bank/

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n bankcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n businesscc -v 1.0 -p github.com/chaincode/Business/

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n businesscc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n txncc -v 1.0 -p github.com/chaincode/Transactions/

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n txncc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n disbursementcc -v 1.0 -p github.com/chaincode/Transactions/Disbursement

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n disbursementcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n txnbalcc -v 1.0 -p github.com/chaincode/TxnBalance

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n txnbalcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n walletcc -v 1.0 -p github.com/chaincode/Wallet

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n walletcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

peer chaincode install -n programcc -v 1.0 -p github.com/chaincode/Program

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n programcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n pprcc -v 1.0 -p github.com/chaincode/PPR

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n pprcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n accrualcc -v 1.0 -p github.com/chaincode/Transactions/Accrual

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n accrualcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

peer chaincode install -n chargescc -v 1.0 -p github.com/chaincode/Transactions/Charges

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n chargescc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

peer chaincode install -n interestAccruedChargecc -v 1.0 -p github.com/chaincode/Transactions/InterestAccruedCharge

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n interestAccruedChargecc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n interestInAdvancecc -v 1.0 -p github.com/chaincode/Transactions/InterestInAdvance

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n interestInAdvancecc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n interestRefundcc -v 1.0 -p github.com/chaincode/Transactions/InterestRefund

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n interestRefundcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n marginRefundcc -v 1.0 -p github.com/chaincode/Transactions/MarginRefund

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n marginRefundcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n penalChargescc -v 1.0 -p github.com/chaincode/Transactions/PenalCharges

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n penalChargescc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n penalInterestCollectioncc -v 1.0 -p github.com/chaincode/Transactions/PenalInterestCollection

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n penalInterestCollectioncc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n repaymentcc -v 1.0 -p github.com/chaincode/Transactions/Repayment

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n repaymentcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


peer chaincode install -n tdscc -v 1.0 -p github.com/chaincode/Transactions/TDS

peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n tdscc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"


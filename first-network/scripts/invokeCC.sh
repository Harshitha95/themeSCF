echo "invoking bank"
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n bankcc -c '{"Args":["writeBankInfo","1bank","kvb","chennai","40A","2333s673sxx78","sdr3cfgtdui3","23rfs6vhj148b","897vhessety","86zs0lhtd"]}' -C myc

echo "invoking business"
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n businesscc -c '{"Args":["putNewBusinessInfo","1bus","tata","12348901","4000000","23hhnx56s673sxx78","sdr32123d3","23rfs148b","12.4","8.09","0","1000000"]}' -C myc

echo "invoking wallets"
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n walletcc -c '{"Args":["newWallet","23hhnx56s673sxx78","1000"]}' -C myc
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n walletcc -c '{"Args":["newWallet","23rfs148b","1000"]}' -C myc
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n walletcc -c '{"Args":["newWallet","sdr32123d3","1000"]}' -C myc
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n walletcc -c '{"Args":["newWallet","2333s673sxx78","1000"]}' -C myc
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n walletcc -c '{"Args":["newWallet","sdr3cfgtdui3","1000"]}' -C myc
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n walletcc -c '{"Args":["newWallet","23rfs6vhj148b","1000"]}' -C myc
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n walletcc -c '{"Args":["newWallet","897vhessety","1000"]}' -C myc
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n walletcc -c '{"Args":["newWallet","86zs0lhtd","1000"]}' -C myc

echo "invoking loan"
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n loancc -c '{"Args":["newLoanInfo","1loan","1ins","1eb","1prg","900","23/04/2018:12:45:20","pragadeesh","5.6","23/10/2018","25/09/2018:20:45:01","sanctioned","900"]}' -C myc

echo "invoking loanBalance"
peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n loanbalcc -c '{"Args":["putLoanBalInfo","1loanbal","1loan","1txn","23/04/2018","disbursement","900","0","0","900","sanctioned"]}' -C myc



peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n instrumentcc -c '{"Args":["enterInstrument","5051","12/09/2018","B002","B001","10000","23/07/2019","123456","TZQ2 486","32645","31/12/2019"]}' -C myc


peer chaincode invoke -o orderer.example.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem  -C myc -n instrumentcc -c '{"Args":["getInstrument","3ins"]}' -C myc

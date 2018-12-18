#Kill the existing network
./avadaKedavra.sh 
# Start the network
./startFabric.sh
npm install
# enroll admin
node enrollAdmin.js
# register user
node registerUser.js
#get uuid package
docker exec cli go get github.com/google/uuid
# install chaincodes
./installCC.sh

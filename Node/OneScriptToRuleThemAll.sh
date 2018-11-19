./startFabric.sh;
echo "###########################" 
echo "Established Network"
echo "###########################"
node enrollAdmin.js
node registerUser.js
echo "###########################" 
echo "Enrolled Admin and Registered User"
echo "###########################"
./installCC.sh
echo "###########################" 
echo "Installed and Instantiated Chaincodes"
echo "###########################"
./invokeCC.sh
echo "###########################" 
echo "Invoked necessary ledgers for handling txns"
echo "###########################"
node app.js

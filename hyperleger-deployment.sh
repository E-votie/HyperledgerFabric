sudo apt update
sudo apt install docker.io
sudo apt install docker-compose
sudo apt install openssl
wget https://github.com/hyperledger/firefly-cli/releases/download/v1.3.0/firefly-cli_1.3.0_Linux_x86_64.tar.gz
sudo tar -zxf ~/Downloads/firefly-cli_*.tar.gz -C /usr/local/bin ff && rm ~/Downloads/firefly-cli_*.tar.gz
go install github.com/hyperledger/firefly-cli/ff@latest
ff version
ff init fab
#give the stack name as fab and number of nodes as 3
ff start fab #what ever the name given above
cd voting-chaincode
ff deploy fabric fab e_votie.zip firefly e_votie 1.0
curl -X POST http://localhost:5000/api/v1/namespaces/default/contracts/interfaces?publish=true \
     -H "Content-Type: application/json" \
     -d @interfaces.json
# Step 1: Send the first POST request and capture the response
response=$(curl -s -X POST http://localhost:5000/api/v1/namespaces/default/contracts/interfaces?publish=true \
                -H "Content-Type: application/json" \
                -d @interfaces.json)

# Extract the 'id' from the response using 'jq'
id=$(echo $response | jq -r '.id')

# Step 2: Construct the next JSON using the extracted 'id' and send the second POST request
json_data=$(cat <<EOF
{
  "name": "e_votie",
  "interface": {
    "id": "$id"
  },
  "location": {
    "channel": "firefly",
    "chaincode": "e_votie"
  }
}
EOF
)

# Send the second POST request
curl -X POST http://localhost:5000/api/v1/namespaces/default/contracts/interfaces?publish=true \
     -H "Content-Type: application/json" \
     -d "$json_data"

#then goto the http://127.0.0.1:5108 and go to the contact under API you can find the api



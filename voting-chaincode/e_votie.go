/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"voting-chaincode/chaincode"
)

func main() {
	ElectionChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		fmt.Printf("Error creating voting chaincode: %s", err.Error())
		return
	}

	if err := ElectionChaincode.Start(); err != nil {
		fmt.Printf("Error starting voting chaincode: %s", err.Error())
	}
}

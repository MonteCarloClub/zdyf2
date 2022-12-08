/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SimpleAsset struct {
}

const IdPrefix = "UUID:"
const serialNumberPrefix = "SRIALNUMBER:"

func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// args := stub.GetStringArgs()
	// if len(args) != 2 {
	// 	return shim.Error("Incorrect arguments. Expecting a key and a value")
	// }
	// err := stub.PutState(args[0], []byte(args[1]))
	// if err != nil {
	// 	return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	// }
	return shim.Success(nil)
}

func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()	

	var result string
	var err error
	if fn == "setCertificate" {
		result, err = setCertificate(stub, args)
	} else if fn == "getCertificate" { 
		result, err = getCertificate(stub, args)
	}else if fn == "getCertificateByUUID"{
		result, err = getCertificateByUUID(stub, args)
	}else if fn == "revokeCertificate"{
		result, err = revokeCertificate(stub, args)
	}else if fn == "revokeCertificateByUUID"{
		result, err = revokeCertificateByUUID(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

func setCertificate(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a UUID, serialNumber and a Certificate")
	}
	var err error
	// var value []byte
	err = stub.PutState(IdPrefix + args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set UUID: %s", args[0])
	}
	err = stub.PutState(serialNumberPrefix + args[1], []byte(args[2]))
	if err != nil {
		return "", fmt.Errorf("Failed to set serialNumber: %s", args[0])
	}
	_, err = stub.GetState(serialNumberPrefix + args[1])
	if err != nil {
		return "", fmt.Errorf("Failed to get serialNumber: %s with error: %s", args[0], err)
	}
	return args[1], nil
}

func getCertificate(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a serialNumber")
	}

	value, err := stub.GetState(serialNumberPrefix + args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get serialNumber: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Certificate not found for serialNumber: %s", args[0])
	}
	return string(value), nil
}


func getCertificateByUUID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a UUID")
	}

	value, err := stub.GetState(IdPrefix + args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get UUID: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Certificate not found for UUID: %s", args[0])
	}
	return string(value), nil
}

func revokeCertificate(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a serialNumber")
	}

	err := stub.DelState(serialNumberPrefix + args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete serialNumber: %s with error: %s", args[0], err)
	}

	return "", nil
}

func revokeCertificateByUUID(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a UUID")
	}

	err := stub.DelState(IdPrefix + args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to delete UUID: %s with error: %s", args[0], err)
	}
	return "", nil
}


// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}

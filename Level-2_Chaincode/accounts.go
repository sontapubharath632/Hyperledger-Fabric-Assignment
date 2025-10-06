package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Asset struct {
	DealerID    string  `json:"dealerId"`
	MSISDN      string  `json:"msisdn"`
	MPIN        string  `json:"mpin"`
	Balance     float64 `json:"balance"`
	Status      string  `json:"status"`
	TransAmount float64 `json:"transAmount"`
	TransType   string  `json:"transType"`
	Remarks     string  `json:"remarks"`
}

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface,
	dealerId, msisdn, mpin string, balance float64, status string) error {

	exists, _ := ctx.GetStub().GetState(dealerId)
	if exists != nil {
		return fmt.Errorf("asset %s already exists", dealerId)
	}

	asset := Asset{DealerID: dealerId, MSISDN: msisdn, MPIN: mpin, Balance: balance, Status: status}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("error marshalling asset: %v", err)
	}

	return ctx.GetStub().PutState(dealerId, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, dealerId string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(dealerId)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset %s: %v", dealerId, err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("asset %s not found", dealerId)
	}

	var asset Asset
	if err := json.Unmarshal(assetJSON, &asset); err != nil {
		return nil, err
	}
	return &asset, nil
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dealerId string, balance float64, status string) error {
	asset, err := s.ReadAsset(ctx, dealerId)
	if err != nil {
		return err
	}

	asset.Balance = balance
	asset.Status = status

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("error marshalling updated asset: %v", err)
	}

	return ctx.GetStub().PutState(dealerId, assetJSON)
}

func (s *SmartContract) GetHistory(ctx contractapi.TransactionContextInterface, dealerId string) ([]string, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(dealerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get history for %s: %v", dealerId, err)
	}
	defer resultsIterator.Close()

	var history []string
	for resultsIterator.HasNext() {
		modification, _ := resultsIterator.Next()
		history = append(history, string(modification.Value))
	}
	return history, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(fmt.Sprintf("error creating chaincode: %v", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("error starting chaincode: %v", err))
	}
}

/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a proposal
type SmartContract struct {
	contractapi.Contract
}

// Proposal describes basic details of what makes up a proposal
type Proposal struct {
	ApplicantName         string `json:"applicantName"`
	LoanAmount            uint64 `json:"loanAmount"`
	IsApproved            bool   `json:"isApproved"`
	IsCIBILVerified       bool   `json:"isCIBILVerified"`
	IsTrackRecordVerified bool   `json:"isTrackRecordVerified"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Proposal
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Proposal{
		Proposal{ApplicantName: "Akshay Kurhekar", LoanAmount: 1000, IsApproved: false, IsCIBILVerified : false, IsTrackRecordVerified: false},
		Proposal{ApplicantName: "Rama", LoanAmount: 9999, IsApproved: false, IsCIBILVerified : false,  IsTrackRecordVerified: false},
		Proposal{ApplicantName: "Krishna", LoanAmount: 108, IsApproved: true, IsCIBILVerified : true, IsTrackRecordVerified: true},		
	}

	for i, proposal := range cars {
		proposalAsByte, _ := json.Marshal(proposal)
		err := ctx.GetStub().PutState("Proposal"+strconv.Itoa(i), proposalAsByte)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateProposal adds a new proposal to the world state with given details
func (s *SmartContract) CreateProposal(ctx contractapi.TransactionContextInterface, proposalNumber string, applicantName string, loanAmount uint64) error {
	proposal := Proposal{
		ApplicantName:         applicantName,
		LoanAmount:            loanAmount,
		IsApproved:            false,
		IsCIBILVerified:       false,
		IsTrackRecordVerified: false,
	}

	proposalAsBytes, _ := json.Marshal(proposal)

	return ctx.GetStub().PutState(proposalNumber, proposalAsBytes)
}

// QueryProposal returns the proposal stored in the world state with given id
func (s *SmartContract) QueryProposal(ctx contractapi.TransactionContextInterface, proposalNumber string) (*Proposal, error) {
	proposalAsByte, err := ctx.GetStub().GetState(proposalNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if proposalAsByte == nil {
		return nil, fmt.Errorf("%s does not exist", proposalNumber)
	}

	proposal := new(Proposal)
	_ = json.Unmarshal(proposalAsByte, proposal)

	return proposal, nil
}

// QueryAllProposal returns all cars found in world state
func (s *SmartContract) QueryAllProposal(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		proposal := new(Proposal)
		_ = json.Unmarshal(queryResponse.Value, proposal)

		queryResult := QueryResult{Key: queryResponse.Key, Record: proposal}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of proposal with given id in world state
func (s *SmartContract) SetCIBILTrack(ctx contractapi.TransactionContextInterface, proposalNumber string, cibil bool, track bool) error {
	proposal, err := s.QueryProposal(ctx, proposalNumber)

	if err != nil {
		return err
	}

	proposal.IsCIBILVerified = cibil;
	proposal.IsTrackRecordVerified = track;

	proposalAsByte, _ := json.Marshal(proposal)

	return ctx.GetStub().PutState(proposalNumber, proposalAsByte)
}

// ApproveProposal updates the owner field of proposal with given id in world state
func (s *SmartContract) ApproveProposal(ctx contractapi.TransactionContextInterface, proposalNumber string) error {
	proposal, err := s.QueryProposal(ctx, proposalNumber)

	if err != nil {
		return err
	}

	if proposal.IsCIBILVerified == true && proposal.IsTrackRecordVerified == true {
		
		proposalAsByte, _ := json.Marshal(proposal)
	
		return ctx.GetStub().PutState(proposalNumber, proposalAsByte)
	}else {
		return fmt.Errorf("CIBIL OR Track is Not verified")
	}
}

// GetHistoryOfProposal
func (s *SmartContract) GetHistoryOfProposal(ctx contractapi.TransactionContextInterface, proposalId string) ([]*Proposal, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(proposalId)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal history: %v", err)
	}
	defer resultsIterator.Close()

	var history []*Proposal
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate history results: %v", err)
		}

		asset := new(Proposal)
		err = asset.FromJSON(response.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize asset JSON: %v", err)
		}

		history = append(history, asset)
	}

	return history, nil
}

func (a *Proposal) FromJSON(data []byte) error {
	return json.Unmarshal(data, a)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create loan proposal chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting loan proposal chaincode: %s", err.Error())
	}
}

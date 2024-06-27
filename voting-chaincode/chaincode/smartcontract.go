package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Voter struct {
	NIC               string `json:"nic"`
	Name              string `json:"name"`
	VoterID           string `json:"voterId"`
	BiometricTemplate string `json:"biometricTemplate"`
}

type Election struct {
	ElectionID      string    `json:"electionId"`
	StartTimestamp  time.Time `json:"startTimestamp"`
	EndTimestamp    time.Time `json:"endTimestamp"`
	ContractAddress string    `json:"contractAddress"`
}

type VoterElectionRelation struct {
	VoterID        string `json:"voterId"`
	ElectionID     string `json:"electionId"`
	Eligibility    bool   `json:"eligibility"`
	Voted          bool   `json:"voted"`
	PollingStation string `json:"pollingStation"`
}

func (s *SmartContract) CreateVoter(ctx contractapi.TransactionContextInterface, nic string, name string, voterID string, biometricTemplate string) error {
	voter := Voter{
		NIC:               nic,
		Name:              name,
		VoterID:           voterID,
		BiometricTemplate: biometricTemplate,
	}

	voterJSON, err := json.Marshal(voter)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(voterID, voterJSON)
}

func (s *SmartContract) CreateElection(ctx contractapi.TransactionContextInterface, electionID string, startTime time.Time, endTime time.Time, contractAddress string) error {
	election := Election{
		ElectionID:      electionID,
		StartTimestamp:  startTime,
		EndTimestamp:    endTime,
		ContractAddress: contractAddress,
	}

	electionJSON, err := json.Marshal(election)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(electionID, electionJSON)
}

func (s *SmartContract) AssignVoterToElection(ctx contractapi.TransactionContextInterface, voterID string, electionID string, eligibility bool, pollingStation string) error {
	relation := VoterElectionRelation{
		VoterID:        voterID,
		ElectionID:     electionID,
		Eligibility:    eligibility,
		Voted:          false,
		PollingStation: pollingStation,
	}

	relationJSON, err := json.Marshal(relation)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("relation_%s_%s", voterID, electionID)
	return ctx.GetStub().PutState(key, relationJSON)
}

func (s *SmartContract) RecordVote(ctx contractapi.TransactionContextInterface, voterID string, electionID string) error {
	key := fmt.Sprintf("relation_%s_%s", voterID, electionID)
	relationJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return fmt.Errorf("failed to read relation: %v", err)
	}
	if relationJSON == nil {
		return fmt.Errorf("relation does not exist")
	}

	var relation VoterElectionRelation
	err = json.Unmarshal(relationJSON, &relation)
	if err != nil {
		return err
	}

	if !relation.Eligibility {
		return fmt.Errorf("voter is not eligible for this election")
	}

	if relation.Voted {
		return fmt.Errorf("voter has already voted in this election")
	}

	relation.Voted = true

	updatedRelationJSON, err := json.Marshal(relation)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(key, updatedRelationJSON)
}

func (s *SmartContract) GetVoter(ctx contractapi.TransactionContextInterface, voterID string) (*Voter, error) {
	voterJSON, err := ctx.GetStub().GetState(voterID)
	if err != nil {
		return nil, fmt.Errorf("failed to read voter: %v", err)
	}
	if voterJSON == nil {
		return nil, fmt.Errorf("voter does not exist")
	}

	var voter Voter
	err = json.Unmarshal(voterJSON, &voter)
	if err != nil {
		return nil, err
	}

	return &voter, nil
}

func (s *SmartContract) GetElection(ctx contractapi.TransactionContextInterface, electionID string) (*Election, error) {
	electionJSON, err := ctx.GetStub().GetState(electionID)
	if err != nil {
		return nil, fmt.Errorf("failed to read election: %v", err)
	}
	if electionJSON == nil {
		return nil, fmt.Errorf("election does not exist")
	}

	var election Election
	err = json.Unmarshal(electionJSON, &election)
	if err != nil {
		return nil, err
	}

	return &election, nil
}

func (s *SmartContract) GetVoterElectionRelation(ctx contractapi.TransactionContextInterface, voterID string, electionID string) (*VoterElectionRelation, error) {
	key := fmt.Sprintf("relation_%s_%s", voterID, electionID)
	relationJSON, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, fmt.Errorf("failed to read relation: %v", err)
	}
	if relationJSON == nil {
		return nil, fmt.Errorf("relation does not exist")
	}

	var relation VoterElectionRelation
	err = json.Unmarshal(relationJSON, &relation)
	if err != nil {
		return nil, err
	}

	return &relation, nil
}

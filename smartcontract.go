package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Donation
type SmartContract struct {
	contractapi.Contract
}

// Donation describes basic details of what makes up a simple donation
// Insert struct field in alphabetic order => to achieve determinism across languages
// rder when marshal golang keeps the oto json but doesn't order automatically
type Donation struct {
	AppraisedValue int    `json:"AppraisedValue"`
	DonationType   string `json:"DonationType"`
	ID             string `json:"ID"`
	Donor          string `json:"Donor"`
	Size           int    `json:"Size"`
	Timestamp      int    `jason:"Timestamp"`
	Status         string `jason:"Status"`
}

// InitLedger adds a base set of donations to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Donation{
		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
	}

	for _, donation := range donations {
		donationJSON, err := json.Marshal(donation)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(donation.ID, donationJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the donation %s already exists", id)
	}

	donation := Donation{
		ID:             id,
		DonationType:   donationType,
		Size:           size,
		Donor:          donor,
		AppraisedValue: appraisedValue,
		Timestamp:      timestamp,
		Status:         status,
	}
	donationJSON, err := json.Marshal(donation)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, donationJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if donationJSON == nil {
		return nil, fmt.Errorf("the donation %s does not exist", id)
	}

	var donation Donation
	err = json.Unmarshal(donationJSON, &donation)
	if err != nil {
		return nil, err
	}

	return &donation, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the donation %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Donation{
		ID:             id,
		DonationType:   donationType,
		Size:           size,
		Donor:          donor,
		AppraisedValue: appraisedValue,
		Timestamp:      timestamp,
		Status:         status,
	}
	donationJSON, err := json.Marshal(donation)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, donationJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the donation %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return donationJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldDonor := donation.Donor
	donation.Donor = newDonor

	donationJSON, err := json.Marshal(donation)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, donationJSON)
	if err != nil {
		return "", err
	}

	return oldDonor, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all donations in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var donations []*Donation
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var donation Donation
		err = json.Unmarshal(queryResponse.Value, &donation)
		if err != nil {
			return nil, err
		}
		donations = append(donations, &donation)
	}

	return donations, nil
}

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
}

// InitLedger adds a base set of donations to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	donations := []Donation{
		{ID: "donation1", DonationType: "money", Size: 0, Donor: "Tomoko", AppraisedValue: 300},
		{ID: "donation2", DonationType: "money", Size: 0, Donor: "Brad", AppraisedValue: 400},
		{ID: "donation3", DonationType: "money", Size: 0, Donor: "Jin Soo", AppraisedValue: 500},
		{ID: "donation4", DonationType: "money", Size: 0, Donor: "Max", AppraisedValue: 600},
		{ID: "donation5", DonationType: "ssd", Size: 1, Donor: "Adriana", AppraisedValue: 700},
		{ID: "donation6", DonationType: "keyboard", Size: 5, Donor: "Michel", AppraisedValue: 800},
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

// CreateDonation issues a new donation to the world state with given details.
func (s *SmartContract) CreateDonation(ctx contractapi.TransactionContextInterface, id string, donationType string, size int, donor string, appraisedValue int) error {
	exists, err := s.DonationExists(ctx, id)
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
	}
	donationJSON, err := json.Marshal(donation)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, donationJSON)
}

// ReadDonation returns the donation stored in the world state with given id.
func (s *SmartContract) ReadDonation(ctx contractapi.TransactionContextInterface, id string) (*Donation, error) {
	donationJSON, err := ctx.GetStub().GetState(id)
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

// UpdateDonation updates an existing donation in the world state with provided parameters.
func (s *SmartContract) UpdateDonation(ctx contractapi.TransactionContextInterface, id string, donationType string, size int, donor string, appraisedValue int) error {
	exists, err := s.DonationExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the donation %s does not exist", id)
	}

	// overwriting original donation with new donation
	donation := Donation{
		ID:             id,
		DonationType:   donationType,
		Size:           size,
		Donor:          donor,
		AppraisedValue: appraisedValue,
	}
	donationJSON, err := json.Marshal(donation)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, donationJSON)
}

// DeleteDonation deletes an given donation from the world state.
func (s *SmartContract) DeleteDonation(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.DonationExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the donation %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// DonationExists returns true when donation with given ID exists in world state
func (s *SmartContract) DonationExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	donationJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return donationJSON != nil, nil
}

// TransferDonation updates the donor field of donation with given id in world state, and returns the old donor.
func (s *SmartContract) TransferDonation(ctx contractapi.TransactionContextInterface, id string, newDonor string) (string, error) {
	donation, err := s.ReadDonation(ctx, id)
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

// GetAllDonations returns all donations found in world state
func (s *SmartContract) GetAllDonations(ctx contractapi.TransactionContextInterface) ([]*Donation, error) {
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

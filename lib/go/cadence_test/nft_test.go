package test

import (
	"testing"

	"github.com/onflow/cadence"
	"github.com/stretchr/testify/assert"
)

func TestDeployingContracts(t *testing.T) {
	blockchain := newBlockchain(t)
	standardContractsAddress := deployStandardContracts(t, blockchain)
	deployRedSquirrelNftContract(t, blockchain, standardContractsAddress)
}

func TestGettingTotalSupply(t *testing.T) {
	blockchain := newBlockchain(t)
	standardContractAddress := deployStandardContracts(t, blockchain)
	redSquirrelNftAddress, _ := deployRedSquirrelNftContract(t, blockchain, standardContractAddress)
	totalSupply := getTotalSupply(t, blockchain, redSquirrelNftAddress)
	assert.Equal(t, totalSupply, cadence.NewUInt64(0))
}

func TestMintingAnRedSquirrelNFT(t *testing.T) {
	blockchain := newBlockchain(t)
	standardContractAddress := deployStandardContracts(t, blockchain)
	redSquirrelNftAddress, signer := deployRedSquirrelNftContract(t, blockchain, standardContractAddress)

	setUpTestAccountToReceiveRedSquirrelNfts(t, blockchain, standardContractAddress, redSquirrelNftAddress, signer)

	totalSupply := getTotalSupply(t, blockchain, redSquirrelNftAddress)
	assert.EqualValues(t, cadence.NewUInt64(0), totalSupply)

	mint(t, blockchain, standardContractAddress, redSquirrelNftAddress, signer, redSquirrelNftAddress, "test name", "test description", "test thumbnail")

	totalSupply = getTotalSupply(t, blockchain, redSquirrelNftAddress)
	assert.EqualValues(t, cadence.NewUInt64(1), totalSupply)
}

func TestGetRedSquirrel(t *testing.T) {
	type RedSquirrelMetadata struct {
		RedSquirrelID uint64
		ResourceID    uint64
		Name          string
		Description   string
		Thumbnail     string
		Owner         string
		NftType       string
	}

	blockchain := newBlockchain(t)
	standardContractAddress := deployStandardContracts(t, blockchain)
	redSquirrelNftAddress, signer := deployRedSquirrelNftContract(t, blockchain, standardContractAddress)

	setUpTestAccountToReceiveRedSquirrelNfts(t, blockchain, standardContractAddress, redSquirrelNftAddress, signer)

	testName := "test name"
	testDescription := "test description"
	testThumbnail := "test thumbnail"
	mint(t, blockchain, standardContractAddress, redSquirrelNftAddress, signer, redSquirrelNftAddress, testName, testDescription, testThumbnail)

	redSquirrelID := uint64(0)
	redSquirrel := getRedSquirrel(t, blockchain, standardContractAddress, redSquirrelNftAddress, redSquirrelNftAddress, redSquirrelID)
	redSquirrelAsCadenceStruct := redSquirrel.(cadence.Struct)

	var actualMetadata RedSquirrelMetadata
	for fieldIndex, field := range redSquirrelAsCadenceStruct.StructType.Fields {
		fieldValue := redSquirrelAsCadenceStruct.Fields[fieldIndex].ToGoValue()
		switch field.Identifier {
		case "redSquirrelID":
			actualMetadata.RedSquirrelID = fieldValue.(uint64)
		case "resourceID":
			actualMetadata.ResourceID = fieldValue.(uint64)
		case "name":
			actualMetadata.Name = fieldValue.(string)
		case "description":
			actualMetadata.Description = fieldValue.(string)
		case "thumbnail":
			actualMetadata.Thumbnail = fieldValue.(string)
		case "owner":
			addressBytes := fieldValue.([8]uint8)
			addressString := cadence.NewAddress(addressBytes).Hex()
			actualMetadata.Owner = addressString
		case "type":
			actualMetadata.NftType = fieldValue.(string)
		}
	}

	expectedMetadata := RedSquirrelMetadata{
		RedSquirrelID: uint64(0),
		ResourceID:  uint64(31),
		Name:        testName,
		Description: testDescription,
		Thumbnail:   "ipfs://" + testThumbnail,
		Owner:       redSquirrelNftAddress.Hex(),
		NftType:     "A." + redSquirrelNftAddress.Hex() + ".RedSquirrelNFT.NFT",
	}

	assert.Equal(t, actualMetadata, expectedMetadata)
}

package test

import (
	"red-squirrel-nft/utils"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-emulator/types"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
	sdkTestHelper "github.com/onflow/flow-go-sdk/test"
	"github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/onflow/cadence/encoding/json"
)

func newBlockchain(t *testing.T) *emulator.Blockchain {
	blockchain, err := emulator.NewBlockchain()
	assert.NoError(t, err)
	return blockchain
}

func deployStandardContracts(t *testing.T, blockchain *emulator.Blockchain) flow.Address {
	nftContract := templates.Contract{
		Name:   "NonFungibleToken",
		Source: string(contracts.NonFungibleToken()),
	}
	metadataViewsCode, getContractErr := utils.GetMetadataContractCode()
	assert.NoError(t, getContractErr)
	metadataViewsContract := templates.Contract{
		Name:   "MetadataViews",
		Source: string(metadataViewsCode),
	}

	contracts := []templates.Contract{nftContract, metadataViewsContract}

	address, accountErr := blockchain.CreateAccount(nil, contracts)
	assert.NoError(t, accountErr)
	_, commitBlockErr := blockchain.CommitBlock()
	assert.NoError(t, commitBlockErr)

	return address
}

func deployRedSquirrelNftContract(t *testing.T, blockchain *emulator.Blockchain, standardContractsAddress flow.Address) (flow.Address, crypto.Signer) {
	redSquirrelNftContractCode, getContractErr := utils.GetRedSquirrelNftContract(standardContractsAddress)
	assert.NoError(t, getContractErr)
	redSquirrelNftContract := templates.Contract{
		Name:   "RedSquirrelNFT",
		Source: string(redSquirrelNftContractCode),
	}

	accountKeys := sdkTestHelper.AccountKeyGenerator()
	accountKey, signer := accountKeys.NewWithSigner()
	address, createAccountErr := blockchain.CreateAccount(
		[]*flow.AccountKey{accountKey},
		[]templates.Contract{redSquirrelNftContract},
	)
	assert.NoError(t, createAccountErr)
	_, commitBlockErr := blockchain.CommitBlock()
	assert.NoError(t, commitBlockErr)

	return address, signer
}

func getTotalSupply(t *testing.T, blockchain *emulator.Blockchain, address flow.Address) cadence.Value {
	totalSupplyScript, totalSupplyScriptErr := utils.GetTotalSupplyScript(address)
	assert.NoError(t, totalSupplyScriptErr)

	return executeScript(
		t,
		blockchain,
		totalSupplyScript,
		[][]byte{},
	)
}

func setUpTestAccountToReceiveRedSquirrelNfts(
	t *testing.T,
	blockchain *emulator.Blockchain,
	standardContractsAddress,
	redSquirrelNftAddress flow.Address,
	userSigner crypto.Signer,
) {
	txCode, txCodeErr := utils.GetSetUpAccountTransactionCode(standardContractsAddress, redSquirrelNftAddress)
	assert.NoError(t, txCodeErr)

	tx := flow.NewTransaction().
		SetScript([]byte(txCode)).
		SetGasLimit(100).
		SetProposalKey(blockchain.ServiceKey().Address, blockchain.ServiceKey().Index, blockchain.ServiceKey().SequenceNumber).
		SetPayer(blockchain.ServiceKey().Address).
		AddAuthorizer(redSquirrelNftAddress)

	signAndSubmitTransaction(
		t,
		blockchain,
		tx,
		[]flow.Address{blockchain.ServiceKey().Address, redSquirrelNftAddress},
		[]crypto.Signer{blockchain.ServiceKey().Signer(), userSigner},
		false,
	)
}

func mint(
	t *testing.T,
	blockchain *emulator.Blockchain,
	standardContractsAddress,
	redSquirrelNftAddress flow.Address,
	redSquirrelNftSigner crypto.Signer,
	recipientAddress flow.Address,
	name,
	description,
	thumbnailPath string,
) {
	mintTxCode, mintTxCodeErr := utils.GetMintTransactionCode(standardContractsAddress, redSquirrelNftAddress)
	assert.NoError(t, mintTxCodeErr)

	cadenceName, nameErr := cadence.NewString(name)
	assert.NoError(t, nameErr)
	cadenceDescription, descriptionErr := cadence.NewString(description)
	assert.NoError(t, descriptionErr)
	cadenceThumbnailPath, thumbnailPathErr := cadence.NewString(thumbnailPath)
	assert.NoError(t, thumbnailPathErr)

	tx := flow.NewTransaction().
		SetScript(mintTxCode).
		SetGasLimit(100).
		SetProposalKey(blockchain.ServiceKey().Address, blockchain.ServiceKey().Index, blockchain.ServiceKey().SequenceNumber).
		SetPayer(blockchain.ServiceKey().Address).
		AddAuthorizer(redSquirrelNftAddress)

	tx.AddArgument(cadence.NewAddress(redSquirrelNftAddress))
	tx.AddArgument(cadenceName)
	tx.AddArgument(cadenceDescription)
	tx.AddArgument(cadenceThumbnailPath)

	signAndSubmitTransaction(
		t,
		blockchain,
		tx,
		[]flow.Address{blockchain.ServiceKey().Address, redSquirrelNftAddress},
		[]crypto.Signer{blockchain.ServiceKey().Signer(), redSquirrelNftSigner},
		false,
	)
}

func getRedSquirrel(
	t *testing.T,
	blockchain *emulator.Blockchain,
	standardContractsAddress,
	redSquirrelNftAddress,
	ownerAddress flow.Address,
	redSquirrelID uint64,
) cadence.Value {
	script, scriptErr := utils.GetRedSquirrelScript(standardContractsAddress, redSquirrelNftAddress)
	assert.NoError(t, scriptErr)

	return executeScript(
		t,
		blockchain,
		script,
		[][]byte{
			json.MustEncode(cadence.NewAddress(ownerAddress)),
			json.MustEncode(cadence.NewUInt64(redSquirrelID)),
		},
	)
}

func executeScript(t *testing.T, blockchain *emulator.Blockchain, script []byte, arguments [][]byte) cadence.Value {
	result, err := blockchain.ExecuteScript(script, arguments)
	assert.NoError(t, err)
	assert.NoError(t, result.Error)

	return result.Value
}

func signAndSubmitTransaction(
	t *testing.T,
	b *emulator.Blockchain,
	tx *flow.Transaction,
	signerAddresses []flow.Address,
	signers []crypto.Signer,
	shouldRevert bool,
) *types.TransactionResult {
	for i := len(signerAddresses) - 1; i >= 0; i-- {
		signerAddress := signerAddresses[i]
		signer := signers[i]

		if i == 0 {
			err := tx.SignEnvelope(signerAddress, 0, signer)
			assert.NoError(t, err)
		} else {
			err := tx.SignPayload(signerAddress, 0, signer)
			assert.NoError(t, err)
		}
	}

	return submitTransaction(t, b, tx, shouldRevert)
}

func submitTransaction(
	t *testing.T,
	blockchain *emulator.Blockchain,
	tx *flow.Transaction,
	shouldRevert bool,
) *types.TransactionResult {
	err := blockchain.AddTransaction(*tx)
	assert.NoError(t, err)

	result, err := blockchain.ExecuteNextTransaction()
	assert.NoError(t, err)

	if shouldRevert {
		assert.True(t, result.Reverted())
	} else {
		if !assert.True(t, result.Succeeded()) {
			t.Log(result.Error.Error())
		}
	}

	_, err = blockchain.CommitBlock()
	assert.NoError(t, err)

	return result
}

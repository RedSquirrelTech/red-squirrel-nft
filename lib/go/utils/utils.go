package utils

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/onflow/flow-go-sdk"
)

var (
	_, file, _, _    = runtime.Caller(0)
	basepath         = filepath.Dir(file)
	cadenceRootPath  = basepath + "/../../../cadence"
	contractsPath    = cadenceRootPath + "/contracts"
	transactionsPath = cadenceRootPath + "/transactions"
	scriptsPath      = cadenceRootPath + "/scripts"
)

func GetMetadataContractCode() ([]byte, error) {
	metadataContractPath := "../../../cadence/contracts/standard/MetadataViews.cdc"
	return ioutil.ReadFile(metadataContractPath)
}

func GetRedSquirrelNftContract(standardContractsAddress flow.Address) ([]byte, error) {
	redSquirrelNftContractPath := "../../../cadence/contracts/RedSquirrelNFT.cdc"
	contractCode, err := ioutil.ReadFile(redSquirrelNftContractPath)
	if err != nil {
		return nil, err
	}
	codeWithStandardContractsAddress := ReplaceStandardContractsAddress(standardContractsAddress, contractCode)

	return codeWithStandardContractsAddress, nil
}

func ReplaceStandardContractsAddress(address flow.Address, cadenceCode []byte) []byte {
	var MetadataViewsAddressPlaceholder = regexp.MustCompile(`"[^"\s].*/MetadataViews.cdc"`)
	var NonFungibleTokenAddressPlaceholder = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)

	var addressAsString string
	if !strings.Contains(address.Hex(), "0x") {
		addressAsString = "0x" + address.Hex()
	}

	codeWithNonFungibleTokenAddress := NonFungibleTokenAddressPlaceholder.ReplaceAllString(string(cadenceCode), addressAsString)
	codeWithMetadataViewsAddress := MetadataViewsAddressPlaceholder.ReplaceAllString(string(codeWithNonFungibleTokenAddress), addressAsString)
	return []byte(codeWithMetadataViewsAddress)
}

func ReplaceRedSquirrelNFTContractAddress(address flow.Address, cadenceCode []byte) []byte {
	var RedSquirrelNFTPlaceholder = regexp.MustCompile(`"[^"\s].*/RedSquirrelNFT.cdc"`)

	var addressAsString string
	if !strings.Contains(address.Hex(), "0x") {
		addressAsString = "0x" + address.Hex()
	}

	codeWithRedSquirrelNFTAddress := RedSquirrelNFTPlaceholder.ReplaceAllString(string(cadenceCode), addressAsString)
	return []byte(codeWithRedSquirrelNFTAddress)
}

func GetTotalSupplyScript(address flow.Address) ([]byte, error) {
	scriptCode, err := ioutil.ReadFile(scriptsPath + "/get_total_supply.cdc")
	if err != nil {
		return nil, err
	}
	scriptWithAddress := ReplaceRedSquirrelNFTContractAddress(address, scriptCode)
	return scriptWithAddress, nil
}

func GetSetUpAccountTransactionCode(standardContractsAddress, redSquirrelNftAddress flow.Address) ([]byte, error) {
	txCode, err := ioutil.ReadFile(transactionsPath + "/set_up_account.cdc")
	if err != nil {
		return nil, err
	}
	txCodeWithStandardAddresses := ReplaceStandardContractsAddress(standardContractsAddress, txCode)
	txCodeWithRedSquirrelNftAddress := ReplaceRedSquirrelNFTContractAddress(redSquirrelNftAddress, txCodeWithStandardAddresses)

	return txCodeWithRedSquirrelNftAddress, nil
}

func GetMintTransactionCode(standardContractsAddress, redSquirrelNftAddress flow.Address) ([]byte, error) {
	txCode, err := ioutil.ReadFile(transactionsPath + "/mint.cdc")
	if err != nil {
		return nil, err
	}
	txCodeWithStandardAddresses := ReplaceStandardContractsAddress(standardContractsAddress, txCode)
	txCodeWithRedSquirrelNftAddress := ReplaceRedSquirrelNFTContractAddress(redSquirrelNftAddress, txCodeWithStandardAddresses)
	return txCodeWithRedSquirrelNftAddress, nil
}

func GetRedSquirrelScript(standardContractsAddress, redSquirrelNftAddress flow.Address) ([]byte, error) {
	scriptCode, err := ioutil.ReadFile(scriptsPath + "/get_red_squirrel.cdc")
	if err != nil {
		return nil, err
	}
	scriptCodeWithStandardAddresses := ReplaceStandardContractsAddress(standardContractsAddress, scriptCode)
	scriptCodeWithRedSquirrelNftAddress := ReplaceRedSquirrelNFTContractAddress(redSquirrelNftAddress, scriptCodeWithStandardAddresses)
	return scriptCodeWithRedSquirrelNftAddress, nil
}

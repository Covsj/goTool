package evm

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

func CreateNewWallet() (priKey string, address string, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	priKey = hexutil.Encode(privateKeyBytes)[2:]
	priKey = strings.TrimPrefix(priKey, "0x")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return priKey, address, nil
}

func CreateNewWalletByMnemonic() (mnemonic, priKey, address string, err error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", "", "", err
	}
	mnemonic, err = bip39.NewMnemonic(entropy)
	if err != nil {
		return "", "", "", err
	}

	seed := bip39.NewSeed(mnemonic, "")
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return "", "", "", err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, true)
	if err != nil {
		return "", "", "", err
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return "", "", "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	priKey = hexutil.Encode(privateKeyBytes)[2:]
	priKey = strings.TrimPrefix(priKey, "0x")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return mnemonic, priKey, address, nil
}

func CalculateLastWord(mnemonicWords []string) (string, error) {
	if len(mnemonicWords) != 11 {
		return "", errors.New("mnemonicWords not 11 length")
	}

	var wordList = bip39.GetWordList()

	var binaryString string

	for _, word := range mnemonicWords {
		index, found := findIndex(wordList, word)
		if !found {
			return "", fmt.Errorf("word '%s' not found in BIP-39 word list", word)
		}
		binaryString += fmt.Sprintf("%011b", index)
	}
	hash := sha256.Sum256(binaryStringToBytes(binaryString))
	checksum := fmt.Sprintf("%08b", hash[0])[:4]
	fullBinaryString := binaryString + checksum
	lastWordIndex, err := binaryToDecimal(fullBinaryString[len(fullBinaryString)-11:])
	if err != nil {
		return "", err
	}
	return wordList[lastWordIndex], nil
}

func GetWalletNonce(h Handler, address string) (uint64, error) {
	cli := h.Client()
	if cli == nil {
		return 0, errors.New("client empty")
	}
	return cli.PendingNonceAt(context.Background(), common.HexToAddress(address))
}

func GetAddressByPrivateKey(hexPrivateKey string) (string, error) {
	hexPrivateKey = strings.TrimPrefix(hexPrivateKey, "0x")

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex(), nil
}

func GetAddressByMnemonic(mnemonic string) (address, priKey string, err error) {
	seed := bip39.NewSeed(mnemonic, "")
	wallet, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		return "", "", err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, true)
	if err != nil {
		return "", "", err
	}

	privateKey, err := wallet.PrivateKey(account)
	if err != nil {
		return "", "", err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	priKey = hexutil.Encode(privateKeyBytes)[2:]
	priKey = strings.TrimPrefix(priKey, "0x")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return priKey, address, nil
}

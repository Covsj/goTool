package basic

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

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
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", "", err
	}

	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		return "", "", "", err
	}

	key := masterKey
	for _, n := range path {
		key, err = key.DeriveNonStandard(n)
		if err != nil {
			return "", "", "", err
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return "", "", "", err
	}

	privateKeyECDSA := privateKey.ToECDSA()

	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	priKey = strings.TrimPrefix(hexutil.Encode(privateKeyBytes)[2:], "0x")

	address = crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()
	return mnemonic, priKey, address, nil
}

func CalculateLastWord(mnemonicWords []string) (string, error) {
	if len(mnemonicWords) != 11 {
		return "", errors.New("mnemonicWords not 11 length")
	}

	var wordList = bip39.GetWordList()

	var binaryString string

	for _, word := range mnemonicWords {
		index, found := FindIndex(wordList, word)
		if !found {
			return "", fmt.Errorf("word '%s' not found in BIP-39 word list", word)
		}
		binaryString += fmt.Sprintf("%011b", index)
	}
	hash := sha256.Sum256(BinaryStringToBytes(binaryString))
	checksum := fmt.Sprintf("%08b", hash[0])[:4]
	fullBinaryString := binaryString + checksum
	lastWordIndex, err := BinaryToDecimal(fullBinaryString[len(fullBinaryString)-11:])
	if err != nil {
		return "", err
	}
	return wordList[lastWordIndex], nil
}

func FindIndex(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func BinaryToDecimal(binaryStr string) (int, error) {
	num := new(big.Int)
	num, ok := num.SetString(binaryStr, 2)
	if !ok {
		return 0, fmt.Errorf("invalid binary number: %s", binaryStr)
	}
	return int(num.Int64()), nil
}

func BinaryStringToBytes(s string) []byte {
	bi := new(big.Int)
	bi.SetString(s, 2)
	return bi.Bytes()
}

func ConvertBigIntToString(balance *big.Int, decimals uint8) string {
	if decimals == 0 {
		decimals = 18
	}
	// Convert balance to float
	balanceFloat := new(big.Float).SetInt(balance)

	// Calculate divisor
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))

	// Divide balance by divisor
	balanceFloat.Quo(balanceFloat, divisor)

	// Convert float balance to string
	balanceStr := balanceFloat.Text('f', int(decimals))

	return balanceStr
}

func ConvertStringToBigInt(balanceStr string, decimals uint8) (*big.Int, error) {
	if decimals == 0 {
		decimals = 18
	}
	// Convert balance string to big.Float
	balanceFloat, ok := new(big.Float).SetString(balanceStr)
	if !ok {
		return nil, fmt.Errorf("failed to convert balance to big.Float")
	}

	// Multiply balance by 10^decimals
	exp := big.NewInt(int64(decimals))
	multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), exp, nil))
	balanceFloat.Mul(balanceFloat, multiplier)

	// Convert balance to big.Int
	balanceInt := new(big.Int)
	balanceFloat.Int(balanceInt)

	return balanceInt, nil
}

func SignText(privateKeyHex, message string) (string, error) {
	// Ensure the private key is in the correct format
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	// Decode the hex string to a byte slice
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %v", err)
	}

	// Convert to ECDSA private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("cannot parse private key: %v", err)
	}

	// Format the message like MetaMask does
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	// Adjust the v value for Ethereum's replay-protected transaction format (EIP-155)
	if signature[64] < 27 {
		signature[64] += 27
	}

	return hexutil.Encode(signature), nil
}

package basic

import (
	"encoding/hex"
	"fmt"
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

func SignText(privateKeyHex, message string) (string, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0X")

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

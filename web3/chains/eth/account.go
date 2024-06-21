package eth

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/Covsj/goTool/web3/chains/basic"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

type Account struct {
	*Util
	privateKeyECDSA *ecdsa.PrivateKey
	address         string
}

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		return nil, err
	}

	key := masterKey
	for _, n := range path {
		key, err = key.DeriveNonStandard(n)
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	privateKeyECDSA := privateKey.ToECDSA()
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()

	return &Account{
		Util:            NewUtil(),
		privateKeyECDSA: privateKeyECDSA,
		address:         address,
	}, nil
}

func NewAccountWithMnemonicIndex(mnemonic string, accountIndex string) (*Account, error) {
	if accountIndex == "" {
		accountIndex = "0"
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/" + accountIndex)
	if err != nil {
		return nil, err
	}

	key := masterKey
	for _, n := range path {
		key, err = key.DeriveNonStandard(n)
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	privateKeyECDSA := privateKey.ToECDSA()
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()

	return &Account{
		Util:            NewUtil(),
		privateKeyECDSA: privateKeyECDSA,
		address:         address,
	}, nil
}

func NewAccountWithPrivateKey(privateKey string) (*Account, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	priData, err := types.HexDecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	privateKeyECDSA, err := crypto.ToECDSA(priData)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()

	return &Account{
		Util:            NewUtil(),
		privateKeyECDSA: privateKeyECDSA,
		address:         address,
	}, nil
}

func (a *Account) PrivateKey() ([]byte, error) {
	return crypto.FromECDSA(a.privateKeyECDSA), nil
}

func (a *Account) PrivateKeyHex() (string, error) {
	bytes := crypto.FromECDSA(a.privateKeyECDSA)
	return types.HexEncodeToString(bytes), nil
}

func (a *Account) PublicKey() []byte {
	return crypto.FromECDSAPub(&a.privateKeyECDSA.PublicKey)
}

func (a *Account) PublicKeyHex() string {
	bytes := crypto.FromECDSAPub(&a.privateKeyECDSA.PublicKey)
	return types.HexEncodeToString(bytes)
}

func (a *Account) Address() string {
	return a.address
}

func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	dataString := string(message)
	hashBytes := SignHashForMsg(dataString)
	return a.SignHash(hashBytes)
}

func (a *Account) SignHex(messageHex string, password string) (*basic.OptionalString, error) {
	data, err := types.HexDecodeString(messageHex)
	if err != nil {
		return nil, err
	}
	signed, err := a.Sign(data, password)
	if err != nil {
		return nil, err
	}
	signedString := types.HexEncodeToString(signed)
	return &basic.OptionalString{Value: signedString}, nil
}

func (a *Account) SignHash(hash []byte) ([]byte, error) {
	signature, err := crypto.Sign(hash, a.privateKeyECDSA)
	if err != nil {
		return nil, err
	}
	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// SignHashForMsg 以太坊的 hash 专门在数据前面加上了一段话
func SignHashForMsg(data string) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func VerifySignature(pubkey, message, signedMsg string) bool {
	pubBytes, err := types.HexDecodeString(pubkey)
	if err != nil {
		return false
	}
	originBytes, err := types.HexDecodeString(message)
	if err != nil {
		return false
	}
	originMsgHash := SignHashForMsg(string(originBytes))

	signedBytes, err := types.HexDecodeString(signedMsg)
	if err != nil {
		return false
	}
	signedBytes = signedBytes[:len(signedBytes)-1]
	return crypto.VerifySignature(pubBytes, originMsgHash, signedBytes)
}

func IsValidSignature(publicKey, msg, signature []byte) bool {
	originMsgHash := SignHashForMsg(string(msg))
	signature = signature[:len(signature)-1]
	return crypto.VerifySignature(publicKey, originMsgHash, signature)
}

func AsEthereumAccount(account basic.Account) *Account {
	if r, ok := account.(*Account); ok {
		return r
	} else {
		return nil
	}
}

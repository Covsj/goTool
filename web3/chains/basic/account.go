package basic

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Handler interface {
	ChainId() int64
	GasLimit() uint64
	ChainName() string
	Client() *ethclient.Client
	RpcClient() *rpc.Client
	NativeBalance(address common.Address) (*big.Int, error)
	TokenBalance(tokenAddress common.Address, ownerAddress common.Address) (*big.Int, error)
	SendTransaction(toAddress common.Address, amount *big.Int, gasLimit uint64, data []byte, result interface{}) (common.Hash, error)
	SendEth(toAddress common.Address, amount *big.Int, gasLimit uint64) (common.Hash, error)
	SendToken(tokenAddress common.Address, toAddress common.Address, amount *big.Int) (common.Hash, error)
	TokenDecimals(tokenAddress common.Address) (uint8, error)
	CallContract(contractAddress common.Address, abiStr, method string, results interface{},
		amount *big.Int, params ...interface{}) (common.Hash, error)
	QueryTransactionStatus(txHash common.Hash) (isPending bool, status bool, err error)
}

type AddressUtil interface {
	// EncodePublicKeyToAddress @param publicKey can start with 0x or not.
	EncodePublicKeyToAddress(publicKey string) (string, error)
	// DecodeAddressToPublicKey @return publicKey that will start with 0x.
	DecodeAddressToPublicKey(address string) (string, error)

	IsValidAddress(address string) bool
}

type Account interface {
	// PrivateKey @return privateKey data
	PrivateKey() ([]byte, error)
	// PrivateKeyHex @return privateKey string that will start with 0x.
	PrivateKeyHex() (string, error)

	// PublicKey @return publicKey data
	PublicKey() []byte
	// PublicKeyHex @return publicKey string that will start with 0x.
	PublicKeyHex() string

	// Address @return address string
	Address() string

	Sign(message []byte, password string) ([]byte, error)
	SignHex(messageHex string, password string) (*OptionalString, error)
}

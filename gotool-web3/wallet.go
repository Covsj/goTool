package gotool_web3

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"

	ERC20 "github.com/Covsj/goTool/gotool-lib/erc20"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetTokenBalance(chain string, address string, token string) (res string, err error) {
	var client *ethclient.Client

	client = GetRpcClientByName(chain)

	// 获取原生token
	if token == "" {
		balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
		if err != nil {
			return "", err
		}
		return balance.String(), nil
	} else {
		tokenAddress := common.HexToAddress(token)
		instance, err := ERC20.NewUtils(tokenAddress, client)
		if err != nil {
			return "", err
		}
		addr := common.HexToAddress(address)

		balance, err := instance.BalanceOf(&bind.CallOpts{}, addr)
		if err != nil {
			return "", err
		}
		return balance.String(), nil
	}
}

func CreateNewWallet() (priKey string, address string, err error) {
	privateKey, err := crypto.GenerateKey()

	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	priKey = hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return priKey, address, nil
}

func SignText(privateKeyHex, text string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}

	// Hash the text
	hash := crypto.Keccak256Hash([]byte(text))

	// Sign the hash
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(signature), nil
}

func CreateTransaction(fromAddress, chainId string, to, gasLimit, gasPrice, value, data string) (*types.Transaction, error) {
	client := GetRpcClientByChainId(chainId)
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(fromAddress))
	if err != nil {
		return nil, err
	}
	//    gasPrice := "20000000000" // 例如：20 Gwei
	//    gasLimit := "21000"       // 例如：21000
	//    value := "1000000000000000000" // 例如：1 ETH
	if gasLimit == "" {
		gasLimit = "21000"
	}
	if gasPrice == "" {
		gasPrice = "20000000000"
	}
	valueBigInt, _ := new(big.Int).SetString(value, 10)
	gasLimitUint64, _ := strconv.ParseUint(gasLimit, 10, 64)
	gasPriceBigInt, _ := new(big.Int).SetString(gasPrice, 10)

	tx := types.NewTransaction(nonce, common.HexToAddress(to), valueBigInt, gasLimitUint64, gasPriceBigInt, []byte(data))

	return tx, nil
}

func SignTransactionData(privateKeyHex string, tx *types.Transaction, chainID string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}

	chainIDBigInt, _ := new(big.Int).SetString(chainID, 10)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainIDBigInt), privateKey)
	if err != nil {
		return "", err
	}

	rlpEncodedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(rlpEncodedTx), nil
}

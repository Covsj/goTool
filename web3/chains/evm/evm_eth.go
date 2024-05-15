package evm

import (
	"context"
	"errors"
	"math/big"
	"strings"

	ERC20 "github.com/Covsj/goTool/lib/erc20"
	"github.com/Covsj/goTool/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthHandler struct {
	privateKey string
	mnemonic   string
	Address    common.Address
	client     *ethclient.Client
	rpcClient  *rpc.Client
	Rpc        string
}

func NewEthHandler(rpc string, privateKey, mnemonic string) (*EthHandler, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")

	if rpc == "" {
		rpc = RpcEth
	}
	p := &EthHandler{
		Rpc:        rpc,
		privateKey: privateKey,
		mnemonic:   mnemonic,
	}
	if privateKey != "" {
		p.privateKey = privateKey
		address, err := GetAddressByPrivateKey(privateKey)
		if err != nil {
			log.Error("初始化EthHandler获取地址失败", "error", err.Error())
			return nil, err
		} else {
			p.Address = common.HexToAddress(address)
		}
	} else if mnemonic != "" {
		address, priKey, err := GetAddressByMnemonic(mnemonic)
		if err != nil {
			log.Error("初始化EthHandler获取地址失败", "error", err.Error())
			return nil, err
		} else {
			p.Address = common.HexToAddress(address)
			p.privateKey = priKey
		}
	} else {
		mnemonic, key, address, err := CreateNewWalletByMnemonic()
		if err != nil {
			log.Error("初始化EthHandler随机生成钱包账户失败", "error", err.Error())
			return nil, err
		} else {
			p.Address = common.HexToAddress(address)
			p.privateKey = key
			p.mnemonic = mnemonic
		}
	}

	client, rpcClient, err := CreateClient(rpc)
	if err != nil {
		log.Error("初始化EthHandler获取客户端失败", "error", err.Error())
		return nil, err
	}
	p.rpcClient = rpcClient
	p.client = client
	return p, err
}

func (h *EthHandler) ChainId() int64 {
	id, err := h.Client().ChainID(context.Background())
	if err != nil {
		return 1
	}
	return id.Int64()
}

func (h *EthHandler) DefaultGasLimit() uint64 {
	return 21_0000
}

func (h *EthHandler) ChainName() string {
	return "ETH"
}

func (h *EthHandler) Client() *ethclient.Client {
	return h.client
}

func (h *EthHandler) RpcClient() *rpc.Client {
	return h.rpcClient
}

// GetEthBalance 获取以太币余额
func (h *EthHandler) GetEthBalance(address common.Address) (*big.Int, error) {
	return h.Client().BalanceAt(context.Background(), address, nil)
}

// GetTokenBalance 获取代币余额
func (h *EthHandler) GetTokenBalance(tokenAddress common.Address, ownerAddress common.Address) (*big.Int, error) {
	instance, err := ERC20.NewUtils(tokenAddress, h.Client())
	if err != nil {
		return nil, err
	}
	return instance.BalanceOf(nil, ownerAddress)
}

// SendTransaction 发送交易
func (h *EthHandler) SendTransaction(toAddress common.Address,
	amount *big.Int, gasLimit uint64, data []byte, result interface{}) (common.Hash, error) {
	privateKey, err := crypto.HexToECDSA(h.privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	gasPrice, err := h.Client().SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	nonce, err := h.Client().PendingNonceAt(context.Background(), h.Address)
	if err != nil {
		return common.Hash{}, err
	}

	if gasLimit <= 0 {
		gasLimit = h.DefaultGasLimit()
	}

	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(h.ChainId())), privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	signedTxData, err := signedTx.MarshalBinary()
	if err != nil {
		return common.Hash{}, err
	}
	err = h.RpcClient().CallContext(context.Background(), result,
		"eth_sendRawTransaction", hexutil.Encode(signedTxData))
	if err != nil {
		return common.Hash{}, err
	}
	return signedTx.Hash(), nil
}

// SendEth 发送以太币
func (h *EthHandler) SendEth(toAddress common.Address, amount *big.Int, gasLimit uint64) (common.Hash, error) {
	return h.SendTransaction(toAddress, amount, gasLimit, nil, nil)
}

// SendToken 发送代币
func (h *EthHandler) SendToken(tokenAddress common.Address, toAddress common.Address, amount *big.Int) (common.Hash, error) {
	privateKey, err := crypto.HexToECDSA(h.privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	tokenInstance, err := ERC20.NewUtils(tokenAddress, h.Client())
	if err != nil {
		return common.Hash{}, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	tx, err := tokenInstance.Transfer(auth, toAddress, amount)
	if err != nil {
		return common.Hash{}, err
	}

	return tx.Hash(), nil
}

func (h *EthHandler) GetTokenDecimals(tokenAddress common.Address) (uint8, error) {
	opts := &bind.CallOpts{}
	token, err := ERC20.NewUtils(tokenAddress, h.Client())
	if err != nil {
		return 18, err
	}
	decimals, err := token.Decimals(opts)
	if err != nil {
		return 18, err
	}
	return decimals, nil
}

func (h *EthHandler) CallContract(contractAddress common.Address, abiStr,
	method string, results interface{}, amount *big.Int, params ...interface{}) (common.Hash, error) {

	// 加载合约ABI
	contractABI, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return common.Hash{}, err
	}

	txData, err := contractABI.Pack(method, params...)
	if err != nil {
		return common.Hash{}, err
	}

	return h.SendTransaction(contractAddress, amount, 0, txData, results)
}

func (h *EthHandler) QueryTransactionStatus(txHash common.Hash) (isPending bool, status bool, err error) {
	// Use Ethereum client to get pending transaction
	tx, isPending, err := h.Client().TransactionByHash(context.Background(), txHash)
	if err != nil {
		return false, false, err
	}

	// If transaction is pending, return true
	if isPending {
		return true, false, nil
	}

	// If transaction is not pending, check if it has receipt
	if tx != nil {
		receipt, err := h.Client().TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return false, false, err
		}

		// Check if the transaction receipt exists
		if receipt != nil {
			// Transaction receipt found, check if the transaction was successful
			return false, receipt.Status == types.ReceiptStatusSuccessful, nil
		}
	}

	// If neither pending nor receipt found, return false
	return false, false, errors.New("NOT FOUND TRANSACTION")
}

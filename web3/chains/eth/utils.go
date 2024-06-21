package eth

import (
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"math/big"
	"strconv"
	"strings"

	"github.com/Covsj/goTool/web3/chains/basic"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

func OptsTobiInt(opts *CallMethodOpts) *CallMethodOptsBigInt {
	GasPrice, _ := new(big.Int).SetString(opts.GasPrice, 10)
	GasLimit, _ := strconv.Atoi(opts.GasLimit)
	MaxPriorityFeePerGas, _ := new(big.Int).SetString(opts.MaxPriorityFeePerGas, 10)
	Value, _ := new(big.Int).SetString(opts.Value, 10)

	return &CallMethodOptsBigInt{
		Nonce:                uint64(opts.Nonce),
		Value:                Value,
		GasPrice:             GasPrice,
		MaxPriorityFeePerGas: MaxPriorityFeePerGas,
		GasLimit:             uint64(GasLimit),
	}
}

func PrivateKeyToAddress(privateKey string) (string, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}
	return crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex(), nil
}

func EncodeErc20Transfer(toAddress, amount string) ([]byte, error) {
	if !common.IsHexAddress(toAddress) {
		return nil, errors.New("invalid receiver address")
	}
	amountInt, valid := big.NewInt(0).SetString(amount, 10)
	if !valid {
		return nil, errors.New("invalid transfer amount")
	}
	return EncodeContractData(Erc20AbiStr, Erc20MethodTransfer, common.HexToAddress(toAddress), amountInt)
}

func EncodeErc20Approve(spender string, amount *big.Int) ([]byte, error) {
	if !common.IsHexAddress(spender) {
		return nil, errors.New("Invalid receiver address")
	}
	return EncodeContractData(Erc20AbiStr, Erc20MethodApprove, common.HexToAddress(spender), amount)
}

func EncodeErc721TransferFrom(sender, receiver, nftId string) ([]byte, error) {
	if !common.IsHexAddress(sender) || !common.IsHexAddress(receiver) {
		return nil, basic.ErrInvalidAddress
	}
	amtInt, ok := big.NewInt(0).SetString(nftId, 10)
	if !ok {
		return nil, basic.ErrInvalidAmount
	}
	return EncodeContractData(Erc721Abi_TransferOnly, "transferFrom",
		common.HexToAddress(sender),
		common.HexToAddress(receiver),
		amtInt)
}

func EncodeContractData(abiString, method string, params ...interface{}) ([]byte, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}
	return parsedAbi.Pack(method, params...)
}

func DecodeContractParams(abiString string, data []byte) (string, []interface{}, error) {
	if len(data) <= 4 {
		return "", nil, nil
	}
	parsedAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return "", nil, err
	}

	method, err := parsedAbi.MethodById(data[:4])
	if err != nil {
		return "", nil, err
	}

	params, err := method.Inputs.Unpack(data[4:])
	return method.RawName, params, err
}

func CreateNewWalletByMnemonic(index int) (mnemonic, priKey, address string, err error) {
	if index <= 0 {
		index = 0
	}
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

	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/" + strconv.Itoa(index))
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

// CalculateLastWord 根据传入的11个助记词，计算最后助记词
func CalculateLastWord(mnemonicWords []string) (string, error) {
	if len(mnemonicWords) != 11 {
		return "", errors.New("mnemonicWords not 11 length")
	}
	// found own morning
	wordList := bip39.GetWordList()
	m := []string{}
	for _, word := range wordList {
		mnemonic := fmt.Sprintf("%s %s", strings.Join(mnemonicWords, " "), word)
		if bip39.IsMnemonicValid(mnemonic) {
			m = append(m, word)
		}
	}

	if len(m) <= 0 {
		return "", errors.New("not found")
	}
	stringToNumber := func(str string, length int) int {
		h := fnv.New32a()
		_, err := h.Write([]byte(str))
		if err != nil {
			return 0
		}
		hashValue := h.Sum32()
		number := int(hashValue) % length
		return number
	}
	index := stringToNumber(mnemonicWords[len(mnemonicWords)-2]+mnemonicWords[len(mnemonicWords)-1], len(m))
	mnemonic := fmt.Sprintf("%s %s", strings.Join(mnemonicWords, " "), m[index])
	return mnemonic, nil
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

func BalanceBigIntToString(balance *big.Int, decimals uint8) string {
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

	for strings.HasSuffix(balanceStr, "0") {
		balanceStr = strings.TrimSuffix(balanceStr, "0")
	}
	return balanceStr
}

func BalanceStringToBigInt(balanceStr string, decimals uint8) (*big.Int, error) {
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

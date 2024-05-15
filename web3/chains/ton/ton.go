package ton

import (
	"context"
	"encoding/base64"

	log "github.com/Covsj/goTool/log"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

type ChainHandler struct {
	Api ton.APIClientWrapped `json:"api"`
	Ctx context.Context      `json:"ctx"`
}

func NewTonChainHandler() (*ChainHandler, error) {
	client := liteclient.NewConnectionPool()
	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/global.config.json")
	if err != nil {
		log.ErrorF("NewTonChainHandler get config err: %s", err.Error())
		return nil, err
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.ErrorF("NewTonChainHandler connection err: %s", err.Error())
		return nil, err
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicySecure).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	// bound all requests to single ton node
	ctx := client.StickyContext(context.Background())

	return &ChainHandler{
		Api: api,
		Ctx: ctx,
	}, nil
}

func (p *ChainHandler) CreateNewWallet(words []string) (*wallet.Wallet, []string, error) {
	if len(words) == 0 {
		words = wallet.NewSeed()
	}

	w, err := wallet.FromSeed(p.Api, words, wallet.V4R2)
	if err != nil {
		log.Error("CreateNewWallet FromSeed err: %s", err.Error())
		return nil, words, err
	}
	log.InfoF("CreateNewWallet wallet address: %s", w.WalletAddress())

	return w, words, nil

}

func (p *ChainHandler) GetTonBalance(checkAddr string) (balance tlb.Coins, err error) {

	log.InfoF("GetWalletBalance fetching and checking proofs since config init block, it may take near a minute...")
	// we need fresh block info to run get methods
	b, err := p.Api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.ErrorF("GetTonBalance get master chain info err: %s", err.Error())
		return
	}
	log.InfoF("GetWalletBalance get master chain info success, Work chain:%d, SeqNo:%d", b.Workchain, b.SeqNo)

	addr := address.MustParseAddr(checkAddr)

	// we use WaitForBlock to make sure block is ready,
	// it is optional but escapes us from liteserver block not ready errors
	//res, err := p.Api.WaitForBlock(b.SeqNo).GetAccount(p.Ctx, b, addr)
	res, err := p.Api.GetAccount(p.Ctx, b, addr)
	if err != nil {
		log.ErrorF("GetTonBalance get account err: %s", err.Error())
		return
	}
	if res.IsActive {
		//fmt.Printf("Status: %s\n", res.State.Status)
		//fmt.Printf("Balance: %s TON\n", res.State.Balance.String())
		//if res.Data != nil {
		//	fmt.Printf("Data: %s\n", res.Data.Dump())
		//}
		balance = res.State.Balance
	}
	return

}

func (p *ChainHandler) Transfer(w *wallet.Wallet, amount tlb.Coins, to string) error {

	addr := address.MustParseAddr(to)

	// if destination wallet is not initialized (or you don't care)
	// you should set bounce to false to not get money back.
	// If bounce is true, money will be returned in case of not initialized destination wallet or smart-contract error
	// bounce := false

	transfer, err := w.BuildTransfer(addr, amount, true, "")
	if err != nil {
		log.ErrorF("Transfer err: %s", err.Error())
		return err
	}

	tx, block, err := w.SendWaitTransaction(p.Ctx, transfer)
	if err != nil {
		log.ErrorF("Transfer SendWaitTransaction err: %s", err.Error())
		return err
	}

	balance, err := w.GetBalance(p.Ctx, block)
	if err != nil {
		log.ErrorF("Transfer GetBalance err: %s", err.Error())
		return err
	}

	log.InfoF("Transfer transaction confirmed at block %d, hash: %s balance left: %s", block.SeqNo,
		base64.StdEncoding.EncodeToString(tx.Hash), balance.String())

	return nil
}

package http

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	ss := &T{}
	resp, err := DoRequest(&ReqOpt{Method: "GET", Url: "https://gmgn.ai/defi/quotation/v1/signals?size=10", RespOut: &ss})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp.String())
}

type T struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Signals []struct {
			Id              int         `json:"id"`
			Timestamp       int         `json:"timestamp"`
			Maker           interface{} `json:"maker"`
			TokenAddress    string      `json:"token_address"`
			TokenPrice      float64     `json:"token_price"`
			FromTimestamp   int         `json:"from_timestamp,omitempty"`
			UpdatedAt       int         `json:"updated_at,omitempty"`
			BuyDuration     int         `json:"buy_duration,omitempty"`
			BuyUsd          float64     `json:"buy_usd,omitempty"`
			TxCount         int         `json:"tx_count,omitempty"`
			SignalType      string      `json:"signal_type"`
			SmartBuy        int         `json:"smart_buy"`
			SmartSell       int         `json:"smart_sell"`
			Signal1HCount   int         `json:"signal_1h_count"`
			FirstEntryPrice float64     `json:"first_entry_price"`
			PriceChange     float64     `json:"price_change"`
			Token           struct {
				Id                         int     `json:"id"`
				Chain                      string  `json:"chain"`
				Address                    string  `json:"address"`
				AntiWhaleModifiable        int     `json:"anti_whale_modifiable"`
				BuyTax                     int     `json:"buy_tax"`
				CannotBuy                  int     `json:"cannot_buy"`
				CanTakeBackOwnership       int     `json:"can_take_back_ownership"`
				CreatorAddress             string  `json:"creator_address"`
				CreatorBalance             float64 `json:"creator_balance"`
				CreatorPercent             float64 `json:"creator_percent"`
				ExternalCall               int     `json:"external_call"`
				HiddenOwner                int     `json:"hidden_owner"`
				HolderCount                int     `json:"holder_count"`
				HoneypotWithSameCreator    int     `json:"honeypot_with_same_creator"`
				IsAntiWhale                int     `json:"is_anti_whale"`
				IsBlacklisted              int     `json:"is_blacklisted"`
				IsHoneypot                 int     `json:"is_honeypot"`
				IsInDex                    int     `json:"is_in_dex"`
				IsMintable                 int     `json:"is_mintable"`
				IsOpenSource               int     `json:"is_open_source"`
				IsProxy                    int     `json:"is_proxy"`
				IsWhitelisted              int     `json:"is_whitelisted"`
				OwnerAddress               string  `json:"owner_address"`
				OwnerBalance               int     `json:"owner_balance"`
				OwnerChangeBalance         int     `json:"owner_change_balance"`
				OwnerPercent               int     `json:"owner_percent"`
				PersonalSlippageModifiable int     `json:"personal_slippage_modifiable"`
				Selfdestruct               int     `json:"selfdestruct"`
				SellTax                    float64 `json:"sell_tax"`
				SlippageModifiable         int     `json:"slippage_modifiable"`
				TokenName                  string  `json:"token_name"`
				TokenSymbol                string  `json:"token_symbol"`
				TotalSupply                int64   `json:"total_supply"`
				TradingCooldown            int     `json:"trading_cooldown"`
				TransferPausable           int     `json:"transfer_pausable"`
				HasResult                  int     `json:"has_result"`
				UpdatedAt                  int     `json:"updated_at"`
				LpHolders                  []struct {
					Tag        string      `json:"tag"`
					Value      interface{} `json:"value"`
					Address    string      `json:"address"`
					Balance    string      `json:"balance"`
					Percent    string      `json:"percent"`
					NFTList    interface{} `json:"NFT_list"`
					IsLocked   int         `json:"is_locked"`
					IsContract int         `json:"is_contract"`
				} `json:"lp_holders"`
				LpTotalSupply      float64     `json:"lp_total_supply"`
				FakeToken          interface{} `json:"fake_token"`
				CannotSellAll      int         `json:"cannot_sell_all"`
				LpHolderCount      int         `json:"lp_holder_count"`
				Symbol             string      `json:"symbol"`
				Name               string      `json:"name"`
				Decimals           int         `json:"decimals"`
				Logo               *string     `json:"logo"`
				CmcId              interface{} `json:"cmc_id"`
				CoingeckoId        interface{} `json:"coingecko_id"`
				CirculatingSupply  interface{} `json:"circulating_supply"`
				MaxSupply          interface{} `json:"max_supply"`
				CreationBlock      interface{} `json:"creation_block"`
				CreationTimestamp  interface{} `json:"creation_timestamp"`
				Liquidity          float64     `json:"liquidity"`
				OpenTimestamp      int         `json:"open_timestamp"`
				Price              float64     `json:"price"`
				Price1M            float64     `json:"price_1m"`
				Price5M            float64     `json:"price_5m"`
				Price1H            float64     `json:"price_1h"`
				Price6H            float64     `json:"price_6h"`
				Price24H           float64     `json:"price_24h"`
				HighPrice          float64     `json:"high_price"`
				HighPriceTimestamp int         `json:"high_price_timestamp"`
				LowPrice           interface{} `json:"low_price"`
				LowPriceTimestamp  interface{} `json:"low_price_timestamp"`
				Volume1M           float64     `json:"volume_1m"`
				Volume5M           float64     `json:"volume_5m"`
				Volume1H           float64     `json:"volume_1h"`
				Volume6H           float64     `json:"volume_6h"`
				Volume24H          float64     `json:"volume_24h"`
				BuyVolume1M        float64     `json:"buy_volume_1m"`
				BuyVolume5M        float64     `json:"buy_volume_5m"`
				BuyVolume1H        float64     `json:"buy_volume_1h"`
				BuyVolume6H        float64     `json:"buy_volume_6h"`
				BuyVolume24H       float64     `json:"buy_volume_24h"`
				SellVolume1M       float64     `json:"sell_volume_1m"`
				SellVolume5M       float64     `json:"sell_volume_5m"`
				SellVolume1H       float64     `json:"sell_volume_1h"`
				SellVolume6H       float64     `json:"sell_volume_6h"`
				SellVolume24H      float64     `json:"sell_volume_24h"`
				Swaps1M            int         `json:"swaps_1m"`
				Swaps5M            int         `json:"swaps_5m"`
				Swaps1H            int         `json:"swaps_1h"`
				Swaps6H            int         `json:"swaps_6h"`
				Swaps24H           int         `json:"swaps_24h"`
				Buys1M             int         `json:"buys_1m"`
				Buys5M             int         `json:"buys_5m"`
				Buys1H             int         `json:"buys_1h"`
				Buys6H             int         `json:"buys_6h"`
				Buys24H            int         `json:"buys_24h"`
				Sells1M            int         `json:"sells_1m"`
				Sells5M            int         `json:"sells_5m"`
				Sells1H            int         `json:"sells_1h"`
				Sells6H            int         `json:"sells_6h"`
				Sells24H           int         `json:"sells_24h"`
				BiggestPoolAddress string      `json:"biggest_pool_address"`
				Standard           interface{} `json:"standard"`
				DextScore          interface{} `json:"dext_score"`
				Renounced          int         `json:"renounced"`
				HotLevel           int         `json:"hot_level"`
				IsShowAlert        bool        `json:"is_show_alert"`
				TokenAddress       string      `json:"token_address"`
				Fdv                float64     `json:"fdv"`
			} `json:"token"`
			Link struct {
				TwitterUsername string  `json:"twitter_username,omitempty"`
				Website         *string `json:"website,omitempty"`
				Telegram        string  `json:"telegram,omitempty"`
			} `json:"link"`
			RecentBuys struct {
				SmartWallets     int           `json:"smart_wallets"`
				SmartBuyUsd      int           `json:"smart_buy_usd"`
				FollowingWallets int           `json:"following_wallets"`
				FollowingBuyUsd  int           `json:"following_buy_usd"`
				BuyTimestamp     int           `json:"buy_timestamp"`
				BuyList          []interface{} `json:"buy_list"`
			} `json:"recent_buys"`
			PreviousSignals []struct {
				Id                int         `json:"id"`
				Timestamp         int         `json:"timestamp"`
				Maker             interface{} `json:"maker"`
				TokenAddress      string      `json:"token_address"`
				TokenPrice        float64     `json:"token_price"`
				FromTimestamp     int         `json:"from_timestamp,omitempty"`
				UpdatedAt         int         `json:"updated_at,omitempty"`
				BuyDuration       int         `json:"buy_duration,omitempty"`
				BuyUsd            float64     `json:"buy_usd,omitempty"`
				TxCount           int         `json:"tx_count,omitempty"`
				SignalType        string      `json:"signal_type"`
				PriceReview       float64     `json:"price_review,omitempty"`
				PriceChangeReview float64     `json:"price_change_review,omitempty"`
				ReviewToTimestamp int         `json:"review_to_timestamp,omitempty"`
				TurnoverRate      float64     `json:"turnover_rate,omitempty"`
				MarketCap         string      `json:"market_cap,omitempty"`
			} `json:"previous_signals"`
			IsFirst           bool    `json:"is_first"`
			PriceReview       float64 `json:"price_review,omitempty"`
			PriceChangeReview float64 `json:"price_change_review,omitempty"`
			ReviewToTimestamp int     `json:"review_to_timestamp,omitempty"`
		} `json:"signals"`
		Next string `json:"next"`
	} `json:"data"`
}

package eth

import (
	"github.com/Covsj/goTool/web3/chains/basic"
)

var (
	_ basic.Account     = (*Account)(nil)
	_ basic.Chain       = (*Chain)(nil)
	_ basic.Token       = (*Token)(nil)
	_ basic.Token       = (*Erc20Token)(nil)
	_ basic.Transaction = (*Transaction)(nil)

	_ basic.Token = (*Src20Token)(nil)
)

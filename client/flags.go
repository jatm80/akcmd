package client

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gookit/gcli/v3"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

func AddQueryFlagsToCmd(cmd *gcli.Command) {
	cmd.Str(flags.FlagChainID, "", "", "The network chain ID")
	cmd.Str(flags.FlagNode, "", "tcp://localhost:26657",
		"<host>:<port> to Tendermint RPC interface for this chain")
	cmd.Int64(flags.FlagHeight, "", 0,
		"Use a specific height to query state at (this can error if the node is pruning state)")
	cmd.Str(tmcli.OutputFlag, "o", "text", "Output format (text|json)")

	cmd.Required(flags.FlagChainID)
}

func AddPaginationFlagsToCmd(cmd *gcli.Command, query string) {
	cmd.Uint64(flags.FlagPage, "", 1,
		fmt.Sprintf("pagination page of %s to query for. This sets offset to a multiple of limit", query))
	cmd.Str(flags.FlagPageKey, "", "",
		fmt.Sprintf("pagination page-key of %s to query for", query))
	cmd.Uint64(flags.FlagOffset, "", 0,
		fmt.Sprintf("pagination offset of %s to query for", query))
	cmd.Uint64(flags.FlagLimit, "", 100,
		fmt.Sprintf("pagination limit of %s to query for", query))
	cmd.Bool(flags.FlagCountTotal, "", false,
		fmt.Sprintf("count total number of records in %s to query for", query))
}

package client

import (
	"fmt"
	"math"

	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gookit/gcli/v3"
	dtypes "github.com/ovrclk/akash/x/deployment/types"
	mtypes "github.com/ovrclk/akash/x/market/types"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

var (
	qryOpts              = QueryOpts{}
	pageOpts             = PaginationOpts{}
	deploymentFilterOpts = dtypes.DeploymentFilters{}
	deploymentIDOpts     = dtypes.DeploymentID{}
	gSeq                 uint64
	oSeq                 uint64
)

type QueryOpts struct {
	ChainID string
	Node    string
	Height  int64
	Output  string
}

type PaginationOpts struct {
	Page       uint64
	PageKey    string
	Offset     uint64
	Limit      uint64
	CountTotal bool
	Reverse    bool
}

func AddQueryFlagsToCmd(cmd *gcli.Command) {
	cmd.StrOpt(&qryOpts.ChainID, flags.FlagChainID, "", "", "The network chain ID")
	cmd.StrOpt(&qryOpts.Node, flags.FlagNode, "", "tcp://localhost:26657",
		"<host>:<port> to Tendermint RPC interface for this chain")
	cmd.Int64Opt(&qryOpts.Height, flags.FlagHeight, "", 0,
		"Use a specific height to query state at (this can error if the node is pruning state)")
	cmd.StrOpt(&qryOpts.Output, tmcli.OutputFlag, "o", "text", "Output format (text|json)")

	cmd.Required(flags.FlagChainID)
}

func QueryFlagsFromCmd() QueryOpts {
	return qryOpts
}

func AddPaginationFlagsToCmd(cmd *gcli.Command, query string) {
	cmd.Uint64Opt(&pageOpts.Page, flags.FlagPage, "", 1,
		fmt.Sprintf("pagination page of %s to query for. This sets offset to a multiple of limit", query))
	cmd.StrOpt(&pageOpts.PageKey, flags.FlagPageKey, "", "",
		fmt.Sprintf("pagination page-key of %s to query for", query))
	cmd.Uint64Opt(&pageOpts.Offset, flags.FlagOffset, "", 0,
		fmt.Sprintf("pagination offset of %s to query for", query))
	cmd.Uint64Opt(&pageOpts.Limit, flags.FlagLimit, "", 100,
		fmt.Sprintf("pagination limit of %s to query for", query))
	cmd.BoolOpt(&pageOpts.CountTotal, flags.FlagCountTotal, "", false,
		fmt.Sprintf("count total number of records in %s to query for", query))
	cmd.BoolOpt(&pageOpts.Reverse, flags.FlagReverse, "", false,
		"results are sorted in descending order")
}

func PaginationOptsFromCmd() PaginationOpts {
	return pageOpts
}

func AddDeploymentFilterFlags(cmd *gcli.Command) {
	cmd.StrOpt(&deploymentFilterOpts.Owner, "owner", "", "", "deployment owner address to filter")
	cmd.StrOpt(&deploymentFilterOpts.State, "state", "", "", "deployment state to filter (active,closed)")
	cmd.Uint64Opt(&deploymentFilterOpts.DSeq, "dseq", "", 0, "deployment sequence to filter")
}

func DepFiltersFromFlags() dtypes.DeploymentFilters {
	return deploymentFilterOpts
}

func AddDeploymentIDFlags(cmd *gcli.Command) {
	cmd.StrOpt(&deploymentIDOpts.Owner, "owner", "", "", "Deployment Owner")
	cmd.Uint64Opt(&deploymentIDOpts.DSeq, "dseq", "", 0, "Deployment Sequence")
}

func MarkReqDeploymentIDFlags(cmd *gcli.Command) {
	cmd.Required("owner", "dseq")
}

func DeploymentIDFromFlags() dtypes.DeploymentID {
	return deploymentIDOpts
}

func AddGroupIDFlags(cmd *gcli.Command) {
	AddDeploymentIDFlags(cmd)
	cmd.Uint64Opt(&gSeq, "gseq", "", 1, "Group Sequence")
}

func MarkReqGroupIDFlags(cmd *gcli.Command) {
	MarkReqDeploymentIDFlags(cmd)
}

func GroupIDFromFlags() (dtypes.GroupID, error) {
	dID := DeploymentIDFromFlags()
	val, err := getGSeq()
	if err != nil {
		return dtypes.GroupID{}, err
	}
	return dtypes.MakeGroupID(dID, val), nil
}

func getGSeq() (uint32, error) {
	if gSeq > math.MaxUint32 {
		return 0, errors.New("gseq out of uint32 range")
	}
	return uint32(gSeq), nil
}

func AddOrderFilterFlags(cmd *gcli.Command) {
	cmd.StrOpt(&deploymentFilterOpts.Owner, "owner", "", "", "order owner address to filter")
	cmd.StrOpt(&deploymentFilterOpts.State, "state", "", "", "order state to filter (open,matched,closed)")
	cmd.Uint64Opt(&deploymentFilterOpts.DSeq, "dseq", "", 0, "deployment sequence to filter")
	cmd.Uint64Opt(&gSeq, "gseq", "", 1, "group sequence to filter")
	cmd.Uint64Opt(&oSeq, "oseq", "", 1, "order sequence to filter")
}

func OrderFiltersFromFlags() (mtypes.OrderFilters, error) {
	dFilters := DepFiltersFromFlags()
	filter := mtypes.OrderFilters{
		Owner: dFilters.Owner,
		State: dFilters.State,
		DSeq:  dFilters.DSeq,
	}
	var err error
	if filter.GSeq, err = getGSeq(); err != nil {
		return filter, err
	}
	if filter.OSeq, err = getOSeq(); err != nil {
		return filter, err
	}
	return filter, nil
}

func getOSeq() (uint32, error) {
	if oSeq > math.MaxUint32 {
		return 0, errors.New("oseq out of uint32 range")
	}
	return uint32(oSeq), nil
}

func AddOrderIDFlags(cmd *gcli.Command) {
	AddGroupIDFlags(cmd)
	cmd.Uint64Opt(&oSeq, "oseq", "", 1, "Order Sequence")
}

func MarkReqOrderIDFlags(cmd *gcli.Command) {
	MarkReqGroupIDFlags(cmd)
}

func OrderIDFromFlags() (mtypes.OrderID, error) {
	gID, err := GroupIDFromFlags()
	if err != nil {
		return mtypes.OrderID{}, err
	}
	val, err := getOSeq()
	if err != nil {
		return mtypes.OrderID{}, err
	}
	return mtypes.MakeOrderID(gID, val), nil
}

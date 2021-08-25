package client

import (
	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func createContext() client.Context {
	iregistry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(iregistry)
	banktypes.RegisterInterfaces(iregistry)
	stakingtypes.RegisterInterfaces(iregistry)
	vestingtypes.RegisterInterfaces(iregistry)
	cryptocodec.RegisterInterfaces(iregistry)

	cctx := client.Context{}
	cctx = cctx.WithOffline(false)
	cctx = cctx.WithInterfaceRegistry(iregistry)

	return cctx
}

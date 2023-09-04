package pool

import (
	"context"

	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

type (
	CommunityPoolKeeper    = distkeeper.Keeper
	CommunityPoolMsgServer = disttypes.MsgServer
)

func FundCommunityPool(ctx context.Context, ms CommunityPoolMsgServer, msg *disttypes.MsgFundCommunityPool) (*disttypes.MsgFundCommunityPoolResponse, error) {
	return ms.FundCommunityPool(ctx, msg)
}

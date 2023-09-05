package keeper

import (
	"context"

	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/pool/types"
)

type Keeper struct {
	storeService storetypes.KVStoreService

	authority string
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService, authority string) Keeper {
	return Keeper{
		storeService: storeService,
		authority:    authority,
	}
}

// GetAuthority returns the x/pool module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With(log.ModuleKey, "x/"+types.ModuleName)
}

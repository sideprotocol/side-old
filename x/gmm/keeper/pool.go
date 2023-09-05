package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	utils "github.com/sideprotocol/side/sideutils"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k Keeper) initializePool(ctx sdk.Context, msg *types.MsgCreatePool) (*string, error) {
	poolCreator := msg.PoolCreator()
	pool := msg.CreatePool()
	totalShares := sdk.NewInt(0)

	poolShareBaseDenom := pool.PoolId //types.GetPoolShareDenom(pool.PoolId)
	poolShareDisplayDenom := fmt.Sprintf("GAMM-%s", pool.PoolId)

	assets := make(map[string]types.PoolAsset)
	for _, liquidity := range msg.Liquidity {
		assets[liquidity.Token.Denom] = liquidity
		totalShares = totalShares.Add(liquidity.Token.Amount)
	}

	// Check pool already created or not
	if _, found := k.GetPool(ctx, pool.PoolId); found {
		return nil, types.ErrAlreadyCreatedPool
	}

	// Check balance.
	for _, liquidity := range msg.Liquidity {
		balance := k.bankKeeper.GetBalance(ctx, poolCreator, liquidity.Token.Denom)
		if balance.Amount.LT(liquidity.Token.Amount) {
			return nil, types.ErrInsufficientBalance
		}
	}

	// Create Module Account
	escrowAccount := types.GetEscrowAddress(pool.PoolId)
	if err := utils.CreateModuleAccount(ctx, k.accountKeeper, escrowAccount); err != nil {
		return nil, err
	}

	// Move asset from creator to module account

	if err := k.LockTokens(ctx, pool.PoolId, msg.PoolCreator(), msg.InitialLiquidity()); err != nil {
		return nil, err
	}

	// Register metadata to bank keeper

	k.bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
		Description: fmt.Sprintf("The share token of the gamm pool %s", pool.GetPoolId()),
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    poolShareBaseDenom,
				Exponent: 0,
				Aliases: []string{
					"attopoolshare",
				},
			},
			{
				Denom:    poolShareDisplayDenom,
				Exponent: types.OneShareExponent,
				Aliases:  nil,
			},
		},
		Base:    poolShareBaseDenom,
		Display: poolShareDisplayDenom,
	})

	// Mint share to creator
	if err := k.MintPoolShareToAccount(ctx, poolCreator, sdk.NewCoin(
		poolShareBaseDenom,
		totalShares,
	)); err != nil {
		return nil, err
	}

	// Save pool to chain
	k.AppendPool(ctx, pool)
	return &pool.PoolId, nil
}

// RemoveInterchainLiquidityPool removes a interchainLiquidityPool from the store
func (k Keeper) RemoveInterchainLiquidityPool(
	ctx sdk.Context,
	poolID string,
) {
	// Get current pool count
	poolCount, found := k.GetCountByPoolID(ctx, poolID)
	if !found {
		return
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolsPrefix)
	store.Delete(GetInterchainLiquidityPoolKey(poolCount))
}

func (k Keeper) GetPoolCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolsPrefix)
	b := store.Get(types.KeyCurrentPoolCountPrefix)
	if b == nil {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func (k Keeper) SetPoolCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolsPrefix)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, count)
	store.Set(types.KeyCurrentPoolCountPrefix, b)
}

func GetInterchainLiquidityPoolKey(count uint64) []byte {
	return []byte(fmt.Sprintf("%020d", count))
}

// GetAllInterchainLiquidityPool returns all interchainLiquidityPool
func (k Keeper) GetAllPool(ctx sdk.Context) (list []types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(string(types.KeyPoolsPrefix)))

	// Start from the latest pool and move to the oldest
	poolCount := k.GetPoolCount(ctx)
	for i := poolCount; i >= 1 && (poolCount-i) < types.MaxPoolCount; i-- {
		b := store.Get(GetInterchainLiquidityPoolKey(i))
		if b == nil {
			continue
		}
		var val types.Pool
		k.cdc.MustUnmarshal(b, &val)
		list = append(list, val)
	}
	return
}

// Sets the mapping between poolId and its count index
func (k Keeper) SetPoolIDToCountMapping(ctx sdk.Context, poolID string, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolIDToCountPrefix)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, count)
	store.Set([]byte(poolID), b)
}

// Gets the count index of the poolId
func (k Keeper) GetCountByPoolID(ctx sdk.Context, poolID string) (count uint64, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolIDToCountPrefix)
	b := store.Get([]byte(poolID))
	if b == nil {
		return 0, false
	}
	return binary.BigEndian.Uint64(b), true
}

// Modified SetInterchainLiquidityPool
func (k Keeper) AppendPool(ctx sdk.Context, pool types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolsPrefix)

	// Get current pool count
	poolCount := k.GetPoolCount(ctx)

	// Increment the count
	poolCount++

	// Set the new count
	k.SetPoolCount(ctx, poolCount)

	// Set the poolId to count mapping
	k.SetPoolIDToCountMapping(ctx, pool.PoolId, poolCount)

	// Marshal the pool and set in store
	b := k.cdc.MustMarshal(&pool)
	store.Set(GetInterchainLiquidityPoolKey(poolCount), b)

	// Check if we exceed max pools
	if poolCount > types.MaxPoolCount {
		// Remove the oldest pool
		store.Delete(GetInterchainLiquidityPoolKey(poolCount - types.MaxPoolCount))
	}
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(string(types.KeyPoolsPrefix)))
	// Get current pool count
	poolCount, found := k.GetCountByPoolID(ctx, pool.PoolId)
	if !found {
		return
	}
	// Marshal the pool and set in store
	b := k.cdc.MustMarshal(&pool)
	store.Set(GetInterchainLiquidityPoolKey(poolCount), b)
}

// Modified GetInterchainLiquidityPool
func (k Keeper) GetPool(ctx sdk.Context, poolID string) (val types.Pool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolsPrefix)

	count, found := k.GetCountByPoolID(ctx, poolID)
	if !found {
		return val, false
	}

	b := store.Get(GetInterchainLiquidityPoolKey(count))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

package keeper

type UTXOViewKeeper interface {
}

type UTXOKeeper interface {
	UTXOViewKeeper
}

// var _ UTXOKeeper = (*BaseUTXOKeeper)(nil)

// type BaseUTXOViewKeeper struct {
// 	UTXOs *collections.IndexedMap[collections.Pair[sdk.AccAddress, string], math.Int, BalancesIndexes]
// }

// func NewBaseUTXOViewKeeper() *BaseUTXOViewKeeper {
// 	// return &BaseUTXOViewKeeper{
// 	// 	UTXOs: &collections.NewIndexedMap(sb, types.BalancesPrefix, "utxos", collections.PairKeyCodec(sdk.AccAddressKey, collections.StringKey), types.BalanceValueCodec, newBalancesIndexes(sb)),
// 	// }
// 	return nil
// }

// type BaseUTXOKeeper struct {
// }

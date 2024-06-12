package keeper

type UTXOViewKeeper interface {
}

type UTXOKeeper interface {
	UTXOViewKeeper
}

var _ UTXOKeeper = (*BaseUTXOKeeper)(nil)

type BaseUTXOViewKeeper struct {
}

func NewBaseUTXOViewKeeper() *BaseUTXOViewKeeper {
	return &BaseUTXOViewKeeper{}
}

type BaseUTXOKeeper struct {
}

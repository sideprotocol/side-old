package bindings

// SideQuery contains osmosis custom queries.
type SideQuery struct {
	// Pool   *Pool   `json:"pool,omitempty"`
	Params *Params `json:"params,omitempty"`
}

// type Pool struct {
// 	PoolId string `json:"pool_id"`
// }

type Params struct{}

// type PoolResponse struct {
// 	Admin string `json:"admin"`
// }

type ParamsResponse struct {
	Params *ParamsRes `json:"params"`
}

type ParamsRes struct {
	PoolCreationFee uint64 `json:"pool_creation_fee"`
}

package types

func FormatICAAccountOwner(chainID string, accountType string) (result string) {
	return chainID + "." + accountType
}

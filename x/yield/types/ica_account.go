package types

func FormatICAAccountOwner(chainId string, accountType string) (result string) {
	return chainId + "." + accountType
}

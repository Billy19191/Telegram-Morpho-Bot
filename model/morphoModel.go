package model

type MorphoResponseEntity struct {
	Data []VaultEntity `json:"data"`
}

type VaultEntity struct {
	VaultName     string  `json:"vaultName"`
	TotalAssetUsd float64 `json:"totalAssetUsd"`
	Liquidity     float64 `json:"liquidity"`
	MyAssetUsd    float64 `json:"myAssetUsd"`
	AvgApy        float64 `json:"avgApy"`
	SharedInVault float64 `json:"sharedInVault"`
}

type MorphoResponseModel struct {
	Data []VaultModel `json:"data"`
}

type VaultModel struct {
	VaultName     string  `json:"vaultName"`
	TotalAssetUsd float64 `json:"totalAssetUsd"`
	Liquidity     float64 `json:"liquidity"`
	MyAssetUsd    float64 `json:"myAssetUsd"`
	AvgApy        float64 `json:"avgApy"`
	SharedInVault float64 `json:"sharedInVault"`
}

package service

import (
	"fmt"

	"github.com/Billy19191/Telegram-Morpho-Bot/model"
)

func EvaluateVaultRisk(vault model.VaultEntity) model.RiskReport {
	metrics := []model.MetricEvaluation{
		evaluateLiquidityRatio(vault),
		evaluateShareConcentration(vault),
		evaluateAPY(vault),
		evaluateVaultUtilization(vault),
	}

	overall := model.StatusSafe
	for _, m := range metrics {
		if m.Status == model.StatusCritical {
			overall = model.StatusCritical
			break
		}
		if m.Status == model.StatusMonitor {
			overall = model.StatusMonitor
		}
	}

	return model.RiskReport{
		OverallStatus: overall,
		Metrics:       metrics,
	}
}

func evaluateLiquidityRatio(vault model.VaultEntity) model.MetricEvaluation {
	var ratio float64
	if vault.MyAssetUsd > 0 {
		ratio = vault.Liquidity / vault.MyAssetUsd
	}

	var status model.RiskStatus

	switch {
	case ratio < 3:
		status = model.StatusCritical
	case ratio < 10:
		status = model.StatusMonitor
	default:
		status = model.StatusSafe
	}

	return model.MetricEvaluation{
		Name:   "Liquidity Ratio",
		Value:  fmt.Sprintf("%.1fx", ratio),
		Status: status,
	}
}

func evaluateShareConcentration(vault model.VaultEntity) model.MetricEvaluation {
	share := vault.SharedInVault

	var status model.RiskStatus

	switch {
	case share > 5:
		status = model.StatusCritical
	case share > 2:
		status = model.StatusMonitor
	default:
		status = model.StatusSafe
	}

	return model.MetricEvaluation{
		Name:   "Share in Vault",
		Value:  fmt.Sprintf("%.2f%%", share),
		Status: status,
	}
}

func evaluateAPY(vault model.VaultEntity) model.MetricEvaluation {
	apy := vault.AvgApy

	var status model.RiskStatus

	switch {
	case apy >= 6 && apy <= 8:
		status = model.StatusSafe
	case (apy >= 5 && apy < 6) || (apy > 8 && apy <= 10):
		status = model.StatusMonitor
	default:
		status = model.StatusCritical
	}

	return model.MetricEvaluation{
		Name:   "Avg APY",
		Value:  fmt.Sprintf("%.2f%%", apy),
		Status: status,
	}
}

func evaluateVaultUtilization(vault model.VaultEntity) model.MetricEvaluation {
	var utilization float64
	if vault.TotalAssetUsd > 0 {
		utilization = ((vault.TotalAssetUsd - vault.Liquidity) / vault.TotalAssetUsd) * 100
	}

	var status model.RiskStatus

	switch {
	case utilization > 90:
		status = model.StatusCritical
	case utilization > 80:
		status = model.StatusMonitor
	default:
		status = model.StatusSafe

	}

	return model.MetricEvaluation{
		Name:   "Vault Utilization",
		Value:  fmt.Sprintf("%.1f%%", utilization),
		Status: status,
	}
}

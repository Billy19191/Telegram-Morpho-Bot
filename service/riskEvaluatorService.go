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
	reason := ""

	for _, m := range metrics {
		if m.Status == model.StatusCritical {
			overall = model.StatusCritical
			if reason != "" {
				reason += "\n"
			}
			reason += m.Reason
			continue
		}
		if m.Status == model.StatusMonitor {
			if overall != model.StatusCritical {
				overall = model.StatusMonitor
			}
			if reason != "" {
				reason += "\n"
			}
			reason += m.Reason
		}
	}

	return model.RiskReport{
		OverallStatus: overall,
		Metrics:       metrics,
		Reason:        reason,
	}
}

func evaluateLiquidityRatio(vault model.VaultEntity) model.MetricEvaluation {
	var ratio float64
	if vault.MyAssetUsd > 0 {
		ratio = vault.Liquidity / vault.MyAssetUsd
	}

	var status model.RiskStatus
	var reason string

	switch {
	case ratio < 3:
		status = model.StatusCritical
		reason = "Liquidity is less than 3x your position — withdrawal may be difficult"
	case ratio < 10:
		status = model.StatusMonitor
		reason = "Liquidity is below 10x your position — keep monitoring"
	default:
		status = model.StatusSafe
		reason = "Sufficient liquidity available for your position"
	}

	return model.MetricEvaluation{
		Name:   "Liquidity Ratio",
		Value:  fmt.Sprintf("%.1fx", ratio),
		Status: status,
		Reason: reason,
	}
}

func evaluateShareConcentration(vault model.VaultEntity) model.MetricEvaluation {
	share := vault.SharedInVault

	var status model.RiskStatus
	var reason string

	switch {
	case share > 5:
		status = model.StatusCritical
		reason = "You hold over 5% of the vault — high concentration risk"
	case share > 2:
		status = model.StatusMonitor
		reason = "You hold over 2% of the vault — moderate concentration"
	default:
		status = model.StatusSafe
		reason = "Your share in the vault is well diversified"
	}

	return model.MetricEvaluation{
		Name:   "Share in Vault",
		Value:  fmt.Sprintf("%.2f%%", share),
		Status: status,
		Reason: reason,
	}
}

func evaluateAPY(vault model.VaultEntity) model.MetricEvaluation {
	apy := vault.NetApy

	var status model.RiskStatus
	var reason string

	switch {
	case apy >= 6 && apy <= 8:
		status = model.StatusSafe
		reason = "APY is within the expected 6-8% range"
	case (apy >= 5 && apy < 6) || (apy > 8 && apy <= 10):
		status = model.StatusMonitor
		reason = "APY is slightly outside the normal range — could indicate changing conditions"
	default:
		status = model.StatusCritical
		reason = "APY is significantly abnormal — may signal unsustainable yield or issues"
	}

	return model.MetricEvaluation{
		Name:   "Current APY",
		Value:  fmt.Sprintf("%.2f%%", apy),
		Status: status,
		Reason: reason,
	}
}

func evaluateVaultUtilization(vault model.VaultEntity) model.MetricEvaluation {
	var utilization float64
	if vault.TotalAssetUsd > 0 {
		utilization = ((vault.TotalAssetUsd - vault.Liquidity) / vault.TotalAssetUsd) * 100
	}

	var status model.RiskStatus
	var reason string

	switch {
	case utilization > 90:
		status = model.StatusCritical
		reason = "Over 90% of vault assets are utilized — very low available liquidity"
	case utilization > 80:
		status = model.StatusMonitor
		reason = "Utilization exceeds 80% — liquidity may tighten"
	default:
		status = model.StatusSafe
		reason = "Vault utilization is at a healthy level"
	}

	return model.MetricEvaluation{
		Name:   "Vault Utilization",
		Value:  fmt.Sprintf("%.2f", utilization),
		Status: status,
		Reason: reason,
	}
}

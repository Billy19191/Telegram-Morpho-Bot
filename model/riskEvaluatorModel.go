package model

type RiskStatus string

const (
	StatusSafe     RiskStatus = "🟢 Safe"
	StatusMonitor  RiskStatus = "🟡 Keep Monitor"
	StatusCritical RiskStatus = "🔴 Critical"
)

type MetricEvaluation struct {
	Name   string
	Value  string
	Status RiskStatus
}

type RiskReport struct {
	OverallStatus RiskStatus
	Metrics       []MetricEvaluation
}

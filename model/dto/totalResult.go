package dto

type TotalResult struct {
	OverviewResult   Overview   `json:"overview"`
	StatisticsResult Statistics `json:"statistics"`
	AnalysisResult   Analysis   `json:"analysis"`
}

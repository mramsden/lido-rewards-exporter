package lido

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Event struct {
	ID                     string `json:"id"`
	TotalPooledEtherBefore string `json:"totalPooledEtherBefore"`
	TotalPooledEtherAfter  string `json:"totalPooledEtherAfter"`
	TotalSharesBefore      string `json:"totalSharesBefore"`
	Block                  string `json:"block"`
	BlockTime              string `json:"blockTime"`
	LogIndex               string `json:"logIndex"`
	Type                   string `json:"type"`
	ReportShares           string `json:"reportShares"`
	Balance                string `json:"balance"`
	Rewards                string `json:"rewards"`
	Change                 string `json:"change"`
	CurrencyChange         string `json:"currencyChange"`
	EpochDays              string `json:"epochDays"`
	EpochFullDays          string `json:"epochFullDays"`
	APR                    string `json:"apr"`
}

type RewardsReport struct {
	Events []Event `json:"events"`
	Totals struct {
		EthRewards      string `json:"ethRewards"`
		CurrencyRewards string `json:"currencyRewards"`
	} `json:"totals"`
	AverageApr string `json:"averageApr"`
	TotalItems uint64 `json:"totalItems"`
}

type Currency string

const (
	CurrencyGBP Currency = "GBP"
	CurrencyUSD          = "USD"
	CurrencyEUR          = "EUR"
)

type RewardsReportParams struct {
	Currency    Currency
	ArchiveRate bool
	OnlyRewards bool
}

var DefaultRewardsReportParams = RewardsReportParams{
	Currency:    CurrencyGBP,
	ArchiveRate: true,
	OnlyRewards: true,
}

var ErrInvalidAddress = fmt.Errorf("invalid address")
var ErrInvalidCurrency = fmt.Errorf("invalid currency")

// FetchRewardsReport returns a report of all rewards for the given address
func FetchRewardsReport(address string) (RewardsReport, error) {
	return FetchRewardsReportParams(address, DefaultRewardsReportParams)
}

func FetchRewardsReportParams(address string, params RewardsReportParams) (RewardsReport, error) {
	var report RewardsReport
	var errorPayload struct {
		Message string `json:"message"`
	}

	url := fmt.Sprintf("https://stake.lido.fi/api/rewards?address=%s&currency=%s&archiveRate=%v&onlyRewards=%v", address, params.Currency, params.ArchiveRate, params.OnlyRewards)

	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return report, err
	}

	d := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		err = d.Decode(&errorPayload)
	} else if resp.StatusCode == http.StatusBadRequest {
		err = d.Decode(&errorPayload)
	} else {
		err = d.Decode(&report)
	}
	if err != nil {
		return report, err
	}

	switch errorPayload.Message {
	case "Address is invalid.":
		return report, ErrInvalidAddress
	case "Currency is invalid.":
		return report, ErrInvalidCurrency
	}
	if errorPayload.Message != "" {
		return report, fmt.Errorf(errorPayload.Message)
	}

	return report, nil
}

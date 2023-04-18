package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/mramsden/lido-rewards-exporter/lido"
)

var logger = log.New(os.Stdout, "[Handler] ", log.LstdFlags)

func convertRewardsReportToKoinlyCSV(report lido.RewardsReport) (string, error) {
	b := bytes.NewBufferString("")
	w := csv.NewWriter(b)

	w.Write([]string{"Koinly Date", "Amount", "Currency", "Label", "TxHash"})
	for _, event := range report.Events {
		if blockTime, err := strconv.ParseInt(event.BlockTime, 10, 64); err == nil {

			date := time.Unix(blockTime, 0).UTC().Format("2006-01-02 15:04:05")
			if err != nil {
				log.Printf("error parsing rewards: %v", err)
				return "", err
			}
			w.Write([]string{
				date,
				toDecimal(event.Rewards, 18).Round(8).String(),
				"stETH",
				"reward",
				event.ID,
			})
		}
	}

	return b.String(), nil
}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	_, file := path.Split(r.URL.Path)
	dir := path.Dir(r.URL.Path)
	_, address := path.Split(dir)
	if !strings.HasPrefix(dir, "/rewards") || file != "download" || address == "" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		errorHandler(w, r, http.StatusMethodNotAllowed)
		return
	}

	report, err := lido.FetchRewardsReport(address)
	if err != nil {
		log.Printf("Error fetching report: %s", err)
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	reportCsv, err := convertRewardsReportToKoinlyCSV(report)
	if err != nil {
		log.Printf("error generating CSV: %s", err)
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=report.csv")
	fmt.Fprint(w, reportCsv)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, http.StatusText(status))
}

func newServeMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/rewards/", reportHandler)
	return mux
}

func main() {
	addr := "localhost:8080"
	s := &http.Server{
		Addr:         addr,
		Handler:      newServeMux(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Listening on %s", addr)
	log.Fatal(s.ListenAndServe())
}

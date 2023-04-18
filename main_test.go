package main

import (
	"encoding/csv"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestReportDownload(t *testing.T) {
	expectedContentType := "text/csv"
	expectedContentDisposition := "attachment; filename=report.csv"
	expectedCsvHeader := []string{"Koinly Date", "Amount", "Currency", "Label", "TxHash"}

	ts := httptest.NewServer(newServeMux())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/rewards/0x03d04a5F3cc050aB69A84eB0Da3242cd84DBf724/download")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != expectedContentType {
		t.Errorf("expected content type %s, got %s", expectedContentType, contentType)
	}

	contentDisposition := res.Header.Get("Content-Disposition")
	if contentDisposition != expectedContentDisposition {
		t.Errorf("expected content disposition %s, got %s", expectedContentDisposition, contentDisposition)
	}

	r := csv.NewReader(res.Body)
	header, err := r.Read()
	if err != nil {
		if err == io.EOF {
			t.Error("expected csv to have a header")
		} else {
			t.Fatal(err)
		}
	}

	if !reflect.DeepEqual(header, expectedCsvHeader) {
		t.Errorf("expected csv header %v, got %v", expectedCsvHeader, header)
	}
}

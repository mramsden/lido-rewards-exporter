package lido

import (
	"errors"
	"testing"
)

func TestLido(t *testing.T) {
	examples := []struct {
		description        string
		address            string
		expectedTotalItems uint
		expectedError      error
	}{
		{"address with rewards", "0x03d04a5F3cc050aB69A84eB0Da3242cd84DBf724", 570, nil},
		{"address without rewards", "0x03d04a5F3cc050aB69A84eB0Da3242cd84DBf725", 0, nil},
		{"invalid address", "0x03d04a5F3cc050aB69A84eB0Da3242cd84DBf72", 0, ErrInvalidAddress},
	}

	for _, e := range examples {
		r, err := FetchRewardsReport(e.address)
		if e.expectedError != nil {
			if !errors.Is(err, e.expectedError) {
				t.Errorf("Expected error %v, got %v for %s", e.expectedError, err, e.description)
			}
		} else if err != nil {
			t.Fatal(err)
		}
		if r.TotalItems != uint64(e.expectedTotalItems) {
			t.Errorf("Expected %v rewards, got %v for %s", e.expectedTotalItems, r.TotalItems, e.description)
		}
	}
}

package expenses

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type MonthlyTotalData struct {
	Total string
}

func (d *MonthlyTotalData) String() string {
	if d == nil || d.Total == "" {
		return ""
	}
	return "Monthly Expenses: " + d.Total + "€"
}

// GetMonthlyTotal fetches the monthly total from the URL specified in the MONTHLY_TOTAL_URL env var.
func GetMonthlyTotal() *MonthlyTotalData {
	total, err := getMonthlyTotal()
	if err != nil {
		log.Printf("Error fetching monthly total: %v", err)
		return nil
	}

	return &MonthlyTotalData{Total: total}
}

func getMonthlyTotal() (string, error) {
	url := os.Getenv("MONTHLY_TOTAL_URL")
	if url == "" {
		return "", errors.New("environment variable MONTHLY_TOTAL_URL is not set")
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch monthly total: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return strings.TrimSpace(string(bodyBytes)), nil
}

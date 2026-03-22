package expenses

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetMonthlyTotal_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("12345.67\n"))
	}))
	defer ts.Close()

	if err := os.Setenv("MONTHLY_TOTAL_URL", ts.URL); err != nil {
		t.Fatal(err)
	}

	total := GetMonthlyTotal()
	if total == nil {
		t.Fatal("expected data, got nil")
	}
	if total.Total != "12345.67" {
		t.Errorf("expected total 12345.67, got %q", total.Total)
	}
}

func TestGetMonthlyTotal_MissingEnv(t *testing.T) {
	_ = os.Unsetenv("MONTHLY_TOTAL_URL")

	if total := GetMonthlyTotal(); total != nil {
		t.Fatal("expected nil for missing env var")
	}
}

func TestGetMonthlyTotal_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	}))
	defer ts.Close()

	if err := os.Setenv("MONTHLY_TOTAL_URL", ts.URL); err != nil {
		t.Fatal(err)
	}

	if total := GetMonthlyTotal(); total != nil {
		t.Fatal("expected nil due to non-OK status")
	}
}

func TestGetMonthlyTotal_HTTPError(t *testing.T) {
	if err := os.Setenv("MONTHLY_TOTAL_URL", "http://127.0.0.1:0"); err != nil {
		t.Fatal(err)
	}

	if total := GetMonthlyTotal(); total != nil {
		t.Fatal("expected nil for unreachable server")
	}
}

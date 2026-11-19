package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestServer starts a small web server on this computer for the tests,
// so the tests do not need the internet.
func newTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
	})
	return httptest.NewServer(mux) // any other path gives 404
}

func TestCheckStatusCodes(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	client := &http.Client{Timeout: time.Second}

	tests := []struct {
		name string
		path string
		want int
	}{
		{"page that exists", "/ok", 200},
		{"page that does not exist", "/missing", 404},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := check(client, server.URL+tt.path)
			if r.err != nil {
				t.Fatalf("unexpected error: %v", r.err)
			}
			if r.code != tt.want {
				t.Errorf("got %d, want %d", r.code, tt.want)
			}
		})
	}
}

func TestCheckTimeout(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	client := &http.Client{Timeout: 50 * time.Millisecond} // shorter than /slow takes

	r := check(client, server.URL+"/slow")
	if r.err == nil {
		t.Fatal("expected a timeout error, got none")
	}
}

func TestCheckUnreachable(t *testing.T) {
	client := &http.Client{Timeout: time.Second}

	r := check(client, "http://127.0.0.1:1") // nothing listens on port 1
	if r.err == nil {
		t.Fatal("expected an error for an unreachable site, got none")
	}
}

func TestTallyIsSafe(t *testing.T) {
	var stats tally
	done := make(chan bool)

	for i := 0; i < 100; i++ {
		go func() {
			stats.add(i%2 == 0)
			done <- true
		}()
	}
	for i := 0; i < 100; i++ {
		<-done
	}

	if stats.up != 50 || stats.down != 50 {
		t.Errorf("got up=%d down=%d, want 50 and 50", stats.up, stats.down)
	}
}

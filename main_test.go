package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompressOnlyText(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte(strings.Repeat("a", 4096)), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "video.mp4"), []byte(strings.Repeat("a", 4096)), 0o644); err != nil {
		t.Fatal(err)
	}
	app := newApp(dir)

	for path, want := range map[string]string{"/app.js": "gzip", "/video.mp4": ""} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Accept-Encoding", "gzip")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if got := resp.Header.Get("Content-Encoding"); got != want {
			t.Fatalf("%s: Content-Encoding = %q, want %q", path, got, want)
		}
	}
}

func TestNewApp(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("home"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte("console.log(1)"), 0o644); err != nil {
		t.Fatal(err)
	}
	app := newApp(dir)

	tests := []struct {
		name   string
		path   string
		header map[string]string
		status int
		body   string
	}{
		{"root serves index", "/", nil, 200, "home"},
		{"asset served", "/app.js", nil, 200, "console.log(1)"},
		{"range request", "/app.js", map[string]string{"Range": "bytes=0-6"}, 206, "console"},
		{"spa route falls back", "/products/42", nil, 200, "home"},
		{"missing asset is 404", "/missing.png", nil, 404, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			for k, v := range tt.header {
				req.Header.Set(k, v)
			}
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != tt.status {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tt.status)
			}
			if tt.body != "" && string(body) != tt.body {
				t.Fatalf("body = %q, want %q", body, tt.body)
			}
		})
	}
}

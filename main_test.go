package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMux(t *testing.T) {
	r := newRouter()

	ts := httptest.NewTLSServer(r)
	defer ts.Close()

	client := ts.Client()

	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status should be OK, but got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	got := string(b)
	want := "Hello, World."

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestMuxForNonexistentRoute(t *testing.T) {
	r := newRouter()
	ts := httptest.NewTLSServer(r)
	defer ts.Close()
	client := ts.Client()

	resp, err := client.Post(ts.URL+"/", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code: got %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := ""

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	ts := httptest.NewTLSServer(r)
	defer ts.Close()
	client := ts.Client()

	resp, err := client.Get(ts.URL + "/static/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("response status code %d, wanted %d", resp.StatusCode, http.StatusOK)
	}

	got := resp.Header.Get("Content-Type")
	want := "text/html; charset=utf-8"

	if got != want {
		t.Errorf("content type: %s, wanted: %s", got, want)
	}
}

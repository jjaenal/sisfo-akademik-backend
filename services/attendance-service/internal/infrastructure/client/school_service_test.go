package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetLocation_Success(t *testing.T) {
	res := struct {
		Success bool `json:"success"`
		Data    struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"data"`
	}{Success: true}
	res.Data.Latitude = -6.2
	res.Data.Longitude = 106.8
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(res)
	}))
	defer ts.Close()
	svc := NewSchoolService(ts.URL, 2*time.Second)
	loc, err := svc.GetLocation(context.Background(), "tenant-x")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if loc == nil || loc.Latitude != -6.2 || loc.Longitude != 106.8 || loc.Radius != 100 {
		t.Fatalf("unexpected location")
	}
}

func TestGetLocation_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	svc := NewSchoolService(ts.URL, time.Second)
	_, err := svc.GetLocation(context.Background(), "t")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetLocation_DecodeError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("{invalid"))
	}))
	defer ts.Close()
	svc := NewSchoolService(ts.URL, time.Second)
	_, err := svc.GetLocation(context.Background(), "t")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetLocation_ApiFailure(t *testing.T) {
	res := struct {
		Success bool `json:"success"`
		Data    struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"data"`
	}{Success: false}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(res)
	}))
	defer ts.Close()
	svc := NewSchoolService(ts.URL, time.Second)
	_, err := svc.GetLocation(context.Background(), "t")
	if err == nil {
		t.Fatalf("expected error")
	}
}


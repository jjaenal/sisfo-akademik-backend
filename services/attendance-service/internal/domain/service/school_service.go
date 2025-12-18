package service

import (
	"context"
)

type SchoolLocation struct {
	Latitude  float64
	Longitude float64
	Radius    float64 // Allowed radius in meters
}

type SchoolService interface {
	GetLocation(ctx context.Context, tenantID string) (*SchoolLocation, error)
}

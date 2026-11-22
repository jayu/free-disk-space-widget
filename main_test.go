package main

import (
	"math"
	"testing"
)

func TestStatusIcon(t *testing.T) {
	for _, tt := range []struct {
		name    string
		color   bool
		percent float64
		want    string
	}{
		{"disabled", false, 0, diskIcon},
		{"above warning", true, 20.01, diskIcon},
		{"at warning", true, 20, "🟡"},
		{"above error", true, 10.01, "🟡"},
		{"at error", true, 10, "🔴"},
		{"empty", true, 0, "🔴"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts := options{color: tt.color, warn: 20, critical: 10}
			if got := opts.icon(tt.percent); got != tt.want {
				t.Errorf("icon(%g) = %q, want %q", tt.percent, got, tt.want)
			}
		})
	}
}

func TestThresholds(t *testing.T) {
	for _, tt := range []struct {
		name     string
		warn     float64
		critical float64
		valid    bool
	}{
		{"range endpoints", 100, 0, true},
		{"decimals", 19.5, 8.25, true},
		{"equal", 10, 10, false},
		{"reversed", 10, 20, false},
		{"above 100", 101, 10, false},
		{"below zero", 20, -1, false},
		{"NaN", math.NaN(), 10, false},
		{"infinite", 20, math.Inf(-1), false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opts := options{warn: tt.warn, critical: tt.critical}
			if got := opts.validThresholds(); got != tt.valid {
				t.Errorf("validThresholds() = %v, want %v", got, tt.valid)
			}
		})
	}
}

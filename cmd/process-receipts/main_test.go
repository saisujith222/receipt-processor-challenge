package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCalculateRetailerPoints(t *testing.T) {
	tests := []struct {
		name     string
		retailer string
		expected int
	}{
		{
			name:     "Valid retailer name with letters and numbers",
			retailer: "Target123",
			expected: 9,
		},
		{
			name:     "Retailer with only letters",
			retailer: "Wal mart",
			expected: 7,
		},
		{
			name:     "Retailer with spaces and special characters",
			retailer: "Best Buy!",
			expected: 7,
		},
		{
			name:     "Empty retailer name",
			retailer: "",
			expected: 0,
		},
		{
			name:     "Retailer with all non-alphanumeric characters",
			retailer: "!@#$%^&*()_+",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateRetailerPoints(tt.retailer)
			assert.Equal(t, tt.expected, result)
		})
	}
}


func TestCalculateItemPairPoints(t *testing.T) {
	tests := []struct {
		name     string
		items    []Item
		expected int
	}{
		{"No items", []Item{}, 0},
		{"One item", []Item{{"Item1", "1.00"}}, 0},
		{"Two items", []Item{{"Item1", "1.00"}, {"Item2", "2.00"}}, 5},
		{"Three items", []Item{{"Item1", "1.00"}, {"Item2", "2.00"}, {"Item3", "3.00"}}, 5},
		{"Four items", []Item{{"Item1", "1.00"}, {"Item2", "2.00"}, {"Item3", "3.00"}, {"Item4", "4.00"}}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := calculateItemPairPoints(tt.items); result != tt.expected {
				assert.Equal(t,tt.expected,result)
			}
		})
	}
}


func TestCalculateTimePoints(t *testing.T) {
	tests := []struct {
		name         string
		purchaseTime string
		expected     int
		expectError  bool
	}{
		{"Valid Time in Range", "15:30", 10, false},
		{"Valid Time Edge Case", "14:30", 10, false},
		{"Valid Time Not in Range", "13:59", 0, false},
		{"Invalid Hour Format", "25:00", 0, true},
		{"Invalid Minute Format", "14:61", 0, true},
		{"Empty Time", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculateTimePoints(tt.purchaseTime)
			if tt.expectError {
				assert.NotNil(t,err);
			} else {
				assert.Equal(t,tt.expected,result)
			}
		})
	}
}


func TestCalculateOddDayPoints(t *testing.T) {
	tests := []struct {
		name          string
		purchaseDate  string
		expectedPoints int
	}{
		{"Odd Day", "2022-01-01", 6},
		{"Even Day", "2022-01-02", 0},
		{"Invalid Date", "2022-01-XX", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points, _ := calculateOddDayPoints(tt.purchaseDate)
			assert.Equal(t,points,tt.expectedPoints)
		})
	}
}

func TestCalculateTotalPoints(t *testing.T) {
	tests := []struct {
		name       string
		total      string
		expectedPoints int
	}{
		{"Round Total", "2", 75},
		{"Multiple of 0.25", "50.25", 25},
		{"Invalid Total", "invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points, _ := calculateTotalPoints(tt.total)
			assert.Equal(t, points, tt.expectedPoints)
		})
	}
}


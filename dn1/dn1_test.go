package main

import (
	"math"
	"testing"
)

const floatTolerance = 1e-9

func almostEqual(first float64, second float64) bool {
	return math.Abs(first-second) <= floatTolerance
}

func TestStatistikaCalculations(t *testing.T) {
	postaje := map[string]Postaja{
		"LJ01": {
			ime:     "Ljubljana 01",
			kraj:    "Ljubljana",
			meritve: []float64{10.0, 20.0, 30.0},
		},
		"MB01": {
			ime:     "Maribor 01",
			kraj:    "Maribor",
			meritve: []float64{-10.5, 0.0, 15.5, 20.0},
		},
		"KP01": {
			ime:     "Koper 01",
			kraj:    "Koper",
			meritve: []float64{12.0, 12.0, 12.0},
		},
	}

	testCases := []struct {
		stationIdentifier string
		expectedMinimum   float64
		expectedMaximum   float64
		expectedMean      float64
		expectedSuccess   bool
	}{
		{
			stationIdentifier: "LJ01",
			expectedMinimum:   10.0,
			expectedMaximum:   30.0,
			expectedMean:      20.0,
			expectedSuccess:   true,
		},
		{
			stationIdentifier: "MB01",
			expectedMinimum:   -10.5,
			expectedMaximum:   20.0,
			expectedMean:      6.25,
			expectedSuccess:   true,
		},
		{
			stationIdentifier: "KP01",
			expectedMinimum:   12.0,
			expectedMaximum:   12.0,
			expectedMean:      12.0,
			expectedSuccess:   true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.stationIdentifier, func(t *testing.T) {
			minimum, maximum, mean, success := statistika(postaje, testCase.stationIdentifier)

			if success != testCase.expectedSuccess {
				t.Fatalf("expected success %v, got %v", testCase.expectedSuccess, success)
			}
			if !almostEqual(minimum, testCase.expectedMinimum) {
				t.Errorf("expected minimum %f, got %f", testCase.expectedMinimum, minimum)
			}
			if !almostEqual(maximum, testCase.expectedMaximum) {
				t.Errorf("expected maximum %f, got %f", testCase.expectedMaximum, maximum)
			}
			if !almostEqual(mean, testCase.expectedMean) {
				t.Errorf("expected mean %f, got %f", testCase.expectedMean, mean)
			}
		})
	}
}

func TestStatistikaInsufficientMeasurements(t *testing.T) {
	postaje := map[string]Postaja{
		"EMPTY": {
			ime:     "Prazna 01",
			kraj:    "Prazna",
			meritve: []float64{},
		},
		"ONE": {
			ime:     "Ena 01",
			kraj:    "Ena",
			meritve: []float64{15.0},
		},
		"TWO": {
			ime:     "Dve 01",
			kraj:    "Dve",
			meritve: []float64{15.0, 16.0},
		},
	}

	identifiers := []string{"EMPTY", "ONE", "TWO", "NON_EXISTENT"}

	for _, identifier := range identifiers {
		t.Run(identifier, func(t *testing.T) {
			_, _, _, success := statistika(postaje, identifier)
			if success {
				t.Errorf("expected statistika to fail for identifier %s", identifier)
			}
		})
	}
}

func TestOcenaTemperature(t *testing.T) {
	testCases := []struct {
		temperatureValue float64
		expectedRating   string
	}{
		{35.0, "vroče"},
		{25.0, "vroče"},
		{24.9, "toplo"},
		{15.0, "toplo"},
		{14.9, "hladno"},
		{5.0, "hladno"},
		{4.9, "mrzlo"},
		{0.0, "mrzlo"},
		{-15.0, "mrzlo"},
	}

	for _, testCase := range testCases {
		actualRating := ocena_temperature(testCase.temperatureValue)
		if actualRating != testCase.expectedRating {
			t.Errorf("for temperature %f expected %s, got %s", testCase.temperatureValue, testCase.expectedRating, actualRating)
		}
	}
}

func TestDodajMeritev(t *testing.T) {
	postaje := map[string]Postaja{
		"LJ01": {
			ime:     "Ljubljana 01",
			kraj:    "Ljubljana",
			meritve: []float64{10.0},
		},
	}

	t.Run("ValidWithinRange", func(t *testing.T) {
		dodaj_meritev(postaje, "LJ01", 20.0)
		station := postaje["LJ01"]
		if len(station.meritve) != 2 || station.meritve[1] != 20.0 {
			t.Errorf("measurement 20.0 was not correctly added: %v", station.meritve)
		}
	})

	t.Run("BoundaryValuesAccepted", func(t *testing.T) {
		dodaj_meritev(postaje, "LJ01", TEMP_MIN)
		dodaj_meritev(postaje, "LJ01", TEMP_MAX)
		station := postaje["LJ01"]
		if len(station.meritve) != 4 {
			t.Fatalf("expected 4 measurements, got %d", len(station.meritve))
		}
		if station.meritve[2] != TEMP_MIN || station.meritve[3] != TEMP_MAX {
			t.Errorf("boundary values were not properly added: %v", station.meritve)
		}
	})

	t.Run("OutOfRangeRejected", func(t *testing.T) {
		initialLength := len(postaje["LJ01"].meritve)
		dodaj_meritev(postaje, "LJ01", TEMP_MIN-0.1)
		dodaj_meritev(postaje, "LJ01", TEMP_MAX+0.1)
		if len(postaje["LJ01"].meritve) != initialLength {
			t.Errorf("measurements outside temperature limits were incorrectly added")
		}
	})

	t.Run("MissingStationIgnored", func(t *testing.T) {
		dodaj_meritev(postaje, "NON_EXISTENT", 10.0)
		if _, exists := postaje["NON_EXISTENT"]; exists {
			t.Errorf("non-existent station was unexpectedly created")
		}
	})
}

func TestFormattedAppendix(t *testing.T) {
	testCases := []struct {
		inputNumber    int
		expectedOutput string
	}{
		{1, "01"},
		{9, "09"},
		{10, "10"},
		{99, "99"},
	}

	for _, testCase := range testCases {
		output := formatted_appendix(testCase.inputNumber)
		if output != testCase.expectedOutput {
			t.Errorf("expected formatted_appendix(%d) = %s, got %s", testCase.inputNumber, testCase.expectedOutput, output)
		}
	}
}

package gvalidator_test

import (
	"testing"

	globalvalidator "github.com/Phi-S/cs-server-manager/gvalidator"
)

func TestGlobalValidator_PortTag_Ok(t *testing.T) {
	type testInput struct {
		port        uint32
		shouldError bool
	}
	testData := []testInput{
		{0, true},
		{1, false},
		{65535, false},
		{65536, true},
		{65538, true},
	}

	for _, td := range testData {
		err := globalvalidator.Instance().Var(td.port, "port")
		if err == nil {
			if td.shouldError {
				t.Fatalf("error expected but nill returned. port: %v", td.port)
			}
		} else {
			if !td.shouldError {
				t.Fatalf("Test failed for data: %v Error: %v", td.port, err)
			}
		}
	}
}

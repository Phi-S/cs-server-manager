package globalvalidator_test

import (
	globalvalidator "cs-server-manager/global_validator"
	"testing"
)

func init() {
	globalvalidator.Init()
}

func TestPortTag_Ok(t *testing.T) {

	type testInput struct {
		port             uint32
		shouldThrowError bool
	}
	testData := []testInput{
		{0, true},
		{1, false},
		{65535, false},
		{65536, true},
		{65538, true},
	}

	for _, td := range testData {
		err := globalvalidator.Instance.Var(td.port, "port")
		if err == nil {
			if td.shouldThrowError {
				t.Fatalf("error expected but nill returned. port: %v", td.port)
			}
		} else {
			if !td.shouldThrowError {
				t.Fatalf("Test failed for data: %v Error: %v", td.port, err)
			}
		}
	}
}

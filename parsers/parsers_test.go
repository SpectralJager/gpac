package parsers

import "testing"

func TestInteger(t *testing.T) {
	testCases := []struct {
		input  string
		result string
	}{
		{
			input:  "0",
			result: "0",
		},
		{
			input:  "12",
			result: "12",
		},
		{
			input:  "1344",
			result: "1344",
		},
		{
			input:  "012301984",
			result: "012301984",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.result, func(t *testing.T) {
			result := Integer()([]byte(tC.input))
			if result.Error != nil {
				t.Fatal(result.Error)
			}
			if result.Ok != tC.result {
				t.Fail()
			}
		})
	}
}

func TestSignedInteger(t *testing.T) {
	testCases := []struct {
		input  string
		result string
	}{
		{
			input:  "0",
			result: "0",
		},
		{
			input:  "12",
			result: "12",
		},
		{
			input:  "1344",
			result: "1344",
		},
		{
			input:  "012301984",
			result: "012301984",
		},
		{
			input:  "-0",
			result: "-0",
		},
		{
			input:  "-12",
			result: "-12",
		},
		{
			input:  "-1344",
			result: "-1344",
		},
		{
			input:  "-012301984",
			result: "-012301984",
		},
		{
			input:  "+0",
			result: "0",
		},
		{
			input:  "+12",
			result: "12",
		},
		{
			input:  "+1344",
			result: "1344",
		},
		{
			input:  "+012301984",
			result: "012301984",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.result, func(t *testing.T) {
			result := SignedInteger()([]byte(tC.input))
			if result.Error != nil {
				t.Fatal(result.Error)
			}
			if result.Ok != tC.result {
				t.Fail()
			}
		})
	}
}

func TestMatch(t *testing.T) {
	testCases := []struct {
		input  string
		result string
	}{
		{
			input:  "anyn,mn,am",
			result: "any",
		},
		{
			input:  "forafsln",
			result: "for",
		},
		{
			input:  "struct_lkasfn",
			result: "struct",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.result, func(t *testing.T) {
			result := Match(tC.result)([]byte(tC.input))
			if result.Error != nil {
				t.Fatal(result.Error)
			}
			if result.Ok != tC.result {
				t.Fail()
			}
		})
	}
}

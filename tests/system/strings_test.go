package system_test

import (
	"service-recordingStorage/system"
	"testing"
)

type removeSpecialCharsData struct {
	OriginalStr       string
	SpecialCharacters string
	ExpectedStr       string
}

type truncateByMaxLenData struct {
	OriginalStr string
	MaxLen      int
	ExpectedStr string
}

func TestRemoveSpecialCharsOK(t *testing.T) {

	testData := []removeSpecialCharsData{
		{
			OriginalStr:       `S/ome\ t\ext <wi>th spe::cial" *|*characters?`,
			SpecialCharacters: `/\<>:"|?*`,
			ExpectedStr:       "Some text with special characters",
		},
		{
			OriginalStr:       "Some text with! special (characters)",
			SpecialCharacters: `/\<>:"|?*`,
			ExpectedStr:       "Some text with! special (characters)",
		},
		{
			OriginalStr:       "Some text without special characters",
			SpecialCharacters: `/\<>:"|?*`,
			ExpectedStr:       "Some text without special characters",
		},
		{
			OriginalStr:       "Some text without special characters",
			SpecialCharacters: ``,
			ExpectedStr:       "Some text without special characters",
		},
	}

	for _, testRow := range testData {
		result := system.RemoveSpecialChars(testRow.OriginalStr, testRow.SpecialCharacters)
		if result != testRow.ExpectedStr {
			t.Fatalf("Result string %s is not equal to expected string %s. Original string is %s, special characters are %s", result, testRow.ExpectedStr, testRow.OriginalStr, testRow.SpecialCharacters)
		}
	}
}

func TestTruncateByMaxLenOK(t *testing.T) {

	testData := []truncateByMaxLenData{
		{
			OriginalStr: `Some long long long text`,
			MaxLen:      10,
			ExpectedStr: "Some long ",
		},
		{
			OriginalStr: "Some short text",
			MaxLen:      20,
			ExpectedStr: "Some short text",
		},
	}

	for _, testRow := range testData {
		result := system.TruncateByMaxLen(testRow.OriginalStr, testRow.MaxLen)
		if result != testRow.ExpectedStr {
			t.Fatalf("Result string %s is not equal to expected string %s. Original string is %s, maxLen is %d", result, testRow.ExpectedStr, testRow.OriginalStr, testRow.MaxLen)
		}
	}
}

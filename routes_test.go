package main

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

// TestValidateClean tests the validate function, checking
// that is returns no errors for a known valid receipts
func TestValidateClean(t *testing.T) {
	// test valid receipt
	filetitles := []string{"ex1", "ex2", "ex3", "ex4"}
	receipts := getTestReceipts(filetitles)
	for i, receipt := range receipts {
		validation := validate(receipt)
		if len(validation) > 0 {
			t.Errorf("Validation of %v now has %v number of Error; want 0", filetitles[i], len(validation))
			t.Logf("Error Given: %v", validation)
		}
	}
}

// TestValidateError tests the validate function, checking
// that is returns any errors for a known invalid receipts
func TestValidateError(t *testing.T) {
	// test invalid receipt
	filetitles := []string{"exbreak1", "exbreak2", "exbreak3"}
	receipts := getTestReceipts(filetitles)
	for i, receipt := range receipts {
		validation := validate(receipt)
		t.Logf("Number of Errors in %v: %v", filetitles[i], len(validation))
		if len(validation) == 0 {
			t.Errorf("validate(%v) now has %v number of Error; want more than 0", filetitles[i], len(validation))
		}
	}
}

// parse local json file from ExampleReceipts folder for use in tests
func getTestReceipt(file_title string) Receipt {
	file_title = "ExampleReceipts/" + file_title + ".json"
	data, err := os.Open(file_title)
	if err != nil {
		panic(err)
	}
	defer data.Close()
	byteValue, _ := io.ReadAll(data)
	var receipt Receipt
	json.Unmarshal(byteValue, &receipt)
	return receipt
}

func getTestReceipts(file_titles []string) []Receipt {
	var receipts []Receipt
	for _, file_title := range file_titles {
		receipts = append(receipts, getTestReceipt(file_title))
	}
	return receipts
}

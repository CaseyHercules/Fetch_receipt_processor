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
	receipts := getTestReceipts(filetitles, t)
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
	filetitles := []string{"exbreak1", "exbreak2", "exbreak3", "exbreak_testbadTotal"}
	receipts := getTestReceipts(filetitles, t)
	for i, receipt := range receipts {
		validation := validate(receipt)
		t.Logf("Number of Errors in %v: %v", filetitles[i], len(validation))
		if len(validation) == 0 {
			t.Errorf("validate(%v) now has %v number of Error; want more than 0", filetitles[i], len(validation))
		}
	}
}

// Rule 1: One point for every alphanumeric character in the retailer name.
func TestValidateRule1(t *testing.T) {
	filetitles := []string{"rule1_test1", "rule1_test2", "rule1_test3"}
	expected := []int{11, 12, 9}
	receipts := getTestReceipts(filetitles, t)
	for i, receipt := range receipts {
		result := calculatePointsRule1(receipt).Points
		if result != expected[i] {
			t.Errorf("calculatePointsRule1(%v) = %v; want %v", filetitles[i], result, expected[i])
		}
	}
}

// Rule 2: 50 points if the total is a round dollar amount with no cents.
func TestValidateRule2(t *testing.T) {
	filetitles := []string{"rule2_test1", "rule2_test2", "ex1"}
	expected := []int{50, 0, 0}
	receipts := getTestReceipts(filetitles, t)
	for i, receipt := range receipts {
		result := calculatePointsRule2(receipt).Points
		if result != expected[i] {
			t.Errorf("calculatePointsRule2(%v) = %v; want %v", filetitles[i], result, expected[i])
		}
	}
}

// Rule 3: 25 points if the total is a multiple of 0.25.
func TestValidateRule3(t *testing.T) {
	filetitles := []string{"rule2_test1", "rule2_test2", "ex1"}
	expected := []int{25, 25, 0}
	receipts := getTestReceipts(filetitles, t)
	for i, receipt := range receipts {
		result := calculatePointsRule3(receipt).Points
		if result != expected[i] {
			t.Errorf("calculatePointsRule3(%v) = %v; want %v", filetitles[i], result, expected[i])
		}
	}
}

// Rule 4: 5 points for every two items on the receipt.
func TestValidateRule4(t *testing.T) {
	filetitles := []string{"ex1", "ex2", "ex3"}
	expected := []int{5, 10, 10}
	receipts := getTestReceipts(filetitles, t)
	for i, receipt := range receipts {
		result := calculatePointsRule4(receipt).Points
		if result != expected[i] {
			t.Errorf("calculatePointsRule4(%v) = %v; want %v", filetitles[i], result, expected[i])
		}
	}
}

// Rule 5: If the trimmed length of the item description is a multiple of 3,
//
//	multiply the price by 0.2 and round up to the nearest integer.
//	The result is the number of points earned.
func TestValidateRule5(t *testing.T) {
	filetitles := []string{"ex1", "ex2", "ex3", "ex4"}
	expected_pts := []int{1, 0, 6, 26}
	expected_count := []int{1, 0, 3, 2}
	receipts := getTestReceipts(filetitles, t)
	for i, receipt := range receipts {
		pts := 0
		count := len(calculatePointsRule5(receipt))
		for _, item := range calculatePointsRule5(receipt) {
			pts += item.Points
		}
		if pts != expected_pts[i] {
			t.Errorf("calculatePointsRule5(%v) Count of points = %v; want %v", filetitles[i], pts, expected_pts[i])
		}
		if count != expected_count[i] {
			t.Errorf("calculatePointsRule5(%v) Count of items = %v; want %v", filetitles[i], count, expected_count[i])
		}
	}
}

// Rule 6: 6 points if the day in the purchase date is odd.
func TestValidateRule6(t *testing.T) {
	filetitles := []string{"ex1", "ex2", "ex3", "ex4"}
	expected := []int{0, 0, 6, 0}
	receipts := getTestReceipts(filetitles, t)
	for i, receipt := range receipts {
		result := calculatePointsRule6(receipt).Points
		if result != expected[i] {
			t.Errorf("calculatePointsRule6(%v) = %v; want %v", filetitles[i], result, expected[i])
		}
	}
}

// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func TestValidateRule7(t *testing.T) {
	filetitles := []string{"ex1", "ex2", "ex3", "ex4"}
	expected := []int{0, 0, 0, 10}
	receipts := getTestReceipts(filetitles, t)
	for i, receipt := range receipts {
		result := calculatePointsRule7(receipt).Points
		if result != expected[i] {
			t.Errorf("calculatePointsRule7(%v) = %v; want %v", filetitles[i], result, expected[i])
		}
	}
}

// parse local json file from ExampleReceipts folder for use in tests
func getTestReceipt(file_title string, t *testing.T) Receipt {
	file_title = "ExampleReceipts/" + file_title + ".json"
	data, err := os.Open(file_title)
	if err != nil {
		t.Fatal(err)
	}
	defer data.Close()
	byteValue, _ := io.ReadAll(data)
	var receipt Receipt
	json.Unmarshal(byteValue, &receipt)
	return receipt
}

func getTestReceipts(file_titles []string, t *testing.T) []Receipt {
	var receipts []Receipt
	for _, file_title := range file_titles {
		receipts = append(receipts, getTestReceipt(file_title, t))
	}
	return receipts
}

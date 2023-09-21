package main

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total,string"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price,string"`
}

type Points struct {
	Points int `json:"points"`
}

type Breakdown struct {
	Breakdown int `json:"breakdown"`
}

var rx = regexp.MustCompile("[^\\w]+")

func setupRoutes(app *fiber.App) {

	var receipts = make(map[string]Receipt)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the Fetch_receipt_processor_challenge_API! \n\n" +
			"Please use the following endpoints to interact with the API: \n\n" +
			"POST /receipts/process \n" +
			"GET /receipts/:id/points specific Receipt's point total via the rule set\n" +
			"GET /receipts/:id/breakdown For a breakdown of a specific Receipt via returned Receipt UUID\n" +
			"GET /debug for a rough breakdown of recent processed receipts \n\n" +
			"Please see the README.md file for more information.")
	})
	// Endpoint for processing receipts
	app.Post("/receipts/process", func(c *fiber.Ctx) error {
		// Parse the JSON payload from the request
		receipt := new(Receipt)
		if err := c.BodyParser(receipt); err != nil {
			return err
		}
		// Validate the receipt
		validation := validate(*receipt)
		if len(validation) > 0 {
			return c.JSON(fiber.Map{"error": validation})
		}

		// Generate a unique ID for the receipt
		id := uuid.New().String()

		// Store the receipt in memory using the ID as the key
		receipts[id] = *receipt

		// Return the ID as a JSON response
		return c.JSON(fiber.Map{"id": id})
	})

	// Endpoint for retrieving points for a receipt
	app.Get("/receipts/:id/points", func(c *fiber.Ctx) error {
		// Get the ID from the URL parameter
		id := c.Params("id")

		// Look up the receipt by ID and calculate the points
		receipt, ok := receipts[id]
		if !ok {
			return fiber.NewError(fiber.StatusNotFound, "Receipt not found")
		}
		points := calculatePoints(receipt)

		// Return the points as a JSON response
		return c.JSON(fiber.Map{"points": points})
	})
	app.Get("/receipts/:id/breakdown", func(c *fiber.Ctx) error {
		// Get the ID from the URL parameter
		id := c.Params("id")

		// Look up the receipt by ID and calculate the points
		receipt, ok := receipts[id]
		if !ok {
			return fiber.NewError(fiber.StatusNotFound, "Receipt not found")
		}
		breakdown := calculateBreakdown(receipt, true)

		// Return the points as a JSON response
		return c.JSON(fiber.Map{"breakdown": breakdown})
	})
	app.Get("/debug", func(c *fiber.Ctx) error {
		// update the landing page to display all receipts ids and total points
		// for each receipt and their breakdown

		if len(receipts) < 1 {
			return c.SendString("No receipts yet")
		}
		var receipt_ids []string
		for id := range receipts {
			receipt_ids = append(receipt_ids, id, strings.Join(calculateBreakdown(receipts[id], true), "\n"), "\n")
		}
		return c.SendString(strings.Join(receipt_ids, "\n"))

	})
}

// calculatePoints calculates the number of points awarded for a given receipt
func calculatePoints(receipt Receipt) int {
	// Rule 1: One point for every alphanumeric character in the retailer name.
	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	// Rule 3: 25 points if the total is a multiple of 0.25.
	// Rule 4: 5 points for every two items on the receipt.
	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// 		multiply the price by 0.2 and round up to the nearest integer.
	//		The result is the number of points earned.
	// Rule 6: 6 points if the day in the purchase date is odd.
	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	var pts int = 0
	pts += len(rx.ReplaceAllString(receipt.Retailer, ""))
	if receipt.Total == float64(int(receipt.Total)) {
		pts += 50
	}
	if receipt.Total == float64(int(receipt.Total*4))/4 {
		pts += 25
	}
	pts += (len(receipt.Items) / 2) * 5
	for _, item := range receipt.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			//if value is negative, set to 0
			temp := int(math.Ceil(item.Price * 0.2))
			if temp < 0 {
				temp = 0
			}
			pts += temp
		}
	}
	if receipt.PurchaseDate[len(receipt.PurchaseDate)-1]%2 == 1 {
		pts += 6
	}
	if receipt.PurchaseTime > "14:00" && receipt.PurchaseTime < "16:00" {
		pts += 10
	}
	return pts
}
func calculateBreakdown(receipt Receipt, view_breakdown bool) []string {
	// Function to define the breakdown of points
	// Rule 1: One point for every alphanumeric character in the retailer name.
	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	// Rule 3: 25 points if the total is a multiple of 0.25.
	// Rule 4: 5 points for every two items on the receipt.
	// Rule 5: If the trimmed length of the item description is a multiple of 3,
	// 		multiply the price by 0.2 and round up to the nearest integer.
	//		The result is the number of points earned.
	// Rule 6: 6 points if the day in the purchase date is odd.
	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	var breakdown []string
	rule1_pts := len(rx.ReplaceAllString(receipt.Retailer, ""))
	breakdown = append(breakdown, strconv.Itoa(rule1_pts)+" points - retailer name has "+strconv.Itoa(len(rx.ReplaceAllString(receipt.Retailer, "")))+" alphanumeric characters")
	rule2_pts := 0
	if receipt.Total == float64(int(receipt.Total)) {
		rule2_pts = 50
	}
	breakdown = append(breakdown, strconv.Itoa(rule2_pts)+" points - total is a round dollar amount with no cents")
	rule3_pts := 0
	if receipt.Total == float64(int(receipt.Total*4))/4 {
		rule3_pts = 25
	}
	breakdown = append(breakdown, strconv.Itoa(rule3_pts)+" points - total is a multiple of 0.25")
	rule4_pts := (len(receipt.Items) / 2) * 5
	breakdown = append(breakdown, strconv.Itoa(rule4_pts)+" points - "+strconv.Itoa(len(receipt.Items))+" items ("+strconv.Itoa(len(receipt.Items)/2)+" pairs @ 5 points per pair)")
	for _, item := range receipt.Items {
		rule5_pts_temp := 0
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			rule5_pts_temp = int(math.Ceil(item.Price * 0.2))
			if rule5_pts_temp < 0 {
				rule5_pts_temp = 0
			}
			breakdown = append(breakdown, strconv.Itoa(rule5_pts_temp)+" points - \""+strings.TrimSpace(item.ShortDescription)+"\" has "+strconv.Itoa(len(strings.TrimSpace(item.ShortDescription)))+" characters and is a multiple of 3")
		}
	}
	rule6_pts := 0
	if receipt.PurchaseDate[len(receipt.PurchaseDate)-1]%2 == 1 {
		rule6_pts = 6
	}
	breakdown = append(breakdown, strconv.Itoa(rule6_pts)+" points - day in the purchase date is odd")
	rule7_pts := 0
	if receipt.PurchaseTime > "14:00" && receipt.PurchaseTime < "16:00" {
		rule7_pts = 10
	}
	breakdown = append(breakdown, strconv.Itoa(rule7_pts)+" points - time of purchase is after 2:00pm and before 4:00pm")
	breakdown = append(breakdown, strconv.Itoa(calculatePoints(receipt))+" points - total points")
	return breakdown
}
func validate(receipt Receipt) []string {
	var output []string
	// Validate the receipt

	// Validate all items have a price and short description AND
	// Validate total matches sum of items
	var temp_total float64 = 0
	for _, item := range receipt.Items {
		if item.Price == 0 {
			output = append(output, "Item is missing or has an invalid price")
		}
		if len(item.ShortDescription) == 0 {
			output = append(output, "Item is missing or has an invalid short description")
		}
		temp_total += item.Price
	}
	if temp_total != receipt.Total {
		output = append(output, "Total does not match sum of items. Total from items: "+
			strconv.FormatFloat(temp_total, 'f', 2, 64)+
			"   Total from receipt: "+
			strconv.FormatFloat(receipt.Total, 'f', 2, 64))
	}

	// Validate PurchaseDate is a valid date
	_, date_err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date_err != nil {
		output = append(output, "Date is a invalid date. Please use the format YYYY-MM-DD")
	}
	//validate PurchaseTime is a valid time
	_, time_err := time.Parse("15:04", receipt.PurchaseTime)
	if time_err != nil {
		output = append(output, "Time is a invalid time Please use the format HH:MM")
	}
	//validate Retailer is not empty after removing non-alphanumeric characters
	if len(rx.ReplaceAllString(receipt.Retailer, "")) < 1 {
		output = append(output, "Retailer is empty")
	}
	return output
}

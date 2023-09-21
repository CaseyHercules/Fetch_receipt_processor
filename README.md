# Fetch_receipt_processor
The Fetch_receipt_processor is a Go-based application designed to process and score receipts based on a set of rules. This README provides information on how to use the application, its endpoints, project description, rules, and considerations.

## Usage
1. Clone the repo
2. Navigate to the repo locally
3. Build the container via Docker Compose   `docker-compose up -d`
4. <a href="http://localhost:3000/">Head to Localhost:3000</a>
5. Send a POST request to port 3000 at endpoint `http://localhost:3000/receipts/process` to get started.
6. In the background a score is generated based on the rules listed below, with a few considerations made due to unknown answers to specific edge cases.
7. Use the UUID you received as a response and navigate to `http://localhost:3000/receipts/:id/points`, replacing :id with your UUID.
8. Received a JSON response of the points that received.
## Endpoints
Fetch_receipt_processor exposes the following endpoints:

- POST `/receipts/process`
    - Accepts formatted receipts. On success, it returns a UUID. Invalid receipts are not accepted, and responses include reasons for receipt rejection.
- GET  `/receipts/:id/points`
    - Retrieves the total points for a specific receipt, calculated using the rule set. Accessed via Receipt UUID.
- GET  `/receipts/:id/breakdown`
    - Provides a breakdown of point scoring for a specific receipt. Accessed via Receipt UUID.
- GET  `/debug`
    - Offers a human-readable breakdown of recently processed receipts.
## Project Description
This project serves as a take-home exercise for the **Backend Engineer** role at Fetch. The primary functionality of this program is to accept a JSON-formatted receipt via the POST endpoint /receipts/process, score that receipt and return a UUID.

This is my first time using Go as a programming language, it was quite nice and quite verbose, coming from Ruby that is. If you have any questions or issue, please contact me or create an Issue on this project.

### Rules
The scoring of receipts in Fetch_receipt_processor follows a set of rules:

- One point for every alphanumeric character in the retailer name.
- 50 points if the total is a round dollar amount with no cents.
- 25 points if the total is a multiple of 0.25.
- 5 points for every two items on the receipt.
- If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
- 6 points if the day in the purchase date is odd.
- 10 points if the time of purchase is after 2:00pm and before 4:00pm.

### Consideration
While adhering to the defined rules, This Fetch_receipt_processor has made one consideration outside of these rules, specifically regarding negative item values and negative item totals. Due to there being legitimate cases where items have negative values, resulting in a negative total, I made the deliberate choice not to assign negative points. This decision was rooted in the belief that negative points might not align with the intended goals of the scoring system.

Allowing negative points could potentially create opportunities for individuals to manipulate the system by strategically purchasing and returning items with character lengths divisible by 3. While this type of manipulation can be addressed through other means, for simplicity and consistency, I opted to set individual items that would otherwise contribute negative points to a minimum of 0 points.

This approach maintains the integrity of the scoring system while addressing potential vulnerabilities without overly complicating the logic.
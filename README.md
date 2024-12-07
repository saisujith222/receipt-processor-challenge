# Process Receipts API

This project provides an API to process retail receipts, calculate reward points, and retrieve those points using a unique receipt ID.

## Features

- Process receipts to calculate reward points based on retailer details, purchase date/time, items, and total price.
- Retrieve reward points for a given receipt ID.

## Installation

- Clone the Repository
   ```bash
   git clone https://github.com/yourusername/process-receipts.git
   cd cmd/process-receipts
	```
- Install Dependencies: Initialize the Go module and install the required dependencies:
  ```bash
  go mod tidy
  ```

- Run the Application
  ```bash
  go run main.go
  ```

## Testing Endpoints
- Here are the sample requests for testing the endpoints
  ### Process Receipts
```bash
curl -X POST http://localhost:8080/receipts/process \
     -H "Content-Type: application/json" \
     -d '{
         "retailer": "Target",
         "purchaseDate": "2022-01-01",
         "purchaseTime": "13:01",
         "items": [
             {
                 "shortDescription": "Mountain Dew 12PK",
                 "price": "6.49"
             },
             {
                 "shortDescription": "Emils Cheese Pizza",
                 "price": "12.25"
             },
             {
                 "shortDescription": "Knorr Creamy Chicken",
                 "price": "1.26"
             },
             {
                 "shortDescription": "Doritos Nacho Cheese",
                 "price": "3.35"
             },
             {
                 "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
                 "price": "12.00"
             }
         ],
         "total": "35.35"
     }'
```

### Get Points
```bash
curl -X GET http://localhost:8080/receipts/<receipt-id>/points
```

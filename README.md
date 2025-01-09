Receipt Processor Challenge
===========================

This is my solution for the [Fetch Rewards - Receipt Processor
Challenge](https://github.com/fetch-rewards/receipt-processor-challenge),
implemented with Go 1.23.3.

After cloning this repo, be sure to get the lone project dependency,
[uuid](https://github.com/google/uuid).

```
go get github.com/google/uuid
```

## Testing

If desired, run the test suite which includes a number of tests for the `Item`
and `Receipt` structs. These tests check the validity of the contents of the
fields in each struct and calculate expected points for each rule defined in
[the
rules](https://github.com/fetch-rewards/receipt-processor-challenge?tab=readme-ov-file#rules).

> Also points for an additional rule if an item's description starts with a G/g

```
go test
```

> Optionally pass the `-v` flag to show the results of each individual test.

There are also tests to get the total number of points for the three example
receipts provided in the challenge repo:

- `example1.json` is the
  [morning-receipt.json](https://github.com/fetch-rewards/receipt-processor-challenge/blob/main/examples/morning-receipt.json)
- `example2.json` and `example3.json` are the example JSON objects from the
  [challenge
  README](https://github.com/fetch-rewards/receipt-processor-challenge/blob/main/README.md#examples)

## Run the API server

The API server will run on localhost:8080

```
go run main.go
```

Endpoints:

- POST `/receipts/process` with receipt JSON as the payload (see `example*.json`
  files)
    - 200 response: JSON with 'id' field for the stored receipt
    - 400 response: JSON with 'error' field containing any validation errors
      with the payload
- GET `/receipts/{id}/points`
    - 200 response: JSON with 'points' field containing integer number of points
      awarded
    - 404 response: JSON with 'error' field if receipt not found
- GET `/receipts/{id}/breakdown`
    - 200 response: JSON with 'breakdown' field containing array of the points
      breakdown
    - 404 response: JSON with 'error' field if receipt not found

## (Optional) Using the python-webclient

The Python [webclient-helper](https://pypi.org/project/webclient-helper) package
is a wrapper I maintain to the popular
[requests](https://pypi.org/project/requests) package for making HTTP requests
to resources online.

The python-webclient directory contains a `client.py` file that defines a
`ReceiptClient` class to easily interact with the receipt-processor API written
in Go. The `ReceiptClient` has a number of methods:

- the `post_and_save` method will POST a given Python dict as JSON to the
  `/receipts/process` endpoint and save the UUID returned to an internal list of
  UUIDs if the endpoint gave a 200 response
    - the `post_receipt1` method uses the
      [morning-receipt.json](https://github.com/fetch-rewards/receipt-processor-challenge/blob/main/examples/morning-receipt.json)
      from the challenge examples directory
    - the `post_receipt2` and `post_receipt3` methods use the two example JSON
      objects from the [challenge
      README](https://github.com/fetch-rewards/receipt-processor-challenge/blob/main/README.md#examples)
- the `get_points` method will GET the response from the points endpoint for a
  given UUID (`/receipts/{_id}/points`)
- the `get_breakdown` method will GET the response from the breakdown endpoint
  for a given UUID (`/receipts/{_id}/breakdown`)
- the `get_all_points_and_breakdowns` method will iterate over the internal list
  of UUIDs returned from successful POST calls and hit the points and breakdown
  endpoints for each.

### Setup

Create a virtual environment and install the requirements.txt to it.

```
python3 -m venv venv

venv/bin/pip install -r requirements.txt
```

> Note: The only requirement is webclient-helper

### Using the client interactively

```
>>> from client import client

>>> from pprint import pprint

>>> response = client.post_receipt1()

>>> response.json()
{'id': '5710602b-a949-4a1f-b41e-68bf82c2fb7d'}

>>> response2 = client.get_points(response.json()['id'])

>>> response2.json()
{'points': 15}

>>> response3 = client.get_breakdown(response.json()['id'])

>>> pprint(response3.json())
{'breakdown': ['9 points for retailer name (Walgreens)',
               '0 points for round dollar amount (2.65)',
               '0 points for being multiple of 0.25 (2.65)',
               '5 points for number of items (2)',
               '0 point(s) for item (Pepsi - 12-oz | 1.25)',
               '1 point(s) for item (Dasani | 1.40)',
               '0 points for purchase day being odd (2022-01-02)',
               '0 points for time of purchase between 2pm and 4pm (08:13)',
               '15 points total']}

>>> client.post_and_save(json={}).json()
{'error': 'Validation errors: retailer cannot be empty | purchaseDate cannot be empty | purchaseTime cannot be empty | total cannot be empty | sum of item prices (0.00) != given total ()'}

>>> exit
```

### Running the `client.py` file as a script

This calls a number of methods on an instance of the ReceiptClient class, named
`client`. It makes some bad requests, POSTS the 3 good examples, POSTS a couple
bad examples, gets the points and breakdown from a random good UUID, tries to
get points and breakdown from a bad UUID, then gets the points and breakdowns
for all UUIDs.

```
% venv/bin/python client.py
(404) GET http://localhost:8080/receipts/None/points
{'error': 'receipt None not found'}

(405) GET http://localhost:8080/receipts/process
{'error': 'Only POST is allowed'}

(404) GET http://localhost:8080/receipts/bad
b'404 page not found\n'

(200) POST http://localhost:8080/receipts/process
{'id': 'e64f7946-6791-4536-ad94-c0ccda111f48'}

(200) POST http://localhost:8080/receipts/process
{'id': '2c0dc78d-0c7e-4a91-97cf-296e1c4911dc'}

(200) POST http://localhost:8080/receipts/process
{'id': '72f3b285-e192-49d0-940a-ed819a575166'}

(400) POST http://localhost:8080/receipts/process
{'error': 'Validation errors: purchaseDate cannot be parsed (2022-13-20) | '
          'purchaseTime cannot be parsed (24:33)'}

(400) POST http://localhost:8080/receipts/process
{'error': 'Validation errors: sum of item prices (4.50) != given total (4.60)'}

Selected _id='e64f7946-6791-4536-ad94-c0ccda111f48'
(200) GET http://localhost:8080/receipts/e64f7946-6791-4536-ad94-c0ccda111f48/points
{'points': 15}

Selected _id='e64f7946-6791-4536-ad94-c0ccda111f48'
(200) GET http://localhost:8080/receipts/e64f7946-6791-4536-ad94-c0ccda111f48/breakdown
{'breakdown': ['9 points for retailer name (Walgreens)',
               '0 points for round dollar amount (2.65)',
               '0 points for being multiple of 0.25 (2.65)',
               '5 points for number of items (2)',
               '0 point(s) for item (Pepsi - 12-oz | 1.25)',
               '1 point(s) for item (Dasani | 1.40)',
               '0 points for purchase day being odd (2022-01-02)',
               '0 points for time of purchase between 2pm and 4pm (08:13)',
               '15 points total']}

(404) GET http://localhost:8080/receipts/abc-123/points
{'error': 'receipt abc-123 not found'}

(404) GET http://localhost:8080/receipts/abc-123/breakdown
{'error': 'receipt abc-123 not found'}

(200) GET http://localhost:8080/receipts/e64f7946-6791-4536-ad94-c0ccda111f48/points
{'points': 15}

(200) GET http://localhost:8080/receipts/e64f7946-6791-4536-ad94-c0ccda111f48/breakdown
{'breakdown': ['9 points for retailer name (Walgreens)',
               '0 points for round dollar amount (2.65)',
               '0 points for being multiple of 0.25 (2.65)',
               '5 points for number of items (2)',
               '0 point(s) for item (Pepsi - 12-oz | 1.25)',
               '1 point(s) for item (Dasani | 1.40)',
               '0 points for purchase day being odd (2022-01-02)',
               '0 points for time of purchase between 2pm and 4pm (08:13)',
               '15 points total']}

(200) GET http://localhost:8080/receipts/2c0dc78d-0c7e-4a91-97cf-296e1c4911dc/points
{'points': 28}

(200) GET http://localhost:8080/receipts/2c0dc78d-0c7e-4a91-97cf-296e1c4911dc/breakdown
{'breakdown': ['6 points for retailer name (Target)',
               '0 points for round dollar amount (35.35)',
               '0 points for being multiple of 0.25 (35.35)',
               '10 points for number of items (5)',
               '0 point(s) for item (Mountain Dew 12PK | 6.49)',
               '3 point(s) for item (Emils Cheese Pizza | 12.25)',
               '0 point(s) for item (Knorr Creamy Chicken | 1.26)',
               '0 point(s) for item (Doritos Nacho Cheese | 3.35)',
               '3 point(s) for item (   Klarbrunn 12-PK 12 FL OZ   | 12.00)',
               '6 points for purchase day being odd (2022-01-01)',
               '0 points for time of purchase between 2pm and 4pm (13:01)',
               '28 points total']}

(200) GET http://localhost:8080/receipts/72f3b285-e192-49d0-940a-ed819a575166/points
{'points': 109}

(200) GET http://localhost:8080/receipts/72f3b285-e192-49d0-940a-ed819a575166/breakdown
{'breakdown': ['14 points for retailer name (M&M Corner Market)',
               '50 points for round dollar amount (9.00)',
               '25 points for being multiple of 0.25 (9.00)',
               '10 points for number of items (4)',
               '0 point(s) for item (Gatorade | 2.25)',
               '0 point(s) for item (Gatorade | 2.25)',
               '0 point(s) for item (Gatorade | 2.25)',
               '0 point(s) for item (Gatorade | 2.25)',
               '0 points for purchase day being odd (2022-03-20)',
               '10 points for time of purchase between 2pm and 4pm (14:33)',
               '109 points total']}
```

> You can choose to remain in the Python interpreter after loading the client.py
> file with the `-i` flag (`venv/bin/python -i client.py`)

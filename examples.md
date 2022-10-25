# Examples
This is example API requests
## GET /user
### Request:
```localhost:8080/v1/user?id=1```

### Response:
```json
{
    "id": 1,
    "amount": "600.00"
}
```

## POST /user

### Request:
```localhost:8080/v1/user```

### Request body:
```json
{
    "id": 1,
    "amount": "600.00"
}
```

### Response:
```json
{}
```

## POST /order

### Request:
```localhost:8080/v1/order```

### Request body:

#### Create
```json
{
  "action": "create",
  "order_id": 1,
  "service_id": 1,
  "user_id": 1,
  "sum": "200"
}
```
#### Approve
```json
{
  "action": "approve",
  "order_id": 1,
  "service_id": 1,
  "user_id": 1,
  "sum": "200"
}
```
#### Cancel
```json
{
  "action": "cancel",
  "order_id": 1,
  "service_id": 1,
  "user_id": 1,
  "sum": "200"
}
```

### Response:
```json
{}
```
## GET /history

### Request:
```localhost:8080/v1/history?id=1&order_by=date&desc=true&limit=2&page=1```

### Response:
```json
{
  "orders": [
    {
      "sum": "200.00",
      "service": "Replenishment",
      "status": "Approved",
      "time": "13:19 24 Oct 22 UTC"
    },
    {
      "sum": "200.00",
      "service": "Replenishment",
      "status": "Approved",
      "time": "13:17 24 Oct 22 UTC"
    }
  ]
}
```

#### Note: Pagination and sort params are optional. If you omit them, you will get whole user's history sorted by date in ascending order. Example:

### Request:
```localhost:8080/v1/history?id=1```

### Response:
```json
{
  "orders": [
    {
      "sum": "200.00",
      "service": "Replenishment",
      "status": "Approved",
      "time": "13:17 24 Oct 22 UTC"
    },
    {
      "sum": "200.00",
      "service": "Replenishment",
      "status": "Approved",
      "time": "13:19 24 Oct 22 UTC"
    },
    {
      "sum": "200.00",
      "service": "Good bought",
      "status": "Approved",
      "time": "13:19 24 Oct 22 UTC"
    },
    {
      "sum": "200.00",
      "service": "Advertisement bought",
      "status": "Approved",
      "time": "13:21 24 Oct 22 UTC"
    }
  ]
}
```

## GET /report

### Request:
```localhost:8080/v1/report?year=2022&month=10```

### Response:
```json
{
  "link": "localhost:8080/v1/reports/2022-10.csv"
}
```
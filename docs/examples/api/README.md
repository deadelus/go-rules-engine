# REST API Example

This example demonstrates how to wrap the `go-rules-engine` in a simple REST API using the Go standard library (`net/http`). 

## Why a REST API?

While the engine is primarily a library, exposing it as a service allows for:
- **Centralized Rules**: One service calculating rules for multiple clients.
- **Cross-language Support**: Non-Go services can evaluate rules via JSON over HTTP.
- **Stateless Evaluation**: Facts are sent in the request, results are returned in the response.

## Features of this Example

1. **Stateless Handlers**: Every request creates a fresh `Almanac`.
2. **Smart Skip Integration**: Rules are automatically skipped if required facts are missing from the request, preventing evaluation errors.
3. **Audit Trace**: The response includes the full evaluation path (if enabled).
4. **Execution Headers**: Returns `X-Execution-Time` in response headers.
5. **Error Handling**: Proper JSON error responses for malformed inputs.

## How to run

```bash
go run docs/examples/api/main.go
```

## How to test

Open another terminal and use `curl` to send facts:

```bash
curl -X POST http://localhost:8080/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "facts": {
      "totalSpend": 1200,
      "accountAgeDays": 400,
      "amount": 100,
      "isFirstPurchase": false
    }
  }'
```

The output will be a JSON object containing the `results` for all evaluated rules.

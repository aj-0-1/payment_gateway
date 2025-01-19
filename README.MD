# Documentation

## How to run the solution
- Clone the repository
- `cd` into the root directory and run `go mod tidy` and `make api-test`
- This will start the server and run curl commands to verify that the API can process a payment and also retrieve the processed payment's details
- Run `make unit-tests` to run various unit tests

```
=== RUN   TestHTTPHandlerSuite
=== RUN   TestHTTPHandlerSuite/TestGetPayment
=== RUN   TestHTTPHandlerSuite/TestGetPayment/payment_found
=== RUN   TestHTTPHandlerSuite/TestGetPayment/payment_not_found
=== RUN   TestHTTPHandlerSuite/TestProcessPayment
=== RUN   TestHTTPHandlerSuite/TestProcessPayment/successful_payment
=== RUN   TestHTTPHandlerSuite/TestProcessPayment/invalid_request
--- PASS: TestHTTPHandlerSuite (0.00s)
    --- PASS: TestHTTPHandlerSuite/TestGetPayment (0.00s)
        --- PASS: TestHTTPHandlerSuite/TestGetPayment/payment_found (0.00s)
        --- PASS: TestHTTPHandlerSuite/TestGetPayment/payment_not_found (0.00s)
    --- PASS: TestHTTPHandlerSuite/TestProcessPayment (0.00s)
        --- PASS: TestHTTPHandlerSuite/TestProcessPayment/successful_payment (0.00s)
        --- PASS: TestHTTPHandlerSuite/TestProcessPayment/invalid_request (0.00s)
=== RUN   TestServiceSuite
=== RUN   TestServiceSuite/TestGetPayment
=== RUN   TestServiceSuite/TestGetPayment/payment_found
=== RUN   TestServiceSuite/TestGetPayment/payment_not_found
=== RUN   TestServiceSuite/TestProcessPayment
=== RUN   TestServiceSuite/TestProcessPayment/successful_payment
=== RUN   TestServiceSuite/TestProcessPayment/processing_failure
--- PASS: TestServiceSuite (0.00s)
    --- PASS: TestServiceSuite/TestGetPayment (0.00s)
        --- PASS: TestServiceSuite/TestGetPayment/payment_found (0.00s)
        --- PASS: TestServiceSuite/TestGetPayment/payment_not_found (0.00s)
    --- PASS: TestServiceSuite/TestProcessPayment (0.00s)
        --- PASS: TestServiceSuite/TestProcessPayment/successful_payment (0.00s)
        --- PASS: TestServiceSuite/TestProcessPayment/processing_failure (0.00s)
PASS
ok      payment_gateway/internal/gateway
amirzan@amirzan-1-0:~/code/go/payment_gateway$ make api-test
Starting payment gateway server...
2025/01/09 22:52:51 Starting payment gateway server

Testing payment processing...
{"payment_id":"34506e4a-2237-4eae-895e-860a9b359bf7","merchant_id":"merchant1","card_details":{"number":"","masked_number":"************1111","expiry_month":11,"expiry_year":2026,"cvv":""},"amount":{"value":100,"currency":"GBP"},"status":"success","created_at":"2025-01-09T22:52:51.875836478Z"}

Testing payment retrieval...
{"payment_id":"34506e4a-2237-4eae-895e-860a9b359bf7","merchant_id":"merchant1","card_details":{"number":"","masked_number":"************1111","expiry_month":11,"expiry_year":2026,"cvv":""},"amount":{"value":100,"currency":"GBP"},"status":"success","created_at":"2025-01-09T22:52:51.875836478Z"}
```

## Approach
- Given the time constraints, we have opted for a simple application architecture that meets our requirements
- We have used interfaces to define the behaviour of components such as `store` and `bank` in order to have flexibility. If the scope/scale of the application increased, we can swap out current components e.g. A relational database or multiple bank APIs without changing business logic
- This follows the Dependency Inversion Principle and allows us to reduce coupling between components and makes it easier to test with mocks
- Testing approach with mocks and table-driven tests that scale well as our tests verify behaviour not implementation details

## Assumptions
Here are some of the assumptions that were made given the scale of the application:
- Basic payment input validation using struct tags for card details + amount to simulate cases for a failed/successful payment process
- An in-memory storage solution is suitable for this scale but we need other tools to help solve problems that occur at greater scales in the real world
- No need for authentication, encryption etc.
- A simple bank simulator as opposed to a service that allows integrations with multiple banks/payment providers

## Areas for improvement + use of cloud technologies

### Development / Testing / Deployment / Monitoring
- Auth + Encryption
- Automated testing across multiple levels
- Metrics + Logging
- Containerisation + CI/CD

### Scalability
Here are some improvements that could be made, leveraging cloud technologies, in order to design a scalable solution:

**Requirements**
- Data consistency + high availability, no accidental double payments, accurate records + audit trail
- Data access, payment history + access, aggregated reports/dashboards

**Payments**
- Accurate payment status management
- Handle retries with expotential backoff
- Idempotency key prevents duplicate payments

**Storage**
- A cloud-native distributed database such as CockroachDB which uses the raft consensus algorithm for strong consistency
- Node-based horizontal scaling, automated fallover, partition by merchant ID, ACID compliance etc.
- Acts as our source of truth table from which we can derive other data sources/actions

**Event Store**
- Can use Kafka to track changes and create an audit trail
- Supports multiple consumers reading from same source e.g. Analytics service
- Loose coupling means we can easily add more integrations in the future

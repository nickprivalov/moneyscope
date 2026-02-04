# MoneyScope

_Money under a scope, to manage your finances._

---

Project is under the `Open Software License v. 3.0`.

---

## How to Run

#### 1. Run the Database via Docker

```bash
docker compose up -d
```

#### 2. Run backend services

```bash
# Ledger
go run backend/ledger/cmd/server/main.go
# Ingest
go run backend/ingest/cmd/server/main.go
# Gateway
go run backend/gateway/cmd/server/main.go
```

#### 3. Uploading a sample file

Given a test file saved as `test.csv` looking like:
```csv
Date,Description,Amount
2026-01-01,Company Salary,3500.00
2026-01-02,Starbucks Coffee,-6.50
2026-01-02,Amazon Refund,65.99
2026-01-03,Landlord Rent,-1200.00
2026-01-04,Uber Trip,-22.40
```

You can upload the test file and see what happens:
```bash
curl -X POST http://localhost:8080/api/upload -F "file=@test.csv"
```

A JSON response should be received:
```json
{
  "transactions_processed": 5,
  "transactions_failed": 0, 
  "message": "Successfully processed 5 transactions"
}
```

You can verify the transactions were entered into the database by:
```bash
docker exec -it moneyscope-db psql -U moneyscope_user -d moneyscope_db -c "SELECT * FROM ledger.transactions;"
```

#### 4. Development / Changes to Protobuf Definitions

For any changes to the protobuf definitions, run:
```bash
make proto
```

If new protobuf definitions are added in a new module (e.g. for `auth` and the `proto/auth` folder doesn't exist), it must be manually added to the `Makefile` command.

If GoLand doesn't register/pick up on the new definitions, hit `Ctrl + Alt + Y` to reload from disk (or see if the "Sync Dependencies" notification is there)

Then, changes to protobuf definitions should be reflected in any Go-code struct that is meant to be them (usually in the `models` folder).

Then, also update any functions that copy data from the Protobuf version to the go-lang struct.

GORM's automigrate will easily handle new columns, but changing existing fields and removing fields will need double-checking of the database.

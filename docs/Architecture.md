# MoneyScope: System Architecture Overview

---

### Other Architecture

#### Infrastructure / Running (Docker & AWS)
- AWS CDK, Docker
- Postgres database
- CDK sub-project for ease of creation/tear-down
  - ECS
  - RDS
  - Secrets Manager
  - etc. needed infrastructure (VPC, subnets, etc.)

#### Frontend (React)
- React, Vite, Material UI, Recharts
- Possibles:
  - Angular
    - Material, charts library TBD
  - [Flutter](https://flutter.dev/)
  - [Wails](https://github.com/wailsapp/wails)
- JWT auth/management
- Allows uploading for `Ingest` service
- Visualize financial data
  - Calendar view to see per-day/week spending patterns
- Operates exclusively with `Gateway` service

---

### Data Flows

#### The Upload Flow
1.  User uploads `statements.csv` via `Frontend`.
2.  `Gateway` receives HTTP POST stream.
3.  `Gateway` opens gRPC stream to `Ingest`.
4.  `Ingest` reads stream, buffers chunks, and parses lines via `encoding/csv`.
5.  `Ingest` creates batches of 50 transactions.
6.  `Ingest` calls `Ledger.CreateTransactionsBatch` (gRPC).
7.  `Ledger` starts a DB transaction, inserts rows into `ledger.transactions`, and commits.
8.  `Ingest` returns a summary report (e.g., "50 processed, 0 failed") up the chain to the User.

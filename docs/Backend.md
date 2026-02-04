# MoneyScope: Backend Architecture

---

### Backend

#### Common (`backend/common`)
- Shared code
- Protobuf definitions
- Go `struct` definitions
- Database connectivity (incl. GORM)

#### Authentication Service (`backend/auth`)
- General authentication/identity/credentials handling
- Registration
- Login/Logout
- Token validation middleware for `Gateway` service
- Handles operations for `Profile` service to manage actual user profile

#### Profile Service (`backend/profile`)
- Handles user's profile and personal settings (post-auth)

#### Gateway Service (`backend/gateway`)
- Gin, gRPC
- Endpoint management
    - Routes endpoints to appropriate services/functions
- CORS management
- File streaming proxy for `Ingest` (HTTP file stream piped directly into gRPC client stream)
- Translate JSON to structs
- **[Idea]:** batch responses from across services into response for frontend

#### Ingest Service (`backend/ingest`)
- gRPC
- Processing file uploads and "translating" them into appropriate data structures
- Reads file stream, turns line items into `Transaction` objects, batches them to be sent to `Ledger` service
- Passes "cleaned" data to `Ledger`
- (best-guess try?) categorize incoming transactions
- **[Idea]:** options for (on upload) tweaking the CSV-to-upload's format (e.g. ignore first _n_ rows, etc.)

#### Ledger Service (`backend/ledger`)
- gRPC, GORM
- "Source of truth" for account, balances, transactions, etc.
- Performs CRUD operations for database
    - transaction history en masse (including pagination, filtering, searching, etc.)
    - modifying transactions post-upload
    - manage accounts, savings buckets, debts (secured and unsecured)
- Ensures financial data consistency/correctness/integrity

#### Analytics Service (`backend/analytics`)
- _Unknown what other tech may be used; potentially AI plugins/calls? (Gemini?)_
- Provides analytics/insights on financial data
- Performs CRUD operations for analytics (creating categories of spending, etc.)
- Calculate spending patterns for front-end chart displays
- Handling of budget tracking feature (create user-defined budget plans, savings buckets, etc.)
    - Maybe `Ledger` responsibility?
# MoneyScope Backend

---

The backend operates in a quite **monolithic** way, but with great modularity for organization.

The server (`backend/server`) creates and sets up the various service objects (`backend/db`, `backend/auth`, etc.).

When needed, services may run functions and features in goroutines.

### Main Server (`backend/server`)
- Actual main runnable
- Initializations, connections, etc.

### Common (`backend/common`)
- Shared code
- Shared structs 
  - JSON "translating" for these structs should occur here
  - For structs in other services/modules, their "JSON translating" should occur with the struct's definition

### Database (`backend/db`)
- Database-related functionalities
- Connectivity
- DAOs, entities

### Authentication Service (`backend/auth`)
- General authentication, identity and credentials management
- Registration
- Login, logout
- Token validation middleware for all routes

### Profile Service (`backend/profile`)
- User's profile
- Settings and preferences management

### Gateway Service (`backend/gateway`)
- Endpoint routing
- CORS management
- File streaming proxy for `Ingest` service
- Can batch function responses from other services into JSON responses for frontend

### Ingest Service (`backend/ingest`)
- Handles file uploading and "translating" them into respective data structures
- Passes data to `Ledger` service
- Can take configuration options for CSV uploading

### Ledger Service (`backend/ledger`)
- "Source of truth" for accounts, balances, transactions, debts, etc.
- Performs CRUD operations on database
  - Accounts
    - Bank accounts (checking, savings)
    - Health Savings Accounts (HSAs)
    - Retirement Accounts (401K, Roth IRAs)
    - Investment/Trading Accounts
    - Debts (secured/unsecured)
  - Budgets
    - Buckets within accounts
    - Savings goals
    - Payoff plans
  - Transactions
  - Spending categories
  - etc.

### Analytics Service (`backend/analytics`)
- _This is a "maybe" service, consider the whole service to be an optional feature_
- AI-powered analytics
  - Everything that entails configuring it (API key, etc.)
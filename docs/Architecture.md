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

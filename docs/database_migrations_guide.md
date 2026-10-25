# Database Migrations & Seed Data Guide

This document describes how to run and manage database migrations and seed data for the **NC Business Management Software** backend.

---

## 1. Overview of Artifacts

| Component | Path | Description |
|---|---|---|
| **Schema Migration Up** | [`migrations/000001_init_schema.up.sql`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/migrations/000001_init_schema.up.sql) | Creates all 33 tables, PostgreSQL extensions, constraints, and indexes. |
| **Schema Migration Down** | [`migrations/000001_init_schema.down.sql`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/migrations/000001_init_schema.down.sql) | Cleanly tears down all tables and extensions in reverse dependency order. |
| **Seed Data Up** | [`migrations/000002_seed_baseline_data.up.sql`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/migrations/000002_seed_baseline_data.up.sql) | Comprehensive seed data (Roles, Permissions, Admin User, Store, Settings, Categories, Products, Counters). |
| **Seed Data Down** | [`migrations/000002_seed_baseline_data.down.sql`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/migrations/000002_seed_baseline_data.down.sql) | Reverses all seeded rows cleanly. |
| **Embedded Go Migrations** | [`migrations/migrations.go`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/migrations/migrations.go) | Embeds all `.sql` files into the compiled binary using `embed.FS`. |
| **Migration Engine** | [`pkg/database/migrate.go`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/pkg/database/migrate.go) | Transactional Go migration runner tracking applied versions in `schema_migrations`. |
| **Database Pool** | [`pkg/database/db.go`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/pkg/database/db.go) | PostgreSQL connection pool and env configuration parser. |
| **CLI Runner** | [`cmd/migrate/main.go`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/cmd/migrate/main.go) | CLI tool with `up`, `down`, `reset`, and `status` commands. |
| **Docker Compose** | [`docker-compose.yml`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/docker-compose.yml) | Spawns a PostgreSQL 16 container with health check and persistent storage. |
| **Env Template** | [`.env.example`](file:///e:/NC%20Business%20Management%20Software/business-mgmt-backend/.env.example) | Database and application connection settings. |

---

## 2. Starting PostgreSQL with Docker

To start a clean PostgreSQL 16 instance locally:

```bash
docker compose up -d
```

To stop PostgreSQL:

```bash
docker compose down
```

---

## 3. Running Migrations via CLI

The project includes a pure Go migration tool that executes transactions and tracks state in `schema_migrations`:

### Check Migration Status
```bash
go run ./cmd/migrate status
```

### Apply Migrations (Schema + Baseline Seed)
```bash
go run ./cmd/migrate up
```

### Roll Back Last Migration Step
```bash
go run ./cmd/migrate down
```

### Reset Database from Scratch (Roll Back All & Re-apply)
```bash
go run ./cmd/migrate reset
```

---

## 4. Config Environment Variables

You can configure database credentials either in your environment or via `.env`:

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | *(optional)* | Full URL: `postgres://user:pass@host:5432/dbname?sslmode=disable` |
| `DB_HOST` | `localhost` | PostgreSQL server hostname |
| `DB_PORT` | `5432` | PostgreSQL server port |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `business_mgmt`| Database name |
| `DB_SSLMODE` | `disable` | SSL mode (`disable`, `require`, etc.) |
| `AUTO_MIGRATE` | `false` | When set to `true`, the API server runs `MigrateUp` on boot. |

---

## 5. Seed Data Summary

The baseline seed migration (`000002_seed_baseline_data.up.sql`) includes:

### Roles & Permissions
- **Roles**: `SUPER_ADMIN`, `STORE_ADMIN`, `FINANCE_ADMIN`, `STAFF`
- **Granular Permissions**: 25+ permissions across IAM, Stores, Customers, Products, Suppliers, Purchases, Inventory, Sales, Payments, Expenses, Notifications, and Reports.
- **Role-Permission Mappings**: Configured for all 4 roles.

### Default Admin Account
- **Email**: `admin@example.com`
- **Password**: `Password123!`
- **Role**: `SUPER_ADMIN`
- **Store Access**: Assigned to `Main Branch - Central Store` (`STR-001`)

### Master Data & Settings
- **Default Store**: `STR-001` (*Main Branch - Central Store*, Bengaluru)
- **Company Settings**: Currency `INR`, Timezone `Asia/Kolkata`, default credit `15 days`, overdue threshold `1` (yellow) & `7` (red) days.
- **Product Categories**: Cement & Concrete, Steel & TMT Rebars, Paints & Primers, Plumbing & Pipes, Electrical & Wiring, Hardware & Fasteners.
- **Sample Products**: UltraTech Cement 50KG, Tata Tiscon TMT Bar 12mm, Asian Paints Apex 20L, Astral PVC Pipe 110mm, Havells Wire 2.5 sq mm.
- **Initial Inventory**: Seeded for `STR-001` with average unit costs and quantities.
- **Sample Contacts**: Supplier (*Tata Steel & Building Materials Ltd*), Customer (*Rahul Sharma*).
- **Expense Categories**: Rent, Electricity, Staff Refreshments, Freight, Marketing, Legal.
- **Notification Templates**: WhatsApp payment reminder, invoice generated, payment received.
- **Reminder Rules**: Automatic follow-up on Day 7, Day 9, Day 11.
- **Document Counters**: Initialized for FY 2026-27 (`INVOICE`, `PURCHASE`, `PAYMENT`, `EXPENSE`).

---
description: Guardrails for writing or modifying SQL upgrade scripts under server/sql/. Use before creating or editing any incremental SQL file.
---

# SQL Upgrade Guardrails

Enforce safe, idempotent, MySQL-compatible DDL and DML for incremental upgrade scripts.

## Pre-Flight Checks

Before writing or editing any file under `server/sql/`:

// turbo
1. Read the baseline snapshot `go-base.sql` (search for the relevant table/column, do not read the whole file)
// turbo
2. Read the nearest related upgrade scripts under `server/sql/` to understand existing patterns
3. Identify the exact MySQL version target: treat this repository as **Oracle MySQL 8.x** by default, not MariaDB

## Mandatory Rules

### No Foreign Keys

- **Do NOT** add `FOREIGN KEY` or `REFERENCES` clauses in `CREATE TABLE` or `ALTER TABLE`
- Use plain columns with indexes for cross-table references
- Referential integrity is enforced in application code only

### Idempotent & Rerunnable

- Every script must be safe to execute multiple times
- Use `information_schema` checks or prepared-statement guards for DDL changes
- Use `INSERT ... SELECT ... WHERE NOT EXISTS` or `ON DUPLICATE KEY UPDATE` for seed data
- Use `WHERE` guards on old values for `UPDATE` statements

### MySQL DDL Compatibility

Do **NOT** use these unsupported MySQL incremental DDL patterns:

- `ALTER TABLE ... ADD COLUMN IF NOT EXISTS ...`
- `ALTER TABLE ... ADD INDEX IF NOT EXISTS ...`
- Other unverified `IF [NOT] EXISTS` forms inside `ALTER TABLE`

Instead, use idempotent guards via `information_schema` + dynamic SQL:

```sql
-- Example: add column if missing
SELECT COUNT(*) INTO @col_exists
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'target_table'
  AND COLUMN_NAME = 'new_column';

SET @sql = IF(@col_exists = 0,
    'ALTER TABLE `target_table` ADD COLUMN `new_column` varchar(100) DEFAULT NULL COMMENT ''描述''',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
```

For `CREATE TABLE`, use the standard `IF NOT EXISTS`:

```sql
SELECT COUNT(*) INTO @tbl_exists
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = DATABASE()
  AND TABLE_NAME = 'new_table';

SET @sql = IF(@tbl_exists = 0,
    'CREATE TABLE `new_table` ( ... ) ENGINE=InnoDB CHARACTER SET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC',
    'SELECT 1');
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
```

### Seed Data Safety

- Menu rows, API rows, permission rows, config rows must be duplicate-safe
- Use `INSERT IGNORE`, `ON DUPLICATE KEY UPDATE`, or `NOT EXISTS` subquery
- Only update rows that actually need updating — use `WHERE` guards on old values

### Script Scope

- Keep incremental SQL limited to the current feature
- Do not mix unrelated schema cleanup into the same script
- Do not rewrite `go-base.sql` for normal feature delivery

### Script Header

Every script must start with:

```sql
-- <brief description of what this script does>
-- 执行方法: mysql -u root -p go-base < server/sql/<filename>.sql

SET NAMES utf8mb4;
```

## Self-Check

Before finalizing the script, verify:

- [ ] No `FOREIGN KEY` or `REFERENCES` clauses
- [ ] All DDL is idempotent (safe to rerun)
- [ ] No unsupported `IF NOT EXISTS` inside `ALTER TABLE`
- [ ] Seed data inserts are duplicate-safe
- [ ] `UPDATE` statements have `WHERE` guards on old values
- [ ] Script is scoped to the current feature only
- [ ] Filename is descriptive: `server/sql/<action>_<feature>_<detail>.sql`
- [ ] Script is referenced in commit/PR description

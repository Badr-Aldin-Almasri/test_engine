# Database Migration Guide

## Overview
This migration separates nodes from the `flows` table into a dedicated `flow_nodes` table for better normalization and query performance.

## Migration Steps

### Option 1: Using the migration script (Recommended)

1. **Navigate to the database directory:**
   ```bash
   cd server/database
   ```

2. **Set environment variables (if needed):**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5433
   export DB_NAME=visual_testing_db
   export DB_USER=visual_testing
   export DB_PASSWORD=visual_testing_pass
   ```

3. **Run the migration:**
   ```bash
   ./migrate.sh
   ```

### Option 2: Manual migration using psql

1. **Connect to the database:**
   ```bash
   psql -h localhost -p 5433 -U visual_testing -d visual_testing_db
   ```

2. **Run the migration SQL:**
   ```bash
   \i migration_001_separate_nodes.sql
   ```

### Option 3: Using Docker

If you're using Docker Compose:

1. **Find the database container:**
   ```bash
   docker ps | grep postgres
   ```

2. **Copy the migration file into the container:**
   ```bash
   docker cp migration_001_separate_nodes.sql <container_name>:/tmp/
   ```

3. **Execute the migration:**
   ```bash
   docker exec -i <container_name> psql -U visual_testing -d visual_testing_db < /tmp/migration_001_separate_nodes.sql
   ```

Or directly:
```bash
docker exec -i <container_name> psql -U visual_testing -d visual_testing_db -f /tmp/migration_001_separate_nodes.sql
```

## What the Migration Does

1. **Creates `flow_nodes` table** with all necessary columns
2. **Migrates existing data** from `flows.nodes` JSONB column to `flow_nodes` table
3. **Creates indexes** for better query performance
4. **Sets up triggers** for automatic timestamp updates

## Verification

After running the migration, verify the data was migrated correctly:

```sql
-- Check if nodes were migrated
SELECT COUNT(*) FROM flow_nodes;

-- Check a specific flow's nodes
SELECT fn.* 
FROM flow_nodes fn 
WHERE fn.flow_id = '<some-flow-id>';

-- Compare with original data (if nodes column still exists)
SELECT id, jsonb_array_length(nodes) as node_count 
FROM flows 
WHERE nodes IS NOT NULL;
```

## Final Step (After Verification)

Once you've verified the migration worked correctly, remove the old `nodes` column:

```sql
ALTER TABLE flows DROP COLUMN IF EXISTS nodes;
```

## Rollback (If Needed)

If you need to rollback the migration:

1. **Restore nodes to flows table:**
   ```sql
   -- Add nodes column back
   ALTER TABLE flows ADD COLUMN nodes JSONB DEFAULT '[]'::jsonb;
   
   -- Re-populate nodes from flow_nodes table
   UPDATE flows f
   SET nodes = (
       SELECT jsonb_agg(
           jsonb_build_object(
               'id', fn.node_id,
               'type', fn.type,
               'position', jsonb_build_object('x', fn.position_x, 'y', fn.position_y),
               'data', jsonb_build_object(
                   'id', fn.data_id,
                   'type', fn.data_type,
                   'label', fn.data_label,
                   'status', fn.data_status,
                   'config', COALESCE(fn.data_config, '{}'::jsonb),
                   'output', fn.data_output,
                   'error', fn.data_error
               )
           )
       )
       FROM flow_nodes fn
       WHERE fn.flow_id = f.id
   );
   ```

2. **Drop the flow_nodes table:**
   ```sql
   DROP TABLE IF EXISTS flow_nodes CASCADE;
   ```

## Notes

- The migration is **idempotent** - you can run it multiple times safely
- Existing data is preserved during migration
- The `nodes` column remains in the `flows` table until you manually remove it (for safety)
- All foreign key constraints are preserved with CASCADE delete


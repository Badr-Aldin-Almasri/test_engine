#!/bin/bash

# Database migration script
# This script runs the migration to separate nodes into their own table

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database connection details
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5433}"
DB_NAME="${DB_NAME:-visual_testing_db}"
DB_USER="${DB_USER:-visual_testing}"
DB_PASSWORD="${DB_PASSWORD:-visual_testing_pass}"

echo -e "${YELLOW}Starting database migration...${NC}"
echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo ""

# Check if migration file exists
if [ ! -f "migration_001_separate_nodes.sql" ]; then
    echo -e "${RED}Error: migration_001_separate_nodes.sql not found${NC}"
    exit 1
fi

# Run migration
echo -e "${YELLOW}Running migration_001_separate_nodes.sql...${NC}"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migration_001_separate_nodes.sql

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Migration completed successfully!${NC}"
    echo ""
    echo -e "${YELLOW}Note: The 'nodes' column is still in the flows table for safety.${NC}"
    echo -e "${YELLOW}After verifying the migration worked, you can remove it manually:${NC}"
    echo -e "${YELLOW}ALTER TABLE flows DROP COLUMN IF EXISTS nodes;${NC}"
else
    echo -e "${RED}Migration failed!${NC}"
    exit 1
fi


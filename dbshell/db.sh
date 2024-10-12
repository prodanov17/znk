#!/bin/bash
set -e
export PGPASSWORD=postgres123;
echo "Init postgres db..."
psql -v ON_ERROR_STOP=1 --username "postgres" <<-EOSQL
  CREATE DATABASE znk;
  GRANT ALL PRIVILEGES ON DATABASE znk TO "postgres";
EOSQL
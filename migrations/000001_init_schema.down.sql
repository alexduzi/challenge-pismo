-- migrations/000001_init_schema.down.sql

DROP TABLE IF EXISTS transactions CASCADE;
DROP TABLE IF EXISTS operation_types CASCADE;
DROP TABLE IF EXISTS accounts CASCADE;
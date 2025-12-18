-- Clean old data and reset database
-- Run this SQL script in your PostgreSQL database if you have old data

-- Drop tables in correct order (reverse of creation)
DROP TABLE IF EXISTS balais CASCADE;
DROP TABLE IF EXISTS regencies CASCADE;
DROP TABLE IF EXISTS provinces CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Tables will be recreated by GORM AutoMigrate

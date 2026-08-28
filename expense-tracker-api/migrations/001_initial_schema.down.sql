-- migrations/001_initial_schema.down.sql

DROP INDEX IF EXISTS idx_audit_logs_created;
DROP INDEX IF EXISTS idx_audit_logs_entity;
DROP INDEX IF EXISTS idx_otp_verifications_expires;
DROP INDEX IF EXISTS idx_otp_verifications_user_id;
DROP INDEX IF EXISTS idx_transactions_family_date;
DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_transactions_family_id;
DROP INDEX IF EXISTS idx_transactions_user_id;
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_family_invitations_email;
DROP INDEX IF EXISTS idx_family_invitations_family_id;
DROP INDEX IF EXISTS idx_family_members_family_role;
DROP INDEX IF EXISTS idx_family_members_user_id;
DROP INDEX IF EXISTS idx_families_created_by;
DROP INDEX IF EXISTS idx_families_unique_code;
DROP INDEX IF EXISTS idx_users_google_id;
DROP INDEX IF EXISTS idx_users_email;

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS upi_apps;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS transaction_categories;
DROP TABLE IF EXISTS family_invitations;
DROP TABLE IF EXISTS family_members;
DROP TABLE IF EXISTS families;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS transaction_type;
package config

import (
	"context"
	"log"
)

func UserTable() error {
	query := `
	-- Extensions
	CREATE EXTENSION IF NOT EXISTS pgcrypto;
	CREATE EXTENSION IF NOT EXISTS citext;

	-- Departments Table
	CREATE TABLE IF NOT EXISTS departments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT UNIQUE NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMPTZ,
		created_by UUID,
		updated_by UUID,
		deleted_by UUID
	);

	-- Users Table
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email CITEXT UNIQUE NOT NULL,
		name TEXT NOT NULL CHECK (char_length(name) >= 2),
		password_hash TEXT,
		role TEXT NOT NULL DEFAULT 'applicant'
			CHECK (role IN ('admin','checker','approver','applicant','finance')),
		department_id UUID,
		is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
		email_verified BOOLEAN NOT NULL DEFAULT FALSE,
		email_verified_at TIMESTAMPTZ,
		force_password_change BOOLEAN NOT NULL DEFAULT FALSE,
		last_login_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMPTZ,
		created_by UUID,
		updated_by UUID,
		deleted_by UUID
	);

	-- Add columns if they don't exist (Migration support)
	DO $$ BEGIN
		-- Departments
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='departments' AND column_name='deleted_at') THEN
			ALTER TABLE departments ADD COLUMN deleted_at TIMESTAMPTZ;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='departments' AND column_name='created_by') THEN
			ALTER TABLE departments ADD COLUMN created_by UUID;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='departments' AND column_name='updated_by') THEN
			ALTER TABLE departments ADD COLUMN updated_by UUID;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='departments' AND column_name='deleted_by') THEN
			ALTER TABLE departments ADD COLUMN deleted_by UUID;
		END IF;

		-- Users
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='department_id') THEN
			ALTER TABLE users ADD COLUMN department_id UUID;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='created_by') THEN
			ALTER TABLE users ADD COLUMN created_by UUID;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='updated_by') THEN
			ALTER TABLE users ADD COLUMN updated_by UUID;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='deleted_by') THEN
			ALTER TABLE users ADD COLUMN deleted_by UUID;
		END IF;

		IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='department') THEN
			ALTER TABLE users DROP COLUMN department;
		END IF;

		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='force_password_change') THEN
			ALTER TABLE users ADD COLUMN force_password_change BOOLEAN NOT NULL DEFAULT FALSE;
		END IF;

		-- Seed departments if none exist
		IF NOT EXISTS (SELECT 1 FROM departments) THEN
			INSERT INTO departments (name) VALUES 
			('Implementation'),
			('Application & Systems'),
			('Finance');
		END IF;

		-- Link existing users to first dept if needed
		UPDATE users SET department_id = (SELECT id FROM departments ORDER BY created_at LIMIT 1) WHERE department_id IS NULL;
	END $$;

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
	CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);
	`

	_, err := DB.Exec(context.Background(), query)
	return err
}

func UserTrigger() error {
	query := `
	CREATE OR REPLACE FUNCTION set_updated_at()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
	CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();

	DROP TRIGGER IF EXISTS trg_departments_updated_at ON departments;
	CREATE TRIGGER trg_departments_updated_at BEFORE UPDATE ON departments FOR EACH ROW EXECUTE FUNCTION set_updated_at();

	DROP TRIGGER IF EXISTS trg_overtimes_updated_at ON overtimes;
	CREATE TRIGGER trg_overtimes_updated_at BEFORE UPDATE ON overtimes FOR EACH ROW EXECUTE FUNCTION set_updated_at();
	`

	_, err := DB.Exec(context.Background(), query)
	return err
}

func OvertimeTable() error {
	query := `
	CREATE EXTENSION IF NOT EXISTS pgcrypto;

	-- Enum types (safe + reusable)
	DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'overtime_status') THEN
			CREATE TYPE overtime_status AS ENUM ('pending', 'checked', 'approved', 'rejected');
		END IF;
	END $$;

	DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'overtime_program') THEN
			CREATE TYPE overtime_program AS ENUM ('night', 'weekend', 'holiday');
		END IF;
	END $$;

	-- Table
	CREATE TABLE IF NOT EXISTS overtimes (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		date DATE NOT NULL,
		start_time TIME NOT NULL,
		end_time TIME NOT NULL,
		job_done TEXT NOT NULL CHECK (char_length(job_done) >= 3),
		status overtime_status NOT NULL DEFAULT 'pending',
		department_id UUID NOT NULL,
		program overtime_program NOT NULL,
		duration NUMERIC(5,2) NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMPTZ,
		created_by UUID,
		updated_by UUID,
		deleted_by UUID,

		-- Constraints
		CONSTRAINT fk_overtime_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		CONSTRAINT valid_time_range CHECK (start_time < end_time)
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_overtimes_user_id ON overtimes(user_id);
	CREATE INDEX IF NOT EXISTS idx_overtimes_status ON overtimes(status);
	CREATE INDEX IF NOT EXISTS idx_overtimes_date ON overtimes(date);
	CREATE INDEX IF NOT EXISTS idx_overtimes_program ON overtimes(program);

	DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='overtimes' AND column_name='department_id') THEN
			ALTER TABLE overtimes ADD COLUMN department_id UUID;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='overtimes' AND column_name='created_by') THEN
			ALTER TABLE overtimes ADD COLUMN created_by UUID;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='overtimes' AND column_name='updated_by') THEN
			ALTER TABLE overtimes ADD COLUMN updated_by UUID;
		END IF;
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='overtimes' AND column_name='deleted_by') THEN
			ALTER TABLE overtimes ADD COLUMN deleted_by UUID;
		END IF;

		IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='overtimes' AND column_name='department') THEN
			ALTER TABLE overtimes DROP COLUMN department;
		END IF;

		-- Link existing overtimes to user's current dept, or first dept as fallback
		UPDATE overtimes o SET department_id = u.department_id FROM users u WHERE o.user_id = u.id AND o.department_id IS NULL;
		UPDATE overtimes SET department_id = (SELECT id FROM departments ORDER BY created_at LIMIT 1) WHERE department_id IS NULL;
	END $$;
	`

	_, err := DB.Exec(context.Background(), query)
	return err
}

func Migrate() error {
	if err := UserTable(); err != nil {
		return err
	}

	if err := OvertimeTable(); err != nil {
		return err
	}

	if err := UserTrigger(); err != nil {
		return err
	}

	return nil
}

func DropTables() {
	_, err := DB.Exec(context.Background(), `DROP TABLE IF EXISTS users;`)
	if err != nil {
		log.Fatal(err)
	}
	// _, err = DB.Exec(context.Background(), `DROP TABLE IF EXISTS overtimes;`)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

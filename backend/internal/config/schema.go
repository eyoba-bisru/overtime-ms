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

	-- Users Table
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

		email CITEXT UNIQUE NOT NULL,

		name TEXT NOT NULL CHECK (char_length(name) >= 2),

		password_hash TEXT,

		role TEXT NOT NULL DEFAULT 'applicant'
			CHECK (role IN ('admin','checker','approver','applicant','finance')),

		is_blocked BOOLEAN NOT NULL DEFAULT FALSE,

		email_verified BOOLEAN NOT NULL DEFAULT FALSE,
		email_verified_at TIMESTAMPTZ,

		force_password_change BOOLEAN NOT NULL DEFAULT FALSE,

		last_login_at TIMESTAMPTZ,

		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMPTZ
	);

	-- Add columns if they don't exist (Migration support)
	DO $$ BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='force_password_change') THEN
			ALTER TABLE users ADD COLUMN force_password_change BOOLEAN NOT NULL DEFAULT FALSE;
		END IF;
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

	CREATE TRIGGER trg_users_updated_at
	BEFORE UPDATE ON users
	FOR EACH ROW
	EXECUTE FUNCTION set_updated_at();
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

		program overtime_program NOT NULL,
		duration NUMERIC(5,2) NOT NULL DEFAULT 0,

		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		deleted_at TIMESTAMPTZ,

		-- Constraints
		CONSTRAINT fk_overtime_user
			FOREIGN KEY (user_id)
			REFERENCES users(id)
			ON DELETE CASCADE,

		CONSTRAINT valid_time_range
			CHECK (start_time < end_time)
	);

	-- Indexes (VERY IMPORTANT for performance)
	CREATE INDEX IF NOT EXISTS idx_overtimes_user_id ON overtimes(user_id);
	CREATE INDEX IF NOT EXISTS idx_overtimes_status ON overtimes(status);
	CREATE INDEX IF NOT EXISTS idx_overtimes_date ON overtimes(date);
	CREATE INDEX IF NOT EXISTS idx_overtimes_program ON overtimes(program);
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

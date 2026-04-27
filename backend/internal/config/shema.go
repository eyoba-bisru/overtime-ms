package config

import "log"

func UserTable() error {
	query := `
	CREATE EXTENSION IF NOT EXISTS pgcrypto;

	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

		email CITEXT UNIQUE NOT NULL,

		name TEXT NOT NULL
			CHECK (char_length(name) >= 2),

		password TEXT NOT NULL
			CHECK (char_length(password) >= 6),

		role TEXT NOT NULL DEFAULT 'applicant'
			CHECK (role IN ('admin','checker','approver','applicant')),

		is_blocked BOOLEAN NOT NULL DEFAULT FALSE,

		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	);
	`

	_, err := DB.Exec(query)
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

	_, err := DB.Exec(query)
	return err
}

// func OvertimeRequestTable() {
// 	_, err := DB.Exec(`
// 	CREATE TABLE IF NOT EXISTS overtime_requests (
// 		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
// 		user_id UUID NOT NULL CHECK (user_id IN (SELECT id FROM users)),
// 		request_date DATE NOT NULL,
// 		hours DECIMAL(10, 2) NOT NULL,
// 		reason VARCHAR(255) NOT NULL,
// 		status VARCHAR(50) NOT NULL,
// 		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
// 		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
// 		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
// 	);
// `)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func Migrate() error {
	if err := UserTable(); err != nil {
		return err
	}

	if err := UserTrigger(); err != nil {
		return err
	}

	return nil
}

func DropTables() {
	_, err := DB.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(`DROP TABLE IF EXISTS overtime_requests;`)
	if err != nil {
		log.Fatal(err)
	}
}

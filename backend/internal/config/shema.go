package config

func EnableExtensions() {
	_, err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)
	if err != nil {
		panic(err)
	}
}

func UserTable() {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
		name VARCHAR(255) NOT NULL CHECK (LENGTH(name) >= 2),
		password VARCHAR(255) NOT NULL CHECK (LENGTH(password) >= 6),
		role VARCHAR(50) NOT NULL DEFAULT 'applicant' CHECK (role IN ('admin', 'checker', 'approver', 'applicant')),
		is_blocked BOOLEAN NOT NULL DEFAULT FALSE CHECK (is_blocked IN (TRUE, FALSE)),
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE (email)
	);
`)
	if err != nil {
		panic(err)
	}
}

// func OvertimeRequestTable() {
// 	_, err := DB.Exec(`
// 	CREATE TABLE IF NOT EXISTS overtime_requests (
// 		id UUID PRIMARY KEY,
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

func CreateTables() {
	EnableExtensions()
	UserTable()
	// OvertimeRequestTable()
}

func DropTables() {
	_, err := DB.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec(`DROP TABLE IF EXISTS overtime_requests;`)
	if err != nil {
		panic(err)
	}
}

package config

func UserTable() {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
`)
	if err != nil {
		panic(err)
	}
}

func OvertimeRequestTable() {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS overtime_requests (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		request_date DATE NOT NULL,
		hours DECIMAL(10, 2) NOT NULL,
		reason VARCHAR(255) NOT NULL,
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
`)
	if err != nil {
		panic(err)
	}
}

func CreateTables() {
	UserTable()
	OvertimeRequestTable()
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

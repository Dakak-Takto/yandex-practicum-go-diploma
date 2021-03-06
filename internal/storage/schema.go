package storage

const schema string = `
CREATE TABLE IF NOT EXISTS users (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	login VARCHAR(128) NOT NULL, 
	password VARCHAR(64) NOT NULL,
	balance FLOAT NOT NULL DEFAULT (0)
);

CREATE TABLE IF NOT EXISTS orders (
	number VARCHAR(15) NOT NULL PRIMARY KEY, 
	status VARCHAR(20) NOT NULL DEFAULT ('NEW'),
	accrual FLOAT NOT NULL DEFAULT (0), 
	user_id INT NOT NULL REFERENCES users(id),
	uploaded_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS withdrawals (
	order_number VARCHAR(15) NOT NULL PRIMARY KEY,
	sum FLOAT NOT NULL DEFAULT(0),
	user_id INT NOT NULL REFERENCES users(id),
	processed_at TIMESTAMP DEFAULT NOW()
);`

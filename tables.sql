-- Base Account Table
CREATE TABLE Account (
	ID INT IDENTITY(1,1) PRIMARY KEY NOT NULL,
	FullName NVARCHAR(50) NOT NULL,
	Email NVARCHAR(100) NOT NULL,
	AccountType INT NOT NULL,
	Password NVARCHAR(500) NOT NULL,
	CreatedAt DATETIME NOT NULL,
	LastLogin DATETIME NOT NULL
);
package health

import (
	"context"
	"flyhorizons-userservice/repositories"
)

type DatabaseCheck struct {
	Repository *repositories.BaseRepository
}

func (c DatabaseCheck) Name() string {
	return "mssql-db"
}

func (c DatabaseCheck) Pass() bool {
	db, err := c.Repository.CreateConnection()
	if err != nil {
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	ctx := context.TODO()
	return sqlDB.PingContext(ctx) == nil
}

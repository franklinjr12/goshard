package config

import (
	"database/sql"
	"fmt"
	"goshard/internal/database"
	"goshard/internal/dbmapper"
)

const configDbConnectionString = "host=localhost port=5432 user=postgres password=postgres dbname=goshardconfig sslmode=disable"

func LoadDbMappings() error {
	db, err := database.Connect(configDbConnectionString)
	if err != nil {
		return err
	}
	defer database.Close(db)

	rows, err := db.Query("SELECT shardid, sharduid, dsn FROM database_mappings")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var shardid int
		var sharduid sql.NullString
		var dsn string
		err = rows.Scan(&shardid, &sharduid, &dsn)
		if err != nil {
			return err
		}
		if shardid != 0 {
			dbmapper.AddDbMapId(uint64(shardid), dsn)
		} else {
			dbmapper.AddDbMapUid(sharduid.String, dsn)
		}
	}

	return nil
}

func QueryUserIdFromDbConfig(userToken string) (uint64, error) {
	db, err := database.Connect(configDbConnectionString)
	if err != nil {
		return 0, err
	}
	defer database.Close(db)
	var userId uint64
	err = db.QueryRow("SELECT id FROM users WHERE token = $1", userToken).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("failed to query user id from db: %w", err)
	}
	return userId, nil
}

func ReadSchemaFromDbConfig(userId uint64) (string, error) {
	db, err := database.Connect(configDbConnectionString)
	if err != nil {
		return "", err
	}
	defer database.Close(db)
	var schema string
	err = db.QueryRow("SELECT schema FROM user_schemas WHERE user_id = $1", userId).Scan(&schema)
	if err != nil {
		return "", fmt.Errorf("failed to read schema from db: %w", err)
	}
	return schema, nil
}

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

	rows, err := db.Query("SELECT user_id, shardid, sharduid, dsn FROM database_mappings")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var userId uint64
		var shardid uint64
		var sharduid sql.NullString
		var dsn string
		err = rows.Scan(&userId, &shardid, &sharduid, &dsn)
		if err != nil {
			return err
		}
		if userId != 0 {
			dbmapper.AddDbMapWithUserId(userId, dbmapper.DbMap{
				Shardid:  shardid,
				Sharduid: sharduid.String,
				Dsn:      dbmapper.DbConnectionString(dsn),
			})
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

func SchemaExists(userId uint64) (bool, error) {
	db, err := database.Connect(configDbConnectionString)
	if err != nil {
		return false, err
	}
	defer database.Close(db)
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_schemas WHERE user_id = $1)", userId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if schema exists: %w", err)
	}
	return exists, nil
}

func UpdateSchemaInDbConfig(userId uint64, schema string) error {
	db, err := database.Connect(configDbConnectionString)
	if err != nil {
		return err
	}
	defer database.Close(db)
	_, err = db.Exec("UPDATE user_schemas SET schema = $1 WHERE user_id = $2", schema, userId)
	if err != nil {
		return fmt.Errorf("failed to update schema in db: %w", err)
	}
	return nil
}

func WriteSchemaToDbConfig(userId uint64, schema string) error {
	db, err := database.Connect(configDbConnectionString)
	if err != nil {
		return err
	}
	defer database.Close(db)
	_, err = db.Exec("INSERT INTO user_schemas (user_id, schema) VALUES ($1, $2)", userId, schema)
	if err != nil {
		return fmt.Errorf("failed to write schema to db: %w", err)
	}
	return nil
}

func WriteNewMapping(userId uint64, shardid uint64, sharduid string, dsn string) error {
	db, err := database.Connect(configDbConnectionString)
	if err != nil {
		return err
	}
	defer database.Close(db)
	_, err = db.Exec("INSERT INTO database_mappings (user_id, shardid, sharduid, dsn) VALUES ($1, $2, $3, $4)", userId, shardid, sharduid, dsn)
	if err != nil {
		return fmt.Errorf("failed to write new mapping to db: %w", err)
	}
	return nil
}

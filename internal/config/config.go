package config

import (
	"database/sql"
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

package dbmapper

import "errors"

type DbMap struct {
	Shardid  uint64
	Sharduid string
	Dsn      string
}

type DbConnectionString string

var DbMapsId = make(map[uint64]DbConnectionString)
var DbMapsUid = make(map[string]DbConnectionString)

func AddDbMapId(shardid uint64, dsn string) error {
	if shardid == 0 {
		return errors.New("shardid cannot be 0")
	}
	DbMapsId[shardid] = DbConnectionString(dsn)
	return nil
}

func AddDbMapUid(sharduid string, dsn string) error {
	if sharduid == "" {
		return errors.New("shardid cannot be empty string")
	}
	DbMapsUid[sharduid] = DbConnectionString(dsn)
	return nil
}

func GetDbConnectionString(shardid uint64, sharduid string) (DbConnectionString, error) {
	if shardid != 0 {
		if dsn, ok := DbMapsId[shardid]; ok {
			return dsn, nil
		}
	}
	if sharduid != "" {
		if dsn, ok := DbMapsUid[sharduid]; ok {
			return dsn, nil
		}
	}
	return "", errors.New("dbMap not found")
}

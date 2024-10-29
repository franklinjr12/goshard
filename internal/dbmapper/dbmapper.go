package dbmapper

import "errors"

type DbMap struct {
	Shardid  uint64
	Sharduid string
	Dsn      DbConnectionString
}

type DbConnectionString string

const UserIdZeroStr = "userId cannot be 0"
const DbMapNotFoundStr = "dbMap not found"

var DbMapsId = make(map[uint64]DbConnectionString)
var DbMapsUid = make(map[string]DbConnectionString)
var DbMapsByUserId = make(map[uint64][]DbMap)

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
	return "", errors.New(DbMapNotFoundStr)
}

func AddDbMapWithUserId(userId uint64, dbMap DbMap) error {
	if userId == 0 {
		return errors.New(UserIdZeroStr)
	}
	DbMapsByUserId[userId] = append(DbMapsByUserId[userId], dbMap)
	return nil
}

func GetDbConnectionStringByUserId(userId uint64, shardId uint64, shardUid string) (DbConnectionString, error) {
	if userId != 0 {
		if dbMaps, ok := DbMapsByUserId[userId]; ok {
			for _, dbMap := range dbMaps {
				if dbMap.Shardid != 0 && dbMap.Shardid == shardId {
					return dbMap.Dsn, nil
				}
				if dbMap.Sharduid != "" && dbMap.Sharduid == shardUid {
					return dbMap.Dsn, nil
				}
			}
		}
	}
	return "", errors.New(DbMapNotFoundStr)
}

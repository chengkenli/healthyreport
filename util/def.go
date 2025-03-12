/**
 * @title: def
 * @Author ChengKen
 * @Date: 10/11/2022 13:35
 * @Version 1.0
 */

package util

import (
	"go.uber.org/zap"
	"time"
)

const ENCKEY = "**************"

var (
	Logger   *zap.Logger
	P        ArvgParms
	Domain   []map[string]string
	Obsole   []string
	MetaLink []map[string]interface{}
)

type ArvgParms struct {
	App        string
	ReplicaNum int
	Bucket     int
	Help       bool
	Thread     int
	Visits     bool
	Database   string
	File       string
}

type ConnectParms struct {
	Host string
	Port int
	User string
	Pass string
	Base string
	Uri  string
}

type SchemaTables []struct {
	TABLE_CATALOG   string    `bson:"TABLE_CATALOG"`
	TABLE_SCHEMA    string    `bson:"TABLE_SCHEMA"`
	TABLE_NAME      string    `bson:"TABLE_NAME"`
	TABLE_TYPE      string    `bson:"TABLE_TYPE"`
	ENGINE          string    `bson:"ENGINE"`
	VERSION         int64     `bson:"VERSION"`
	ROW_FORMAT      string    `bson:"ROW_FORMAT"`
	TABLE_ROWS      int64     `bson:"TABLE_ROWS"`
	AVG_ROW_LENGTH  int64     `bson:"AVG_ROW_LENGTH"`
	DATA_LENGTH     int64     `bson:"DATA_LENGTH"`
	MAX_DATA_LENGTH int64     `bson:"MAX_DATA_LENGTH"`
	INDEX_LENGTH    int64     `bson:"INDEX_LENGTH"`
	DATA_FREE       int64     `bson:"DATA_FREE"`
	AUTO_INCREMENT  int64     `bson:"AUTO_INCREMENT"`
	CREATE_TIME     time.Time `bson:"CREATE_TIME"`
	UPDATE_TIME     time.Time `bson:"UPDATE_TIME"`
	CHECK_TIME      time.Time `bson:"CHECK_TIME"`
	TABLE_COLLATION string    `bson:"TABLE_COLLATION"`
	CHECKSUM        int64     `bson:"CHECKSUM"`
	CREATE_OPTIONS  string    `bson:"CREATE_OPTIONS"`
	TABLE_COMMENT   string    `bson:"TABLE_COMMENT"`
}

package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	utils "github.com/HasanNugroho/starter-golang/pkg/utlis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	RDBMS RDBMS // relational database
	REDIS REDIS // redis database
}

// RDBMS - relational database variables
type RDBMS struct {
	Activate bool
	Env      struct {
		Driver      string
		Host        string
		Port        string
		TimeZone    string
		Synchronize bool
		LogLevel    int
	}
	Access struct {
		DbName string
		User   string
		Pass   string
	}
	Ssl struct {
		Sslmode    string
		MinTLS     string
		RootCA     string
		ServerCert string
		ClientCert string
		ClientKey  string
	}
	Conn struct {
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime time.Duration
	}
	Client *gorm.DB
}

// REDIS - redis database variables
type REDIS struct {
	Activate bool
	Env      struct {
		Host     string
		Port     string
		Password string
		DB       string
	}
	Conn struct {
		PoolSize int
		ConnTTL  int
	}
	Client *redis.Client
}

// loadDatabaseConfig loads all database-related configuration
func loadDatabaseConfig() (databaseConfig DatabaseConfig, err error) {
	var dbConfig DatabaseConfig

	if activateRDBMS, err := strconv.ParseBool(os.Getenv("ACTIVATE_RDBMS")); err == nil && activateRDBMS {
		dbRDBMS, err := loadRDBMSConfig()
		if err != nil {
			return dbConfig, err
		}
		dbConfig.RDBMS = dbRDBMS.RDBMS
		dbConfig.RDBMS.Activate = true
	}

	if activateRedis, err := strconv.ParseBool(os.Getenv("ACTIVATE_REDIS")); err == nil && activateRedis {
		dbRedis, err := loadRedisConfig()
		if err != nil {
			return dbConfig, err
		}
		dbConfig.REDIS = dbRedis.REDIS
		dbConfig.REDIS.Activate = true
	}

	return dbConfig, nil
}

func loadRDBMSConfig() (databaseConfig DatabaseConfig, err error) {
	// set Env
	databaseConfig.RDBMS.Env = struct {
		Driver      string
		Host        string
		Port        string
		TimeZone    string
		Synchronize bool
		LogLevel    int
	}{
		Driver:      strings.ToLower(strings.TrimSpace(os.Getenv("DBDRIVER"))),
		Host:        strings.TrimSpace(os.Getenv("DBHOST")),
		Port:        strings.TrimSpace(os.Getenv("DBPORT")),
		TimeZone:    strings.TrimSpace(os.Getenv("DBTIMEZONE")),
		Synchronize: utils.ToBool(os.Getenv("DBSYNCHRONIZE"), false),
		LogLevel:    utils.ToInt("DBLOGLEVEL", 1),
	}

	// set Env Access
	databaseConfig.RDBMS.Access = struct {
		DbName string
		User   string
		Pass   string
	}{
		DbName: strings.TrimSpace(os.Getenv("DBNAME")),
		User:   strings.TrimSpace(os.Getenv("DBUSER")),
		Pass:   strings.TrimSpace(os.Getenv("DBPASS")),
	}

	// set Env SSL
	databaseConfig.RDBMS.Ssl = struct {
		Sslmode    string
		MinTLS     string
		RootCA     string
		ServerCert string
		ClientCert string
		ClientKey  string
	}{
		Sslmode:    strings.TrimSpace(os.Getenv("DBSSLMODE")),
		MinTLS:     strings.TrimSpace(os.Getenv("DBSSL_TLS_MIN")),
		RootCA:     strings.TrimSpace(os.Getenv("DBSSL_ROOT_CA")),
		ServerCert: strings.TrimSpace(os.Getenv("DBSSL_SERVER_CERT")),
		ClientCert: strings.TrimSpace(os.Getenv("DBSSL_CLIENT_CERT")),
		ClientKey:  strings.TrimSpace(os.Getenv("DBSSL_CLIENT_KEY")),
	}

	// set Env Connection
	databaseConfig.RDBMS.Conn = struct {
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime time.Duration
	}{
		MaxIdleConns:    utils.ToInt("DBMAXIDLECONNS", 10),
		MaxOpenConns:    utils.ToInt("DBMAXOPENCONNS", 100),
		ConnMaxLifetime: utils.ToDuration("DBCONNMAXLIFETIME", 30*time.Minute),
	}

	return
}

// databaseRedis - all REDIS DB variables
func loadRedisConfig() (databaseConfig DatabaseConfig, err error) {
	defaultPoolSize := 10
	defaultConnTTL := 60

	databaseConfig.REDIS.Env.Host = utils.ToString(os.Getenv("REDISHOST"), "localhost")
	databaseConfig.REDIS.Env.Port = utils.ToString(os.Getenv("REDISPORT"), "6379")
	databaseConfig.REDIS.Env.Password = utils.ToString(os.Getenv("REDISPASSWORD"), "")
	databaseConfig.REDIS.Conn.PoolSize = utils.ToInt(os.Getenv("POOLSIZE"), defaultPoolSize)
	databaseConfig.REDIS.Conn.ConnTTL = utils.ToInt(os.Getenv("CONNTTL"), defaultConnTTL)

	return
}

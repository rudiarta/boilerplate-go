package libgorm

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm(
	username string,
	password string,
	host string,
	port string,
	dbname string,
) (*gorm.DB, error) {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // Disable color
		},
	)

	gormConfig := &gorm.Config{
		// enhance performance config
		// PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 dbLogger,
	}

	// username, password, host, port, database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		username,
		password,
		host,
		port,
		dbname)
	dsn += `&loc=Asia%2FJakarta&charset=utf8`

	sql, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error DB Connection: ", err)
		return nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sql.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sql.SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sql.SetConnMaxLifetime(time.Minute)

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sql,
	}), gormConfig)

	if err != nil {
		fmt.Println("Error DB Connection: ", err)
		return nil, err
	}

	return db, nil
}

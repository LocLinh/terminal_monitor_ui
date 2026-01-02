package database

import (
	"fmt"
	"terminal_monitor_ui/config"
	"terminal_monitor_ui/logger"

	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func GetDB(cfg *config.AppConfig) (*gorm.DB, error) {
	// dsn := ".://sa:123456@localhost:1433?database=HasakiDb"
	// dsn := fmt.Sprintf("%s://%s:%s?database=%s",
	// 	cfg.SqlServer.ServerName,
	// 	cfg.SqlServer.User,
	// 	cfg.SqlServer.Password,
	// 	cfg.SqlServer.Database)

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s",
		cfg.SqlServer.User,
		cfg.SqlServer.Password,
		cfg.SqlServer.ServerName,
		cfg.SqlServer.Database)
	fmt.Println("dsn", dsn)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = logger.WriteFile("./logger/error.txt", fmt.Sprintf("[Database] Error: %v\n\n", err))
		zap.S().Errorf("Error connecting to database:%v \n", err)
		return nil, err
	}
	return db, nil

}

package database

import (
	"content-alchemist/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dsnThinkRoot = config.DB_CONNECTION +
	"/think-root?charset=utf8mb4&parseTime=True&loc=Local"
var DBThinkRoot *gorm.DB

type GithubRepositories struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement"`
	URL        string     `gorm:"column:url;type:text;not null;index"`
	Text       string     `gorm:"column:text;type:text;not null"`
	Posted     bool       `gorm:"column:posted;type:tinyint;not null;default:0;index"`
	DateAdded  *time.Time `gorm:"column:date_added;type:datetime;index"`
	DatePosted *time.Time `gorm:"column:date_posted;type:datetime;index"`
}

func (GithubRepositories) TableName() string {
	return "github_repositories"
}

func init() {
	var err error

	gormConfig := &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Error),
	}

	dsnWithoutDB := config.DB_CONNECTION + "/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsnWithoutDB), gormConfig)
	if err != nil {
		log.Printf("Error opening MySQL connection: %v", err)
	}

	db.Exec("CREATE DATABASE IF NOT EXISTS `think-root`")

	DBThinkRoot, err = gorm.Open(mysql.Open(dsnThinkRoot), gormConfig)
	if err != nil {
		log.Printf("Error opening think-root database connection: %v", err)
	}

	sqlDB, err := DBThinkRoot.DB()
	if err != nil {
		log.Printf("Error getting DB instance: %v", err)
	} else {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	}

	log.Println("Successfully connected to think-root database with optimized connection pool")

	if !DBThinkRoot.Migrator().HasTable(&GithubRepositories{}) {
		err = DBThinkRoot.AutoMigrate(&GithubRepositories{})
		if err != nil {
			log.Printf("Error creating github_repositories table: %v", err)
		}
		log.Println("Created github_repositories table")
	}
}

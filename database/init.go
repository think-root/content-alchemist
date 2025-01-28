package database

import (
	"chappie/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsnThinkRoot = config.DB_CONNECTION +
	"/think-root?charset=utf8mb4&parseTime=True&loc=Local"
var DBThinkRoot *gorm.DB

type GithubRepositories struct {
	ID     int    `gorm:"column:id;primaryKey;autoIncrement"`
	URL    string `gorm:"column:url;type:text;not null"`
	Text   string `gorm:"column:text;type:text;not null"`
	Posted bool   `gorm:"column:posted;type:tinyint;not null;default:0"`
}

func (GithubRepositories) TableName() string {
	return "github_repositories"
}

func init() {
	var err error

	dsnWithoutDB := config.DB_CONNECTION + "/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})
	if err != nil {
		log.Printf("Error opening MySQL connection: %v", err)
	}

	db.Exec("CREATE DATABASE IF NOT EXISTS `think-root`")

	DBThinkRoot, err = gorm.Open(mysql.Open(dsnThinkRoot), &gorm.Config{})
	if err != nil {
		log.Printf("Error opening think-root database connection: %v", err)
	}
	log.Println("Successfully connected to think-root database")

	if !DBThinkRoot.Migrator().HasTable(&GithubRepositories{}) {
		err = DBThinkRoot.AutoMigrate(&GithubRepositories{})
		if err != nil {
			log.Printf("Error creating github_repositories table: %v", err)
		}
		log.Println("Created github_repositories table")
	}
}

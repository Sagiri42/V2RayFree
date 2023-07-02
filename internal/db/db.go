package db

import (
	"V2RayFree/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

var (
	dbFilePath = "configs/nodes.db"
	DB         *gorm.DB
)

func connectDB() {
	var err error
	if DB, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{}); err != nil {
		log.Fatalf("打开nodes.db失败: %v", err)
	}
	DB.AutoMigrate(&model.Node{})
}

func init() {
	connectDB()
}

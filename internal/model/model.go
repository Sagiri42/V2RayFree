package model

import (
	"gorm.io/gorm"
	"time"
)

type Node struct {
	ID      uint   `gorm:"column:id;primaryKey;autoIncrement"`
	Link    string `gorm:"column:link;unique;not null"`
	TCPTest uint16 `gorm:"column:tcpTest"`
	//URLTest   uint8          `gorm:"column:urlTest"`
	Failures  uint16         `gorm:"column:failures"`
	CreatedAt time.Time      `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updatedAt;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deletedAt;index"`
}

func (n *Node) TableName() string {
	return "node"
}

func NewNodes(nodes []string) []Node {
	ns := make([]Node, len(nodes))
	for i := range nodes {
		ns[i].Link = nodes[i]
	}
	return ns
}

func InsetNode(db *gorm.DB, link string, tcpTest uint16) {
	db.Model(Node{}).Create(&Node{Link: link, TCPTest: tcpTest})
}

func DeleteNode(db *gorm.DB, id uint) {
	db.Model(Node{}).Delete(&Node{ID: id})
}

func UpdateTCPLink(db *gorm.DB, id uint, tcpTest uint16) {
	db.Model(Node{}).Updates(&Node{ID: id, TCPTest: tcpTest})
}

//func UpdateUrlLink(db *gorm.DB, id uint, urlTest uint8) {
//	db.Model(Node{}).Updates(&Node{ID: id, URLTest: urlTest})
//}

func UpdateFailures(db *gorm.DB, id uint) {
	var n Node
	db.Model(Node{}).Find(&n)
	n.Failures++
	db.Model(Node{}).Updates(&n)
}

func QueryAllNode(db *gorm.DB) []Node {
	var ns []Node
	db.Model(Node{}).Find(&ns)
	return ns
}

package main

import (
	"V2RayFree/internal/config"
	"V2RayFree/internal/db"
	"V2RayFree/internal/model"
	"V2RayFree/internal/server"
	"V2RayFree/internal/service"
	"V2RayFree/pkg/utils"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// TODO 公共节点定时获取
	utils.Ticker(config.Config.UpdateAt, func() {
		_ = config.Config.ReadConfig()

		var nodes []model.Node

		nodes = append(nodes, model.NewNodes(config.Config.Nodes)...)

		for i := range config.Config.SubscribeLink {
			nodes = append(nodes, server.NewFreeConverter(server.SubscribeLink(config.Config.SubscribeLink[i]))...)
		}

		for i := range config.Config.SubscribeBase64 {
			nodes = append(nodes, server.NewFreeConverter(server.SubscribeBase64(config.Config.SubscribeBase64[i]))...)
		}

		// 过滤
		nodes = server.NodeType.Filter(nodes)
		nodes = server.NodeTcpTest.Filter(nodes)

		for i := range nodes {
			if nodes[i].TCPTest > 0 && nodes[i].TCPTest < 1000 {
				db.DB.Model(model.Node{}).Create(&nodes[i])
			}
		}

	})

	// TODO 数据库节点定时检查
	utils.Ticker(config.Config.UpdateAt/2, func() {
		nodes := model.QueryAllNode(db.DB)
		server.NodeTcpTest.Filter(nodes)
		server.NodeFailures.Filter(nodes)
		var ns []model.Node
		for i := range nodes {
			if nodes[i].Failures <= 5 {
				ns = append(ns, nodes[i])
			} else {
				log.Printf("%v节点已失效", nodes[i].Link)
				db.DB.Model(model.Node{}).Delete(&nodes[i])
			}
		}
		db.DB.Model(model.Node{}).Updates(&ns)
	})

	// TODO 订阅API接口
	http.HandleFunc("/api/subscribe", service.SubscribeHandle)
	log.Printf("服务已启动: http://%s:%d/api/subscribe", config.Config.Host, config.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port), nil); err != nil {
		return
	}
}

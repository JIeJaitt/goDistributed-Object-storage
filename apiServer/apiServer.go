package main

import (
	"goDistributed-Object-storage/apiServer/heartbeat"
	"goDistributed-Object-storage/apiServer/locate"
	"goDistributed-Object-storage/apiServer/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	// 接收来自数据服务的心跳信息
	go heartbeat.ListenHeartbeat()
	// 将对象请求转发给数据服务
	http.HandleFunc("/objects/", objects.Handler)
	// 发送定位文件名的消息以及处理数据服务locate反馈回来的消息
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
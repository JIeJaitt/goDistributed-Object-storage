package main

import (
	"goDistributed-Object-storage/dataServer/heartbeat"
	"goDistributed-Object-storage/dataServer/locate"
	"goDistributed-Object-storage/dataServer/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	// 每隔一段时间就会发送一个心跳信息给到这个apiServer的exchange
	// 所有连接到apiServer exchange的那些接口服务节点就会收到当前数据服务节点的心跳信息
	// 从而就会知道当前这个数据节点是存活状态
	// 当接口服务接收到来自客户端的请求以后，它就有可能会选中当前这个数据服务节点进行具体的逻辑操作
	go heartbeat.StartHeartbeat()
	// 开启一个locate协程监听来自接口服务的locate请求
	// 接收到每一个请求以后就会在本地磁盘上查找有没有这个文件
	// 如果有就返回当前节点的本机地址
	// 我这个dataServer有你要的文件，你可以直接来我这里下载
	// 如果没有找到任何结果就什么都不发，假如所有数据节点都没发送任何结果
	// 等一秒以后那个接口节点就超时了
	go locate.StartLocate()
	// objects Handler 用于处理对象的 API
	// 用来给接口服务提供一个存取对象的服务
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
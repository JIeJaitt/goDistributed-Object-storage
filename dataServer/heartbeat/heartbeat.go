package heartbeat

import (
	"goDistributed-Object-storage/src/lib/rabbitmq"
	"os"
	"time"
)

func StartHeartbeat() {
	// 创建一个新的rabbitmq连接
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	//
	for {
		// 不停的向apiServers这个exchange发送心跳信息
		// 不停的把自己的本机地址LISTEN_ADDRESS 发送到apiServers里面去
		q.Publish("apiServers", os.Getenv("LISTEN_ADDRESS"))
		// 每五秒钟发送一次心跳信息
		time.Sleep(5 * time.Second)
	}
}
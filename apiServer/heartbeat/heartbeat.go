package heartbeat

import (
	"goDistributed-Object-storage/src/lib/rabbitmq"
	"os"
	"strconv"
	"sync"
	"time"
)

// 用来缓存所有的数据服务节点
// key：数据服务节点的本机地址
// value：数据服务节点最后一次发送心跳信息的时间
var dataServers = make(map[string]time.Time)

// 用来保护dataServers这个map的互斥锁
// 用来对dataServers这个map做访问的时候的一个并发的保护
var mutex sync.Mutex

// 接收来自数据服务的心跳信息
// apiServer：go heartbeat.ListenHeartbeat()
func ListenHeartbeat() {
	// 创建rabbitmq消息队列结构体，控制跟rabbitmq相关的API
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	// 把 q 绑定apiServers上面
	// 绑定之后所有发给apiServers的消息都会转发给我一份
	q.Bind("apiServers")
	// 生成一个 channel
	c := q.Consume()
	// 并行的检查哪些dataServer已经很久没收到心跳信息了
	// 移除超时的dataServer
	go removeExpiredDataServer()
	// 从消息队列里面不断的接收消息
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

// 获取所有的数据服务节点
func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}



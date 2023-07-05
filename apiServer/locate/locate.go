package locate

import (
	"goDistributed-Object-storage/src/lib/rabbitmq"
	"os"
	"strconv"
	"time"
)

// 从消息队列里面接收消息
func Locate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

// 判断文件是否存在
func Exist(name string) bool {
	return Locate(name) != ""
}
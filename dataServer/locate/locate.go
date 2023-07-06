package locate

import (
	"goDistributed-Object-storage/src/lib/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool {
	// 利用 err 是否存在来判断name存在的状态
	// 如果文件状态正常返回，则err的值为nil
	_, err := os.Stat(name)
	// name 存在则err为nil 则os.IsNotExist为false ，取反返回true
	return !os.IsNotExist(err)
}

func StartLocate() {

	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	// 自己作为一个dataServer会绑定到dataServers这个exchange
	// 对于这个exchange所收到的所有消息都会转发给我一份
	q.Bind("dataServers")
	// 从dataServers这个exchange里面接收消息
	c := q.Consume()
	// 消息队列收到了以后
	for msg := range c {
		// 不断的从消息队列里面取出消息
		// 每次取到的一个消息就意味着接受到了一个对象的定位请求
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		// 调用Locate函数来判断这个对象是否存在
		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object) {
			// 如果Locate返回true则证明文件存在，就把当前节点的本机地址发送给接口服务
			// 如果locate返回false则证明文件不存在，就什么都不做，接口服务就什么也收不到
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}
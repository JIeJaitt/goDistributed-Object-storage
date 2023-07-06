package objects

import (
	"goDistributed-Object-storage/apiServer/heartbeat"
	"fmt"
	"goDistributed-Object-storage/src/lib/objectstream"
)

// 用来给接口服务提供一个存取对象的服务
func putStream(object string) (*objectstream.PutStream, error) {
	// 获得一个随机数据服务节点
	server := heartbeat.ChooseRandomDataServer()
	// 没有可用的数据服务节点
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}
	// 返回一个objectstream.NewPutStream的指针
	return objectstream.NewPutStream(server, object), nil
}
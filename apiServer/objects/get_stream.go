package objects

import (
	"fmt"
	"goDistributed-Object-storage/apiServer/locate"
	"goDistributed-Object-storage/src/lib/objectstream"
	"io"
)

func getStream(object string) (io.Reader, error) {
	// 调用locate.Locate函数来获得这个object具体存储在dataServer的什么地方
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}

package objects

import (
	"fmt"
	"goDistributed-Object-storage/apiServer/heartbeat"
	"goDistributed-Object-storage/src/lib/objectstream"
)

func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any dataServer")
	}

	return objectstream.NewTempPutStream(server, hash, size)
}

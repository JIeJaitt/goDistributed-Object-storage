package temp

import (
	"goDistributed-Object-storage/dataServer/locate"
	"os"
)

func commitTempObject(datFile string, tempinfo *tempInfo) {
	os.Rename(datFile, os.Getenv("STORAGE_ROOT")+"/objects/"+tempinfo.Name)
	locate.Add(tempinfo.Name)
}

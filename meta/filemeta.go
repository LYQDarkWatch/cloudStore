package meta

import (
	mydb "filestore-server/db"
)

//文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

//init方法会在程序启动时执行一遍
func init() {
	fileMetas = make(map[string]FileMeta)
}

//新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnfileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

//通过sha1获取文件元信息对象
func GetFileMeta(filesha1 string) FileMeta {
	return fileMetas[filesha1]
}

// GetLastFileMetas : 获取批量的文件元信息列表
// func GetLastFileMetas(count int) []FileMeta {
// 	fMetaArray := make([]FileMeta, len(fileMetas))
// 	for _, v := range fileMetas {
// 		fMetaArray = append(fMetaArray, v)
// 	}

// 	sort.Sort(ByUploadTime(fMetaArray))
// 	return fMetaArray[0:count]
// }

// RemoveFileMeta : 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

//从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}

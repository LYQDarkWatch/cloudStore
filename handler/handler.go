package handler

import (
	"encoding/json"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "Internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Println("faild to get data, err : %s", err.Error())
			return
		}
		defer file.Close()
		filemeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/Users/lyq/filesave/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(filemeta.Location)
		if err != nil {
			fmt.Println("faild to create file. err: %s", err.Error())
			return
		}
		defer file.Close()

		filemeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Println("faild to save data into file, err : %s", err.Error())
			return
		}
		newFile.Seek(0, 0)
		filemeta.FileSha1 = util.FileSha1(newFile)
		//meta.UpdateFileMeta(filemeta)
		_ = meta.UpdateFileMetaDB(filemeta)
		//更新用户文件记录表
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

//上传已完成
func FileuploadSuccHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload finish")
}

//获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	//fmeta := meta.GetFileMeta(filehash)
	fMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
func DonwloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	date, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fm.FileName+"")
	w.Write(date)
}

//更新文件元信息(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析请求的参数列表
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//删除文件
func FileDelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}

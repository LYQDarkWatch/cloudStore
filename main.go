package main

import (
	"filestore-server/handler"
	"log"
	"net/http"
)

func StartServer() {
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.FileuploadSuccHandler)
	http.HandleFunc("/file/download", handler.DonwloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDelHandler)
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/login", handler.Login)
	http.HandleFunc("/user/home", handler.Home)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))
	log.Println("Now listening...")
	http.ListenAndServe(":8085", nil)
}

func main() {
	StartServer()
}

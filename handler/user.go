package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwdsalt = "#890"
)

//登陆页面跳转
func Login(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./static/view/signin.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
func Home(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile("./static/view/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	//r.ParseForm()方法获取用户提交的表单参数
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}
	//使用util里的Sha1哈希的加密方法
	enc_passwd := util.Sha1([]byte(passwd + pwdsalt))
	suc := dblayer.Usersigup(username, enc_passwd)
	if suc {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("faild"))
	}
}

//登陆接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password + pwdsalt))

	//校验用户提交的用户名和密码
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		fmt.Println("11")
		w.Write([]byte("FAILED"))
		return
	}
	//登录成功后获得访问凭证(token)
	token := GenToken(username)
	err := dblayer.UpdateToken(username, token)
	if !err {
		fmt.Println("22ç")
		w.Write([]byte("FAILED"))
		return
	}
	//登录成功后重定向到首页
	//w.Write([]byte("success"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/user/home",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

//生成token
func GenToken(username string) string {
	// md5(username+timestamp+token_salt) + timestamp[:8]  40位字符
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

//查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	//验证token是否有效
	IsTokenValid := IsTokenValid(token)
	if !IsTokenValid {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//组装并响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

//token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	//判断token时效性，是否过期
	//从数据库查询是否有此token
	//对比两个token是否一致
	return true
}

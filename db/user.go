package db

import (
	mybd "filestore-server/db/mysql"
	"fmt"
)

//通过用户名及密码完成user表的注册操作
func Usersigup(username string, passwd string) bool {
	//Prepare返回两个参数  stmt和err
	stmt, err := mybd.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`,`user_pwd`)value (?,?)")
	if err != nil {
		fmt.Println("failed to insert:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("faild to insert: " + err.Error())
		return false
	}
	if RowsAffected, err := ret.RowsAffected(); nil == err && RowsAffected > 0 {
		return true
	}
	return false
}

//判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	stmt, err := mybd.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("uername not found: " + username)
		return false
	}
	pRows := mybd.ParseRows(rows)
	//go的断言  []byte
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

//刷新用户登陆的token
func UpdateToken(username string, token string) bool {
	stmt, err := mybd.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`,`user_token`) values (?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mybd.DBConn().Prepare("select user_name,signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()
	//执行查询操作,scan方法将读出的数据转换到给一个外部变量
	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

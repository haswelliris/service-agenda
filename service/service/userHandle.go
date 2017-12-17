package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/haswelliris/service-agenda/service/model"
	"github.com/unrolled/render"
)

func userLoginHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		curUser, _ := model.GetCurUser()
		fmt.Println("检测是否有当前用户：" + curUser)
		if curUser != "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户已登录"))
			fmt.Println("警告：检测到用户已经登录，禁止重复登录 ...")
		} else {
			name := req.FormValue("username")
			if name != "" {
				errCode, theUser, err := model.GetUser(name)
				if errCode == -1 || err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("用户不存在"))
					return
				}
				if req.FormValue("password") != theUser.Password {
					fmt.Println("登录密码错误")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("登录密码错误"))
				} else {
					model.UpdateCurUser(name)
					fmt.Println("登录成功")
					formatter.JSON(w, http.StatusOK, name)
				}
			} else {
				fmt.Println("错误：登录用户名不能为空")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("登录用户名不能为"))
			}

		}
	}
}

func userLogoutHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name, err := model.GetCurUser()
		if err != nil || name == "" {
			fmt.Println("错误，没有已经登录的用户")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("没有已经登录的用户"))
		} else {
			model.UpdateCurUser("")
			fmt.Println("用户 " + name + " 退出登录")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func listAllUsersHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name, err := model.GetCurUser()
		if err != nil || name == "" {
			fmt.Println("用户未登录")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户未登录"))
			return
		}
		fmt.Println("当前登录用户是：" + name)
		users := model.GetUsers()
		fmt.Println("列出所有用户")
		for _, user := range users {
			fmt.Println(user)
		}
		formatter.JSON(w, http.StatusOK, users)
	}
}

func userRegisterHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name, err := model.GetCurUser()
		if err != nil || name != "" {
			fmt.Println("用户已经注册过了")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户已注册"))
			return
		}
		defer req.Body.Close()
		data, _ := ioutil.ReadAll(req.Body)
		fmt.Println("收到提交数据" + string(data))
		var theUser model.User
		json.Unmarshal([]byte(data), &theUser)
		fmt.Println(theUser)
		if err != nil {
			fmt.Println("解析提交数据失败")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("非法提交"))
			return
		} else {
			if theUser.UserName != "" && theUser.Password != "" && theUser.Email != "" && theUser.Phone != "" {
				_, _, err := model.GetUser(theUser.UserName)
				if err == nil {
					fmt.Println("用户名重复")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("用户名重复"))
					return
				}
				err = model.AddUser(theUser.UserName, theUser.Password, theUser.Email, theUser.Phone)
				if err != nil {
					fmt.Println(err)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("注册失败"))
					return
				} else {
					model.UpdateCurUser(theUser.UserName)
					fmt.Println("注册成功")
					formatter.JSON(w, http.StatusOK, theUser.UserName)
				}
			} else {
				fmt.Println("非法提交")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("非法提交"))
				return
			}
		}
	}
}

package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haswelliris/service-agenda/service/model"
	"github.com/unrolled/render"
)

func userLoginHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GET /user/login :")
		req.ParseForm()
		if model.getCurUser() {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户已登录"))
			fmt.Println("警告：检测到用户已经登录，禁止重复登录 ...")
		} else {
			var name string := req.FormValue("username")
			if name {
				errCode,theUser,err = model.GetUser(name)
				if errCode == -1 || err {
					fmt.Println(err);
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("用户不存在"))
					return 					
				}
				if req.FormValue("password") != theUser.password {
					fmt.Println("登录密码错误")
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("登录密码错误"))
				} else {
					model.UpdateCurUser(name)
					fmt.Printf("登录成功")
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

func LogoutHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GET /user/logout")
		name,err := model.getCurUser()
		if err || name != "" {
			fmt.Printf("错误，没有已经登录的用户")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("没有已经登录的用户"))			
		} else {
			model.UpdateCurUser("")
			fmt.Printf("退出登录")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func ListAllUserHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GET /users")
		if model.getCurUser() == "" {
			fmt.Println("用户未登录")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户未登录"))	
			return
		}
		users:=model.getUsers()
		fmt.Println(users)
		formatter.JSON(w, http.StatusOK, users)
	}
}

func userRegisterHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GET /user/register")
		if model.getCurUser() != "" {
			fmt.Println("用户已经注册过了")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户已注册"))	
			return
		}
		decoder := json.NewDecoder(req.Body)
		var theUser model.User
		err := decoder.Decode(&theUser)
		if err != nil {
			fmt.Printf(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("非法提交"))	
			return
		} else {
			if theUser.UserName != "" && theUser.Password != "" && theUser.Email != "" && theUser.Phone != "" {
				err = model.AddUser()
				if err {
					fmt.Printf(err)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("非法提交"))	
					return
				} else {
					model.UpdateCurUser(theUser.UserName)
					fmt.Printf("登录成功")
					formatter.JSON(w, http.StatusOK, theUser.UserName)
				}
			}else {
				fmt.Printf(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("非法提交"))	
				return
			}
		}
	}
}

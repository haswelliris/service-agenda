package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type User struct {
	UserName string `xorm:"pk"`
	Password string
	Email    string
	Phone    string
}

type Users []*User

var users Users

func UpdateCurUser(curUser string) error {
	// write current login user into file
	fout, err := os.Create("data/curUser")
	if err != nil {
		return errors.New("update current user fail: \n->" + err.Error())
	}
	defer fout.Close()
	if err != nil {
		return errors.New("update current user fail: \n->" + err.Error())
	}
	fmt.Fprintf(fout, "%s", curUser)
	return nil
}

func GetCurUser() (string, error) {
	// read current login user from file and return it
	fin, err := os.Open("data/curUser")
	if err != nil {
		return "", errors.New("get current user fail: \n->" + err.Error())
	}
	defer fin.Close()
	curUser, err := ioutil.ReadAll(fin)
	if err != nil {
		return "", errors.New("get current user fail: \n->" + err.Error())
	}
	return string(curUser), nil
}

func AddUser(userName string, password string, email string, phone string) error {
	user := &User{userName, password, email, phone}
	users = append(users, user)
	// insert new user infomation into the database
	_, err := Engine.Insert(user)
	if err != nil {
		return errors.New("write to database fail: \n->" + err.Error())
	}
	return nil
}

func DeleteUser(userName string) error {
	index, _, err := GetUser(userName)
	if err != nil {
		return errors.New("delete user fail: " + err.Error())
	} else {
		// delete user with given index
		users = append(users[:index], users[index+1:]...)
		// delete user from the database
		user := new(User)
		Engine.ID(userName).Delete(user)
	}
	return nil
}

func GetUser(userName string) (int, *User, error) {
	for index, user := range users {
		if user.UserName == userName {
			return index, user, nil
		}
	}
	return -1, nil, errors.New("can not get user by " + userName)
}

func GetUsers() Users {
	return users
}

func userReadFromDB() {
	Engine.Find(&users)
}

func init() {
	userReadFromDB()
}

package model

import (
	"testing"
)

func TestAddUser(t *testing.T) {
	addCases := []*User{
		{"alice", "123456", "alice@qq.com", "12345678912"},
		{"bob", "654321", "bob@qq.com", "14785236987"},
	}
	for _, testUser := range addCases {
		AddUser(testUser.UserName, testUser.Password, testUser.Email, testUser.Phone)
		_, addedUser, err := GetUser(testUser.UserName)
		if err != nil || *addedUser != *testUser {
			t.Error(err)
		}
	}
}

func TestDeleteUser(t *testing.T) {
	deleteCases := []struct {
		in, want string
	}{
		{"alice", "success"},
		{"james", "fail"},
	}
	AddUser("alice", "123456", "alice@qq.com", "12345678912")
	for _, c := range deleteCases {
		err := DeleteUser(c.in)
		if err != nil && c.want != "fail" {
			t.Error(err)
		} else if index, _, _ := GetUser(c.in); index != -1 {
			t.Error("test DeleteUser: fail; can not remove user", c.in, "successfully.")
		}
	}
}

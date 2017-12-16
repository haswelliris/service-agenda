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
	if GetUsers() != nil {
		users = nil
	}
	deleteCases := []struct {
		in, want string
	}{
		{"bob", "success"},
		{"james", "fail"},
	}
	AddUser("bob", "123456", "bob@qq.com", "12345678912")
	for _, c := range deleteCases {
		err := DeleteUser(c.in)
		if err != nil && c.want != "fail" {
			t.Error(err)
		} else if index, _, _ := GetUser(c.in); index != -1 {
			t.Error("test DeleteUser: fail; can not remove user", c.in, "successfully.")
		}
	}
}

func TestAddMeeting(t *testing.T) {
	addCases := []*Meeting{
		{"meeting1", "alice", "2017-11-11/11:11", "2017-11-11/11:22", []string{"james", "bob"}},
		{"meeting2", "bob", "2017-11-11/12:12", "2017-12-12/11:22", []string{"james", "alice"}},
	}
	for _, testMeeting := range addCases {
		AddMeeting(testMeeting.Title, testMeeting.Sponsor, testMeeting.Participators, testMeeting.StartTime, testMeeting.EndTime)
		_, addedMeeting, err := GetMeeting(testMeeting.Title)
		if err != nil || addedMeeting.Title != testMeeting.Title {
			t.Error(err)
		}
	}
}

func TestDeleteMeeting(t *testing.T) {
	if GetMeetings() != nil {
		meetings = nil
	}
	deleteCases := []struct {
		in, want string
	}{
		{"meeting1", "success"},
		{"meeting2", "fail"},
	}
	AddMeeting("meeting1", "alice", []string{"james", "bob"}, "2017-11-11/11:11", "2017-11-11/11:22")
	for _, c := range deleteCases {
		err := DeleteMeeting(c.in)
		if err != nil && c.want != "fail" {
			t.Error(err)
		} else if index, _, _ := GetMeeting(c.in); index != -1 {
			t.Error("test DeleteUser: fail; can not remove user", c.in, "successfully.")
		}
	}
}

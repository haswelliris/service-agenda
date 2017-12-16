package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/haswelliris/service-agenda/service/model"

	"github.com/unrolled/render"
)

func ListAllMeetingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GET /meetings")
		if model.getCurUser() == "" {
			fmt.Println("用户未登录")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户未登录"))
			return
		}
		meetings := model.GetMeetings()
		fmt.Println(meetings)
		formatter.JSON(w, http.StatusOK, ,meetings)

	}
}

func CreateMeetingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("GET /meeting")
		fmt.Println("GET /meetings")
		if model.getCurUser() == "" {
			fmt.Println("用户未登录")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户未登录"))
			return
		}
		decoder := json.NewDecoder(req.Body)
		var meeting model.Meeting
		err := decoder.Decode(&meeting)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("非法提交"))
			return
		}
		fmt.Printf(meeting)
		if meeting.StartTime == "" || meeting.EndTime == "" || meeting.title == "" || len(meeting.Participators) <= 0 {
			fmt.Println("非法提交")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("非法提交"))
			return
		}
		meeting.Sponsor = model.getCuUser()

		meetings := model.GetMeetings()

		for _, meeting := range meetings {
			if meeting.Title == title {
				fmt.Println("标题重复")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("标题重复"))
				return
			}
		}
		if !IsTimeValid(meeting.StartTime) || !IsTimeValid(meeting.EndTime) {
			fmt.Println("会议时间格式错误")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("会议时间格式错误"))
			return
		}
		if util.Time2str(meeting.StartTime) > util.Time2str(meeting.EndTime) ||
			util.Time2str(meeting.StartTime) == util.Time2str(meeting.EndTime) {
			fmt.Println("错误：会议结束时间早于开始时间")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("会议结束时间早于开始时间"))
			return
		}
		for _, user := range meeting.Participators {
			if user == model.getCurUser() {
				fmt.Println("错误：会议发起者不能在其它参与者中")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("会议发起者不能在其它参与者中"))
				return
			}
			index, _, err := model.getUser(user)
			if err {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("参与者中有未注册用户"))
				return
			}
			if !IsUserHaveTime(user, startTime, endTime) {
				fmt.Println(user + " 时间冲突")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("会议时间冲突"))
				return
			}
		}
		model.AddMeeting(meeting.Title, meeting.Sponsor, meeting.Participators, meeting.StartTime, meeting.EndTime)
		fmt.Println("创建会议成功")
		_, createdMeeting, _ := model.GetMeeting(title)
		formatter.JSON(w, http.StatusOK, createdMeeting)
	}
}

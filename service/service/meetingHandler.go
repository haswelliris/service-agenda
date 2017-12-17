package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/haswelliris/service-agenda/service/model"

	"github.com/unrolled/render"
)

func listAllMeetingsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name, err := model.GetCurUser()
		if err != nil || name == "" {
			fmt.Println("用户未登录")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户未登录"))
			return
		}
		fmt.Println("当前登录用户是：" + name)
		meetings := model.GetMeetings()
		fmt.Println("列出所有会议")
		fmt.Println("title  sponsor  startTime  endTime  participators")
		for _, meeting := range meetings {
			printMeeting(*meeting)
		}
		formatter.JSON(w, http.StatusOK, meetings)

	}
}

func createNewMeetingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		name, err := model.GetCurUser()
		if err != nil || name == "" {
			fmt.Println("用户未登录")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("用户未登录"))
			return
		}
		fmt.Println("当前登录用户是：" + name)
		defer req.Body.Close()
		data, _ := ioutil.ReadAll(req.Body)
		fmt.Println("收到提交数据" + string(data))
		var meeting model.Meeting
		err = json.Unmarshal([]byte(data), &meeting)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("非法提交"))
			return
		}
		if meeting.StartTime == "" || meeting.EndTime == "" || meeting.Title == "" || len(meeting.Participators) <= 0 {
			fmt.Println("非法提交")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("非法提交"))
			return
		}
		fmt.Println("收到会议配置")
		fmt.Println("title  sponsor  startTime  endTime  participators")
		printMeeting(meeting)
		meeting.Sponsor, err = model.GetCurUser()

		meetings := model.GetMeetings()

		for _, m := range meetings {
			if meeting.Title == m.Title {
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
		if Time2str(meeting.StartTime) > Time2str(meeting.EndTime) ||
			Time2str(meeting.StartTime) == Time2str(meeting.EndTime) {
			fmt.Println("错误：会议结束时间早于开始时间")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("会议结束时间早于开始时间"))
			return
		}
		for _, user := range meeting.Participators {
			if user == meeting.Sponsor {
				fmt.Println("错误：会议发起者不能在其它参与者中")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("会议发起者不能在其它参与者中"))
				return
			}
			_, _, err := model.GetUser(user)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("参与者中有未注册用户"))
				return
			}
			if !IsUserHaveTime(user, meeting.StartTime, meeting.EndTime) {
				fmt.Println(user + " 时间冲突")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("会议时间冲突"))
				return
			}
		}
		model.AddMeeting(meeting.Title, meeting.Sponsor, meeting.Participators, meeting.StartTime, meeting.EndTime)
		fmt.Println("创建会议成功")
		_, createdMeeting, _ := model.GetMeeting(meeting.Title)
		formatter.JSON(w, http.StatusOK, createdMeeting)
	}
}

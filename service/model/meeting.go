package model

import (
	"errors"
)

type Meeting struct {
	Title         string `xorm:"pk"`
	Sponsor       string
	StartTime     string
	EndTime       string
	Participators []string
}

type Meetings []*Meeting

var meetings Meetings

func AddMeeting(title string, sponsor string, participators []string, startTime string, endTime string) error {
	meeting := &Meeting{
		Title:         title,
		Sponsor:       sponsor,
		StartTime:     startTime,
		EndTime:       endTime,
		Participators: participators}
	meetings = append(meetings, meeting)
	// insert new meeting infomation into the database
	_, err := Engine.Insert(meeting)
	if err != nil {
		return errors.New("write to database fail: \n->" + err.Error())
	}
	return nil
}

func DeleteMeeting(title string) error {
	index, _, err := GetMeeting(title)
	if err != nil {
		return errors.New("delete meeting fail: " + err.Error())
	} else {
		meetings = append(meetings[:index], meetings[index+1:]...)
		// delete meeting from the database
		meeting := new(Meeting)
		Engine.ID(title).Delete(meeting)
	}
	return nil
}

func GetMeeting(title string) (int, *Meeting, error) {
	for index, meeting := range meetings {
		if meeting.Title == title {
			return index, meeting, nil
		}
	}
	return -1, nil, errors.New("can not get meeting by " + title)
}

func GetMeetings() Meetings {
	return meetings
}

// func UpdateParticipators(title string, participators []string) error {
// 	index, _, err := GetMeeting(title)
// 	if err != nil {
// 		return errors.New("update participators fail: " + err.Error())
// 	} else {
// 		meetings[index].Participators = participators
// 	}
// 	meeting := new(Meeting)
// 	meeting.Participators = participators
// 	Engine.ID(title).Update(meeting)
// 	return nil
// }

func meetingReadFromDB() {
	Engine.Find(&meetings)
}

func init() {
	meetingReadFromDB()
}

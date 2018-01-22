package utils

import (
	"time"

	"github.com/astaxie/beego"
)

type Task struct {
	TaskID       int64  `json:"task_id"`
	UserID       string `json:"-"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
	IsRun        bool   `json:"task_status"`
	TicketCode   string `json:"ticket_code"`
	UserKyfw     *Kyfw  `json:"-"`
	SecretKey    string `json:"-"`
	TrainNo      string `json:"-"`
	TrainCode    string `json:"train_code"`
	TrainDate    string `json:"date_time"`
	ForMatDate   string `json:"-"`
	StartStation string `json:"start_name"`
	EndStation   string `json:"end_name"`
	StartCode    string `json:"-"`
	EndCode      string `json:"-"`
	TicketStr    string `json:"ticket_str"`
	PassengerStr string `json:"passenger_str"`
}

type TaskList struct {
	UserTask map[string]map[int64]interface{}
}

func InitTaskList() *TaskList {
	return &TaskList{
		UserTask: make(map[string]map[int64]interface{}),
	}
}

func (tasks *TaskList) Set(key string, task *Task) {
	if tasks.UserTask[key] == nil {
		tasks.UserTask[key] = make(map[int64]interface{})
	}
	task.TaskID = time.Now().Unix()
	task.UserID = key
	tasks.UserTask[key][task.TaskID] = task
	task.Start()
}

func (tasks *TaskList) Get(key string) map[int64]interface{} {
	v, ok := tasks.UserTask[key]
	if !ok {
		return nil
	}
	return v
}

func (tasks *TaskList) Del(key string) {
	delete(tasks.UserTask, key)
}

func (tasks *TaskList) CreateTask(kyfw *Kyfw, SecretKey, TrainNo, TrainCode, TrainDate, ForMatDate, StartStation, EndStation, StartCode, EndCode, TicketStr, PassengerStr string) *Task {
	return &Task{
		UserKyfw:     kyfw,
		SecretKey:    SecretKey,
		TrainNo:      TrainNo,
		TrainCode:    TrainCode,
		TrainDate:    TrainDate,
		ForMatDate:   ForMatDate,
		StartStation: StartStation,
		EndStation:   EndStation,
		StartCode:    StartCode,
		EndCode:      EndCode,
		TicketStr:    TicketStr,
		PassengerStr: PassengerStr,
	}
}

func (tasks *TaskList) DeleteTask(key string, taskID int64) {
	delete(tasks.UserTask[key], taskID)
}

func (task *Task) Start() {
	//任务开始
	task.IsRun = true
	task.StartTime = time.Now().Unix()
	go func() {
		for task.TicketCode == "" && task.IsRun {
			time.Sleep(10 * time.Second)
			err := task.UserKyfw.PlaceAnOrder(task.SecretKey, task.TrainNo, task.TrainCode, task.StartStation, task.StartCode, task.EndStation, task.EndCode, task.TrainDate, task.ForMatDate, task.TicketStr, task.PassengerStr)
			//抢票失败,写日志
			if err != nil {
				beego.Debug("Task:", err.Error())
				continue
			}
			beego.Debug("Task: success")
			//抢票成功
			task.TicketCode = task.UserKyfw.OrderTicketCode
			task.IsRun = false
			task.EndTime = time.Now().Unix()
		}
	}()
}

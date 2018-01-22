package controllers

import "strconv"

// PassengerController Operations about object
type TaskController struct {
	BaseController
}

// @Title Post
// @Description 添加任务
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router / [post]
func (t *TaskController) Post() {

	secretStr := t.GetString("secret_key")
	trainNo := t.GetString("train_no")
	trainCode := t.GetString("train_code")
	trainDate := t.GetString("train_date")
	formatDate := t.GetString("format_date")
	startStation := t.GetString("start_station")
	endStation := t.GetString("end_station")
	startCode := t.GetString("start_code")
	endCode := t.GetString("end_code")
	ticketStr := t.GetString("ticket_str")
	passengerStr := t.GetString("passenger_str")
	task := tasks.CreateTask(t.Kyfw, secretStr, trainNo, trainCode, startStation, startCode, endStation, endCode, trainDate, formatDate, ticketStr, passengerStr)
	tasks.Set(t.UserID, task)
	t.Success().SetMsg("任务添加成功").Send()
}

// @Title Get
// @Description 获取任务
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router / [get]
func (t *TaskController) Get() {
	taskMaps := tasks.Get(t.UserID)
	t.Success().SetData(taskMaps).Send()
}

// @Title Get
// @Description 获取任务日志
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /log/:id:int [get]
func (t *TaskController) Log() {
	taskID, _ := strconv.ParseInt(t.Ctx.Input.Param(":id"), 10, 64)
	t.Success().SetData(taskID).Send()
}

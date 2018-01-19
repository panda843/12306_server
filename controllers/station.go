package controllers


// StationController Operations about object
type StationController struct {
	BaseController
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router / [get]
func (s *StationController) Get() {
	data,err := s.Kyfw.GetStations()
	if err != nil {
		s.Fail().SetMsg(err.Error()).Send()
		return
	}
	//缓存12个小时
	//s.Ctx.Output.Header("Cache-Control:", "public,max-age=43200")
	s.Success().SetData(data).Send()
}

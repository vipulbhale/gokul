package controller

import (
	"github.com/sirupsen/logrus"
	"github.com/vipulbhale/gokul/server/util"
)

var log *logrus.Logger = util.GetLogger()

type Controller interface {
	Render()
}

type BaseController struct {
	Controller
}

type ModelAndView struct {
	Model interface{}
	View  string
}

func (modelAndView *ModelAndView) SetModel(model interface{}) {
	modelAndView.Model = model
}

func (modelAndView *ModelAndView) SetView(view string) {
	modelAndView.View = view
}

func (modelAndView *ModelAndView) GetModel() interface{} {
	return modelAndView.Model
}

func (modelAndView *ModelAndView) GetView() string {
	return modelAndView.View
}

func (baseController *BaseController) Render() {
	log.Debugln("Inside render method of BaseController...")
}

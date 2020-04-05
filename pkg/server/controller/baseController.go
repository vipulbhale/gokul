package controller

import (
	"github.com/vipulbhale/gokul/pkg/server/util"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

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

func init() {
	log = util.GetLogger()
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

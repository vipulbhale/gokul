package controller

import (
	"fmt"
)

var (
	MapOfControllers map[string]BaseController
)

func init() {
	MapOfControllers = make(map[string]BaseController)
}

type Controller interface {
	Render()
}

type BaseController struct {
	Controller

}

func (baseController *BaseController) Render() {
	fmt.Println("Inside render method")
}

func (baseController BaseController) RegisterController(name string) {
	MapOfControllers[name] = baseController
	fmt.Println(MapOfControllers)
}

func (baseController BaseController) GetMapOfController() (map[string]BaseController){
	return MapOfControllers
}
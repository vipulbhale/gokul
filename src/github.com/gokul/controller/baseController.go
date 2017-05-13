package controller

import (
	"fmt"
)
var (
	MapOfControllers map[string]BaseController
)


func init() {

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
package controller

import (
	"fmt"
)



type Controller interface {
	Render()
}

type BaseController struct {
	Controller

}

func (baseController *BaseController) Render() {
	fmt.Println("Inside render method")
}


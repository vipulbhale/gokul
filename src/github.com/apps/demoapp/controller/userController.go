package controller

import (
	"fmt"
	controller2 "github.com/gokul/controller"
)

type DemoController struct {
	*controller2.BaseController
}


func (d *DemoController) Demo() {
	fmt.Print("Hi there")
	d.Render()
}

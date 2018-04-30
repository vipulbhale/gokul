package controller

import (
	"fmt"
	controller2 "github.com/gokul/controller"
)

type DemoController1 struct {
	*controller2.BaseController
}


func (d *DemoController1) Demo() {
	fmt.Print("Hi there")
	d.Render()
}

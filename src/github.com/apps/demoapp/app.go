package main

import (
	"github.com/apps/demoapp/controller"
	"fmt"
)
var (
	demoController *controller.DemoController

)

func RegisterAllControllers(){
	demoController = new(controller.DemoController)

}

func main(){
	fmt.Println("hi there how are you")
}
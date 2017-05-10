package controller

import (
	"fmt"
	"net/http"
)

type Controller interface {
}

type BaseController struct {
	Controller
	w http.ResponseWriter
	r *http.Request
}

func (cntrl BaseController) Render() {
	fmt.Println("Inside render method")
}
func NewController(w http.ResponseWriter, r *http.Request) *BaseController {
	return &BaseController{w: w, r: r}
}

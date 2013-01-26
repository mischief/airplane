package controllers

import (
  "github.com/robfig/revel"
//  "github.com/mischief/airplane/app/models"
)

type Application struct {
	GorpController
}

func (c Application) Index() rev.Result {
	return c.Render()
}


package controllers

import (
//  "code.google.com/p/go.crypto/bcrypt"
//	"fmt"
  "github.com/robfig/revel"
//	"github.com/mischief/airplane/app/models"
//	"strings"
)

type Blogs struct {
  Application
}

func (b Blogs) Index() rev.Result {
  return b.Render()
}

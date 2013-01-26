package models

import (
	"fmt"
	"github.com/robfig/revel"
	"github.com/coopernurse/gorp"
//	"regexp"
)

type Post struct {
  PostId  int
  UserId  int
  Content string

  //
  User *User
}

func (p *Post) String() string {
	return fmt.Sprintf("Post(%d,%d)", p.PostId, p.UserId)
}

func (p *Post) Validate(v *rev.Validation) {
  v.Required(p.User)
}

func (p *Post) PreInsert(_ gorp.SqlExecutor) error {
  p.UserId = p.User.UserId
  return nil
}

func (p *Post) PostGet(exe gorp.SqlExecutor) error {
  var (
    obj interface{}
    err error
  )

  obj, err = exe.Get(User{}, p.UserId)
  if err != nil {
    return fmt.Errorf("Error loading post's user (%d): %s", p.UserId, err)
  }

  p.User = obj.(*User)

  return nil
}


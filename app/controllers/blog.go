package controllers

import (
//  "code.google.com/p/go.crypto/bcrypt"
//	"fmt"
  "github.com/robfig/revel"
	"github.com/mischief/airplane/app/models"
	"strings"
)

type Blog struct {
  Application
}

func (b Blog) checkUser() rev.Result {
	if user := b.connected(); user == nil {
		b.Flash.Error("Please log in first")
		return b.Redirect(Application.Index)
	}
	return nil
}

func (b Blog) Index() rev.Result {
  return b.Render()
}

func (b Blog) List(search string, size, page int) rev.Result {
  if page == 0 {
    page = 1
  }

  nextPage := page + 1
  search = strings.TrimSpace(search)

  var posts []*models.Post

  if search == "" {
    posts = loadPosts(b.Txn.Select(models.Post{},
      `select * from Post limit ?, ?`, (page-1)*size, size))
  } else {
    search = strings.ToLower(search)
    posts = loadPosts(b.Txn.Select(models.Post{},
      `select * from Post where lower(Content) like ? limit ?, ?`, "%"+search+"%", (page-1)*size, size))

  }

  return b.Render(posts, search, size, page, nextPage)
}

func loadPosts(results []interface{}, err error) []*models.Post {
  if err != nil {
    panic(err)
  }

  var posts []*models.Post

  for _, r := range results {
    posts = append(posts, r.(*models.Post))
  }

  return posts
}

func (b Blog) Post(post models.Post) rev.Result {
  post.User = b.connected()
  post.Validate(b.Validation)

  if b.Validation.HasErrors() {
    b.Validation.Keep()
    b.FlashParams()
    return b.Redirect(Blog.Index)
  }

  err := b.Txn.Insert(&post)
  if err != nil {
    panic(err)
  }

  b.Flash.Success("Post %d by %s successful", post.PostId, post.User)

  return b.Redirect(Blog.Index)
}

func (b Blog) loadPostById(id int) *models.Post {
  p, err := b.Txn.Get(models.Post{}, id)
  if err != nil {
    panic(err)
  }

  return p.(*models.Post)
}

func (b Blog) Show(id int) rev.Result {
  post := b.loadPostById(id)
  if post == nil {
    return b.NotFound("Post not found.")
  }

  return b.Render(post)
}


package controllers

import (
//	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/db/app"

  "github.com/mischief/airplane/app/models"
)

var (
  dbm *gorp.DbMap
)

type GorpPlugin struct {
  rev.EmptyPlugin
}

func (p GorpPlugin) OnAppStart() {
	db.DbPlugin{}.OnAppStart()
	dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.SqliteDialect{}}

  setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

  t := dbm.AddTable(models.User{}).SetKeys(true, "UserId")
	t.ColMap("Password").Transient = true
	setColumnSizes(t, map[string]int{
		"Username": 20,
		"Name":     100,
	})

  t = dbm.AddTable(models.Post{}).SetKeys(true, "PostId")
  t.ColMap("User").Transient = true
  setColumnSizes(t, map[string]int{
    "Content": 1000,
  })

  dbm.TraceOn("[gorp]", rev.INFO)
  dbm.CreateTables()

//  _, err := dbm.Get(models.User{}, 0)
//  if(err != nil) {
//    // user exists, skip insertion
//  } else {
//  	bcryptPassword, _ := bcrypt.GenerateFromPassword(
//  		[]byte("demo"), bcrypt.DefaultCost)
//	  demoUser := &models.User{0, "Demo User", "demo", "demo", bcryptPassword}
//  	if err := dbm.Insert(demoUser); err != nil {
//  		panic(err)
//  	}
//  }

}

type GorpController struct {
  *rev.Controller
  Txn *gorp.Transaction
}

func (c *GorpController) Begin() rev.Result {
	txn, err := dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() rev.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() rev.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}


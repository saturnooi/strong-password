package tu

import (
	"crypto/rand"
	"database/sql"
	"flag"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

type Context struct {
	database       string
	databaseSource string
	pDB            *sql.DB
	cleanupHooks   []func()

	DB *sql.DB
}

func (ctx *Context) setup() {
	var err error

	defer func() {
		if err != nil {
			panic(err)
		}
	}()

	ctx.pDB, err = sql.Open("postgres", fmt.Sprintf(ctx.databaseSource, "postgres"))
	if err != nil {
		panic(err)
	}
	_, err = ctx.pDB.Exec(`create database ` + ctx.database)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err != nil {
			ctx.Teardown()
		}
	}()

	ctx.DB, err = sql.Open("postgres", fmt.Sprintf(ctx.databaseSource, ctx.database))
	if err != nil {
		panic(err)
	}

	_, err = ctx.DB.ExecContext(ctx.Ctx(), `CREATE TABLE strong_password_log (
    	id BIGSERIAL PRIMARY KEY,
    req JSONB NOT NULL,
    res JSONB NOT NULL
	)
	`)
	if err != nil {
		panic(err)
	}

}

func (ctx *Context) Ctx() *gin.Context {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c
}

func (ctx *Context) Teardown() {
	for _, f := range ctx.cleanupHooks {
		f()
	}

	if ctx.DB != nil {
		ctx.DB.Close()
	}

	if ctx.pDB != nil {
		ctx.pDB.Exec(`drop database if exists ` + ctx.database)
		ctx.pDB.Close()
	}
}

// Setup setups test dependencies
func Setup() *Context {
	ctx := &Context{
		database:       fmt.Sprintf("test_%d", randInt()),
		databaseSource: os.Getenv("TEST_DB_URL"),
	}
	if ctx.databaseSource == "" {
		panic("TEST_DB_URL env required")
	}
	ctx.setup()

	return ctx
}

var (
	inTest     bool
	inTestOnce sync.Once
)

func InTest() bool {
	inTestOnce.Do(func() {
		inTest = flag.Lookup("test.v") != nil
	})
	return inTest
}

func randInt() int {
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

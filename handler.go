package main

import (
	"database/sql"
	"html/template"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber"
	"gopkg.in/russross/blackfriday.v2"
)

type Handler struct {
	conf *Config
	*sync.RWMutex
	db *sql.DB
}

func (hd *Handler) loadConfig() error {
	configMap, err := GetConfigMap(hd.db)
	if err != nil {
		return err
	}
	config, err := NewFromMap(configMap)
	if err != nil {
		return err
	}
	hd.conf = config
	return nil
}

func NewHandler(db *sql.DB) (*Handler, error) {
	hd := &Handler{
		DefaultConfig,
		&sync.RWMutex{},
		db,
	}
	err := hd.loadConfig()
	if err != nil {
		return &Handler{}, err
	}
	return hd, err
}

func (hd *Handler) ReloadConfig() error {
	hd.Lock()
	defer hd.Unlock()
	return hd.loadConfig()
}

func (hd *Handler) Config() *Config {
	hd.RLock()
	defer hd.RUnlock()
	if hd.conf == nil {
		c := *DefaultConfig
		return &c
	}
	c := *hd.conf
	return &c
}

func (hd *Handler) HandleHomeWithPage(ctx *fiber.Ctx) error {
	conf := hd.Config()
	total, err := GetThoughtsNumbers(hd.db)
	if err != nil {
		return err
	}

	tp := int(total) / conf.PageSize
	if int(total)%conf.PageSize != 0 {
		tp++
	}

	page := tp
	if p, err := strconv.Atoi(ctx.Params("pid")); err == nil {
		page = p
	}

	if page > tp {
		return fiber.ErrNotFound
	}

	thoughts, err := GetThoughtsByPage(hd.db, int64(conf.PageSize), int64(tp-page+1))
	if err != nil {
		return err
	}

	return ctx.Render("home", fiber.Map{
		"Thoughts":   thoughts,
		"Config":     conf,
		"Page":       page,
		"Total":      int(total),
		"TotalPages": tp,
		"PageList":   makePageList(tp),
		"LoggedIn":   ctx.Cookies("password") == conf.Password,
	})
}

func makePageList(i int) (arr []int) {
	if i < 1 {
		return []int{}
	}
	for j := i; j > 0; j-- {
		arr = append(arr, j)
	}
	return
}

func (hd *Handler) HandlePOSTNew(ctx *fiber.Ctx) error {
	conf := hd.Config()
	password := ctx.Cookies("password")
	if password != conf.Password {
		return ctx.Redirect("/admin/login", 302)
	}
	title := ctx.FormValue("title")
	markdown := ctx.FormValue("markdown")
	html := template.HTML(blackfriday.Run([]byte(markdown)))
	tht := &Thought{
		Title:    title,
		Date:     time.Now(),
		HTML:     html,
		Markdown: markdown,
	}
	_, err := InsertThought(hd.db, tht)
	if err != nil {
		return err
	}
	return ctx.Redirect("/", 302)
}

func (hd *Handler) HandleAdminLogin(ctx *fiber.Ctx) error {
	conf := hd.Config()
	password := ctx.Cookies("password")
	if password == conf.Password {
		return ctx.Redirect("/", 302)
	}
	return ctx.Render("login", fiber.Map{
		"Config": conf,
	})
}

func (hd *Handler) HandlePOSTAdminLogin(ctx *fiber.Ctx) error {
	ctx.Cookie(
		&fiber.Cookie{
			Name:  "password",
			Value: ctx.FormValue("password"),
		})
	return ctx.Redirect("/", 302)
}

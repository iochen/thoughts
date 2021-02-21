package main

import (
	"database/sql"
	"strconv"
	"sync"

	"github.com/gofiber/fiber"
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

func (hd *Handler) HandleHome(ctx *fiber.Ctx) error {
	page := 1
	if p, err := strconv.Atoi(ctx.Query("p")); err == nil {
		if p > page {
			page = p
		}
	}
	conf := hd.Config()
	thoughts, err := GetThoughtsByPage(hd.db, int64(conf.PageSize), int64(page))
	if err != nil {
		return err
	}

	total, err := GetThoughtsNumbers(hd.db)
	if err != nil {
		return err
	}
	tp := int(total) / conf.PageSize
	if int(total)%conf.PageSize != 0 {
		tp++
	}

	if page > tp {
		return fiber.ErrNotFound
	}

	return ctx.Render("home", fiber.Map{
		"Thoughts":   thoughts,
		"Config":     conf,
		"Page":       page,
		"Total":      int(total),
		"TotalPages": tp,
		"PageList":   makePageList(tp),
	})
}

func makePageList(i int) (arr []int) {
	if i < 1 {
		return []int{}
	}
	for j := 1; j <= i; j++ {
		arr = append(arr, j)
	}
	return
}

func (hd *Handler) HandleAdmin(ctx *fiber.Ctx) error {
	return nil
}

func (hd *Handler) HandleAdminLogin(ctx *fiber.Ctx) error {
	return nil

}

func (hd *Handler) HandlePOSTAdminLogin(ctx *fiber.Ctx) error {
	return nil

}

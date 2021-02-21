package main

import (
	"strconv"

	"github.com/iochen/thoughts/utils"
)

type Config struct {
	Name     string
	URL      string
	Author   string
	Avatar   string
	Icon     string
	Static   string
	Views    string
	Password string
	PageSize     int
}

var DefaultConfig = &Config{
	Name:     "Thoughts",
	URL:      "https://blog.sorakado.com",
	Author:   "Richard Chen",
	Avatar:   "/avatar.png",
	Icon:     "/favicon.ico",
	Static:   "static/",
	Views:    "views/",
	Password: utils.RandStr(32),
	PageSize: 20,
}

func NewFromMap(m map[string]string) (*Config, error) {
	if m == nil {
		return DefaultConfig, nil
	}
	if m["password"] == "" {
		m["password"] = utils.RandStr(32)
	}
	ps,err := strconv.Atoi(m["page_size"])
	if err != nil {
		return &Config{}, err
	}
	return &Config{
		Name:     m["name"],
		URL:      m["url"],
		Author:   m["author"],
		Avatar:   m["avatar"],
		Icon:     m["icon"],
		Static:   m["static"],
		Views:    m["views"],
		Password: m["password"],
		PageSize: ps,
	}, nil
}

func (conf *Config) Map() map[string]string {
	m := make(map[string]string)
	m["name"] = conf.Name
	m["url"] = conf.URL
	m["author"] = conf.Author
	m["avatar"] = conf.Avatar
	m["icon"] = conf.Icon
	m["static"] = conf.Static
	m["views"] = conf.Views
	m["password"] = conf.Password
	m["page_size"] = strconv.Itoa(conf.PageSize)
	return m
}

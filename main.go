package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gofiber/template/html"
	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	flagSQLite := flag.String("db", "thoughts.db", "SQLite path")
	flagAddr := flag.String("addr", "127.0.0.1:9001", "Listen Address")
	flag.Parse()

	logrus.Infoln("Starting Thoughts...")
	logrus.Infoln("Loading database(SQLite)...")
	if _, err := os.Stat(*flagSQLite); err == nil {
		logrus.Infof("Found database at %s.", *flagSQLite)
	} else if os.IsNotExist(err) {
		logrus.Infof("Database not found, creating a new one at %s.", *flagSQLite)
		if err := createDatabase(*flagSQLite); err != nil {
			logrus.Fatalf("An error occurred when creating a new database: %s.", err.Error())
			os.Exit(1)
		}
		logrus.Warnln("Your password:", DefaultConfig.Password)
	} else {
		logrus.Infof("Can not get status of database at %s, with error:%s.", *flagSQLite, err.Error())
	}

	logrus.Infof("Using database at %s.", *flagSQLite)
	db, _ := sql.Open("sqlite3", *flagSQLite)
	defer db.Close()
	if err := db.Ping(); err != nil {
		logrus.Fatalf("An error occurred when pinging database: %s.", err.Error())
		os.Exit(1)
	}

	handler, err := NewHandler(db)
	if err != nil {
		logrus.Fatalf("An error occurred when loading handler: %s.", err.Error())
		os.Exit(1)
	}

	conf := handler.Config()

	engine := html.New(conf.Views, ".html")
	engine.Reload(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", handler.HandleHomeWithPage)
	app.Get("/p/:pid", handler.HandleHomeWithPage)
	//app.Get("/search", handler.HandleSearch)
	app.Get("/admin/login", handler.HandleAdminLogin)
	app.Post("/admin/login", handler.HandlePOSTAdminLogin)
	app.Post("/new", handler.HandlePOSTNew)
	//app.Thought("/admin/new", handler.HandlePOSTAdminNew)
	//app.Get("/admin/config", handler.HandleAdminConfig)
	//app.Thought("/admin/config", handler.HandlePOSTAdminConfig)

	app.Static("/", conf.Static)

	logrus.Infof("Listening at %s.", *flagAddr)
	app.Listen(*flagAddr)
}

func createDatabase(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	db, _ := sql.Open("sqlite3", path)
	defer db.Close()
	err = CreateTable(db)
	if err != nil {
		return err
	}
	err = UpdateConfig(db, DefaultConfig)
	if err != nil {
		return err
	}
	return nil
}

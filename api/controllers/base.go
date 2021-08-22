package controllers

import (
	"fmt"
	"log"
	"net/http"
	"github.com/haikalvidya/goApiBlog/api/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/mysql"    //for mysql database driver
)

type Server struct {
	DB 		*gorm.DB
	Router	*mux.Router
}

// init for db and the router
func (server *Server) Init(DBDriver, DBUser, DBPassword, DBPort, DBHost, DBName string) {
	var err error
	if DBDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DBHost, DBPort, DBUser, DBName, DBPassword)
		server.DB, err = gorm.Open(DBDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", DBDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", DBDriver)
		}
	}
	// if use mysql database
	if DBDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DBUser, DBPassword, DBHost, DBPort, DBName)
		server.DB, err = gorm.Open(DBDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", DBDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", DBDriver)
		}
	}
	// if use sqlite3
	if DBDriver == "sqlite3" {
		server.DB, err = gorm.Open(DBDriver, DBName)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", DBDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", DBDriver)
		}
		server.DB.Exec("PRAGMA foreign_keys = ON")
	}

	// database migration
	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{})
	// new router
	server.Router = mux.NewRouter()

	server.initializeRoute()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
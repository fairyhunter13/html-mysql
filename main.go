package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func init() {
	var err error
	dbConn, err = gorm.Open(mysql.Open(os.Getenv("MYSQL_URI")), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error in openning the connection to mysql: %v", err)
		return
	}
}

func main() {
	f := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		Immutable:     true,
		AppName:       "Test",
	})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	f.Static("/html", "./staticdir")
	f.Get("/mysql", func(c *fiber.Ctx) (err error) {
		db, _ := dbConn.DB()
		err = db.Ping()
		if err != nil {
			return
		}
		resp := map[string]string{
			"status": "Ping MYSQL DB is succeed!",
		}

		err = c.JSON(resp)
		return
	})

	go func() {
		err := f.Listen(":80")
		if err != nil {
			log.Printf("Error in listening to the server: %v", err)
		}
	}()

	<-signalChan
	err := f.Shutdown()
	if err != nil {
		log.Printf("Error in shutdown the server: %v", err)
	}
}

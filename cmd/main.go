package main

import (
	"fmt"
	"github.com/autumnterror/onit/internal/domain"
	"github.com/autumnterror/onit/internal/infra/psql"
	"github.com/autumnterror/onit/internal/net"
	"github.com/autumnterror/onit/internal/repo"
	"github.com/autumnterror/onit/internal/service"
	"github.com/autumnterror/utils_go/pkg/log"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic(err)
	}
	user := os.Getenv("POSTGRES_USER")
	pw := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")

	if user == "" || pw == "" || db == "" {
		log.Panic("missing environment variables")
	}

	psx, err := psql.NewDB(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, pw, "db", 5432, db))
	if err != nil {
		log.Panic(err)
	}

	err = psx.AutoMigrate(&domain.Product{})
	if err != nil {
		log.Panic(err)
	}

	s := service.NewService(repo.NewProductRepository(psx))

	e := net.New(s)
	go e.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	if err := e.Stop(); err != nil {
		log.Error("", "stop echo", err)
	}

	log.Green(fmt.Sprintf("%s:%s", "", sig.String()))
}

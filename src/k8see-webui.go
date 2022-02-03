package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type appServer struct {
	db     *sql.DB
	router *mux.Router
}

func newServer(db *sql.DB) *appServer {
	s := appServer{
		db:     db,
		router: mux.NewRouter().StrictSlash(true),
	}
	return &s
}

func (s *appServer) Router() http.Handler {
	return s.router
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "f", "", "Configuration file")
	flag.Parse()

	if len(fileName) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	db, _ := cnxDB(fileName)
	defer db.Close()

	_, err := readYamlCnxFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	s := newServer(db)
	s.routes()

	err = http.ListenAndServe(":8081", s.Router())
	if err != nil {
		log.Fatal(err)
	}

}

func cnxDB(fileName string) (*sql.DB, error) {
	cfgEnvCnx, err := readYamlCnxFile(fileName)
	if err != nil {
		fmt.Println("Error when reading configuration file")
		os.Exit(1)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfgEnvCnx.Catsdb.Host, cfgEnvCnx.Catsdb.Port, cfgEnvCnx.Catsdb.User, cfgEnvCnx.Catsdb.Password, cfgEnvCnx.Catsdb.Dbname)
	// fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Failed to connect to database")
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to database")
		os.Exit(1)
	}
	return db, err
}

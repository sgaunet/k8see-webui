package main

import (
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

//go:embed templates
var htmlFiles embed.FS

//go:embed static
var staticCSS embed.FS

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
	var err error
	var fileName string
	var cfg AppConfig
	flag.StringVar(&fileName, "f", "", "Configuration file")
	flag.Parse()

	if len(fileName) == 0 {
		fmt.Println("No config file (option -f). Try to get configuration from environment variables.")
		cfg.Catsdb.Host = os.Getenv("DBHOST")
		cfg.Catsdb.Dbname = os.Getenv("DBNAME")
		if cfg.Catsdb.Port, err = strconv.Atoi(os.Getenv("DBPORT")); err != nil {
			fmt.Println("Env var DBPORT empty or not an int. Set port to default value (5432)")
			cfg.Catsdb.Port = 5432
		}
		cfg.Catsdb.User = os.Getenv("DBUSER")
		cfg.Catsdb.Password = os.Getenv("DBPASSWORD")
	} else {
		cfg, err = readYamlCnxFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err := cnxDB(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	s := newServer(db)
	s.routes()

	err = http.ListenAndServe(":8081", s.Router())
	if err != nil {
		log.Fatal(err)
	}
}

func cnxDB(cfgEnvCnx AppConfig) (*sql.DB, error) {
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

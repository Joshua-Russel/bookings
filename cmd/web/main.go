package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/Joshua-Russel/bookings/internal/config"
	"github.com/Joshua-Russel/bookings/internal/driver"
	"github.com/Joshua-Russel/bookings/internal/handlers"
	"github.com/Joshua-Russel/bookings/internal/helpers"
	"github.com/Joshua-Russel/bookings/internal/models"
	"github.com/Joshua-Russel/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const portNum = ":8000"

var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var app config.AppConfig

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()

	fmt.Println(fmt.Sprintf("Serving on port %s", portNum))

	serve := &http.Server{
		Addr:    portNum,
		Handler: routes(&app),
	}

	_ = serve.ListenAndServe()

}
func run() (*driver.DB, error) {

	infoLog = log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR \t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.User{})
	gob.Register(map[string]int{})

	//flags
	inProduction := flag.Bool("production", false, "Application in production mode")
	useCache := flag.Bool("cache", false, "Use template caches")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbSSL := flag.String("dbssl", "disable", "Database SSL mode (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	mailchan := make(chan models.MailData)
	app.MailChan = mailchan

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Secure = app.InProduction
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	app.Session = session

	log.Println("Connecting to database ... ")
	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=%s ",
		*dbHost, *dbUser, *dbName, *dbPass, *dbPort, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatalf("cannot connect to DB ....%v", err)
	}
	log.Println("Connected to database")

	tempSet, err := render.CreateTemplate()
	if err != nil {
		log.Fatal("cannot create template:", err)
		return nil, err
	}
	app.TmplCache = tempSet
	app.UseCache = *useCache
	app.InProduction = *inProduction

	render.NewRenderer(&app)
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	helpers.Newhelpers(&app)
	return db, nil
}

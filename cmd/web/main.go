package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"

	// "net/smtp"
	"os"
	"time"

	"github.com/Kawaeugtkp/bookings/internal/config"
	"github.com/Kawaeugtkp/bookings/internal/driver"
	"github.com/Kawaeugtkp/bookings/internal/handlers"
	"github.com/Kawaeugtkp/bookings/internal/helpers"
	"github.com/Kawaeugtkp/bookings/internal/models"
	"github.com/Kawaeugtkp/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080" // どうやらconstがletにあたるということみたい

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main Application function
func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err) // Fatalはterminalに表示するっていうprint的な役割に加えて
		// アプリケーションをここで停止させるっていう役割を持つ
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	// from := "me@here.com"
	// auth := smtp.PlainAuth("", from, "", "localhost")
	// err = smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("Hello, world"))
	// if err != nil {
	// 	log.Println(err)
	// }

	fmt.Println("Starting application on port", portNumber)
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{}) // アプリを通して保持するものだから
	// ここに書いているらしい、
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = *inProduction
	app.Usercache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo) // handler自体がなんか広大なインスタンスみたいに
	// なっていて、そこの要素を色々な部分で変えていると。だから下のhandleFuncでも
	// 普通にhandlerのRepoが使えているってことだと思う
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	// http.HandleFunc("/", handler.Repo.Home) // リロードしてもmainをもう一回呼び出すのではないみたい。でもここは実行されている
	// http.HandleFunc("/about", handler.Repo.About)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	// 	n, err := fmt.Fprintf(w, "Hello, World!")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(fmt.Sprintf("Number of bytes written: %d", n))
	// })
	return db, nil
}

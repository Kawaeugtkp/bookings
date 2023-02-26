package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kawaeugtkp/bookings/pkg/config"
	"github.com/Kawaeugtkp/bookings/pkg/handlers"
	"github.com/Kawaeugtkp/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080" // どうやらconstがletにあたるということみたい

var app config.AppConfig
var session *scs.SessionManager

// main is the main Application function
func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache() 
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.Usercache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo) // handler自体がなんか広大なインスタンスみたいに
	// なっていて、そこの要素を色々な部分で変えていると。だから下のhandleFuncでも
	// 普通にhandlerのRepoが使えているってことだと思う
	render.Newtemplates(&app)

	// http.HandleFunc("/", handler.Repo.Home) // リロードしてもmainをもう一回呼び出すのではないみたい。でもここは実行されている
	// http.HandleFunc("/about", handler.Repo.About)

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	// 	n, err := fmt.Fprintf(w, "Hello, World!")
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(fmt.Sprintf("Number of bytes written: %d", n))
	// })

	fmt.Println("Starting application on port", portNumber)
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
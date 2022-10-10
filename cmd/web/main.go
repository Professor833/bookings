package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Professor833/bookings/pkg/config"
	"github.com/Professor833/bookings/pkg/handlers"
	"github.com/Professor833/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	// initialize the session cache
	session = scs.New()
	session.Lifetime = 24 * time.Hour // 24 hours
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // set to true when in production

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println("Cannot create template cache")
	}
	app.TemplateCache = tc
	if app.InProduction {
		app.UseCache = true
	} else {
		app.UseCache = false
	}

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = serve.ListenAndServe()
	log.Fatal(err)
}

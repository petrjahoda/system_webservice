package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"net/http"
	"os"
)

const version = "2021.1.2.17"
const serviceName = "System WebService"
const serviceDescription = "System web interface"
const config = "user=postgres password=Zps05..... dbname=version3 host=localhost port=5432 sslmode=disable"

type program struct{}

func main() {
	logInfo("MAIN", serviceName+" ["+version+"] starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("MAIN", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] started")
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	logInfo("MAIN", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	updateProgramVersion()
	router := httprouter.New()
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/mif/*filepath", http.Dir("mif"))
	router.GET("/", basicAuth(index))
	router.GET("/index", basicAuth(index))
	router.GET("/workplaces", basicAuth(workplaces))
	router.GET("/charts", basicAuth(charts))
	router.GET("/statistics", basicAuth(statistics))
	router.GET("/data", basicAuth(data))
	router.GET("/settings", basicAuth(settings))
	router.POST("/update_user_settings", updateUserSettings)
	go cacheUsers()
	go cacheCompanyName()
	go cacheLocales()
	err := http.ListenAndServe(":80", router)
	if err != nil {
		logError("MAIN", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("MAIN", serviceName+" ["+version+"] running")
}

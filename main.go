package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"net/http"
	"os"
)

const version = "2021.1.3.9"
const serviceName = "System WebService"
const serviceDescription = "System web interface"
const config = "user=postgres password=pj79.. dbname=system host=localhost port=5432 sslmode=disable"

type program struct{}

func main() {
	logInfo("SYSTEM", serviceName+" ["+version+"] starting...")
	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		logError("SYSTEM", "Cannot start: "+err.Error())
	}
	err = s.Run()
	if err != nil {
		logError("SYSTEM", "Cannot start: "+err.Error())
	}
}

func (p *program) Start(service.Service) error {
	logInfo("SYSTEM", serviceName+" ["+version+"] started")
	go p.run()
	return nil
}

func (p *program) Stop(service.Service) error {
	logInfo("SYSTEM", serviceName+" ["+version+"] stopped")
	return nil
}

func (p *program) run() {
	updateProgramVersion()
	router := httprouter.New()
	router.ServeFiles("/html/*filepath", http.Dir("html"))
	router.ServeFiles("/css/*filepath", http.Dir("css"))
	router.ServeFiles("/js/*filepath", http.Dir("js"))
	router.ServeFiles("/mif/*filepath", http.Dir("mif"))
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/", basicAuth(index))
	router.GET("/index", basicAuth(index))
	router.GET("/workplaces", basicAuth(workplaces))
	router.GET("/charts", basicAuth(charts))
	router.GET("/statistics", basicAuth(statistics))
	router.GET("/data", basicAuth(data))
	router.GET("/settings", basicAuth(settings))
	router.POST("/update_user_settings", updateUserSettings)
	router.POST("/get_table_data", getTableData)
	router.POST("/get_settings_data", getSettingsData)
	router.POST("/get_chart_data", getChartData)
	router.POST("/get_detail_settings", getDetailSettings)
	go cacheData()
	err := http.ListenAndServe(":82", router)
	if err != nil {
		logError("SYSTEM", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("SYSTEM", serviceName+" ["+version+"] running")
}

func faviconHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("FAVIKON")
	http.ServeFile(w, r, "./icon/favicon.ico")
}

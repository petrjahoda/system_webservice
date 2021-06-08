package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kardianos/service"
	"net/http"
	"os"
)

const version = "2021.2.3.8"
const serviceName = "System WebService"
const serviceDescription = "System web interface"

//const config = "user=postgres password=pj79.. dbname=system host=ec2-3-17-5-15.us-east-2.compute.amazonaws.com port=5432 sslmode=disable application_name=system_webservice"
const config = "user=postgres password=pj79.. dbname=system host=database port=5432 sslmode=disable application_name=system_webservice"
const (
	production = iota + 1
	downtime
	poweroff
)
const (
	digital = iota + 1
	analog
)

const (
	administrator = iota + 1
	poweruser
	user
)

const (
	zapsi = iota + 1
	zapsiTouch
	s7plc
	datamax
	fileBased
	networkBased
)

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
	router.ServeFiles("/icon/*filepath", http.Dir("icon"))
	router.ServeFiles("/fonts/*filepath", http.Dir("fonts"))
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/", basicAuth(index))
	router.GET("/index", basicAuth(index))
	router.GET("/workplaces", basicAuth(workplaces))
	router.GET("/charts", basicAuth(charts))
	router.GET("/statistics", basicAuth(statistics))
	router.GET("/data", basicAuth(data))
	router.GET("/settings", basicAuth(settings))
	router.POST("/update_user_web_settings_from_web", updateUserWebSettingsFromWeb)
	router.POST("/update_workplaces", updateWorkplaces)
	router.POST("/load_index_data", loadIndexData)
	router.POST("/load_table_data", loadTableData)
	router.POST("/load_statistics_data", loadStatisticsData)
	router.POST("/load_chart_data", loadChartData)
	router.POST("/load_settings_data", loadSettingsData)
	router.POST("/load_settings_detail", loadSettingsDetail)
	router.POST("/load_device_port_detail", loadDevicePort)
	router.POST("/load_workplace_port_detail", loadWorkplacePort)
	router.POST("/save_alarm", saveAlarm)
	router.POST("/save_operation", saveOperation)
	router.POST("/save_order", saveOrder)
	router.POST("/save_product", saveProduct)
	router.POST("/save_part", savePart)
	router.POST("/save_state", saveState)
	router.POST("/save_workshift", saveWorkshift)
	router.POST("/save_breakdown", saveBreakdown)
	router.POST("/save_breakdown_type", saveBreakdownType)
	router.POST("/save_downtime", saveDowntime)
	router.POST("/save_downtime_type", saveDowntimeType)
	router.POST("/save_fault", saveFault)
	router.POST("/save_fault_type", saveFaultType)
	router.POST("/save_package", savePackage)
	router.POST("/save_package_type", savePackageType)
	router.POST("/save_user", saveUser)
	router.POST("/save_device", saveDevice)
	router.POST("/save_user_type", saveUserType)
	router.POST("/save_user_settings", saveUserSettings)
	router.POST("/save_system_settings", saveSystemSettingsDetails)
	router.POST("/save_device_port_details", saveDevicePort)
	router.POST("/save_workplace_mode", saveWorkplaceMode)
	router.POST("/save_workplace_section", saveWorkplaceSection)
	router.POST("/save_workplace", saveWorkplace)
	router.POST("/delete_workplace_port", deleteWorkplacePort)
	router.POST("/save_workplace_port_details", saveWorkplacePort)
	go cacheData()
	err := http.ListenAndServe(":80", router)
	if err != nil {
		logError("SYSTEM", "Problem starting service: "+err.Error())
		os.Exit(-1)
	}
	logInfo("SYSTEM", serviceName+" ["+version+"] running")
}

func faviconHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	http.ServeFile(w, r, "./icon/favicon.ico")
}

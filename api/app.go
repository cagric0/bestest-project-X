package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"projectXBackend/hazelcast"
)

type App struct {
	Server *http.Server
	Hz     *hazelcast.HZ
}

func StartApp(hz *hazelcast.HZ) {
	applicationPort := 3000
	httpServer := &http.Server{Addr: fmt.Sprintf(":%d", applicationPort)}

	instance := &App{
		Server: httpServer,
		Hz:     hz,
	}
	instance.RegisterRoutes()
	log.Println("Go service is running on port", applicationPort)
	log.Fatal(instance.Server.ListenAndServe())
}

func (a *App) RegisterRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/", a.home).Methods("GET")
	r.HandleFunc("/test-runs/{test-name}", a.testRunIDs).Methods("GET")
	r.HandleFunc("/test-logs/{log-identifier}", a.testLogs).Methods("GET")
	r.HandleFunc("/log-detail/{log-identifier}/{log-name}", a.testLogDetail).Methods("GET")

	r.HandleFunc("/test", a.testHandler).Methods("GET")
	r.HandleFunc("/push", a.pushHandler).Methods("POST")

	a.Server.Handler = r
}

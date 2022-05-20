package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hazelcast/hazelcast-go-client"
	"log"
	"net/http"
)

type App struct {
	Server *http.Server
	Hz     *hazelcast.Client
}

func StartApp(hz *hazelcast.Client) {
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

	r.HandleFunc("/test", a.testHandler).Methods("GET")
	r.HandleFunc("/push", a.pushHandler).Methods("POST")
	r.HandleFunc("/logs", a.getLogsHandler).Methods("GET")

	a.Server.Handler = r
}

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
	r.HandleFunc("/tests/{class-name}/{method-name}", a.testRunIDs).Methods("GET")
	r.HandleFunc("/tests/{class-name}/{method-name}/{runId}", a.testLogs).Methods("GET")

	r.HandleFunc("/clear", a.clearHandler).Methods("GET")
	r.HandleFunc("/push", a.pushHandler).Methods("POST")

	a.Server.Handler = r
}

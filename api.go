package main 

import (
	//"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)


type API struct {
	router *mux.Router
	cluster *Cluster
}

func NewAPI(cluster *Cluster) *API {
	api := &API{
		router: mux.NewRouter(),
		cluster: cluster,
	}
	api.setupRoutes()
	return api
}

func (api *API) setupRoutes(){
	api.router.HandleFunc("/set/{key}/{value}", api.setHandler).Methods("POST")
	api.router.HandleFunc("/get/{key}", api.getHandler).Methods("GET")
	api.router.HandleFunc("/delete/{key}", api.deleteHandler).Methods("DELETE")
}

func (api *API) setHandler(w http.ResponseWriter, r *http.Request) {
	// handle set request 

}

func (api *API) getHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *API) deleteHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *API) Run(addr string){
	http.ListenAndServe(addr, api.router)
}

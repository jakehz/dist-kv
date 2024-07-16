package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)


type API struct {
	router *mux.Router
	cluster *Cluster
}

type KVPair struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

type ClusterIP struct {
	IpAddress string `json:"ipAddress"`
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
	api.router.HandleFunc("/join", api.joinClusterFromNode).Methods("POST")
}

func (api *API) setHandler(w http.ResponseWriter, r *http.Request) {
	// handle set request 
	
}

func (api *API) getHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *API) deleteHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *API) joinClusterFromNode(w http.ResponseWriter, r *http.Request) {
	var clusterIP ClusterIP
	err := json.NewDecoder(r.Body).Decode(&clusterIP)
	if err != nil {
		log.Printf("Error parsing request: %v\n", err)
	}
	log.Printf("Recieved %v", clusterIP)
	params := []string{clusterIP.IpAddress}
	err = api.cluster.Join(params)
	if err != nil {
		log.Printf("Error joining cluster: %v", err)
	}
}

func (api *API) Run(addr string){
	fmt.Printf("Running server at %v\n", addr)
	http.ListenAndServe(addr, api.router)
}

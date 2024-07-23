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
	api.router.HandleFunc("/get_whole_kv", api.getWholeKV).Methods("GET")
	api.router.HandleFunc("/delete/{key}", api.deleteHandler).Methods("DELETE")
	api.router.HandleFunc("/join", api.joinClusterFromNode).Methods("POST")
	api.router.HandleFunc("/get_memberlist", api.getMemberlist).Methods("GET")
}

func (api *API) setHandler(w http.ResponseWriter, r *http.Request) {
	// Get variables from url
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["value"]
	api.cluster.Set(key, val)
	fmt.Fprintf(w, "{\"success\": true}")
}

func (api *API) getWholeKV(w http.ResponseWriter, r *http.Request){
	regMap := map[string]interface{}{}
	api.cluster.store.data.Range(func(key, value interface{}) bool{
		regMap[fmt.Sprint(key)] = value 
		return true
	})
	fmt.Fprintf(w, "{\"data\": %v}", regMap)
}

func (api *API) getHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *API) deleteHandler(w http.ResponseWriter, r *http.Request) {

}

func (api *API) getMemberlist(w http.ResponseWriter, r *http.Request) {
	m := api.cluster.Memberlist.Members()
	for _, node := range m {
		if node != api.cluster.Memberlist.LocalNode() {
		 api.cluster.pingOtherNode(node)
		}
	}
	fmt.Fprintf(w, "%v", m)
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

func parseKV(r *http.Request) (*KVPair, error) {
	var pair KVPair
	err := json.NewDecoder(r.Body).Decode(&pair)
	if err != nil {
		return nil, err
	}
	return &pair, nil
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

type Site struct {
	Name         string `json:"name"`
	State        string `json:"state"`
	PhysicalPath string `json:"physicalPath"`
}

func getSites() ([]Site, error) {
	cmd := exec.Command("powershell", "Get-Website | Select-Object name, state, physicalPath | ConvertTo-Json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var sites []Site
	err = json.Unmarshal(output, &sites)
	if err != nil {
		return nil, err
	}

	return sites, nil
}

func sitesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	sites, err := getSites()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sites)
}

func main() {
	http.HandleFunc("/sites", sitesHandler)

	log.Println("API en ejecución...")
	log.Fatal(http.ListenAndServe(":18003", nil))
}
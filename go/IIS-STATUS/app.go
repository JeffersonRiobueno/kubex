package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type Site struct {
	Name         string `json:"name"`
	State        string `json:"state"`
	PhysicalPath string `json:"physicalPath"`
}

type SiteUrl struct {
	ResponseUri string `json`
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

func getSiteUrl(sitePath string) ([]SiteUrl, error) {
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Get-WebURL -PSPath '%s' | Select-Object ResponseUri", sitePath))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var siteUrl []SiteUrl
	err = json.Unmarshal(output, &siteUrl)
	if err != nil {
		return nil, err
	}

	return siteUrl, nil
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

func uriHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	sitePath := "IIS:\\SITES\\GRA_LGS_API_JUSTO"
	siteUrl, err := getSiteUrl(sitePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	json.NewEncoder(w).Encode(siteUrl)
}

func main() {
	http.HandleFunc("/sites", sitesHandler)
	http.HandleFunc("/getUrl", uriHandler)
	log.Println("API en ejecución...")
	log.Fatal(http.ListenAndServe(":18003", nil))
}

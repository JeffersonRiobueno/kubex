package main

import (
	"io/ioutil"
    "fmt"
	"log"
	"net/http"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"encoding/json"
)

type DService struct {
	Name  string `json:"name"`
	State string `json:"status"`
}
type MSJ struct {
	Error  int 
	Msg string 
}

func getStateString(state uint32) string {
	switch state {
    case 1:
        return "Stopped"
    case 2:
        return "StartPending"
    case 3:
        return "StopPending"
    case 4:
        return "Running"
    case 5:
        return "ContinuePending"
    case 6:
        return "PausePending"
    case 7:
        return "Paused"
    default:
        return fmt.Sprintf("Unknown (%d)", state)
    }
}

func getServiceStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    birdJson, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatal(err)
    }
	var result map[string]any
	json.Unmarshal([]byte(birdJson), &result)

	JSData := result["service"].(map[string]any)

	a:= make([]DService, 15)
	N := 0
	for key, value := range JSData {
		fmt.Println(key)
		if value.(string) == "" {
			http.Error(w, "El parámetro 'service' es obligatorio", http.StatusBadRequest)
			return
		}
		// Abrir el manejador del administrador de servicios
		scm, err := mgr.Connect()
		if err != nil {
			fmt.Println("No se pudo abrir el manejador del administrador de servicios:", err)
			return
		}
		defer scm.Disconnect()
		// Abrir el servicio a consultar
		service, err := scm.OpenService(value.(string))
		if err != nil {
			fmt.Println("No se pudo abrir el servicio:", err)
			return
		}
		defer service.Close()
		// Obtener el estado del servicio
		status, err := service.Query()
		if err != nil {
			fmt.Println("No se pudo obtener el estado del servicio:", err)
			return
		}
		// Mostrar el estado del servicio
		var NewService = DService{Name: value.(string), State: getStateString(uint32(status.State))}
		a[N]=NewService
		N+=1
		fmt.Println("Estado del servicio:", status.State)
	}
	json.NewEncoder(w).Encode(a)

	

}

func getService(name string) (*mgr.Service, error) {
    m, err := mgr.Connect()
    if err != nil {
        return nil, err
    }
    defer m.Disconnect()

    s, err := m.OpenService(name)
    if err != nil {
        return nil, err
    }
    return s, nil
}

func StopService(name string) MSJ{
	e:= 0
	m:= "Stop Service"
	service, err := getService(name)
	if err != nil {
		fmt.Printf("Could not get service: %v", err)
		e=1
		m="Could not get service"
	}
	defer service.Close()

	status, err := service.Control(svc.Stop)
	if err != nil {
		fmt.Printf("Could not stop service: %v", err)
		e=1
		m="Could not stop service"
	}
	fmt.Printf("Stop status: %v\n", status)
	p := MSJ{
        Error: e,
        Msg: m,
    }
	return p
}


func startService(serviceName string) MSJ {
    e:= 0
	m:= "Start Service"
	service, err := getService(serviceName)
	if err != nil {
		fmt.Printf("Could not get service: %v", err)
		e=1
		m="Could not get service"
	}
    defer service.Close()

    // Iniciar el servicio
    err = service.Start()
    if err != nil {
        fmt.Printf("Could not start service %s: %v", serviceName, err)
		e=1
		m="Could not start service"
    }
	p := MSJ{
        Error: e,
        Msg: m,
    }
    return p
}

func stop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    query := r.URL.Query()
	service := query.Get("service")
	res := StopService(service)
	json.NewEncoder(w).Encode(res)

}

func start(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    query := r.URL.Query()
	service := query.Get("service")
	res := startService(service)
	json.NewEncoder(w).Encode(res)

}

func main() {
	http.HandleFunc("/", getServiceStatus)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/start", start)
    log.Fatal(http.ListenAndServe(":18002", nil))
}





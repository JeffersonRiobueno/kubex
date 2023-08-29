package main
import (
	"net/http"
	"fmt"
	"log"
	"os/exec"
    "math/rand"
	"encoding/json"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "time"
    "strconv"



)

var r *rand.Rand
 
func init() {
    r = rand.New(rand.NewSource(time.Now().Unix()))
}
 
// RandString genera una cadena aleatoria
func RandString(len int) string {
    bytes := make([]byte, len)
    for i := 0; i < len; i++ {
        b := r.Intn(26) + 65
        bytes[i] = byte(b)
    }
    return string(bytes)
}

	type Pais struct {
		Nombre     string
		Habitantes int
		Capital    string
	}
    type Data struct {
        Proyect, Rama, Folder, Papp, Pbd string
        Id, Status int
    }

    func obtenerBaseDeDatos() (db *sql.DB, e error) {
        usuario := "root"
        pass := "admin"
        host := "tcp(172.20.52.93:3304)"
        nombreBaseDeDatos := "kubex"
        // Debe tener la forma usuario:contraseña@host/nombreBaseDeDatos
        db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", usuario, pass, host, nombreBaseDeDatos))
        if err != nil {
            return nil, err
        }
        return db, nil
    }
    
  
func ListarProyect(writer http.ResponseWriter, request *http.Request){
    writer.Header().Set("Content-Type", "application/json")
    writer.Header().Set("Access-Control-Allow-Origin", "*")
    Pro := []Data{}
    db, err := obtenerBaseDeDatos()
    if err != nil {
        fmt.Println(err)

    }
    defer db.Close()
    filas, err := db.Query("SELECT * FROM runner WHERE `status` = '1' order by 1 desc")

    if err != nil {
        fmt.Println(err)

    }
    // Si llegamos aquí, significa que no ocurrió ningún error
    defer filas.Close()

    // Aquí vamos a "mapear" lo que traiga la consulta en el while de más abajo

    // Recorrer todas las filas, en un "while"
    for filas.Next() {
        var proyect, rama, folder, pApp, pBd string
        var id, status int
        err = filas.Scan(&id, &proyect, &rama, &folder, &pApp, &pBd, &status)
        // Al escanear puede haber un error
        if err != nil {
            fmt.Println(err)

        }
        // Y si no, entonces agregamos lo leído al arreglo

        pro := Data{
            Id: id,
            Proyect: proyect,
            Rama: rama,
            Folder: folder,
            Papp: pApp,
            Pbd: pBd,
            Status: status,
        }
        Pro = append(Pro,pro)
    }
    
    json.NewEncoder(writer).Encode(Pro)


}
    
   
func RunPoyect(writer http.ResponseWriter, request *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    rama := query.Get("rama")
    proyect := query.Get("proyect")
    papp := strconv.Itoa(PortApp())
    pbd :=""
    
	folder := RandString(10)
	if proyect == "WebRosatel_Encapsulado_V1" {
        pbd = strconv.Itoa(PortBD())
	}else{
        pbd = "1"
    }
    db, err := obtenerBaseDeDatos()
	if err != nil {
        fmt.Println(err)
	}
	defer db.Close()
	sentenciaPreparada, err := db.Prepare("INSERT INTO `runner` (`proyect`, `rama`, `folder`, `pApp`, `pBd`, `status`) VALUES (?, ?, ?, ?, ?, '1');")
	if err != nil {
        fmt.Println(err)
	}
	defer sentenciaPreparada.Close()
	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec(proyect, rama, folder, papp, pbd)
	if err != nil {
        fmt.Println(err)
	}

    out, err := exec.Command("/bin/sh", "create.sh", rama, papp, pbd, folder, proyect).Output()
    if err != nil {
        fmt.Println(err)
    }

    fmt.Printf("%s\n", out)
    p := Pais{
        Nombre:     "Canada",
        Habitantes: 37314442,
        Capital:    folder,
    }


	json.NewEncoder(writer).Encode(p)
}


func Restart(writer http.ResponseWriter, request *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    rama := query.Get("rama")
    papp := query.Get("papp")
    pbd := query.Get("pbd")
	folder := query.Get("pfolder")
	out, err := exec.Command("/bin/sh", "restart.sh", rama, papp, pbd, folder).Output()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("asas%s\n", out)
	
	p := Pais{
        Nombre:     "Canada",
        Habitantes: 37314442,
        Capital:    folder,
    }


	json.NewEncoder(writer).Encode(p)
//    fmt.Printf(out)
//	return 2

}

func Delete(writer http.ResponseWriter, request *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    ID := query.Get("ProyectID")
    folder := query.Get("pfolder")
	out, err := exec.Command("/bin/sh", "delete.sh", folder).Output()
    if err != nil {
        fmt.Println(err)
    }

    db, err := obtenerBaseDeDatos()
	if err != nil {
        fmt.Println(err)

	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("UPDATE `runner` SET `status` = '0' WHERE (`id` = ?);")
	if err != nil {
        fmt.Println(err)
	}
	defer sentenciaPreparada.Close()
	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec(ID)
	if err != nil {
        fmt.Println(err)
	}
    fmt.Printf("asas%s\n", out)
	
	p := Pais{
        Nombre:     "Canada",
        Habitantes: 37314442,
        Capital:    folder,
    }

	json.NewEncoder(writer).Encode(p)

}

  
func PortApp()(p int){
    db, err := obtenerBaseDeDatos()
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
    filas, err := db.Query("select MIN(portApp) as portApp from portApp where portApp not in(select pApp from runner as r where r.status=1) ")

    if err != nil {
        fmt.Println(err)
    }
    defer filas.Close()
    for filas.Next() {
        err = filas.Scan(&p)
        if err != nil {
            fmt.Println(err)
        }
    }
    return p
}

func PortBD()(p int){
    db, err := obtenerBaseDeDatos()
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
    filas, err := db.Query("select MIN(portBD) as portBD  from portBD where portBD not in(select pBd from runner as r where r.status=1)")

    if err != nil {
        fmt.Println(err)
    }
    defer filas.Close()
    for filas.Next() {
        err = filas.Scan(&p)
        if err != nil {
            fmt.Println(err)
        }
    }
    return p
}



//INSERT INTO `kubex`.`runner` (`proyect`, `rama`, `folder`, `pApp`, `pBd`, `status`) VALUES ('test', 'dev', 'fold', '80', '30', '0');

func main() {
//    bd := NuevaBD()
	http.HandleFunc("/start", RunPoyect)
	http.HandleFunc("/restart", Restart)
	http.HandleFunc("/delete", Delete)
    http.HandleFunc("/listar", ListarProyect)

	fmt.Println("Running ...")
    err := http.ListenAndServe(":8010", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err.Error())
    }
    //http.ListenAndServe(":8080", MuestraLibros(bd))
	
}
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



)
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
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
            log.Fatal(err)

        }
        defer db.Close()
        filas, err := db.Query("SELECT * FROM runner WHERE `status` = '1' order by 1 desc")
    
        if err != nil {
            log.Fatal(err)

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
                log.Fatal(err)

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
    papp := query.Get("papp")
    pbd := query.Get("pbd")
    proyect := query.Get("proyect")
	folder := randSeq(10)
	out, err := exec.Command("/bin/sh", "e.sh", rama, papp, pbd, folder).Output()
    if err != nil {
        log.Fatal(err)
    }
    db, err := obtenerBaseDeDatos()
	if err != nil {
        log.Fatal(err)
	}
	defer db.Close()
	sentenciaPreparada, err := db.Prepare("INSERT INTO `runner` (`proyect`, `rama`, `folder`, `pApp`, `pBd`, `status`) VALUES (?, ?, ?, ?, ?, '1');")
	if err != nil {
        log.Fatal(err)
	}
	defer sentenciaPreparada.Close()
	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec(proyect, rama, folder, papp, pbd)
	if err != nil {
        log.Fatal(err)
	}
    fmt.Printf("asas%s\n", out)
	
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
        log.Fatal(err)
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
        log.Fatal(err)
    }

    db, err := obtenerBaseDeDatos()
	if err != nil {
        log.Fatal(err)
	}
	defer db.Close()

	sentenciaPreparada, err := db.Prepare("UPDATE `runner` SET `status` = '0' WHERE (`id` = ?);")
	if err != nil {
        log.Fatal(err)
	}
	defer sentenciaPreparada.Close()
	// Ejecutar sentencia, un valor por cada '?'
	_, err = sentenciaPreparada.Exec(ID)
	if err != nil {
        log.Fatal(err)
	}
    fmt.Printf("asas%s\n", out)
	
	p := Pais{
        Nombre:     "Canada",
        Habitantes: 37314442,
        Capital:    folder,
    }


	json.NewEncoder(writer).Encode(p)


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
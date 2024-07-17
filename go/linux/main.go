package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "time"

    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
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
    Id, Status                       int
}

func obtenerBaseDeDatos() (*sql.DB, error) {
    // Cargar las variables de entorno desde el archivo .env
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error al cargar el archivo .env: %v", err)
    }
    // Obtener las variables de entorno
    user := os.Getenv("USER_BD")
    pass := os.Getenv("PASS_BD")
    host := os.Getenv("HOST")
    name := os.Getenv("NAME_BD")

    // Formatear la cadena de conexión
    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, host, name)

    // Abrir la conexión a la base de datos
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, fmt.Errorf("error al abrir la base de datos: %w", err)
    }

    // Comprobar la conexión a la base de datos
    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
    }

    return db, nil
}

func ListarProyect(writer http.ResponseWriter, request *http.Request) {
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
    defer filas.Close()
    for filas.Next() {
        var proyect, rama, folder, pApp, pBd string
        var id, status int
        err = filas.Scan(&id, &proyect, &rama, &folder, &pApp, &pBd, &status)
        if err != nil {
            fmt.Println(err)
        }
        pro := Data{
            Id:      id,
            Proyect: proyect,
            Rama:    rama,
            Folder:  folder,
            Papp:    pApp,
            Pbd:     pBd,
            Status:  status,
        }
        Pro = append(Pro, pro)
    }
    json.NewEncoder(writer).Encode(Pro)
}

func RunPoyect(writer http.ResponseWriter, request *http.Request) {
    writer.Header().Set("Content-Type", "application/json")
    writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    rama := query.Get("rama")
    proyect := query.Get("proyect")
    papp := strconv.Itoa(PortApp())
    pbd := ""
    folder := RandString(10)
    Repo := os.Getenv("R_REPO")
    if proyect == "WebRosatel_Encapsulado_V1" {
        pbd = strconv.Itoa(PortBD())
    } else {
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
    _, err = sentenciaPreparada.Exec(proyect, rama, folder, papp, pbd)
    if err != nil {
        fmt.Println(err)
    }
    out, err := exec.Command("/bin/sh", "create.sh", rama, papp, pbd, folder, proyect, Repo).Output()
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

func Restart(writer http.ResponseWriter, request *http.Request) {
    writer.Header().Set("Content-Type", "application/json")
    writer.Header().Set("Access-Control-Allow-Origin", "*")
    query := request.URL.Query()
    rama := query.Get("rama")
    papp := query.Get("papp")
    pbd := query.Get("pbd")
    proyect := query.Get("proyect")
    folder := query.Get("pfolder")
    out, err := exec.Command("/bin/sh", "restart.sh", rama, papp, pbd, folder, proyect).Output()
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

func Delete(writer http.ResponseWriter, request *http.Request) {
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

func PortApp() int {
    db, err := obtenerBaseDeDatos()
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
    filas, err := db.Query("select MIN(portApp) as portApp from portApp where portApp not in(select pApp from runner as r where r.status=1)")
    if err != nil {
        fmt.Println(err)
    }
    defer filas.Close()
    var p int
    for filas.Next() {
        err = filas.Scan(&p)
        if err != nil {
            fmt.Println(err)
        }
    }
    return p
}

func PortBD() int {
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
    var p int
    for filas.Next() {
        err = filas.Scan(&p)
        if err != nil {
            fmt.Println(err)
        }
    }
    return p
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error al cargar el archivo .env: %v", err)
    }
    // Obtener las variables de entorno
    Port := os.Getenv("R_PORT")
    http.HandleFunc("/start", RunPoyect)
    http.HandleFunc("/restart", Restart)
    http.HandleFunc("/delete", Delete)
    http.HandleFunc("/listar", ListarProyect)
    fmt.Println("Running ...")
    err2 := http.ListenAndServe(Port, nil)
    if err2 != nil {
        log.Fatal("ListenAndServe: ", err2.Error())
    }
}
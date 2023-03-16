package main

import (
    "archive/zip"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
	"log"
)

func main() {
    http.HandleFunc("/deploy", deployHandler)
	http.HandleFunc("/", Inicio)
	fmt.Println("Running ...")
    err := http.ListenAndServe(":8081", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err.Error())
    }

}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	appName := r.FormValue("app")
	appFolder := r.FormValue("folder")
    // 1. Subir el archivo comprimido de .NET
    r.ParseMultipartForm(32 << 20) // límite de 32 MB
    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Crear una carpeta temporal para almacenar el archivo
    tempDir, err := os.MkdirTemp("", "miapp")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer os.RemoveAll(tempDir)

    // Guardar el archivo en la carpeta temporal
    zipFilePath := filepath.Join(tempDir, handler.Filename)
    zipFile, err := os.Create(zipFilePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer zipFile.Close()
    io.Copy(zipFile, file)

    // Detener el sitio de IIS existente
    cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Stop-WebSite -Name %s", appName))
    if err := cmd.Run(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Backup dle proyecto 
    cmd2 := exec.Command("powershell", "-Command", fmt.Sprintf("%s/RUN.ps1", appFolder))
    if err := cmd2.Run(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Descomprimir el archivo en la ruta appFolder
    if err := os.MkdirAll(appFolder, os.ModePerm); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    zipR, err := zip.OpenReader(zipFilePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer zipR.Close()
    for _, f := range zipR.File {
        filePath := filepath.Join(appFolder, f.Name)
        if f.FileInfo().IsDir() {
            os.MkdirAll(filePath, os.ModePerm)
            continue
        }
        if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        inFile, err := f.Open()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            outFile.Close()
            return
        }
        io.Copy(outFile, inFile)
        outFile.Close()
        inFile.Close()
    }
    // Run config del proyecto 
    cmd2 := exec.Command("powershell", "-Command", fmt.Sprintf("%s/CONFIG.ps1", appFolder))
    if err := cmd2.Run(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Iniciar el sitio de IIS
    cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Start-WebSite -Name %s", appName))
    if err := cmd.Run(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "La aplicación ha sido desplegada exitosamente.")
}


func Inicio(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
    http.Error(w, "Inicio", http.StatusInternalServerError)

}
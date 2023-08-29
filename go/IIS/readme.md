# API DEPLOY
Esta api permite Ejecutar escript en un ambiente remoto de IIS

Para generar binario es windows ejecutar: 

`GOOS=windows go build -o programa.exe app.go`

Para que no se abra la terminal

`GOOS=windows go build -ldflags "-H windowsgui" -o programa.exe app.go`

Generar binario en linux:

`go build app.go `

Correr desde la terminal
`go run app.go`

### Docker

Para utilizarlce con docker: 

	docker run -dti --name GO -v "$PWD":/app -p 8010:8010 golang:1.18
	docker exec -ti GO bash
	cd /app
	go run app.go

Validar el puerto especificado en el archivo app.go
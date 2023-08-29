# API KUBEX
Esta api cuenta con 3 metodos que ejecutan script.sh en un serviodr remoto.
- create.sh
- update.sh
- delete.sh

Con esto puedes desplegar aplicaciones en un servidor remoto con linux

Para generar binario es windows ejecutar: 

`GOOS=windows go build -o programa.exe main.go`

Para que no se abra la terminal

`GOOS=windows go build -ldflags "-H windowsgui" -o programa.exe main.go`

Generar binario en linux:

`go build main.go `

Correr desde la terminal
`go run main.go`

### Docker

Para utilizarlce con docker: 

	docker run -dti --name GO -v "$PWD":/app -p 8010:8010 golang:1.18
	docker exec -ti GO bash
	cd /app
	go run main.go

Validar el puerto especificado en el archivo main.go
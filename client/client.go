package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	var conexion net.Conn
	var ip, puerto, valor string

	logger := IniciarLogger()
	defer logger.Close()

	config := IniciarConfig()
	//	defer config.Close()

	ip = config["IP"]
	puerto = config["Puerto"]
	valor = config["Valor"]

	LeerConsola(logger)

	conexion, err := CrearConexion(ip, puerto)
	if err != nil {
		fmt.Println("Error al crear la conexión:", err)
		return
	}
	defer LiberarConexion(conexion)

	err = EnviarMensaje(valor, conexion)
	if err != nil {
		fmt.Println("Error al enviar el mensaje:", err)
		return
	}

	EnviarPaquetePrueba(conexion)

	TerminarPrograma(conexion, logger, config)
}

func IniciarLogger() *os.File {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error al crear el archivo de log:", err)
		os.Exit(1)
	}
	fmt.Fprintln(logFile, "Hola! Soy un log")
	return logFile
}

func IniciarConfig() map[string]string {
	config := make(map[string]string)
	config["IP"] = "127.0.0.1"
	config["Puerto"] = "8080"
	config["Valor"] = "valor_clave"
	return config
}

func LeerConsola(logger *os.File) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			break
		}

		fmt.Fprintln(logger, input)
	}
}

func EnviarPaquetePrueba(conexion net.Conn) {
	paquete := CrearPaquete()
	defer EliminarPaquete(paquete)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Paquete> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			break
		}

		AgregarAPaquete(paquete, []byte(input), len(input))
	}

	err := EnviarPaquete(paquete, conexion)
	if err != nil {
		fmt.Println("Error al enviar el paquete:", err)
	}
}

func TerminarPrograma(conexion net.Conn, logger *os.File, config map[string]string) {
	LiberarConexion(conexion)
	logger.Close()
}

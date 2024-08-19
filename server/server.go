package main

import (
	"log"
	"os"
	"utilsServ"
)

var logger *log.Logger

func main() {
	// Configuración del logger
	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	logger = log.New(file, "Servidor: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Iniciar el servidor
	serverFd, err := utilsServ.IniciarServidor()
	if err != nil {
		logger.Fatalf("Error starting server: %v", err)
	}
	logger.Println("Servidor listo para recibir al cliente")

	clienteFd, err := utilsServ.EsperarCliente(serverFd)
	if err != nil {
		logger.Fatalf("Error accepting client: %v", err)
	}

	for {
		codOp, err := utilsServ.RecibirOperacion(clienteFd)
		if err != nil {
			logger.Println("Error receiving operation:", err)
			break
		}

		switch codOp {
		case utilsServ.MENSAJE:
			if err := utilsServ.RecibirMensaje(clienteFd, logger); err != nil {
				logger.Println("Error receiving message:", err)
				return
			}
		case utilsServ.PAQUETE:
			lista, err := utilsServ.RecibirPaquete(clienteFd, logger)
			if err != nil {
				logger.Println("Error receiving package:", err)
				return
			}
			logger.Println("Me llegaron los siguientes valores:")
			for _, valor := range lista {
				logger.Printf("%s", valor)
			}
		case -1:
			logger.Println("El cliente se desconectó. Terminando servidor")
			return
		default:
			logger.Println("Operación desconocida. No quieras meter la pata")
		}
	}
}

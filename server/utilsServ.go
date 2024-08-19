package utilsServ

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

const PUERTO = "8080"
const MENSAJE = 1
const PAQUETE = 2

func IniciarServidor() (net.Listener, error) {
	addr := "0.0.0.0:" + PUERTO
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	log.Println("Listo para escuchar a mi cliente")
	return listener, nil
}

func EsperarCliente(listener net.Listener) (net.Conn, error) {
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	log.Println("Se conectó un cliente!")
	return conn, nil
}

func RecibirOperacion(conn net.Conn) (int, error) {
	var codOp int32
	err := binary.Read(conn, binary.LittleEndian, &codOp)
	if err != nil {
		if err == io.EOF {
			return -1, nil
		}
		return -1, err
	}
	return int(codOp), nil
}

func RecibirBuffer(conn net.Conn) ([]byte, error) {
	var size int32
	err := binary.Read(conn, binary.LittleEndian, &size)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, size)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func RecibirMensaje(conn net.Conn, logger *log.Logger) error {
	buffer, err := RecibirBuffer(conn)
	if err != nil {
		return err
	}
	logger.Printf("Me llegó el mensaje %s", string(buffer))
	return nil
}

func RecibirPaquete(conn net.Conn, logger *log.Logger) ([]string, error) {
	buffer, err := RecibirBuffer(conn)
	if err != nil {
		return nil, err
	}

	var lista []string
	var desplazamiento int
	for desplazamiento < len(buffer) {
		var tamanio int32
		binary.Read(buffer[desplazamiento:], binary.LittleEndian, &tamanio)
		desplazamiento += 4
		valor := string(buffer[desplazamiento : desplazamiento+int(tamanio)])
		lista = append(lista, valor)
		desplazamiento += int(tamanio)
	}

	return lista, nil
}

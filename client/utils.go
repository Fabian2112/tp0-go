package main

import (
	"bytes"
	"encoding/binary"
	"net"
)

// Definición de las estructuras
type Paquete struct {
	CodigoOperacion int
	Buffer          *Buffer
}

type Buffer struct {
	Size   int
	Stream []byte
}

// Función para serializar un paquete
func SerializarPaquete(paquete *Paquete, bytesInt int) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int32(paquete.CodigoOperacion))
	binary.Write(buf, binary.LittleEndian, int32(paquete.Buffer.Size))
	buf.Write(paquete.Buffer.Stream)

	return buf.Bytes()
}

// Función para crear una conexión
func CrearConexion(ip string, puerto string) (net.Conn, error) {
	direccion := ip + ":" + puerto
	var conn net.Conn
	var err error

	for {
		conn, err = net.Dial("tcp", direccion)
		if err == nil {
			break
		}
		// Reintentar conexión
	}

	return conn, err
}

// Función para enviar un mensaje
func EnviarMensaje(mensaje string, conn net.Conn) error {
	paquete := &Paquete{
		CodigoOperacion: 1, // Supongamos que MENSAJE es 1
		Buffer: &Buffer{
			Size:   len(mensaje) + 1,
			Stream: append([]byte(mensaje), 0),
		},
	}

	bytesInt := paquete.Buffer.Size + 2*4 // 2 enteros de 4 bytes cada uno
	aEnviar := SerializarPaquete(paquete, bytesInt)

	_, err := conn.Write(aEnviar)
	if err != nil {
		return err
	}

	return nil
}

// Función para crear un buffer
func CrearBuffer(paquete *Paquete) {
	paquete.Buffer = &Buffer{
		Size:   0,
		Stream: nil,
	}
}

// Función para crear un paquete
func CrearPaquete() *Paquete {
	paquete := &Paquete{
		CodigoOperacion: 2, // Supongamos que PAQUETE es 2
	}
	CrearBuffer(paquete)
	return paquete
}

// Función para agregar datos a un paquete
func AgregarAPaquete(paquete *Paquete, valor []byte, tamanio int) {
	newSize := paquete.Buffer.Size + tamanio + 4
	stream := append(paquete.Buffer.Stream, make([]byte, tamanio+4)...)

	binary.LittleEndian.PutUint32(stream[paquete.Buffer.Size:], uint32(tamanio))
	copy(stream[paquete.Buffer.Size+4:], valor)

	paquete.Buffer.Size = newSize
	paquete.Buffer.Stream = stream
}

// Función para enviar un paquete
func EnviarPaquete(paquete *Paquete, conn net.Conn) error {
	bytesInt := paquete.Buffer.Size + 2*4
	aEnviar := SerializarPaquete(paquete, bytesInt)

	_, err := conn.Write(aEnviar)
	if err != nil {
		return err
	}

	return nil
}

// Función para eliminar un paquete
func EliminarPaquete(paquete *Paquete) {
	paquete.Buffer.Stream = nil
	paquete.Buffer = nil
}

// Función para liberar la conexión
func LiberarConexion(conn net.Conn) {
	conn.Close()
}

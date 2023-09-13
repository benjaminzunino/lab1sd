package main

import (
	"bufio"
	"context"
	"fmt"
	pb "lab1"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

var numeroRecibido int32 //para guardar numero recibido desde la central

func ParametrosDeInicio() int {

	// Leer archivo
	archivo, err := os.Open("regional/parametros_de_inicio.txt")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return -1
	}
	defer archivo.Close()

	// Crea un lector bufio para leer el archivo
	lector := bufio.NewReader(archivo)

	// Lee una línea (suponiendo que el entero está en una sola línea)
	linea, err := lector.ReadString('\n')

	// Convierte la línea a un entero
	valorEntero, err := strconv.Atoi(linea)

	// Ahora tienes el entero en la variable valorEntero
	fmt.Println("El valor leído del archivo es:", valorEntero)
	return valorEntero
}

type numberServer struct {
	pb.UnimplementedNumberServiceServer
}

func (s *numberServer) SendNumber(ctx context.Context, req *pb.NumberRequest) (*pb.NumberResponse, error) {

	//Guarda el numero recibido por la central
	number := req.GetNumber()
	fmt.Println("Interesados: ", interesados)
	fmt.Printf("Número recibido: %d\n", number)
	numeroRecibido = number

	// Se van generando interesados que obtendran la beta
	interesados_final := generarInteresados(interesados)
	fmt.Println("Numero de interesados que obtendran la beta:", interesados_final)
	interesados = interesados - interesados_final
	fmt.Println("Numero de interesados restantes:", interesados)

	// Se envia el valor de interesados_final a la cola asincrona de la central

	//Conexion con cola asincrona
	conn, err := amqp.Dial("amqp://tzoeibjf:BxXQ3VAX_Nyl1hW1D9A28hmrnquflDaG@jackal.rmq.cloudamqp.com/tzoeibjf")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// Enviar un mensaje a la misma cola asincrónica
	numero_a_enviar := interesados_final // El mensaje que deseas enviar
	mensaje := fmt.Sprintf("%d", numero_a_enviar)

	err = ch.Publish(
		"",         // Intercambio (dejar en blanco para usar el intercambio predeterminado)
		"server_1", // Nombre de la cola
		false,      // Mandar de manera inmediata (no esperar confirmación)
		false,      // Publicar como mandato (no es relevante para RabbitMQ)
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(mensaje),
		},
	)
	if err != nil {
		log.Fatalf("No se pudo enviar el mensaje: %s", err)
	}

	fmt.Printf("Mensaje enviado: %s\n", mensaje)

	return &pb.NumberResponse{Acknowledgment: "Número recibido por el servidor"}, nil
}

func generarAleatorio(min int, max int) int {
	// Inicializa el generador de números aleatorios con una semilla única
	rand.Seed(time.Now().UnixNano())
	// Genera y retorna un número aleatorio en el rango [min, max]
	return rand.Intn(max-min+1) + min
}

func generarInteresados(interesados int) int {

	// Calcular el 20% de la mitad de interesados
	veinte_por_ciento := int(float64(interesados/2) * 0.2)

	// Definir el rango mínimo y máximo para el valor aleatorio
	min := (interesados / 2) - veinte_por_ciento
	max := (interesados / 2) + veinte_por_ciento

	// Generar un valor aleatorio dentro del rango
	valorAleatorio := generarAleatorio(min, max)

	// Imprimir el valor aleatorio
	fmt.Printf("Valor aleatorio entre %d y %d: %d\n", min, max, valorAleatorio)
	return valorAleatorio
}

var interesados int //Variable global de interesados

func main() {
	// Primero se guarda la cantidad de interesados del archivo de entrada.
	interesados = ParametrosDeInicio()

	// Iniciar el servidor gRPC en segundo plano
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Error al escuchar: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterNumberServiceServer(s, &numberServer{})
		log.Println("Servidor gRPC escuchando en el puerto 50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Error al servir: %v", err)
		}
	}()

	// Mantén el programa principal en funcionamiento

	select {}

}

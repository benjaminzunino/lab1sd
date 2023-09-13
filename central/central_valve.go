package main

import (
	"bufio"
	"context"
	"fmt"
	pb "lab1"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

var minimo int
var maximo int

func ParametrosDeInicio() (int, int) {
	// Leer el archivo "parametros de inicio" y extraer los valores mínimo y máximo.
	// Retorna un slice de enteros [min, max].
	// Implementa la lógica para leer el archivo según tu estructura de datos y formato.
	// Puedes usar la función ioutil.ReadFile() para leer el contenido del archivo.
	// Puedes separar los valores usando strings.Split() y luego convertirlos a enteros con strconv.Atoi().
	// Retorna [min, max].

	archivo, err := os.Open("central/parametros_de_inicio.txt")
	if err != nil {
		fmt.Println("Error al abrir archivo: ", err)
		return -1, -1
	}
	defer archivo.Close()

	lector := bufio.NewScanner(archivo)

	var min int
	var max int
	var i int = 0
	var iter int

	for lector.Scan() {
		linea := lector.Text()
		if i == 0 {
			values := strings.Split(linea, "-")

			min, err = strconv.Atoi(values[0])
			max, err = strconv.Atoi(values[1])
		}
		if i == 1 {
			iter, err = strconv.Atoi(linea)
		}
		i++

	}

	rand.Seed(time.Now().UnixNano())

	numAleatorio := rand.Intn(max-min+1) + min

	minimo = min
	maximo = max

	return numAleatorio, iter
}

func generarLlaves(min int, max int) int {
	// Inicializa el generador de números aleatorios con una semilla única
	rand.Seed(time.Now().UnixNano())

	// Genera y retorna un número aleatorio en el rango [min, max]
	return rand.Intn(max-min+1) + min
}

func notifyServidoresRegionales(llavesGeneradas int) {
	// Implementa la lógica para notificar a los servidores regionales la cantidad de llaves generadas.
	// Puedes utilizar gRPC o cualquier otro método de comunicación síncrona según los requisitos de tu laboratorio.
	// Asegúrate de manejar la comunicación con los servidores regionales aquí.
	fmt.Printf("Notificar a los servidores regionales: %d llaves generadas\n", llavesGeneradas)
}

func main() {

	// Lee el archivo "parametros de inicio" para obtener los valores mínimo y máximo
	// para generar una cantidad predefinida de llaves al azar.
	var numero_aleatorio int
	var iteraciones int

	// Establece el contador de iteraciones desde el archivo.
	numero_aleatorio, iteraciones = ParametrosDeInicio()
	fmt.Println("Numero generado: ", numero_aleatorio)
	fmt.Println("Iteraciones del programa: ", iteraciones)
	fmt.Println("Minimo: ", minimo)
	fmt.Println("Maximo: ", maximo)

	//establecer conexion
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

	for iteraciones != -1 && iteraciones > 0 {
		// Genera una cantidad aleatoria de llaves dentro del rango definido.
		llavesGeneradas := generarLlaves(minimo, maximo)
		fmt.Println("Numero generado: ", llavesGeneradas)

		// Notifica a los servidores regionales la cantidad de llaves generadas (sincrono)
		conn1, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("No se pudo conectar: %v", err)
		}
		defer conn1.Close()

		c := pb.NewNumberServiceClient(conn1)

		number := int32(llavesGeneradas) // El número que deseas enviar
		req := &pb.NumberRequest{Number: number}

		resp, err := c.SendNumber(context.Background(), req)
		if err != nil {
			log.Fatalf("Error al enviar el número: %v", err)
		}

		log.Printf("Respuesta del servidor: %s", resp.GetAcknowledgment())

		//consumir cola asincrona
		msgs, err := ch.Consume(
			"server_1", // Nombre de la cola
			"",         // Nombre del consumidor (dejar en blanco para un consumidor automático)
			true,       // Autoack (marcando mensajes como entregados automáticamente)
			false,      // Exclusividad (permite múltiples consumidores en la misma cola)
			false,      // No local (no es relevante para RabbitMQ)
			false,      // No esperar confirmación (no es relevante para RabbitMQ)
			nil,        // Argumentos adicionales
		)
		if err != nil {
			log.Fatalf("No se pudo iniciar el consumidor: %s", err)
		}

		go func() {
			for msg := range msgs {
				// Procesa el mensaje aquí (en este ejemplo, simplemente lo imprime)
				//Aca se debe tratar el numero de interesados recibido del servidor regional para luego restarlos del numero de llaves y ver si lograron registrarse
				fmt.Printf("Recibido un mensaje: %s\n", msg.Body)
				interesados_recibidos, err := strconv.Atoi(string(msg.Body))
				if err != nil {
					fmt.Printf("Error al convertir la cadena a entero: %v\n", err)
					return
				}

				llaves_restantes := llavesGeneradas - interesados_recibidos

				if llaves_restantes > 0 {

					fmt.Println("Todos los interesados han obtenido la beta y han sobrado llaves.")

				} else if llaves_restantes == 0 {

					fmt.Println("Se han acabado las llaves.")

				} else {

					fmt.Println("Hay mas interesados que llaves disponibles.")
					exceso := interesados_recibidos - llaves_restantes
					fmt.Println("Se quedaran sin beta: ", exceso)

					//Se debe enviar el exceso al servidor regional y sumarlos nuevamente a los interesados

				}

			}
		}()

		// Decrementa el contador de iteraciones o verifica si se debe continuar.
		iteraciones--
		time.Sleep(time.Second * 10)
	}
}

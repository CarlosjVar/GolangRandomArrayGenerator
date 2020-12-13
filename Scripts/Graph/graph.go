package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pusher "github.com/pusher/pusher-http-go"
)

// Conexión con el cliente de pusher , que comunica go con js
var client = pusher.Client{
	AppID:   "1118979",
	Key:     "7befe6ab035a03a2ada9",
	Secret:  "1605eea01f08a987d955",
	Cluster: "us2",
	Secure:  true,
}

//Comienza la generación de números
func generate(c echo.Context) error {
	cantidad := 0
	if i, err := strconv.Atoi(c.Param("cantidad")); err == nil {
		cantidad = i
	}
	var waitGroup sync.WaitGroup
	randomch := make(chan int, 1)
	go generateRandom(&waitGroup, randomch, 11, cantidad)
	go createArray(&waitGroup, randomch, cantidad)
	waitGroup.Add(1)
	waitGroup.Wait()
	return c.String(http.StatusOK, "Simulation begun")
}

//Normaliza el numero dado en el rango de 0 a 31
func NormalizeRandom(randomNum float64) int {
	RandomNum := 0 + (int)(randomNum*((31-0)+1))
	return RandomNum
}

//Genera un número aleatorio en el intervalo de [0,1[
func generateRandom(wg *sync.WaitGroup, channel chan int, seed int, size int) {
	period := 8192
	constant := 23
	multiplicativeConstant := 3 + (8 * constant)
	for i := 0; i < size; i++ {
		num := ((seed * multiplicativeConstant) + constant) % int(period)
		seed = num
		normalizedNum := float64(num) / float64(period-1)
		randomNum := NormalizeRandom(normalizedNum)
		channel <- randomNum
		//Se pausa
		<-channel

	}
}

//Recibe un número aleatorio por medio de un canal y lo agrega a un array
func createArray(wg *sync.WaitGroup, channel chan int, size int) {
	randomArray := []int{}
	randomNum := 0
	for i := 0; i < size; i++ {
		randomNum = <-channel
		randomArray = append(randomArray, randomNum)
		client.Trigger("randomNum", "addNumber", randomArray)
		channel <- 0

	}
	defer wg.Done()
}
func main() {
	// Instancia de echo nuestro servidor
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Se definen las rutas del frontend
	e.File("/", "Frontend/index.html")
	e.File("/style.css", "Frontend/style.css")
	e.File("/app1.js", "Frontend/app1.js")
	e.GET("/generate/:cantidad", generate)
	// Iniciamos el servidor
	e.Logger.Fatal(e.Start(":9000"))
}

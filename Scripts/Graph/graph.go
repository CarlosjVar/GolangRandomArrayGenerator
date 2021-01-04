package main

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"

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
func generateSeed() int {
	seed := getTime() / 16
	for true {
		if big.NewInt((int64)(seed)).ProbablyPrime(0) {
			break
		}
		seed += 1

	}
	fmt.Print(seed)
	return seed

}

func getTime() int {
	year := time.Now().Year()
	mes := time.Now().Month()
	dia := time.Now().Day()
	hour := time.Now().Hour()
	minutos := time.Now().Minute()
	segundos := time.Now().Second()
	seed := ((((year*100+int(mes))*100+dia)*100+hour)*100 + minutos) / 50 * segundos
	return seed
}

func load(c echo.Context) error {
	return c.String(http.StatusOK, "Simulation begun")
}

func generate(c echo.Context) error {
	seed := generateSeed()
	cantidad := 0
	if i, err := strconv.Atoi(c.Param("cantidad")); err == nil {
		cantidad = i
	}
	var waitGroup sync.WaitGroup
	randomch := make(chan int, 1)
	go generateRandom(&waitGroup, randomch, seed, cantidad)
	waitGroup.Add(1)
	waitGroup.Wait()
	return c.String(http.StatusOK, "Simulation begun")
}

//Normaliza el numero dado en el rango de 0 a 31
func NormalizeRandom(randomNum float64, upperbound int16, lowerbound int16) int {
	RandomNum := 0 + (int)(randomNum*(((float64)(upperbound)-(float64)(lowerbound))+1))
	return RandomNum
}

//Genera un número aleatorio en el intervalo de [0,1[
func generateRandom(wg *sync.WaitGroup, channel chan int, seed int, size int) {
	randomArray := []int{}
	period := 8192
	constant := 23
	multiplicativeConstant := 3 + (8 * constant)
	for i := 0; i < size; i++ {
		num := ((seed * multiplicativeConstant) + constant) % int(period)
		seed = num
		normalizedNum := float64(num) / float64(period-1)
		randomNum := NormalizeRandom(normalizedNum, 31, 0)
		randomArray = append(randomArray, randomNum)

	}
	client.Trigger("ArrayChannel", "addNumber", randomArray)
	fmt.Printf("%v", randomArray)
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
	e.GET("/load", load)
	// Iniciamos el servidor
	e.Logger.Fatal(e.Start(":9000"))
}

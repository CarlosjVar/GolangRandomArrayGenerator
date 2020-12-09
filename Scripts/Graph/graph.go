package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	pusher "github.com/pusher/pusher-http-go"
)

// We register the Pusher client
var client = pusher.Client{
	AppID:   "1118979",
	Key:     "7befe6ab035a03a2ada9",
	Secret:  "1605eea01f08a987d955",
	Cluster: "us2",
	Secure:  true,
}

// visitsData is a struct
type visitsData struct {
	Pages int
	Count int
}

func setInterval(ourFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticker and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine so that the for loop is none blocking
	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					// This won't block
					go ourFunc()
				} else {
					// This will block
					ourFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear
}
func simulateAux(c echo.Context) error {
	setInterval(func() {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		newVisitsData := visitsData{
			Pages: r1.Intn(100),
			Count: r1.Intn(100),
		}
		client.Trigger("visitorsCount", "addNumber", newVisitsData)
	}, 500, true)

	return c.String(http.StatusOK, "Simulation begun")
}

func simulate(c echo.Context) error {

	var waitGroup sync.WaitGroup
	randomch := make(chan int, 1)
	go generateRandom(&waitGroup, randomch, 11, 100)
	go createArray(&waitGroup, randomch, 100)
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
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the HTTP routes
	e.File("/", "Frontend/index.html")
	e.File("/style.css", "Frontend/style.css")
	e.File("/app1.js", "Frontend/app1.js")
	e.GET("/simulate", simulate)
	fmt.Println("hOLA")
	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}

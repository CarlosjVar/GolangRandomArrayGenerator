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

type datos struct {
	intercambios  int
	condicionales int
	ciclos        int
	fechaI        string
	fechaF        string
}

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
func CopyArray(arrOri []int) []int {
	nuevoArray := []int{}
	var num = 0
	for i := 0; i < len(arrOri); i++ { //copiado de Array
		num = arrOri[i]
		nuevoArray = append(nuevoArray, num)
	}
	return nuevoArray
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
func start(c echo.Context) error {
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
	arrayChannel := make(chan []int, 1)
	go generateRandom(&waitGroup, randomch, arrayChannel, seed, cantidad)
	waitGroup.Add(1)
	waitGroup.Wait()
	arr1 := []int{}

	arr1 = <-arrayChannel //saca el array del channel
	arrGlob = CopyArray(arr1)

	return c.String(http.StatusOK, "Simulation Begun")
}

func sortingSimulation(c echo.Context) error {
	var waitGroup sync.WaitGroup
	arr1 := []int{}
	arr2 := []int{}
	arr3 := []int{}
	arr4 := []int{}
	fmt.Println("Unsorted Arrays")
	arr1 = CopyArray(arrGlob)
	arr2 = CopyArray(arr1) //copia el array
	arr3 = CopyArray(arr2) //copia el array
	arr4 = CopyArray(arr3) //copia el array

	quickChannel := make(chan []int, 1)
	insertionChannel := make(chan []int, 1)
	bubbleChannel := make(chan []int, 1)
	heapChannel := make(chan []int, 1)
	ciclosQuick := 0
	condicionalesQuick := 0
	intercambiosQuick := 0
	waitGroup.Add(4)
	go TempGraficarInsertion(insertionChannel)
	go InsertionSort(&waitGroup, arr1, insertionChannel) //InsetionSort al segundo Array
	go TempGraficarBubble(bubbleChannel)
	go BubbleSort(&waitGroup, arr2, bubbleChannel) //BubbleSort al primer Array
	go TempGraficarQuick(quickChannel)

	go QuickSort(&waitGroup, arr3, quickChannel, 0, 0, &ciclosQuick, &condicionalesQuick, &intercambiosQuick) //InsetionSort al segundo Array
	go TempGraficarHeap(heapChannel)
	go heapsort(&waitGroup, arr4, heapChannel)

	waitGroup.Wait()

	// go TempGraficarQuick(quickChannel)
	// go QuickSort(&waitGroup, arr, quickChannel) //InsetionSort al segundo Array
	fmt.Print("ARRRAY \n \n \n ")
	fmt.Print(arr4)
	return c.String(http.StatusOK, "Simulation Begun")
}

//Normaliza el numero dado en el rango de 0 a 31
func NormalizeRandom(randomNum float64, upperbound int16, lowerbound int16) int {
	RandomNum := 0 + (int)(randomNum*(((float64)(upperbound)-(float64)(lowerbound))+1))
	return RandomNum
}

//Genera un número aleatorio en el intervalo de [0,1[
func generateRandom(wg *sync.WaitGroup, channel chan int, arrayChannel chan []int, seed int, size int) {
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
	arrayChannel <- randomArray
}

var arrGlob []int

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
	e.GET("/start", sortingSimulation)
	// Iniciamos el servidor
	e.Logger.Fatal(e.Start(":9000"))
}

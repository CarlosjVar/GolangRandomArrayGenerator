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
	pos1 int
	pos2 int
}

func InsertionSort(wg *sync.WaitGroup, randomArray []int, insertionChannel chan []int) {
	var num1 = 1
	var numAux = 0 //declaracion de variables
	var guardado = 0
	for num1 < len(randomArray) { //recorrida iteratiba sobre el array
		posiciones := []int{}
		posiciones = append(posiciones, num1)
		numAux = num1
		guardado = randomArray[num1]
		var cambio = false
		for numAux > 0 { //iteracion para devolver un indice hasta donde sea necesario en el array
			if guardado < randomArray[numAux-1] {
				randomArray[numAux] = randomArray[numAux-1] //si el numero anterior es mayor inicia cambio de posiciones
				cambio = true                               //condicional usado para saber si el numero debe seguir siendo movido o nunca sera necesario
			} else if cambio {
				randomArray[numAux] = guardado          // si en algun momento hubo algun cambio pero ya no es necesario
				posiciones = append(posiciones, numAux) // significa que es el lugar dode el numero debe permanecer
				insertionChannel <- posiciones
				<-insertionChannel
				break
			} else {
				numAux = -1 //de no haber habido un cambio nunca se cancela el ciclo para el numero actual
			}
			if numAux-1 == 0 {
				randomArray[numAux-1] = guardado
				posiciones = append(posiciones, numAux-1)
				insertionChannel <- posiciones
				<-insertionChannel
			}
			numAux--
		}
		num1++
	}
	posiciones := []int{}
	posiciones = append(posiciones, 0)
	posiciones = append(posiciones, 0)
	insertionChannel <- posiciones
	<-insertionChannel
	defer wg.Done()
}
func TempGraficarInsertion(insertionChannel chan []int) {
	posicCamb := []int{}
	for true {
		posicCamb = <-insertionChannel
		if posicCamb[0] == 0 && posicCamb[1] == 0 {
			insertionChannel <- posicCamb
			break
		}
		fmt.Print(posicCamb[0])
		fmt.Print(" ")
		fmt.Print(posicCamb[1]) //Prueba, en el futuro va a ser el que obtiene para graficar
		fmt.Println(" " + "InsertionSort")
		client.Trigger("ArrayChannel", "insertion", posicCamb)
		insertionChannel <- posicCamb
	}
}

// Función encargade de realizar el quicksort
func QuickSort(wg *sync.WaitGroup, array []int, quickChannel chan []int) { //metodo recursivo con pivote al final
	if len(array) < 2 {
		//fmt.Println("Fin")		//condicion de salida
		posiciones := []int{}
		posiciones = append(posiciones, 0)
		posiciones = append(posiciones, 0)
		quickChannel <- posiciones
		<-quickChannel
		return
	}

	//fmt.Println("Entra")
	ladoIzq := 0
	ladoDer := len(array) - 1 //declaracion de variables
	auxiliar := 0
	posiciones := []int{}

	for i, _ := range array {
		if array[i] < array[ladoDer] { //Si el numero en el indice actual en menor al pivote
			auxiliar = array[ladoIzq] //se mueve el numero del indice al final del subarray
			array[ladoIzq] = array[i] //de menores a la izquierda, y se actualiza la variable
			array[i] = auxiliar       //representando el final del mismo
			posiciones = append(posiciones, 0)
			posiciones = append(posiciones, 0)
			quickChannel <- posiciones
			<-quickChannel
			ladoIzq++

		}
	}

	auxiliar = array[ladoIzq]
	array[ladoIzq] = array[ladoDer] //El pivote se mueve al final del subarray de la izquierda
	array[ladoDer] = auxiliar
	posiciones = append(posiciones, 0)
	posiciones = append(posiciones, 0)
	quickChannel <- posiciones
	<-quickChannel
	QuickSort(wg, array[:ladoIzq], quickChannel) //Llamadas recursivas para ambos subarrays
	QuickSort(wg, array[ladoIzq+1:], quickChannel)
}
func BubbleSort(wg *sync.WaitGroup, randomArray []int, bubbleArray chan []int) {
	fmt.Println("BubbleSort\n")
	for true { //ciclo que atravieza el array multiples veces hasta no necesitar mas cambios
		fmt.Print(randomArray)
		var num1 = 0 //declaracion de variables
		var num2 = 1
		var numAux = 0
		var cambio = false
		for num1 < len(randomArray) && num2 < len(randomArray) { //ciclo de recorrida actual del array
			posiciones := []int{}
			if randomArray[num1] > randomArray[num2] {
				numAux = randomArray[num1]            //si el numero actual es x es mayor al numero x+1
				randomArray[num1] = randomArray[num2] //intercambian lugares
				randomArray[num2] = numAux
				cambio = true //se indica que en la recorrida actual hubo cambios
				posiciones = append(posiciones, num1)
				posiciones = append(posiciones, num2)
				bubbleArray <- posiciones
				<-bubbleArray
			}
			num1++
			num2++
		}
		if cambio == false { //si en la recorrida no hubo cambios se rompe el ciclo
			break
		}
	}
	posiciones := []int{}
	posiciones = append(posiciones, 0)
	posiciones = append(posiciones, 0)
	bubbleArray <- posiciones
	<-bubbleArray
	defer wg.Done()
	//arrChan <- randomArray
}
func TempGraficarBubble(bubbleChannel chan []int) {
	posicCamb := []int{}
	for true {
		posicCamb = <-bubbleChannel
		if posicCamb[0] == 0 && posicCamb[1] == 0 {
			bubbleChannel <- posicCamb
			break
		}
		fmt.Print(posicCamb[0])
		fmt.Print(" ")
		fmt.Print(posicCamb[1]) //Prueba, en el futuro va a ser el que obtiene para graficar
		fmt.Println(" " + "BubbleSort")
		client.Trigger("ArrayChannel", "bubble", posicCamb)
		bubbleChannel <- posicCamb
	}
}

func TempGraficarQuick(quickChannel chan []int) {
	posicCamb := []int{}
	for true {
		posicCamb = <-quickChannel
		if posicCamb[0] == 0 && posicCamb[1] == 0 {
			quickChannel <- posicCamb
			break
		}

		fmt.Print(posicCamb[0])
		fmt.Print(" ")
		fmt.Print(posicCamb[1]) //Prueba, en el futuro va a ser el que obtiene para graficar
		fmt.Println(" " + "QuickSort")
		s := datos{pos1: posicCamb[0], pos2: posicCamb[1]}
		client.Trigger("ArrayChannel", "insertion", s)
		quickChannel <- posicCamb
	}
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
	arr2 := []int{}
	// arr3 := []int{}
	// arr4 := []int{}
	fmt.Println("Unsorted Arrays")
	arr1 = <-arrayChannel  //saca el array del channel
	arr2 = CopyArray(arr1) //copia el array

	// arr3 = CopyArray(arr2) //copia el array

	// arr4 = CopyArray(arr3) //copia el array

	// quickChannel := make(chan []int, 1)
	insertionChannel := make(chan []int, 1)
	bubbleChannel := make(chan []int, 1)
	waitGroup.Add(2)
	go TempGraficarInsertion(insertionChannel)
	go InsertionSort(&waitGroup, arr1, insertionChannel) //InsetionSort al segundo Array
	go TempGraficarBubble(bubbleChannel)
	go BubbleSort(&waitGroup, arr2, bubbleChannel) //BubbleSort al primer Array
	waitGroup.Wait()

	// go TempGraficarQuick(quickChannel)
	// go QuickSort(&waitGroup, arr, quickChannel) //InsetionSort al segundo Array
	client.Trigger("ArrayChannel", "a", "a")
	return c.String(http.StatusOK, "Simulation begun")
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
	e.GET("/start", start)
	// Iniciamos el servidor
	e.Logger.Fatal(e.Start(":9000"))
}

package main

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

func NormalizeRandom(randomNum float64) int {
	RandomNum := 0 + (int)(randomNum*((31-0)+1))
	return RandomNum
}

func generateRandom(wg *sync.WaitGroup, channel chan int, seed int, arrayChannel chan []int, size int) {
	randomArray := []int{}
	period := 8192
	constant := 23
	multiplicativeConstant := 3 + (8 * constant)
	for i := 0; i < size; i++ {
		num := ((seed * multiplicativeConstant) + constant) % int(period)
		seed = num
		normalizedNum := float64(num) / float64(period-1)
		randomNum := NormalizeRandom(normalizedNum)
		randomArray = append(randomArray, randomNum)
	}
	fmt.Printf("%v", randomArray)
	arrayChannel <- randomArray
	defer wg.Done()
}
func BubbleSort(wg *sync.WaitGroup, arrChan chan []int, controller chan int) {
	randomArray := []int{}
	randomArray = <-arrChan
	fmt.Println("BubbleSort")
	for true {
		var num1 = 0
		var num2 = 1
		var numAux = 0
		var cambio = false
		for num1 < len(randomArray) && num2 < len(randomArray) {
			if randomArray[num1] > randomArray[num2] {
				<-controller
				numAux = randomArray[num1]
				randomArray[num1] = randomArray[num2]
				randomArray[num2] = numAux
				TempObtenerIndices(num1, num2, "BubbleSort")
				cambio = true
				controller <- 0
			}
			num1++
			num2++
		}
		if cambio == false {
			break
		}
	}
	defer wg.Done()
	//arrChan <- randomArray
}

func InsertionSort(wg *sync.WaitGroup, arrChan chan []int, controller chan int) {
	var num1 = 1
	var numAux = 0
	var numGua = 0
	var guardado = 0
	randomArray := []int{}
	randomArray = <-arrChan
	fmt.Println("Insertion Sort")
	for num1 < len(randomArray) {
		numAux = num1
		guardado = randomArray[num1]
		var cambio = false
		for numAux > 0 {
			if guardado < randomArray[numAux-1] {
				randomArray[numAux] = randomArray[numAux-1]
				numGua = numAux
				cambio = true
			} else if cambio {
				controller <- 0
				randomArray[numAux] = guardado
				TempObtenerIndices(numGua, numAux, "InsertSort")
				<-controller
				break
			} else {
				numAux = -1
			}
			if numAux-1 == 0 {
				randomArray[numAux-1] = guardado
			}
			numAux--
		}
		num1++
	}
	defer wg.Done()
	//arrChan <- randomArray
}

func TempObtenerIndices(indice1 int, indice2 int, ordenamiento string) {
	fmt.Print(indice1)
	fmt.Print(" ")
	fmt.Print(indice2) //Prueba, en el futuro va a ser el que obtiene para graficar
	fmt.Println(" " + ordenamiento)
}

func CopyArray(arrOri []int) []int {
	nuevoArray := []int{}
	var num = 0
	for i := 0; i < len(arrOri); i++ {
		num = arrOri[i]
		nuevoArray = append(nuevoArray, num)
	}
	return nuevoArray
}

func main() {

	var seed int64
	for true {
		fmt.Println("Ingrese una semilla prima en el intervalo de 11 a 101")
		fmt.Scan(&seed)
		if (seed >= 11) && (seed <= 101) {

			if big.NewInt(seed).ProbablyPrime(0) {
				break

			} else {

				fmt.Println("Debe ingresar un número primo")
			}

		} else {

			fmt.Println("Debe ingresar un número entre 11 y 101")
		}
	}
	var size int
	for true {
		fmt.Println("Ingrese cuantos números desea [10,10000]")
		fmt.Scan(&size)
		if size >= 10 && seed <= 10000 {
			break
		} else {
			fmt.Println("El número debe estar en el intervalo de [10,10000]")
		}
	}
	randomch := make(chan int, 1)
	arrayChannel := make(chan []int, 1)

	arr := []int{}
	arr2 := []int{}

	var waitGroup sync.WaitGroup
	go generateRandom(&waitGroup, randomch, int(seed), arrayChannel, size) //genera el array
	waitGroup.Add(1)
	waitGroup.Wait()

	arr = <-arrayChannel //saca el array del channel
	fmt.Println(arr)
	arr2 = CopyArray(arr) //copia el array
	fmt.Println(arr2)

	go BubbleSort(&waitGroup, arrayChannel, randomch) //se ordena el primer array
	go InsertionSort(&waitGroup, arrayChannel, randomch)
	arrayChannel <- arr //se inserta el array al channel
	<-arrayChannel
	arrayChannel <- arr2 //se inserta el array al channel
	time.Sleep(3 * time.Second)

	fmt.Println(arr)
	fmt.Println(arr2)
	fmt.Println("Terminó") //verificacion
}

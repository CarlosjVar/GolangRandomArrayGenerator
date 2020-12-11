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

func createArray(wg *sync.WaitGroup, channel chan int, arrChan chan []int, size int) {
	randomArray := []int{}
	randomNum := 0
	for i := 0; i < size; i++ {
		randomNum = <-channel
		randomArray = append(randomArray, randomNum)
		//fmt.Printf("%v", randomArray)
		channel <- 0

	}
	arrChan <- randomArray
	arrChan <- randomArray
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
	go generateRandom(&waitGroup, randomch, int(seed), size) //genera el array
	go createArray(&waitGroup, randomch, arrayChannel, size)
	//defer &waitGroup.Done()

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
	//arr = <-arrayChannel  //se saca el array ordenado del channel
	//arr2 = <-arrayChannel //se saca el array ordenado del channel
	//<-arrayChannel        //vacia el channel

	//waitGroup.Add(1)
	//waitGroup.Wait()

	fmt.Println(arr)
	fmt.Println(arr2)
	fmt.Println("Terminó") //verificacion

	// arr := []int{16, 95, 12, 29, 40, 25, 46, 29, 13, 86}
	//arr2 := []int{16, 95, 12, 29, 40, 25, 46, 29, 13, 86}
	// arr := []int{16, 95, 12, 29, 40, 25, 46, 29, 13, 86, 12, 71, 11, 41, 10, 31, 57, 23, 87, 90, 48, 100, 43, 80, 42, 50, 86, 40, 61, 48, 73, 95, 69, 14, 78, 87, 14, 89, 91, 49, 91, 87, 36, 62, 8, 57, 46, 49, 9, 32, 98, 96, 42, 52, 61, 35, 97, 69, 51, 63, 87, 10, 55, 95, 79, 69, 45, 30, 23, 58, 58, 61, 77, 31, 37, 94, 47, 8, 41, 68, 89, 89, 45, 35, 1, 16, 18, 16, 65, 100, 58, 9, 74, 57, 27, 50, 92, 100, 19, 15}
	// arr2 := []int{16, 95, 12, 29, 40, 25, 46, 29, 13, 86, 12, 71, 11, 41, 10, 31, 57, 23, 87, 90, 48, 100, 43, 80, 42, 50, 86, 40, 61, 48, 73, 95, 69, 14, 78, 87, 14, 89, 91, 49, 91, 87, 36, 62, 8, 57, 46, 49, 9, 32, 98, 96, 42, 52, 61, 35, 97, 69, 51, 63, 87, 10, 55, 95, 79, 69, 45, 30, 23, 58, 58, 61, 77, 31, 37, 94, 47, 8, 41, 68, 89, 89, 45, 35, 1, 16, 18, 16, 65, 100, 58, 9, 74, 57, 27, 50, 92, 100, 19, 15}
	// BubbleSort(arr)
	// InsertionSort(arr2)
	// fmt.Println(arr)
	// fmt.Println(arr2)

}

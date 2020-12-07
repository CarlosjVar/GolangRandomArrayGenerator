package main

import (
	"fmt"
	"math/big"
	"sync"
)

func BubbleSort(randomArray []int) {
	for true {
		var num1 = 0
		var num2 = 1
		var numAux = 0
		var cambio = false
		for num1 < len(randomArray) && num2 < len(randomArray) {
			if randomArray[num1] > randomArray[num2] {
				numAux = randomArray[num1]
				randomArray[num1] = randomArray[num2]
				randomArray[num2] = numAux
				cambio = true
			}
			num1++
			num2++
		}
		if cambio == false {
			break
		}
	}
}

func InsertionSort(randomArray []int) {
	var num1 = 1
	var numAux = 0
	var guardado = 0
	for num1 < len(randomArray) {
		numAux = num1
		guardado = randomArray[num1]
		var cambio = false
		for numAux > 0 {
			if guardado < randomArray[numAux-1] {
				randomArray[numAux] = randomArray[numAux-1]
				cambio = true
			} else if cambio {
				randomArray[numAux] = guardado
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
}

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
func createArray(wg *sync.WaitGroup, channel chan int, size int) {
	randomArray := []int{}
	randomNum := 0
	for i := 0; i < size; i++ {
		randomNum = <-channel
		randomArray = append(randomArray, randomNum)
		fmt.Printf("%v", randomArray)
		channel <- 0

	}

	defer wg.Done()

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

	var waitGroup sync.WaitGroup
	go generateRandom(&waitGroup, randomch, int(seed), size)
	go createArray(&waitGroup, randomch, size)
	waitGroup.Add(1)
	waitGroup.Wait()
	fmt.Println("Terminó")

	//arr := []int{16, 95, 12, 29, 40, 25, 46, 29, 13, 86}
	//arr2 := []int{16, 95, 12, 29, 40, 25, 46, 29, 13, 86}
	// arr := []int{16, 95, 12, 29, 40, 25, 46, 29, 13, 86, 12, 71, 11, 41, 10, 31, 57, 23, 87, 90, 48, 100, 43, 80, 42, 50, 86, 40, 61, 48, 73, 95, 69, 14, 78, 87, 14, 89, 91, 49, 91, 87, 36, 62, 8, 57, 46, 49, 9, 32, 98, 96, 42, 52, 61, 35, 97, 69, 51, 63, 87, 10, 55, 95, 79, 69, 45, 30, 23, 58, 58, 61, 77, 31, 37, 94, 47, 8, 41, 68, 89, 89, 45, 35, 1, 16, 18, 16, 65, 100, 58, 9, 74, 57, 27, 50, 92, 100, 19, 15}
	// arr2 := []int{16, 95, 12, 29, 40, 25, 46, 29, 13, 86, 12, 71, 11, 41, 10, 31, 57, 23, 87, 90, 48, 100, 43, 80, 42, 50, 86, 40, 61, 48, 73, 95, 69, 14, 78, 87, 14, 89, 91, 49, 91, 87, 36, 62, 8, 57, 46, 49, 9, 32, 98, 96, 42, 52, 61, 35, 97, 69, 51, 63, 87, 10, 55, 95, 79, 69, 45, 30, 23, 58, 58, 61, 77, 31, 37, 94, 47, 8, 41, 68, 89, 89, 45, 35, 1, 16, 18, 16, 65, 100, 58, 9, 74, 57, 27, 50, 92, 100, 19, 15}
	// BubbleSort(arr)
	// InsertionSort(arr2)
	// fmt.Println(arr)
	// fmt.Println(arr2)

}

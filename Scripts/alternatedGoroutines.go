package main

import (
	"fmt"
	"math/big"
	"sync"
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
	randomch := make(chan int, 3)

	var waitGroup sync.WaitGroup
	go generateRandom(&waitGroup, randomch, int(seed), size)
	go createArray(&waitGroup, randomch, size)
	waitGroup.Add(1)
	waitGroup.Wait()
	fmt.Println("Terminó")

}

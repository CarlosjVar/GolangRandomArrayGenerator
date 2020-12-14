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
	//fmt.Printf("%v", randomArray)
	arrayChannel <- randomArray
	//defer wg.Done()
}
func BubbleSort(wg *sync.WaitGroup, randomArray []int, controller chan int) {
	fmt.Println("BubbleSort")
	for true { //ciclo que atravieza el array multiples veces hasta no necesitar mas cambios
		var num1 = 0 //declaracion de variables
		var num2 = 1
		var numAux = 0
		var cambio = false
		for num1 < len(randomArray) && num2 < len(randomArray) { //ciclo de recorrida actual del array
			if randomArray[num1] > randomArray[num2] {
				//<-controller
				numAux = randomArray[num1]            //si el numero actual es x es mayor al numero x+1
				randomArray[num1] = randomArray[num2] //intercambian lugares
				randomArray[num2] = numAux
				TempObtenerIndices(num1, num2, "BubbleSort")
				cambio = true //se indica que en la recorrida actual hubo cambios
				//controller <- 0
			}
			num1++
			num2++
		}
		if cambio == false { //si en la recorrida no hubo cambios se rompe el ciclo
			break
		}
	}
	defer wg.Done()
	//arrChan <- randomArray
}

func InsertionSort(wg *sync.WaitGroup, randomArray []int, controller chan int) {
	var num1 = 1
	var numAux = 0
	var numGua = 0 //declaracion de variables
	var guardado = 0
	fmt.Println("Insertion Sort")
	for num1 < len(randomArray) { //recorrida iteratiba sobre el array
		numAux = num1
		guardado = randomArray[num1]
		var cambio = false
		for numAux > 0 { //iteracion iterativa para devolver un indice hasta donde sea necesario en el array
			if guardado < randomArray[numAux-1] {
				randomArray[numAux] = randomArray[numAux-1] //si el numero anterior es mayor inicia cambio de posiciones
				numGua = numAux
				cambio = true //condicional usado para saber si el numero debe seguir siendo movido o nunca sera necesario
			} else if cambio {
				//controller <- 0
				randomArray[numAux] = guardado
				TempObtenerIndices(numGua, numAux, "InsertSort") // si en algun momento hubo algun cambio pero ya no es necesario
				//<-controller										// significa que es el lugar dode el numero debe permanecer
				break
			} else {
				numAux = -1 //de no haber habido un cambio nunca se cancela el ciclo para el numero actual
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

func QuickSort(array []int, controller chan int) { //metodo recursivo con pivote al final
	if len(array) < 2 {
		//fmt.Println("Fin")		//condicion de salida
		return
	}

	//fmt.Println("Entra")
	ladoIzq := 0
	ladoDer := len(array) - 1 //declaracion de variables
	auxiliar := 0

	for i, _ := range array {
		if array[i] < array[ladoDer] { //Si el numero en el indice actual en menor al pivote
			auxiliar = array[ladoIzq] //se mueve el numero del indice al final del subarray
			array[ladoIzq] = array[i] //de menores a la izquierda, y se actualiza la variable
			array[i] = auxiliar       //representando el final del mismo
			TempObtenerIndices(i, ladoIzq, "QuickSort")
			ladoIzq++

		}
	}

	auxiliar = array[ladoIzq]
	array[ladoIzq] = array[ladoDer] //El pivote se mueve al final del subarray de la izquierda
	array[ladoDer] = auxiliar
	TempObtenerIndices(ladoIzq, ladoDer, "QuickSort")

	QuickSort(array[:ladoIzq], controller) //Llamadas recursivas para ambos subarrays
	QuickSort(array[ladoIzq+1:], controller)
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
	for i := 0; i < len(arrOri); i++ { //copiado de Array
		num = arrOri[i]
		nuevoArray = append(nuevoArray, num)
	}
	return nuevoArray
}

func main() {

	var seed int64
	for true {
		fmt.Println("Ingrese una semilla prima en el intervalo de 11 a 101") //Menu para parametros a eleccion del usuario
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
	randomch := make(chan int, 1) //creacion de canales para fturo uso
	arrayChannel := make(chan []int, 1)

	arr := []int{}
	arr2 := []int{}
	arr3 := []int{}

	var waitGroup sync.WaitGroup
	go generateRandom(&waitGroup, randomch, int(seed), arrayChannel, size) //genera el array

	arr = <-arrayChannel //saca el array del channel
	fmt.Println(arr)
	arr2 = CopyArray(arr) //copia el array
	fmt.Println(arr2)
	arr3 = CopyArray(arr2) //copia el array
	fmt.Println(arr3)

	go BubbleSort(&waitGroup, arr, randomch)     //BubbleSort al primer Array
	go InsertionSort(&waitGroup, arr2, randomch) //InsetionSort al segundo Array

	waitGroup.Add(2)
	waitGroup.Wait()

	QuickSort(arr3, randomch) //Quicksort sin corrutinas al tercer Array

	fmt.Println(arr)
	fmt.Println(arr2)
	fmt.Println(arr3)
	fmt.Println("Terminó") //verificacion
}

package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func InsertionSort(wg *sync.WaitGroup, randomArray []int, insertionChannel chan []int) {
	var tiempoIni = time.Now()
	var tiempoFinal = time.Now()
	var num1 = 1
	var numAux = 0 //declaracion de variables
	var guardado = 0
	var ciclos = 0
	var condicionales = 0
	var intercambios = 0
	for num1 < len(randomArray) { //recorrida iteratiba sobre el array
		ciclos++
		posiciones := []int{}
		posiciones = append(posiciones, num1)
		numAux = num1
		guardado = randomArray[num1]
		var cambio = false
		for numAux > 0 { //iteracion para devolver un indice hasta donde sea necesario en el array
			ciclos++
			if guardado < randomArray[numAux-1] {
				condicionales++
				randomArray[numAux] = randomArray[numAux-1] //si el numero anterior es mayor inicia cambio de posiciones
				cambio = true                               //condicional usado para saber si el numero debe seguir siendo movido o nunca sera necesario
				intercambios++
			} else if cambio {
				condicionales++
				intercambios++
				randomArray[numAux] = guardado          // si en algun momento hubo algun cambio pero ya no es necesario
				posiciones = append(posiciones, numAux) // significa que es el lugar dode el numero debe permanecer
				insertionChannel <- posiciones
				<-insertionChannel
				break
			} else {
				condicionales++
				numAux = -1 //de no haber habido un cambio nunca se cancela el ciclo para el numero actual
			}
			if numAux-1 == 0 {
				condicionales++
				intercambios++
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
	tiempoFinal = time.Now()
	posiciones = append(posiciones, 0)
	posiciones = append(posiciones, 0)
	dato := []string{}
	dato = append(dato, strconv.Itoa(intercambios), strconv.Itoa(condicionales), strconv.Itoa(ciclos), tiempoIni.Format("2006-01-02 15:04:05"), tiempoFinal.Format("2006-01-02 15:04:05"))

	client.Trigger("ArrayChannel", "insertionStats", dato)
	insertionChannel <- posiciones
	<-insertionChannel
	defer wg.Done()
	//arrChan <- randomArray
}
func TempGraficarInsertion(insertionChannel chan []int) {
	posicCamb := []int{}
	for true {
		posicCamb = <-insertionChannel
		if posicCamb[0] == 0 && posicCamb[1] == 0 {
			insertionChannel <- posicCamb
			break
		}
		time.Sleep(20 * time.Millisecond)
		fmt.Print(posicCamb[0])
		fmt.Print(" ")
		fmt.Print(posicCamb[1]) //Prueba, en el futuro va a ser el que obtiene para graficar
		fmt.Println(" " + "InsertionSort")
		client.Trigger("ArrayChannel", "insertion", posicCamb)
		insertionChannel <- posicCamb
	}
}

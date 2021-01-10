package main

import (
	"strconv"
	"sync"
	"time"
)

// Función encargade de realizar el quicksort
func QuickSort(wg *sync.WaitGroup, array []int, quickChannel chan []int, ultimoIndice int, wgEnd int, ciclos *int, condicionales *int, intercambios *int) { //metodo recursivo con pivote al final
	var tiempoIni = time.Now()
	print(tiempoIni.Format("2006-01-02 15:04:05"))
	if len(array) < 2 {
		*condicionales++
		//fmt.Println("Fin")		//condicion de salida
		return
	}
	wgEnd = wgEnd + 1
	//fmt.Println("Entra")
	ladoIzq := 0
	ladoDer := len(array) - 1 //declaracion de variables
	auxiliar := 0
	posiciones := []int{}

	for i, _ := range array {
		*ciclos++
		if array[i] < array[ladoDer] { //Si el numero en el indice actual en menor al pivote
			*condicionales++
			auxiliar = array[ladoIzq] //se mueve el numero del indice al final del subarray
			array[ladoIzq] = array[i] //de menores a la izquierda, y se actualiza la variable
			array[i] = auxiliar       //representando el final del mismo
			if (i + ultimoIndice) != (ladoIzq + ultimoIndice) {
				*intercambios++
				posiciones = append(posiciones, i+ultimoIndice)
				posiciones = append(posiciones, ladoIzq+ultimoIndice)
				quickChannel <- posiciones
				<-quickChannel
				time.Sleep(50 * time.Millisecond)
				posiciones = nil
			}
			ladoIzq++
		}
	}

	auxiliar = array[ladoIzq]
	array[ladoIzq] = array[ladoDer] //El pivote se mueve al final del subarray de la izquierda
	array[ladoDer] = auxiliar
	*intercambios++
	posiciones = append(posiciones, ladoIzq+ultimoIndice)
	posiciones = append(posiciones, ladoDer+ultimoIndice)
	quickChannel <- posiciones
	<-quickChannel
	time.Sleep(50 * time.Millisecond)
	posiciones = nil
	QuickSort(wg, array[:ladoIzq], quickChannel, ultimoIndice, wgEnd, ciclos, condicionales, intercambios) //Llamadas recursivas para ambos subarrays
	QuickSort(wg, array[ladoIzq+1:], quickChannel, ultimoIndice+ladoIzq+1, wgEnd, ciclos, condicionales, intercambios)
	if wgEnd == 1 {
		*condicionales++
		var tiempoFinal = time.Now()
		dato := []string{}
		dato = append(dato, strconv.Itoa(*intercambios), strconv.Itoa(*condicionales), strconv.Itoa(*ciclos), tiempoIni.Format("2006-01-02 15:04:05"), tiempoFinal.Format("2006-01-02 15:04:05"))

		client.Trigger("ArrayChannel", "quickStats", dato)
		defer wg.Done()
	}
}

func TempGraficarQuick(quickChannel chan []int) {
	posicCamb := []int{}
	for true {
		posicCamb = <-quickChannel
		// if posicCamb[0] == 0 && posicCamb[1] == 0 {
		// 	quickChannel <- posicCamb
		// 	break
		// }
		client.Trigger("ArrayChannel", "quick", posicCamb)
		quickChannel <- posicCamb
	}
}

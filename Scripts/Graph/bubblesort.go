package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func BubbleSort(wg *sync.WaitGroup, randomArray []int, bubbleArray chan []int) {
	var tiempoIni = time.Now()
	var tiempoFinal = time.Now()
	var ciclos = 0
	var condicionales = 0
	var intercambios = 0
	fmt.Println("BubbleSort\n")
	for true { //ciclo que atravieza el array multiples veces hasta no necesitar mas cambios
		ciclos++
		fmt.Print(randomArray)
		var num1 = 0 //declaracion de variables
		var num2 = 1
		var numAux = 0
		var cambio = false
		for num1 < len(randomArray) && num2 < len(randomArray) { //ciclo de recorrida actual del array
			ciclos++
			posiciones := []int{}
			if randomArray[num1] > randomArray[num2] {
				condicionales++
				numAux = randomArray[num1]            //si el numero actual es x es mayor al numero x+1
				randomArray[num1] = randomArray[num2] //intercambian lugares
				randomArray[num2] = numAux
				intercambios++
				cambio = true //se indica que en la recorrida actual hubo cambios
				posiciones = append(posiciones, num1)
				posiciones = append(posiciones, num2)
				bubbleArray <- posiciones
				<-bubbleArray
			}
			num1++
			num2++
		}
		if cambio == false {
			condicionales++ //si en la recorrida no hubo cambios se rompe el ciclo
			break
		}

	}
	posiciones := []int{}
	posiciones = append(posiciones, 0)
	posiciones = append(posiciones, 0)

	bubbleArray <- posiciones
	<-bubbleArray
	tiempoFinal = time.Now()
	dato := []string{}
	dato = append(dato, strconv.Itoa(intercambios), strconv.Itoa(condicionales), strconv.Itoa(ciclos), tiempoIni.Format("2006-01-02 15:04:05"), tiempoFinal.Format("2006-01-02 15:04:05"))

	client.Trigger("ArrayChannel", "bubbleStats", dato)

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

		client.Trigger("ArrayChannel", "bubble", posicCamb)
		bubbleChannel <- posicCamb
	}
}

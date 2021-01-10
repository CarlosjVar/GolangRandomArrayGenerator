package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Estructura del heap
type maxheap struct {
	arr []int
}

// Crea una estructura de heap
func newMaxHeap(arr []int) *maxheap {
	maxheap := &maxheap{
		arr: arr,
	}
	return maxheap
}

// Retorna el índice del hijo izquierdo de un nodo
func (m *maxheap) indiceIzquierdo(index int) int {
	return 2*index + 1
}

// Retorna el índice del hijo derecho de un nodo
func (m *maxheap) indiceDerecho(index int) int {
	return 2*index + 2
}

// Intercambia 2 elementos de un array entre si
func (m *maxheap) swap(primero, segundo int) {
	m.arr[primero], m.arr[segundo] = m.arr[segundo], m.arr[primero]
}

func (m *maxheap) leaf(index int, lenght int) bool {
	if index >= (lenght/2) && index <= lenght {
		return true
	}
	return false
}

// Se encarga de hacer un heap a partir de un array , esto comparando la raíz con sus hijos , si alguno de ellos es más grande que la raíz se intercambia su posición
func (m *maxheap) heapify(current int, lenght int, heapChannel chan []int, intercambios *int, condicionales *int, ciclos *int) {
	*condicionales++
	if m.leaf(current, lenght) {
		return
	}
	mayor := current
	indiceIzquierdo := m.indiceIzquierdo(current)
	rightRightIndex := m.indiceDerecho(current)
	*condicionales++
	if indiceIzquierdo < lenght && m.arr[indiceIzquierdo] > m.arr[mayor] {
		mayor = indiceIzquierdo
	}
	*condicionales++
	if rightRightIndex < lenght && m.arr[rightRightIndex] > m.arr[mayor] {
		mayor = rightRightIndex
	}
	*condicionales++
	if mayor != current {
		*intercambios++
		m.swap(current, mayor)
		posiciones := []int{}
		posiciones = append(posiciones, current)
		posiciones = append(posiciones, mayor)
		heapChannel <- posiciones
		<-heapChannel
		time.Sleep(12 * time.Millisecond)
		m.heapify(mayor, lenght, heapChannel, intercambios, condicionales, ciclos)
	}
	return
}

// Recibe un array y lo convierte en un maxHeap
func (m *maxheap) buildMaxHeap(lenght int, heapChannel chan []int, intercambios *int, condicionales *int, ciclos *int) {
	for index := ((lenght / 2) - 1); index >= 0; index-- {
		*ciclos++
		m.heapify(index, lenght, heapChannel, intercambios, condicionales, ciclos)
	}
}

// Sortea el max heap , esto mediante el método de tomar la raíz , que es el mayor y lo envía al final
func (m *maxheap) sort(lenght int, start time.Time, wg *sync.WaitGroup, heapChannel chan []int, intercambios *int, condicionales *int, ciclos *int) {
	m.buildMaxHeap(lenght, heapChannel, intercambios, condicionales, ciclos)
	for i := lenght - 1; i > 0; i-- {
		*ciclos++
		// Obtiene la raíz y la "elimina " del heap
		posiciones := []int{}
		m.swap(0, i)
		*intercambios++
		posiciones = append(posiciones, 0)
		posiciones = append(posiciones, i)
		heapChannel <- posiciones
		<-heapChannel
		time.Sleep(12 * time.Millisecond)
		m.heapify(0, i, heapChannel, intercambios, condicionales, ciclos)
	}

	defer wg.Done()

	//fmt.Print(time.Since(start))
}

// Función para imprimir el heap
func (m *maxheap) print() {
	for _, val := range m.arr {
		fmt.Println(val)
	}
}

// Función que se encarga del proceso
func heapsort(wg *sync.WaitGroup, array []int, heapChannel chan []int) {
	intercambios := 0
	condicionales := 0
	ciclos := 0
	start := time.Now()
	minHeap := newMaxHeap(array)
	minHeap.sort(len(array), start, wg, heapChannel, &intercambios, &condicionales, &ciclos)
	finish := time.Now()

	dato := []string{}
	dato = append(dato, strconv.Itoa(intercambios), strconv.Itoa(condicionales), strconv.Itoa(ciclos), start.Format("2006-01-02 15:04:05"), finish.Format("2006-01-02 15:04:05"))
	client.Trigger("ArrayChannel", "heapStats", dato)
	// minHeap.print()
}

//Función graficadora del heapsort
func TempGraficarHeap(heapChannel chan []int) {
	posicCamb := []int{}
	for true {
		posicCamb = <-heapChannel
		if posicCamb[0] == 0 && posicCamb[1] == 0 {
			heapChannel <- posicCamb
			break
		}
		fmt.Print(posicCamb[0])
		fmt.Print(" ")
		fmt.Print(posicCamb[1]) //Prueba, en el futuro va a ser el que obtiene para graficar
		fmt.Println(" " + "HeapSort")
		client.Trigger("ArrayChannel", "heap", posicCamb)
		heapChannel <- posicCamb
	}
}

package main

/*
import (
	"fmt"
	"math/rand"
)

func routine(moncanal chan int, synchro chan int, phrase string) {
	for true {
		<-moncanal
		fmt.Println(phrase)
		synchro <- 0
	}
}

func ordonnanceur(canaux []chan int, synchro chan int, fin chan int) {
	mcanaux := canaux
	for j := 0; j < 100; j++ {
		rand.Shuffle(5, func(i, j int) { mcanaux[i], mcanaux[j] = mcanaux[j], mcanaux[i] })
		for i := 0; i < 5; i++ {
			mcanaux[i] <- 0
			<-synchro
		}
	}
	fin <- 0
}

func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	c3 := make(chan int)
	c4 := make(chan int)
	c5 := make(chan int)
	cs := make(chan int)
	cx := []chan int{c1, c2, c3, c4, c5}
	fin := make(chan int)
	go func() { routine(c1, cs, "belle Marquise") }()
	go func() { routine(c2, cs, "vos beaux yeux") }()
	go func() { routine(c3, cs, "me font") }()
	go func() { routine(c4, cs, "mourir") }()
	go func() { routine(c5, cs, "d'amour") }()
	go func() { ordonnanceur(cx, cs, fin) }()
	<-fin
}
*/
/*

package main

import (
	"fmt"
)

func routine(vazy chan int, synchro chan int, phrase string) {
	for true {
		<-vazy
		fmt.Println(phrase)
		synchro <- 0
	}
}

func ordonnanceur(vazy chan int, synchro chan int, fin chan int) {
	for j := 0; j < 100; j++ {
		for i := 0; i < 5; i++ {
			vazy <- 0
		}
		for i := 0; i < 5; i++ {
			<-synchro
		}
	}
	fin <- 0
}

func main() {
	v := make(chan int)
	s := make(chan int)
	fin := make(chan int)
	go func() { routine(v, s, "belle Marquise") }()
	go func() { routine(v, s, "vos beaux yeux") }()
	go func() { routine(v, s, "me font") }()
	go func() { routine(v, s, "mourir") }()
	go func() { routine(v, s, "d'amour") }()
	go func() { ordonnanceur(v, s, fin) }()
	<-fin
}

*/

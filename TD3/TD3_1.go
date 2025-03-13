package main

/*
import "fmt"

func routine(moncanal chan int, mot string) {
	for true {
		<-moncanal
		fmt.Println(mot)
		moncanal <- 0
	}
}

func main() {
	tour := make(chan int)
	go func() { routine(tour, "pong") }()
	go func() { routine(tour, "ping") }()
	tour <- 0
	for true {

	}
}
*/
/*
import "fmt"

func routine(moncanal chan int, soncanal chan int, mot string) {
	for true {
		<-moncanal
		fmt.Println(mot)
		soncanal <- 0
	}
}

func main() {
	ci := make(chan int)
	co := make(chan int)
	go func() { routine(co, ci, "pong") }()
	go func() { routine(ci, co, "ping") }()
	ci <- 0
	for true {

	}
}*/

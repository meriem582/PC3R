package main

import (
	"fmt"
	"time"
)

type Future struct { // Future encapsule le résultat d'un calcul et un canal pour signaler sa disponibilité.
	result int
	ready  chan struct{} // Canal fermé lorsque le résultat est prêt. veut dire que en vrai
}

func NewFuture() *Future { // NewFuture crée une nouvelle future.
	return &Future{
		ready: make(chan struct{}),
	}
}

// Poll permet de consulter le résultat de façon non bloquante. Si le résultat est prêt, il est retourné avec true, sinon false.
func (f *Future) Poll() (int, bool) {
	select {
	case <-f.ready:
		return f.result, true
	default:
		return 0, false
	}
}

// Wait bloque jusqu'à ce que le résultat soit prêt, puis le retourne.
func (f *Future) Wait() int {
	<-f.ready
	return f.result
}

// CalcRequest regroupe la fonction à exécuter et la future associée.
type CalcRequest struct {
	f   func() int
	fut *Future
}

// client envoie une requête de calcul au serveur et retourne la future associée.
func client(calcChan chan CalcRequest, f func() int) *Future {
	fut := NewFuture()
	req := CalcRequest{
		f:   f,
		fut: fut,
	}
	calcChan <- req
	return fut
}

// serveur lit les requêtes sur calcChan et lance leur exécution dans une goroutine.
func serveur(calcChan chan CalcRequest) {
	for req := range calcChan {
		go func(r CalcRequest) {
			res := r.f()       // Exécute la fonction de calcul
			r.fut.result = res // Stocke le résultat dans la future
			close(r.fut.ready) // Signale que le résultat est prêt
		}(req)
	}
}
func main() {
	calcChan := make(chan CalcRequest)   // Création du canal de requêtes de calcul.
	go serveur(calcChan)                 // Lancement du serveur.
	fut := client(calcChan, func() int { // Le client soumet un calcul (ici une fonction qui simule un délai).
		time.Sleep(2 * time.Second)
		return 42
	})

	// Le client effectue un polling toutes les secondes pour voir si le résultat est prêt.
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		if res, ready := fut.Poll(); ready {
			fmt.Printf("Poll: résultat prêt: %d\n", res)
			break
		} else {
			fmt.Println("Poll: résultat non prêt")
		}
	}

	// Ensuite, le client attend le résultat (méthode bloquante).
	res := fut.Wait()
	fmt.Printf("Wait: résultat reçu: %d\n", res)
}

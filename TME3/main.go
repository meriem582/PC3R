/* on s'est inspiré d'un depot git https://github.com/valeeraZ/Sorbonne_PC3R/blob/main/TME3/TME3.go*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func lecteur(canal chan string) {
	file, err := os.Open("./stop_times.txt")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture du fichier")
		return
	}
	/* on ferme le fichier à la fin de la fonction */
	defer file.Close()
	/* on crée un scanner pour lire le fichier */
	scanneur := bufio.NewScanner(file)
	/* on lit le fichier ligne par ligne */
	scanneur.Scan()
	for scanneur.Scan() {
		canal <- scanneur.Text()
	}
}

type paquet struct {
	depart string
	arrive string
	arret  int
}

type coupleAenvoyer struct {
	paquet              paquet
	canalpaquetResultat chan paquet
}

func worker(canal_lec chan string, canal_ser chan coupleAenvoyer, canal_red chan paquet) {
	for {
		ligne_de_donnees := <-canal_lec
		go func(str string) {
			arrive := strings.Split(str, ",")[1]
			depart := strings.Split(str, ",")[2]
			arretInit := 0
			/* conversion de données en un paquet pour l'envoie*/
			paquet0 := paquet{depart: depart, arrive: arrive, arret: arretInit}

			/* envoie sur le canal du seveur*/
			canalpaquetResultat := make(chan paquet)
			canal_ser <- coupleAenvoyer{paquet: paquet0, canalpaquetResultat: canalpaquetResultat}

			paquetResultat := <-canalpaquetResultat
			canal_red <- paquetResultat
		}(ligne_de_donnees)
	}
}

func diff(departArg string, arriveArg string) int {
	arrive, _ := time.Parse("00:00:00", departArg)
	depart, _ := time.Parse("00:00:00", arriveArg)
	diff := depart.Sub(arrive)
	return int(diff.Minutes())
}

func serveur_de_calcul(canal_ser chan coupleAenvoyer) {
	for true {
		couple := <-canal_ser
		go func(c coupleAenvoyer) {
			c.paquet.arret = diff(c.paquet.depart, c.paquet.arrive)
			c.canalpaquetResultat <- c.paquet
		}(couple)
	}
}

func redacteur(canal_red chan paquet, canal_sign_fin chan int) {
	s := 0
	compteur := 0
	for true {
		select {
		case paquet := <-canal_red:
			{
				compteur++
				s += paquet.arret
			}
		case <-canal_sign_fin:
			{
				canal_sign_fin <- s / compteur
				return
			}
		}
	}
}

func main() {
	/*fmt.Println("Hello, Go!")*/

	workers := 10
	canal_lec := make(chan string)
	canal_red := make(chan paquet)
	canal_ser := make(chan coupleAenvoyer)
	canal_principal := make(chan int)

	go func() { lecteur(canal_lec) }()

	for i := 0; i < workers; i++ {
		go func() { worker(canal_lec, canal_ser, canal_red) }()
	}

	go func() { serveur_de_calcul(canal_ser) }()

	go func() { redacteur(canal_red, canal_principal) }()

	time.Sleep(10 * time.Second)

	canal_principal <- 0
	total := <-canal_principal
	fmt.Println("Temps Total : ", total)

}

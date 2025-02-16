package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"

	st "tme4/client/structures"
	tr "tme4/serveur/travaux"
)

var idServ = make(map[int]*personne_serv)

var ADRESSE = "localhost"

var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "M"}

// type d'un paquet de personne stocke sur le serveur, n'implemente pas forcement personne_int (qui n'existe pas ici)
type personne_serv struct {
	// A FAIRE
	statut    string
	tabAFaire []func(st.Personne) st.Personne
	personne  st.Personne
}

// cree une nouvelle personne_serv, est appelé depuis le client, par le proxy, au moment ou un producteur distant
// produit une personne_dist
func creer(id int) *personne_serv {
	// A FAIRE
	p := pers_vide
	tabaFaire := make([]func(st.Personne) st.Personne, 0)
	nvp := personne_serv{statut: "v", tabAFaire: tabaFaire, personne: p}
	idServ[id] = &nvp
	return &nvp
}

// Méthodes sur les personne_serv, on peut recopier des méthodes des personne_emp du client
// l'initialisation peut être fait de maniere plus simple que sur le client
// (par exemple en initialisant toujours à la meme personne plutôt qu'en lisant un fichier)
func (p *personne_serv) initialise() {
	// A FAIRE
	p.personne = st.Personne{Nom: "Doe", Prenom: "John", Age: 42, Sexe: "M"}
	for i := 0; i <= rand.Intn(6); i++ {
		p.tabAFaire = append(p.tabAFaire, tr.UnTravail())
	}
	p.statut = "R"
}

func (p *personne_serv) travaille() {
	// A FAIRE
	p.personne = p.tabAFaire[0](p.personne)
	p.tabAFaire = p.tabAFaire[1:]
	if len(p.tabAFaire) == 0 {
		p.statut = "C"
	}
}

func (p *personne_serv) vers_string() string {
	// A FAIRE
	var sexe string
	if p.personne.Sexe == "M" {
		sexe = "Homme"
	} else {
		sexe = "Femme"
	}
	return fmt.Sprint("Nom: "+p.personne.Nom+" Prénom: "+p.personne.Prenom+" Age: ", p.personne.Age, " Sexe: "+sexe)
}

func (p *personne_serv) donne_statut() string {
	// A FAIRE
	return p.statut
}

// Goroutine qui maintient une table d'association entre identifiant et personne_serv
// il est contacté par les goroutine de gestion avec un nom de methode et un identifiant
// et il appelle la méthode correspondante de la personne_serv correspondante
func mainteneur(f string, id int, canalRetour chan string) {
	// A FAIRE
	if f == "creer" {
		creer(id)
		canalRetour <- "OK"
	} else if f == "initialise" {
		idServ[id].initialise()
		canalRetour <- "OK"
	} else if f == "travaille" {
		idServ[id].travaille()
		canalRetour <- "OK"
	} else if f == "vers_string" {
		canalRetour <- idServ[id].vers_string()
	} else if f == "donne_statut" {
		canalRetour <- idServ[id].donne_statut()
	} else {
		canalRetour <- "Methode inconnue"
	}

}

// Goroutine de gestion des connections
// elle attend sur la socketi un message content un nom de methode et un identifiant et appelle le mainteneur avec ces arguments
// elle recupere le resultat du mainteneur et l'envoie sur la socket, puis ferme la socket
func gere_connection(cnx net.Conn) {
	// A FAIRE
	for {
		m, _ := bufio.NewReader(cnx).ReadString('\n') // lit un message sur la socket
		request := strings.TrimSuffix(m, "\n")        // recupere la requete du client
		tab := strings.Split(request, ",")            // separe la requete en deux parties
		id, _ := strconv.Atoi(tab[0])                 // recupere l'id
		f := tab[1]                                   // recupere la methode
		canalRetour := make(chan string)              // cree un canal de retour
		go func() {
			mainteneur(f, id, canalRetour) // lance le mainteneur
		}()
		result := <-canalRetour          // recupere le resultat du mainteneur
		cnx.Write([]byte(result + "\n")) // envoie le resultat au client

	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Format: client <port>")
		return
	}
	port, _ := strconv.Atoi(os.Args[1]) // doit être le meme port que le client
	addr := ADRESSE + ":" + fmt.Sprint(port)
	// A FAIRE: creer les canaux necessaires, lancer un mainteneur
	ln, _ := net.Listen("tcp", addr) // ecoute sur l'internet electronique
	fmt.Println("Ecoute sur", addr)
	for {
		conn, _ := ln.Accept() // recoit une connection, cree une socket
		fmt.Println("Accepte une connection.")
		go gere_connection(conn) // passe la connection a une routine de gestion des connections
	}
}

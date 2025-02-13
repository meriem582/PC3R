//Benaissa Meriém
//Ahmed Youra

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	st "./structures" // contient la structure Personne
	tr "./travaux"    // contient les fonctions de travail sur les Personnes
)

var ADRESSE string = "localhost"                           // adresse de base pour la Partie 2
var FICHIER_SOURCE string = "./conseillers-municipaux.txt" // fichier dans lequel piocher des personnes
var TAILLE_SOURCE int = 450000                             // inferieure au nombre de lignes du fichier, pour prendre une ligne au hasard
var TAILLE_G int = 5                                       // taille du tampon des gestionnaires
var NB_G int = 2                                           // nombre de gestionnaires
var NB_P int = 2                                           // nombre de producteurs
var NB_O int = 4                                           // nombre d'ouvriers
var NB_PD int = 2                                          // nombre de producteurs distants pour la Partie 2

var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "M"} // une personne vide

type chanDeMessage struct {
	ligne       int
	canalRetour chan string
}

// paquet de personne, sur lequel on peut travailler, implemente l'interface personne_int
type personne_emp struct {
	personne     st.Personne                       // la personne
	ligne        int                               // le numero de ligne dans le fichier source
	tabAFaire    []func(p st.Personne) st.Personne //ici on stocke les fonctions des personnes à appliquer une personne
	statut       string                            // Vide , EnCours ou Fini
	canalLecture chan chanDeMessage
	// A FAIRE
}

// paquet de personne distante, pour la Partie 2, implemente l'interface personne_int
type personne_dist struct {
	// A FAIRE
}

// interface des personnes manipulees par les ouvriers, les
type personne_int interface {
	initialise()          // appelle sur une personne vide de statut V, remplit les champs de la personne et passe son statut à R
	travaille()           // appelle sur une personne de statut R, travaille une fois sur la personne et passe son statut à C s'il n'y a plus de travail a faire
	vers_string() string  // convertit la personne en string
	donne_statut() string // renvoie V, R ou C
}

// fabrique une personne à partir d'une ligne du fichier des conseillers municipaux
// à changer si un autre fichier est utilisé
func personne_de_ligne(l string) st.Personne {
	separateur := regexp.MustCompile("\u0009") // oui, les donnees sont separees par des tabulations ... merci la Republique Francaise
	separation := separateur.Split(l, -1)
	naiss, _ := time.Parse("2/1/2006", separation[7])
	a1, _, _ := time.Now().Date()
	a2, _, _ := naiss.Date()
	agec := a1 - a2
	return st.Personne{Nom: separation[4], Prenom: separation[5], Sexe: separation[6], Age: agec}
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES ***

/*
initialise recupere dans le fichier source la ligne-eme ligne de texte, la convertit en Personne
(le code de personne de ligne est donn´e) et remplit le contenu du paquet avec. Ensuite elle insere
dans le tableau afaire un nombre al´eatoire (1-5) de fonctions de travail (en appelant la fonction
UnTravail de /client/travaux/travaux.go). Puis elle passe le statut du paquet `a R.*/

func (p *personne_emp) initialise() {

	retour := make(chan string)                                          //création d'un canal de message
	p.canalLecture <- chanDeMessage{ligne: p.ligne, canalRetour: retour} //envoie de la ligne à lire
	message := <-retour                                                  //récupération du message
	p.personne = personne_de_ligne(message)                              //création de la personne

	nbFonctions := rand.Intn(5) + 1 //nombre aléatoire de fonctions de travail de 1 à 5
	for i := 0; i < nbFonctions; i++ {
		p.tabAFaire = append(p.tabAFaire, tr.UnTravail()) //ajout des fonctions de travail
	}
	p.statut = "R" //passage du statut à R

}

/*travaille applique `a la personne contenue dans le paquet la fonction de travail en tˆete de afaire
(cela modifie la personne) et retire cette fonction de afaire. Puis, si afaire est vide, elle passe le
statut du paquet `a C.*/

func (p *personne_emp) travaille() {
	if len(p.tabAFaire) > 0 {
		p.personne = p.tabAFaire[0](p.personne) //application de la fonction de travail
		p.tabAFaire = p.tabAFaire[1:]           //suppression de la fonction de travail
	} else {
		p.statut = "C" //passage du statut à C
	}
}

/*vers_string retourne une chaîne de caractères contenant le nom, le prénom, l'âge et le sexe de la personne*/

func (p *personne_emp) vers_string() string {

	var sexe string
	if p.personne.Sexe == "M" {
		sexe = "Homme"
	} else {
		sexe = "Femme"
	}
	return fmt.Sprint("Nom: "+p.personne.Nom+" Prénom: "+p.personne.Prenom+" Age: ", p.personne.Age, " Sexe: "+sexe)
}

func (p *personne_emp) donne_statut() string {
	return p.statut
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES DISTANTES (PARTIE 2) ***
// ces méthodes doivent appeler le proxy (aucun calcul direct)

func (p personne_dist) initialise() {
	// A FAIRE
}

func (p personne_dist) travaille() {
	// A FAIRE
}

func (p personne_dist) vers_string() string {
	// A FAIRE
}

func (p personne_dist) donne_statut() string {
	// A FAIRE
}

// *** CODE DES GOROUTINES DU SYSTEME ***

// Partie 2: contacté par les méthodes de personne_dist, le proxy appelle la méthode à travers le réseau et récupère le résultat
// il doit utiliser une connection TCP sur le port donné en ligne de commande
func proxy() {
	// A FAIRE
}

// Partie 1 : contacté par la méthode initialise() de personne_emp, récupère une ligne donnée dans le fichier source
func lecteur(canal chan chanDeMessage) {
	for {
		message := <-canal
		ligne := message.ligne
		canalRetour := message.canalRetour
		file, err := os.Open(FICHIER_SOURCE)
		if err != nil {
			fmt.Println("Erreur lors de l'ouverture du fichier")
			return
		}
		defer file.Close()
		scanneur := bufio.NewScanner(file)
		//il faut sauter la premiere ligne
		scanneur.Scan()
		for i := 0; i < ligne; i++ {
			scanneur.Scan()
		}
		//faire un test si le scanneur a bien scanné
		if scanneur.Scan() {
			canalRetour <- scanneur.Text()
		} else {
			fmt.Println("Erreur lors de la lecture de la ligne")
		}

		file.Close()

	}
}

// Partie 1: récupèrent des personne_int depuis les gestionnaires, font une opération dépendant de donne_statut()
// Si le statut est V, ils initialise le paquet de personne puis le repasse aux gestionnaires
// Si le statut est R, ils travaille une fois sur le paquet puis le repasse aux gestionnaires
// Si le statut est C, ils passent le paquet au collecteur
func ouvrier(cOutGestion chan personne_int, cInGestion chan personne_int, cInCollecteur chan personne_int) {

	for {
		p := <-cOutGestion
		switch p.donne_statut() {
		case "V":
			p.initialise()
			cInGestion <- p
		case "R":
			p.travaille()
			cInGestion <- p
		case "C":
			cInCollecteur <- p
		}
	}

}

// Partie 1: les producteurs cree des personne_int implementees par des personne_emp initialement vides,
// de statut V mais contenant un numéro de ligne (pour etre initialisee depuis le fichier texte)
// la personne est passée aux gestionnaires
func producteur(cInGestion chan personne_int, lectureDuChan chan chanDeMessage) {
	for {
		ligne := rand.Intn(TAILLE_SOURCE)
		canalRetour := make(chan string)
		lectureDuChan <- chanDeMessage{ligne: ligne, canalRetour: canalRetour}
		personneVide := pers_vide
		aFaire := make([]func(st.Personne) st.Personne, 0)
		p := personne_emp{personne: personneVide, ligne: ligne, tabAFaire: aFaire, statut: "V", canalLecture: lectureDuChan}
		cInGestion <- personne_int(&p)
	}

}

// Partie 2: les producteurs distants cree des personne_int implementees par des personne_dist qui contiennent un identifiant unique
// utilisé pour retrouver l'object sur le serveur
// la creation sur le client d'une personne_dist doit declencher la creation sur le serveur d'une "vraie" personne, initialement vide, de statut V
func producteur_distant() {
	// A FAIRE
}

// Partie 1: les gestionnaires recoivent des personne_int des producteurs et des ouvriers et maintiennent chacun une file de personne_int
// ils les passent aux ouvriers quand ils sont disponibles
// ATTENTION: la famine des ouvriers doit être évitée: si les producteurs inondent les gestionnaires de paquets, les ouvrier ne pourront
// plus rendre les paquets surlesquels ils travaillent pour en prendre des autres
func gestionnaire(cInProdGestion chan personne_int, cOutGestion chan personne_int, cInOuvGestion chan personne_int) {
	file := make([]personne_int, 0)
	for {
		// si la file est vide
		if len(file) == 0 {
			select {
			case p := <-cInProdGestion:
				file = append(file, p)
			case p := <-cInOuvGestion:
				file = append(file, p)
			}
		} else if len(file) == TAILLE_G {
			cOutGestion <- file[0]
			file = file[1:]
		} else if len(file) < TAILLE_G-1 {
			select {
			case p := <-cInProdGestion:
				file = append(file, p)
			case p := <-cInOuvGestion:
				file = append(file, p)
			case cOutGestion <- file[0]:
				file = file[1:]
			}
		} else {
			select {
			case cOutGestion <- file[0]:
				file = file[1:]
			case p := <-cInOuvGestion:
				file = append(file, p)
			}
		}
	}
}

// Partie 1: le collecteur recoit des personne_int dont le statut est c, il les collecte dans un journal
// quand il recoit un signal de fin du temps, il imprime son journal.
func collecteur(cInCollecteur chan personne_int, fintemps chan int) {
	var journal string
	for {
		select {
		case p := <-cInCollecteur:
			journal += p.vers_string() + "\n"
		case <-fintemps:
			fmt.Println("Mon journal d'aujourd'hui : \n", journal)
			fintemps <- 0
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // graine pour l'aleatoire
	if len(os.Args) < 3 {
		fmt.Println("Format: client <port> <millisecondes d'attente>")
		return
	}
	port, _ := strconv.Atoi(os.Args[1])   // utile pour la partie 2
	millis, _ := strconv.Atoi(os.Args[2]) // duree du timeout
	fintemps := make(chan int)
	// A FAIRE
	// creer les canaux
	// lancer les goroutines (parties 1 et 2): 1 lecteur, 1 collecteur, des producteurs, des gestionnaires, des ouvriers
	// lancer les goroutines (partie 2): des producteurs distants, un proxy
	time.Sleep(time.Duration(millis) * time.Millisecond)
	fintemps <- 0
	<-fintemps
}

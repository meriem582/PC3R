//Benaissa Meriem
//Ahmed Youra

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	st "tme4/client/structures" // contient la structure Personne

	tr "tme4/client/travaux" // contient les fonctions de travail sur les Personnes
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
	ligne       int         // ligne à lire
	canalRetour chan string // canal de retour
}

type chanDeMessageDist struct {
	id          int // id de la personne distante
	canalRetour chan string
	methode     string
}

// paquet de personne, sur lequel on peut travailler, implemente l'interface personne_int
type personne_emp struct {
	personne     st.Personne                       // la personne
	ligne        int                               // le numero de ligne dans le fichier source
	tabAFaire    []func(p st.Personne) st.Personne //ici on stocke les fonctions des personnes à appliquer sur une personne
	statut       string                            // V pour vide, R pour rempli, C pour complet
	canalLecture chan chanDeMessage                // canal de communication avec le lecteur
}

// paquet de personne distante, pour la Partie 2, implemente l'interface personne_int
type personne_dist struct {
	// A FAIRE
	id               int
	canalLectureDist chan chanDeMessageDist
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
	fmt.Println("Lecture de la ligne :", l)

	// Séparation par point-virgule
	separation := strings.Split(l, ";")

	// Vérification du nombre de champs
	if len(separation) < 11 {
		fmt.Println("Erreur : nombre de champs insuffisant")
		return st.Personne{}
	}

	// Extraction des données
	nom := separation[6]      // Nom de l'élu
	prenom := separation[7]   // Prénom de l'élu
	sexe := separation[8]     // Code sexe
	naissStr := separation[9] // Date de naissance

	// Conversion de la date de naissance
	naiss, err := time.Parse("02/01/2006", naissStr)
	if err != nil {
		fmt.Println("Erreur lors du parsing de la date :", err)
		return st.Personne{}
	}

	// Calcul de l'âge
	a1, _, _ := time.Now().Date()
	a2, _, _ := naiss.Date()
	agec := a1 - a2
	resultat_lecture := st.Personne{Nom: nom, Prenom: prenom, Age: agec, Sexe: sexe}

	fmt.Println("Resultat de la lecture : ", resultat_lecture)
	return resultat_lecture
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES ***

/*
initialise recupere dans le fichier source la ligne-eme ligne de texte, la convertit en Personne
(le code de personne de ligne est donn´e) et remplit le contenu du paquet avec. Ensuite elle insere
dans le tableau afaire un nombre al´eatoire (1-5) de fonctions de travail (en appelant la fonction
UnTravail de /client/travaux/travaux.go). Puis elle passe le statut du paquet `a R.*/

func (p *personne_emp) initialise() {
	fmt.Println("initialisation pour la ligne  ", p.ligne)

	retour := make(chan string)                                          //création d'un canal de message
	p.canalLecture <- chanDeMessage{ligne: p.ligne, canalRetour: retour} //envoie de la ligne à lire
	message := <-retour                                                  //récupération du message
	p.personne = personne_de_ligne(message)                              //lecture de la ligne et conversion en personne

	nbFonctions := rand.Intn(5) + 1 //nombre aléatoire de fonctions de travail de 1 à 5
	for i := 0; i < nbFonctions; i++ {
		p.tabAFaire = append(p.tabAFaire, tr.UnTravail()) //ajout des fonctions de travail
	}
	p.statut = "R" //passage du statut à R (Rempli)
	fmt.Println("Personne initialisée : ", p.personne, " avec ", nbFonctions, " fonctions de travail", " statut : ", p.statut)

}

/*travaille applique `a la personne contenue dans le paquet la fonction de travail en tˆete de afaire
(cela modifie la personne) et retire cette fonction de afaire. Puis, si afaire est vide, elle passe le
statut du paquet `a C.*/

func (p *personne_emp) travaille() {
	fmt.Println("travail de la personne : ", p.vers_string())
	if len(p.tabAFaire) > 0 {
		p.personne = p.tabAFaire[0](p.personne) //application de la fonction de travail
		p.tabAFaire = p.tabAFaire[1:]           //suppression de la fonction de travail
	} else {
		p.statut = "C" //passage du statut à C
	}
	fmt.Println("Résultat du travaille sur la personne : ", p.vers_string())
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
	retour := make(chan string)
	p.canalLectureDist <- chanDeMessageDist{id: p.id, canalRetour: retour, methode: "initialise"}
	<-retour
}

func (p personne_dist) travaille() {
	// A FAIRE
	retour := make(chan string)
	p.canalLectureDist <- chanDeMessageDist{id: p.id, canalRetour: retour, methode: "travaille"}
	<-retour
}

func (p personne_dist) vers_string() string {
	// A FAIRE
	retour := make(chan string)
	p.canalLectureDist <- chanDeMessageDist{id: p.id, canalRetour: retour, methode: "vers_string"}
	return <-retour
}

func (p personne_dist) donne_statut() string {
	// A FAIRE
	retour := make(chan string)
	p.canalLectureDist <- chanDeMessageDist{id: p.id, canalRetour: retour, methode: "donne_statut"}
	return <-retour
}

// *** CODE DES GOROUTINES DU SYSTEME ***

// Partie 2: contacté par les méthodes de personne_dist, le proxy appelle la méthode à travers le réseau et récupère le résultat
// il doit utiliser une connection TCP sur le port donné en ligne de commande
func proxy(port string, canal chan chanDeMessageDist) {
	// A FAIRE
	address := ADRESSE + ":" + port
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Erreur de connexion au serveur: %v", err)
	}
	for {
		message := <-canal
		request := strconv.Itoa(message.id) + "," + message.methode + "\n"
		fmt.Fprintf(conn, fmt.Sprint(request))
		recu, _ := bufio.NewReader(conn).ReadString('\n')
		response := strings.TrimSuffix(recu, "\n")
		fmt.Println("Reponse du serveur : ", response)
		message.canalRetour <- response
	}
}

// Partie 1 : contacté par la méthode initialise() de personne_emp, récupère une ligne donnée dans le fichier source
func lecteur(canal chan chanDeMessage) {
	fmt.Println("Lecteur du fichier")
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
	fmt.Println("l'ouvrier commence à travailler")

	for {
		//  récuperation des paquets des gestionnaires
		p := <-cOutGestion
		// on examine son statut
		switch p.donne_statut() {
		//si le paquet est Vide, l'ouvrier initialise le paquet puis le renvoient aux gestionnaires,
		case "V":
			p.initialise()
			fmt.Println("initialisation de la personne", p.vers_string())
			cInGestion <- p
			//si le paquet est en cours de traitement, l'ouvrier travaille une fois sur le paquet puis le renvoient aux gestionnaires,
		case "R":
			fmt.Println("avant le travail sur la personne", p.vers_string())
			p.travaille()
			fmt.Println("Résultat du travaille sur la personne", p.vers_string())
			cInGestion <- p
			// si le paquet est Fini, ils lenvoient au collecteur
		case "C":
			fmt.Println("envoie du paquet au collecteur")
			cInCollecteur <- p
		}
	}

}

// Partie 1: les producteurs cree des personne_int implementees par des personne_emp initialement vides,
// de statut V mais contenant un numéro de ligne (pour etre initialisee depuis le fichier texte)
// la personne est passée aux gestionnaires
func producteur(cInGestion chan personne_int, lectureDuChan chan chanDeMessage) {
	for {
		fmt.Println("creation de la personne")
		ligne := rand.Intn(TAILLE_SOURCE)
		personneVide := pers_vide
		aFaire := make([]func(st.Personne) st.Personne, 0)
		// on initialise la personne avec un statut V et on la passe aux gestionnaires
		p := personne_emp{personne: personneVide, ligne: ligne, tabAFaire: aFaire, statut: "V", canalLecture: lectureDuChan}
		fmt.Println("Personne vide crée par le producteur: ", p.vers_string())
		cInGestion <- personne_int(&p)
	}

}

// Partie 2: les producteurs distants cree des personne_int implementees par des personne_dist qui contiennent un identifiant unique
// utilisé pour retrouver l'object sur le serveur
// la creation sur le client d'une personne_dist doit declencher la creation sur le serveur d'une "vraie" personne, initialement vide, de statut V

func producteur_distant(cInGestion chan personne_int, canal chan chanDeMessageDist, idfraischan chan int) {
	for {
		id := <-idfraischan
		nv_personne := personne_dist{id: id, canalLectureDist: canal}
		retour := make(chan string)
		canal <- chanDeMessageDist{id: id, canalRetour: retour, methode: "creer"}
		<-retour
		cInGestion <- nv_personne
	}
}

// Partie 1: les gestionnaires recoivent des personne_int des producteurs et des ouvriers et maintiennent chacun une file de personne_int
// ils les passent aux ouvriers quand ils sont disponibles
// ATTENTION: la famine des ouvriers doit être évitée: si les producteurs inondent les gestionnaires de paquets, les ouvrier ne pourront
// plus rendre les paquets surlesquels ils travaillent pour en prendre des autres
func gestionnaire(cInProdGestion chan personne_int, cOutGestion chan personne_int, cInOuvGestion chan personne_int) {
	fmt.Println("le gestionnaire commence à travailler")
	file := make([]personne_int, 0)
	for {
		// si la file est vide
		if len(file) == 0 {
			fmt.Println("la file est vide")
			// si la file est vide, on peut ajouter un paquet
			select {
			// si le producteur est prêt, on ajoute un paquet
			case p := <-cInProdGestion:
				file = append(file, p)
				// si l'ouvrier est prêt, on ajoute un paquet
			case p := <-cInOuvGestion:
				file = append(file, p)
			}
		} else if len(file) == TAILLE_G {
			fmt.Println("la file est pleine")
			// si la file est pleine, on peut juste defiler un paquet et l'envoyer à l'ouvrier
			cOutGestion <- file[0]
			file = file[1:]
		} else if len(file) < TAILLE_G/2 {
			fmt.Println("la file est moins de la moitié pleine")
			// j'écoute sur les deux canaux quand je suis moins de la moitié pleine
			select {
			// si le producteur est prêt, on ajoute un paquet
			case p := <-cInProdGestion:
				file = append(file, p)
				// si l'ouvrier est prêt, on ajoute un paquet
			case p := <-cInOuvGestion:
				file = append(file, p)
				// si l'ouvrier est prêt, on defile un paquet
			case cOutGestion <- file[0]:
				file = file[1:]
			}
		} else {
			fmt.Println("la file est plus de la moitié pleine")
			// quand je suis plus de la moitié pleine, j'arrête d'écouter sur le canal des producteurs
			// pour éviter la famine des ouvriers
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
	// journal est une variable locale qui contient les personnes traitées
	fmt.Println("le collecteur commence à travailler")
	var journal string
	for {
		// on attend un signal de fin du temps ou une personne à traiter
		select {
		// si on recoit une personne à traiter, on l'ajoute au journal
		case p := <-cInCollecteur:
			journal += p.vers_string() + "\n"
			fmt.Println("journal en cours ... :")
			fmt.Println(journal)
		case <-fintemps:
			// si on recoit un signal de fin du temps, on imprime le journal
			fmt.Println("Mon journal d'aujourd'hui :")
			fmt.Println(journal)
			fintemps <- 0
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // graine pour l'aleatoire

	// partie 1
	/*if len(os.Args) < 1 {
		fmt.Println("Format: client <millisecondes d'attente>")
		return
	}*/
	// partie 2
	if len(os.Args) < 3 {
		fmt.Println("Format: client <port> <millisecondes d'attente>")
		return
	}

	port := os.Args[1]                    // utile pour la partie 2
	millis, _ := strconv.Atoi(os.Args[2]) // duree du timeout
	// partie 1
	//millis, _ := strconv.Atoi(os.Args[1]) // duree du timeout pour la partie 1

	fintemps := make(chan int)
	// A FAIRE
	// creer les canaux
	lectrice := make(chan chanDeMessage)
	cInProdGestion := make(chan personne_int)
	cOutGestion := make(chan personne_int)
	cInOuvGestion := make(chan personne_int)
	cInCollecteur := make(chan personne_int)
	request := make(chan chanDeMessageDist)
	idfraischan := make(chan int)

	// lancer les goroutines (parties 1 et 2): 1 lecteur, 1 collecteur, des producteurs, des gestionnaires, des ouvriers
	go func() {
		lecteur(lectrice)
	}()
	for i := 0; i < NB_P; i++ {
		go func() {
			producteur(cInProdGestion, lectrice)
		}()
	}
	for i := 0; i < NB_G; i++ {
		go func() {
			gestionnaire(cInProdGestion, cOutGestion, cInOuvGestion)
		}()
	}
	for i := 0; i < NB_O; i++ {
		go func() {
			ouvrier(cOutGestion, cInOuvGestion, cInCollecteur)
		}()
	}
	go func() {
		collecteur(cInCollecteur, fintemps)
	}()

	// lancer les goroutines (partie 2): des producteurs distants, un proxy
	go func() { proxy(port, request) }()
	go func() {
		compteur := 0
		for {
			idfraischan <- compteur
			compteur++
		}
	}()
	for i := 0; i < NB_PD; i++ {
		go func() {
			producteur_distant(cInProdGestion, request, idfraischan)
		}()
	}

	time.Sleep(time.Duration(millis) * time.Millisecond)
	fintemps <- 0
	<-fintemps
}

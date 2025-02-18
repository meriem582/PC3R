// MERIEM BENAISSA 
// AHMED YOURA
// on a utiliser copilot
// on s'est inspiré de git pour le "messager" https://github.com/valeeraZ/Sorbonne_PC3R
#include<stdio.h>
#include<stdlib.h>
#include<pthread.h>
#include<string.h>
#include<unistd.h>
#include "include/fthread.h"

int compteur;

struct paquet{ char * nom;};

struct tapis {
    char * nom;
    struct paquet ** fifo;
    int capacite;
    int debut;
    int quantite;
    ft_event_t * event;
};

struct producteur{
    ft_scheduler_t * schedProducteur;
    char * nom;
    int cibleAproduire;
    int nombreProduitTotal;
    struct tapis * tapis;
    FILE * journalProduction;
};

struct consommateur{  
    ft_scheduler_t * schedConsommateur;  
    int id;
    struct tapis * tapis;
    int * compteur;
    FILE * journalConsommation;
    ft_event_t * eventFin;
};

struct messager{
    int id;
    int * compteur;
    ft_scheduler_t * schedProducteur;
    ft_scheduler_t * schedConsommateur;
    struct tapis * tapisProducteur;
    struct tapis * tapisConsommateur;
    FILE * journalMessager;
};

// fonction pour dire que le tapis est vide    
int estVide(struct tapis * tapis){
    return tapis->quantite == 0;
}

// fonction pour dire que le tapis est plein
int estPlein(struct tapis * tapis){
    return tapis->quantite == tapis->capacite ;
}

// fonction pour initialiser le paquet
struct paquet * initPaquet(char * nom){
    struct paquet * paquet = malloc(sizeof(struct paquet));
    paquet->nom = malloc(sizeof(char) * (strlen(nom) + 1));
    strcpy(paquet->nom, nom); 
    return paquet;
}

// fonction pour liberer un paquet 
void libererPaquet(struct paquet * paquet){
    free(paquet->nom);
    free(paquet);
}

// fonction pour initialiser le tapis
struct tapis * initTapis(char * nom,int capacite, ft_event_t * event){ 
    struct tapis * tapis = malloc(sizeof(struct tapis));
    char * nomTapis = malloc(sizeof(char) * (strlen(nom) + 1));
    strcpy(nomTapis, nom);
    tapis->nom = nomTapis;
    tapis->fifo = malloc(capacite * sizeof(struct paquet));
    tapis->capacite = capacite;
    tapis->debut = 0;
    tapis->quantite = 0;
    tapis->event = event;
    return tapis;
}

// fonction pour liberer le tapis
void libererTapis(struct tapis * tapis){
    free(tapis->fifo);
    free(tapis->nom);
}

// fonction pour enfiler un paquet
void enfiler(struct tapis * tapis, struct paquet * paquet){
    while(estPlein(tapis)){
        ft_thread_await(*tapis->event);
        ft_thread_cooperate();
    }
    tapis->fifo[(tapis->debut + tapis->quantite) % tapis->capacite] = paquet;
    tapis->quantite++;
    ft_thread_generate(*tapis->event);
}

// fonction pour defiler un paquet
struct paquet * defiler(struct tapis * tapis, int compteur){
    while(estVide (tapis) && compteur > 0){
        ft_thread_await(*tapis->event);
        ft_thread_cooperate();
    }
    struct paquet * p = tapis->fifo[tapis->debut];
    tapis->quantite--;
    tapis->debut = (tapis->debut + 1) % tapis->capacite;
    compteur --;
    ft_thread_generate(* tapis->event);
    return p;
}


void  production(void * prod)
{
    struct producteur * producteur = (struct producteur *) prod;
    while(producteur->nombreProduitTotal < producteur->cibleAproduire){
        int length = snprintf( NULL, 0, "%d", producteur->nombreProduitTotal );
        char * nb = malloc( length + 1 );
        snprintf( nb, length + 1, "%d", producteur->nombreProduitTotal );
        char * result = malloc(strlen(producteur->nom) + strlen(nb) + 1);
        sprintf(result, "%s%s%s", producteur->nom, " ", nb);
        struct paquet * paquet = initPaquet(result);
        enfiler(producteur->tapis, paquet);
        producteur->nombreProduitTotal++;
        printf("Production de %s\n", paquet->nom);
        fprintf(producteur->journalProduction, "Production de %s\n", paquet->nom);
        ft_thread_cooperate();
    }
    free(producteur);
}

void  consommation(void * cons)
{
    struct consommateur * consommateur = (struct consommateur *) cons;
    while(*(consommateur->compteur)>0){
        struct paquet * paquet = defiler(consommateur->tapis,*(consommateur->compteur));
        if(*(consommateur->compteur) > 0){
            *(consommateur->compteur) = *(consommateur->compteur) - 1;
            printf("C%d mange %s\n", consommateur->id, paquet->nom);
            fprintf(consommateur->journalConsommation, "C%d mange %s\n", consommateur->id, paquet->nom);
            libererPaquet(paquet);
        }
        ft_thread_cooperate();        
    }
    ft_thread_generate(*consommateur->eventFin);
    free(consommateur);
}

void  messagerie(void * mess){
    struct messager * messager = (struct messager *) mess;
    ft_scheduler_t sched = ft_thread_scheduler();
    ft_thread_unlink();
    while(*(messager->compteur) > 0){
        ft_thread_link(*(messager->schedProducteur));
        struct paquet * paquet = defiler(messager->tapisProducteur,*(messager->compteur));
        if(*(messager->compteur) > 0){
            ft_thread_unlink();
            printf("M%d envoie %s\n", messager->id, paquet->nom);
            fprintf(messager->journalMessager, "M%d envoie %s\n", messager->id, paquet->nom);
            ft_thread_link(*(messager->schedConsommateur));
            enfiler(messager->tapisConsommateur, paquet);
        }
        ft_thread_unlink();
    }
    ft_thread_link(sched);
    free(messager);
}

struct terminaison{
    ft_event_t * event;
};

void terminir (void * arg){
    struct terminaison * term =arg;
    ft_thread_await(*term->event);
}
struct producteur * initProducteur(char * nom, int cibleAproduire, struct tapis * tapis, ft_scheduler_t  schedProducteur, FILE * journalProduction){
    struct producteur * p = malloc(sizeof(struct producteur));
    p->nom = nom;
    p->cibleAproduire = cibleAproduire;
    p->nombreProduitTotal = 0;
    p->tapis = tapis;
    p->schedProducteur= &schedProducteur;
    p->journalProduction = journalProduction;
    ft_thread_create(schedProducteur, production, NULL,(void *) p);
    return p;
}

struct consommateur * initConsommateur(int id, struct tapis * tapis, ft_scheduler_t  schedConsommateur, ft_event_t  eventFin, int * compteur, FILE * journalConsommation){
    struct consommateur * c = malloc(sizeof(struct consommateur));
    c->id = id;
    c->compteur = compteur;
    c->tapis = tapis;
    c->schedConsommateur = &schedConsommateur;
    c->eventFin = &eventFin;
    c->journalConsommation = journalConsommation;
    ft_thread_create(schedConsommateur, consommation, NULL,(void *) c);
    return c;
}

struct messager * initMessager(int id, int * compteur, ft_scheduler_t  schedProducteur, ft_scheduler_t  schedConsommateur, struct tapis * tapisProducteur, struct tapis * tapisConsommateur, ft_scheduler_t  schedMess, FILE * journalMessager){
    struct messager * m = malloc(sizeof(struct messager));
    m->id = id;
    m->compteur = compteur;
    m->schedProducteur = &schedProducteur;
    m->schedConsommateur = &schedConsommateur;
    m->tapisProducteur = tapisProducteur;
    m->tapisConsommateur = tapisConsommateur;
    m->journalMessager = journalMessager;
    ft_thread_create(schedMess, messagerie, NULL,(void *) m);
    return m;
}

int main(){
    char * ensProduit[10] = {"produit1", "produit2", "produit3", "produit4", "produit5", "produit6", "produit7", "produit8", "produit9", "produit10"};
    int cibleAproduire = 4, nombreProducteur = 6 ,nombreConsommateur= 7,nbMessage= 5, cpt1 = 0, cpt2 = 0,cpt3 = 0;
    compteur = cibleAproduire * nombreProducteur;
    ft_scheduler_t schedProducteur = ft_scheduler_create();
    ft_scheduler_t schedConsommateur = ft_scheduler_create();
    ft_scheduler_t schedMessager = ft_scheduler_create();

    ft_event_t eventProd = ft_event_create(schedProducteur);
    ft_event_t eventCons = ft_event_create(schedConsommateur);
    ft_event_t eventFin = ft_event_create(schedMessager);

    struct tapis * tapisProd = initTapis("tapisProd", 15, &eventProd);
    struct tapis * tapisCons = initTapis("tapisCons", 15, &eventCons);

    FILE * journalProduction = fopen("journalProduction.txt", "w");
    FILE * journalConsommation = fopen("journalConsommation.txt", "w");
    FILE * journalMessager = fopen("journalMessager.txt", "w");


    while(cpt1 < nombreProducteur){
        struct producteur * prod = initProducteur(ensProduit[cpt1], cibleAproduire, tapisProd, schedProducteur,journalProduction);
        cpt1++;
    }

    while(cpt2 < nombreConsommateur){
        struct consommateur * cons = initConsommateur(cpt2, tapisCons, schedConsommateur, eventFin, &compteur, journalConsommation);
        cpt2++;
    }
    while(cpt3 < nbMessage){
        struct messager * mess = initMessager(cpt3, &compteur, schedProducteur, schedConsommateur, tapisProd, tapisCons, schedMessager, journalMessager);
        cpt3++;
    }
    struct terminaison * term = malloc(sizeof(struct terminaison));
    term->event = &eventFin;
    ft_thread_t termine = ft_thread_create(schedConsommateur, terminir, NULL, (void *) term);
    
    ft_scheduler_start(schedProducteur);
    ft_scheduler_start(schedConsommateur);
    ft_scheduler_start(schedMessager);
    pthread_join(ft_pthread(termine), NULL);

    fclose(journalProduction);
    fclose(journalConsommation);
    fclose(journalMessager);
    

    libererTapis(tapisProd);
    libererTapis(tapisCons);
    return 0;
}
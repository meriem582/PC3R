// MERIEM BENAISSA 
// AHMED YOURA
// on a utiliser copilot
#include<stdio.h>
#include<stdlib.h>
#include<pthread.h>
#include<string.h>
#include<unistd.h>


struct paquet{
    char * nom;
};

struct tapis {
    struct paquet ** fifo;
    int capacite;
    int debut;
    int DernierPaquet;
    // mutex pour la section critique 
    pthread_mutex_t mutex;
    // condition pour que le consomateur ne consomme pas si la file est vide
    pthread_cond_t condCons;
    // condition pour que le producteur ne produise pas si la file est pleine
    pthread_cond_t condProd;
};

struct producteur{
    char * nom;
    int cibleAproduire;
    int nombreProduitTotal;
    // on met le tapis car le produit est enfilé dedans
    struct tapis * tapis;
};

struct consomateur{    
    int id;
    int *cpt;
    struct tapis * tapis;
};

// fonction pour dire que le tapis est vide    
int estVide(struct tapis * tapis){
    return tapis->DernierPaquet == -1;
}

// fonction pour dire que le tapis est plein
int estPlein(struct tapis * tapis){
    return tapis->DernierPaquet == tapis->capacite - 1;
}

// fonction pour initialiser le paquet
struct paquet * initPaquet(char * nom){
    struct paquet * paquet = malloc(sizeof(struct paquet));
    paquet->nom = nom;
    return paquet;
}

// fonction pour liberer un paquet 
void libererPaquet(struct paquet * paquet){
    free(paquet);
}

// fonction pour initialiser le tapis
struct tapis * initTapis(int capacite){
    struct tapis * tapis = malloc(sizeof(struct tapis));
    tapis->fifo = malloc(capacite * sizeof(struct paquet *));
    tapis->capacite = capacite;
    tapis->debut = 0;
    tapis->DernierPaquet = -1;
    pthread_mutex_init(&tapis->mutex, NULL);
    pthread_cond_init(&tapis->condCons, NULL);
    pthread_cond_init(&tapis->condProd, NULL);
    return tapis;
}

// fonction pour liberer le tapis
void libererTapis(struct tapis * tapis){
    free(tapis->fifo);
    free(tapis);
}

// fonction pour enfiler un paquet
void enfiler(struct tapis * tapis, struct paquet * paquet){
    // création de la section critique
    pthread_mutex_lock(&tapis->mutex);
    while(estPlein(tapis)){
        pthread_cond_wait(&tapis->condProd, &tapis->mutex);
    }
    tapis->DernierPaquet++;
    tapis->fifo[tapis->DernierPaquet] = paquet;
    // on signale que le tapis n'est plus vide
    pthread_cond_signal(&tapis->condCons);
    // on sort de la section critique
    pthread_mutex_unlock(&tapis->mutex);
}

// fonction pour defiler un paquet
struct paquet * defiler(struct tapis * tapis){
    pthread_mutex_lock(&tapis->mutex);
    while(estVide(tapis)){
        pthread_cond_wait(&tapis->condCons, &tapis->mutex);
    }
    struct paquet * paquet = tapis->fifo[tapis->debut];
    tapis->debut++;
    tapis->DernierPaquet--;
    if(tapis->debut > tapis->DernierPaquet){
        tapis->debut = 0;
        tapis->DernierPaquet = -1;
    }
    pthread_cond_signal(&tapis->condProd);
    pthread_mutex_unlock(&tapis->mutex);
    return paquet;
}


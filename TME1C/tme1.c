// MERIEM BENAISSA 
// AHMED YOURA
// on a utiliser copilot
#include<stdio.h>
#include<stdlib.h>
#include<pthread.h>
#include<string.h>
#include<unistd.h>



int compteur;

struct paquet{ char * nom;};

struct tapis {
    struct paquet ** fifo;
    int capacite;
    int debut;
    int quantite;
    // on a utiliser les threads POSIX (pthreads) pour la synchronisation {mutex, condCons, condProd}
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
    // on met le tapis car le produit est enfilé dedans (La ressource partagée)
    struct tapis * tapis;
};

struct consomateur{    
    int id;
    struct tapis * tapis;
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
struct tapis * initTapis(int capacite){
    struct tapis * tapis = malloc(sizeof(struct tapis));
    tapis->fifo = malloc(capacite * sizeof(struct paquet));
    tapis->capacite = capacite;
    tapis->debut = 0;
    tapis->quantite = 0;
    pthread_mutex_init(&tapis->mutex, NULL);
    pthread_cond_init(&tapis->condCons, NULL);
    pthread_cond_init(&tapis->condProd, NULL);
    return tapis;
}

// fonction pour liberer le tapis
void libererTapis(struct tapis * tapis){
    free(tapis->fifo);
}

// fonction pour enfiler un paquet
void enfiler(struct tapis * tapis, struct paquet * paquet){
    // verrouiller le mutex 
    pthread_mutex_lock(&tapis->mutex);
    while(estPlein(tapis)){
        pthread_cond_wait(&tapis->condProd, &tapis->mutex); // Les arguments sont le pointeur de la condition et le pointeur du mutex
    }
    if(estVide(tapis)){
        pthread_cond_signal(&tapis->condCons); // on signale que le tapis n'est plus vide à un consomateur qui attend 
    }
    tapis->fifo[ (tapis->debut + tapis->quantite
    ) % tapis->capacite ] = paquet;
    tapis->quantite++;
    // déverrouiller le mutex
    pthread_mutex_unlock(&tapis->mutex);
    // on signale que le tapis n'est plus vide
    pthread_cond_signal(&tapis->condCons);
}

// fonction pour defiler un paquet
struct paquet * defiler(struct tapis * tapis){
    pthread_mutex_lock(&tapis->mutex);
    while(estVide (tapis) && compteur > 0){
        pthread_cond_wait(&tapis->condCons, &tapis->mutex);
    }
    if(estPlein(tapis)){
        pthread_cond_signal(&tapis->condProd);
    }
    struct paquet * p = tapis->fifo[tapis->debut];
    tapis->quantite--;
    tapis->debut = (tapis->debut + 1) % tapis->capacite;
    compteur --;
    pthread_mutex_unlock(&tapis->mutex);
    pthread_cond_signal(&tapis->condProd);
    return p;
}


void * production(void * prod)
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
    }
    free(producteur);
}

void * consommation(void * cons)
{
    struct consomateur * consomateur = (struct consomateur *) cons;
    while(compteur>0){
        struct paquet * paquet = defiler(consomateur->tapis);
        if(compteur > 0){
            printf("C%d mange %s\n", consomateur->id, paquet->nom);
            libererPaquet(paquet);
        }        
    }
    free(consomateur);
}

struct producteur * initProducteur(char * nom, int cibleAproduire, struct tapis * tapis){
    struct producteur * p = malloc(sizeof(struct producteur));
    p->nom = nom;
    p->cibleAproduire = cibleAproduire;
    p->nombreProduitTotal = 0;
    p->tapis = tapis;
    return p;
}

struct consomateur * initConsomateur(int id, struct tapis * tapis){
    struct consomateur * c = malloc(sizeof(struct consomateur));
    c->id = id;
    c->tapis = tapis;
    return c;
}

int main(){
    char * ensProduit[10] = {"produit1", "produit2", "produit3", "produit4", "produit5", "produit6", "produit7", "produit8", "produit9", "produit10"};
    int cibleAproduire = 4, nombreProducteur = 6 ,nombreConsomateur= 7, cpt1 = 0, cpt2 = 0;
    struct tapis * tapis = initTapis(15);
    pthread_t ensProducteur[nombreProducteur], ensConsomateur[nombreConsomateur];
    compteur = cibleAproduire * nombreProducteur;

    while(cpt1 < nombreProducteur){
        struct producteur* prod = initProducteur(ensProduit[cpt1], cibleAproduire, tapis);
        if(pthread_create(&(ensProducteur[cpt1]), NULL, &production, prod) != 0){
            perror("ERREUR DE CREATION DE THREAD");
        }
        cpt1++;
    }

    while(cpt2 < nombreConsomateur){
        struct consomateur * cons = initConsomateur(cpt2, tapis);
        if(pthread_create(&(ensConsomateur[cpt2]), NULL, &consommation, cons) != 0){
            perror("ERREUR DE CREATION DE THREAD");
        }
        cpt2++;
    }

    cpt1 = 0;
    while(cpt1 < nombreProducteur){
        pthread_join(ensProducteur[cpt1], NULL);
        cpt1++;
    }

    cpt2=0;
    while(cpt2 < nombreConsomateur){
        pthread_join(ensConsomateur[cpt2], NULL);
        cpt2++;
    }

    libererTapis(tapis);
    pthread_cond_destroy(&tapis->condProd);
    pthread_cond_destroy(&tapis->condCons);
    return 0;
}
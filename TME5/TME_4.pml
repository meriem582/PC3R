// BENAISSA MERIEM
// AHMED YOURA
// on a utiliser un peu de chat gpt car on connait pas trop la syntaxe de promela

mtype = {Vide, EnCours, Fini}; // vide pour dire que la personne est vide, en cours pour dire que la personne est en train de travailler et fini pour dire que la personne a fini de travailler

chan producteur_gestionnaire = [5] of {int}; // un canal de taille 5 pour que les producteurs puissent envoyer des lignes aux gestionnaires
chan gestionnaire_ouvrier = [5] of {int, mtype}; // un canal de taille 5 pour que les gestionnaires puissent envoyer des personnes avec leurs statuts aux ouvriers
chan ouvrier_collecteur = [5] of {int}; // un canal de taille 5 pour que les ouvriers puissent envoyer des personnes aux collecteurs une fois qu'ils ont fini de travailler sur elles

// Les producteurs envoient des lignes aux gestionnaires, les lignes sont des entiers qui sont des identifiants de lignes et vu que y a pas de fonction random, on a fait une formule pour que les lignes soient différentes 
proctype Producteur(int id) {
    int ligne; 
    do
    :: ligne = (id * 1000) + (id % 450000);
       producteur_gestionnaire!ligne;
    od;
}
// y a deux état qui sont :
// Les gestionnaires recoivent des lignes des producteurs et les envoient aux ouvriers avec un statut vide
// Les gestionnaires recoivent des lignes des ouvriers et les envoient aux collecteurs
proctype Gestionnaire(int id) {
    int ligne;
    mtype statut;
    do
    :: producteur_gestionnaire?ligne -> gestionnaire_ouvrier!ligne, Vide
    :: ouvrier_collecteur?ligne;
    od;
}

// Les ouvriers recoivent des personnes avec leurs statuts des gestionnaires et verifie le statut de la personne, si elle est vide, il la prend et la met en cours et l'envoie au gestionnaire,
// si elle est en cours, il la prend et la met en fini et l'envoie au gestionnaire,
// si elle est fini, il l'envoie au collecteur
proctype Ouvrier(int id) {
    int ligne;
    mtype statut;
    do
    :: gestionnaire_ouvrier?ligne, statut ->
       if
       :: (statut == Vide) -> statut = EnCours; gestionnaire_ouvrier!ligne, statut
       :: (statut == EnCours) -> statut = Fini; gestionnaire_ouvrier!ligne, statut
       :: (statut == Fini) -> ouvrier_collecteur!ligne
       fi;
    od;
}
// skip veut dire que le collecteur ne fait rien, il recoit juste des personnes des ouvriers 
proctype Collecteur() {
    int ligne;
    do
    :: ouvrier_collecteur?ligne -> skip;
    od;
}
// Les observateurs recoivent des personnes avec leurs statuts des gestionnaires et verifie le statut de la personne si il est bon
proctype Observateur() {
    int ligne;
    mtype statut;
    do
    :: gestionnaire_ouvrier?ligne, statut ->
       if
       :: (statut == Vide) -> skip
       :: (statut == EnCours) -> skip
       :: (statut == Fini) -> skip
       fi;
    od;
}

// on a fait un init pour lancer les processus
init {
    atomic {
        run Producteur(1);
        run Producteur(2);
        run Gestionnaire(1);
        run Gestionnaire(2);
        run Ouvrier(1);
        run Ouvrier(2);
        run Ouvrier(3);
        run Ouvrier(4);
        run Collecteur();
        run Observateur();
    }
}

// BENAISSA MERIEM
// AHMED YOURA

mtype {VERT, ORANGE , ROUGE, INDETERMINEE} // type de message pour la couleur du feu
chan observateur_feu = [0] of {mtype,bool} // canal qui transporte la couleur du feu et son clignotement

active proctype Feu_Tricolore() {
    bool clignotement = false; // clignotement du feu (false pour dire que le feu ne clignote pas)
    mtype couleur = INDETERMINEE; // couleur du feu au debut (INDETERMINEE) 
     // les differentes etat du feu
    initial:
        couleur = ORANGE;
        clignotement = true;
        // il y a deux branches possibles
        // true signifie que les deux branches sont toujours executables
        // on choisit une des deux branches de maniere non deterministe
        if 
        // soit on arrete le clignotement du feu et on passe à l'état rouge
        :: true -> clignotement = false; goto rouge;
        // soit on reste dans l'etat init 
        :: true -> goto initial;
        fi
    rouge:
    // on utilise atomic pour dire que les deux instructions doivent etre execute sans interruption
    // la premiere instruction change la couleur du feu a rouge
    // la deuxieme instruction envoie la couleur du feu et son clignotement au canal observateur_feu
        atomic {
            couleur = ROUGE;
            observateur_feu!couleur, clignotement;
        }
        // il y a trois branches possibles
        // on peut passer a l'etat vert, orange ou enpanne
        if
        :: true -> goto vert;
        :: true -> goto enpanne;
        :: true -> goto rouge;
        fi
    vert:
        atomic {
            couleur = VERT;
            observateur_feu!couleur, clignotement;
        }
        if
        :: true -> goto orange;
        :: true -> goto enpanne;
        :: true -> goto vert;
        fi  
    orange:
        atomic {
            couleur = ORANGE;
            observateur_feu!couleur, clignotement;
        }
        if
        :: true -> goto rouge;
        :: true -> goto enpanne;
        :: true -> goto orange;
        fi
    
    enpanne:
        clignotement = true;
    panne_sans_arret :
        couleur = ORANGE;
        observateur_feu!couleur, clignotement;
        if 
        :: true -> goto panne_sans_arret;
        fi
}

active proctype Observateur() {
    mtype couleur_actuelle, couleur_precedente;
    bool clignt;
    couleur_precedente = INDETERMINEE;
    do 
    :: observateur_feu?(couleur_actuelle, clignt) ->
        if
        :: atomic {couleur_actuelle == ORANGE -> assert(clignt == true || couleur_precedente != ROUGE); couleur_precedente = ORANGE}
        :: atomic {couleur_actuelle == ROUGE -> assert(couleur_precedente != VERT); couleur_precedente = ROUGE}
        :: atomic {couleur_actuelle == VERT -> assert(couleur_precedente != ORANGE); couleur_precedente = VERT}
        :: atomic {clignt -> assert(couleur_precedente == ORANGE)}
        fi
    od
}
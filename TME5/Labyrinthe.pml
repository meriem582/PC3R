// BENAISSA MERIEM
// AHMED YOURA

active proctype labyrinthe(){
    initial:
        if
            :: true -> printf("etat initial à la case (ligne = 0 colonne = 4)\n"); goto l0c4
        fi
    l0c4:
        if
            :: true -> printf("Passer à la case (ligne = 1 colonne = 4)\n"); goto l1c4
        fi
    l1c4:
        if
            :: true -> printf("Passer à la case (ligne = 1 colonne = 3)\n"); goto l1c3
            :: true -> printf("Passer à la case (ligne = 2 colonne = 4)\n"); goto l2c4
            :: true -> printf("Passer à la case (ligne = 0 colonne = 4)\n"); goto l0c4
        fi
    l1c3:
        if 
            :: true -> printf("Passer à la case (ligne = 1 colonne = 2)\n"); goto l1c2
            :: true -> printf("Passer à la case (ligne = 1 colonne = 4)\n"); goto l1c4
        fi
    l1c2:
        if
            :: true -> printf("Passer à la case (ligne = 1 colonne = 1)\n"); goto l1c1
            :: true -> printf("Passer à la case (ligne = 1 colonne = 3)\n"); goto l1c3
        fi
    l1c1:
        if
            :: true -> printf("Passer à la case (ligne = 2 colonne = 1)\n"); goto l2c1
            :: true -> printf("Passer à la case (ligne = 0 colonne = 1)\n"); goto l0c1
        fi
    l0c1:
        if
            :: true -> printf("Passer à la case (ligne = 0 colonne = 2)\n"); goto l0c2
            :: true -> printf("Passer à la case (ligne = 1 colonne = 1)\n"); goto l1c1
        fi
    l0c2: 
        if
            :: true -> printf("Passer à la case (ligne = 0 colonne = 3)\n"); goto l0c3
            :: true -> printf("Passer à la case (ligne = 0 colonne = 1)\n"); goto l0c1
        fi
    l0c3:
        if 
            :: true -> printf("Passer à la case (ligne = 0 colonne = 2)\n"); goto l0c2
        fi
    l2c1:
        if
            :: true -> printf("Passer à la case (ligne = 3 colonne = 1)\n"); goto l3c1
            :: true -> printf("Passer à la case (ligne = 2 colonne = 0)\n"); goto l2c0
        fi
    l2c0:
        if
            :: true -> printf("Passer à la case (ligne = 2 colonne = 1)\n"); goto l2c1
        fi
    
    l2c4:
        if
            :: true -> printf("Passer à la case (ligne = 3 colonne = 4)\n"); goto l3c4
            :: true -> printf("Passer à la case (ligne = 1 colonne = 4)\n"); goto l1c4
        fi
    l3c4:
        if
            :: true -> printf("Passer à la case (ligne = 4 colonne = 4)\n"); goto l4c4
            :: true -> printf("Passer à la case (ligne = 2 colonne = 4)\n"); goto l2c4
        fi
    l4c4:
        if
            :: true -> printf("Passer à la case (ligne = 3 colonne = 4)\n"); goto l3c4
        fi
    l3c1:
        if
            :: true -> printf("Passer à la case (ligne = 3 colonne = 2)\n"); goto l3c2
        fi
    l3c2:
        if
            :: true -> printf("Passer à la case (ligne = 3 colonne = 3)\n"); goto l3c3
            :: true -> printf("Passer à la case (ligne = 3 colonne = 1)\n"); goto l3c1
        fi
    l3c3:
        if
            :: true -> printf("Passer à la case (ligne = 4 colonne = 3)\n"); goto l4c3
            :: true -> printf("Passer à la case (ligne = 3 colonne = 2)\n"); goto l3c2
        fi
    l4c3:
        if
            :: true -> printf("Passer à la case (ligne = 4 colonne = 2)\n"); goto l4c2
            :: true -> printf("Passer à la case (ligne = 3 colonne = 3)\n"); goto l3c3
        fi
    l4c2:
        if
            :: true -> printf("Passer à la case (ligne = 4 colonne = 1)\n"); goto l4c1
            :: true -> printf("Passer à la case (ligne = 4 colonne = 3)\n"); goto l4c3
        fi
    l4c1:
        if
            :: true -> printf("Passer à la case (ligne = 4 colonne = 0)\n"); goto l4c0
            :: true -> printf("Passer à la case (ligne = 4 colonne = 2)\n"); goto l4c2
        fi
    l4c0:
        if
            :: true -> goto fin
            :: true -> printf("Passer à la case (ligne = 3 colonne = 0)\n"); goto l3c0
            :: true -> printf("Passer à la case (ligne = 4 colonne = 1)\n"); goto l4c1
        fi
    l3c0:
        if
            :: true -> printf("Passer à la case (ligne = 4 colonne = 0)\n"); goto l4c0
        fi
    fin : 
        printf("Fin du labyrinthe\n"); assert true
}
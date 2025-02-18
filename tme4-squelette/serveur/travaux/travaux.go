package travaux

import (
	"math/rand"

	st "tme4/client/structures" // contient la structure Personne
)

var prenom = []string{"MERIEM", "YOURA", "AHMED", "BENAISSA", "MAYOR"}

// *** LISTES DE FONCTION DE TRAVAIL DE Personne DANS Personne DU SERVEUR ***
// Essayer de trouver des fonctions *différentes* de celles du client

func f1(p st.Personne) st.Personne {
	// A FAIRE
	np := p
	if p.Age > 18 {
		np.Prenom = "MAJOR -" + p.Prenom
	} else {
		np.Prenom = "MINEUR -" + p.Prenom
	}
	return np
}

func f2(p st.Personne) st.Personne {
	// A FAIRE
	np := p
	if p.Age > 18 {
		np.Age = p.Age - 18
	}
	return np
}

func f3(p st.Personne) st.Personne {
	// A FAIRE
	np := p
	if len(p.Prenom) > 5 {
		np.Nom = "TRES LONG -" + p.Nom
	} else {
		np.Nom = "COURT -" + p.Nom
	}
	return np
}

func f4(p st.Personne) st.Personne {
	// A FAIRE
	np := p
	prenom := prenom[rand.Intn(len(prenom))]
	np.Prenom = prenom
	return np
}

func UnTravail() func(st.Personne) st.Personne {
	tableau := make([]func(st.Personne) st.Personne, 0)
	tableau = append(tableau, func(p st.Personne) st.Personne { return f1(p) })
	tableau = append(tableau, func(p st.Personne) st.Personne { return f2(p) })
	tableau = append(tableau, func(p st.Personne) st.Personne { return f3(p) })
	tableau = append(tableau, func(p st.Personne) st.Personne { return f4(p) })
	i := rand.Intn(len(tableau))
	return tableau[i]
}

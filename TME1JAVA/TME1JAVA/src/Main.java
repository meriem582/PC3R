public class Main {

	public static void main(String[] args) {
		// TODO Auto-generated method stub
		String[] produits = { "produit1", "produit2", "produit3", "produit4", "produit5", "produit6", "produit7",
				"produit8", "produit9", "produit10" };
		int cibleAproduire = 4, nombreProducteur = 6, nombreConsomateur = 7;
		Tapis tapis = new Tapis(15);
		int compteur = cibleAproduire * nombreProducteur;
		Thread[] producteurs = new Thread[nombreProducteur];
		Thread[] consommateurs = new Thread[nombreConsomateur];
		for (int i = 0; i < nombreProducteur; i++) {
			Thread producteur = new Thread(new Producteur(produits[i], cibleAproduire, 0, tapis));
			producteurs[i] = producteur;
			producteur.start();
		}

		for (int j = 0; j < nombreConsomateur; j++) {
			Thread consommateur = new Thread(new Consommateur(compteur, j, tapis));
			consommateurs[j] = consommateur;
			consommateur.start();
		}

		for (int i = 0; i < nombreProducteur; i++) {
			try {
				producteurs[i].join();
			} catch (InterruptedException e) {
				e.printStackTrace();
			}
			i++;
		}

		for (int j = 0; j < nombreConsomateur; j++) {
			try {
				consommateurs[j].join();
			} catch (InterruptedException e) {
				e.printStackTrace();
			}
			j++;
		}
	}

}

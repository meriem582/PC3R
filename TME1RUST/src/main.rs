use std::collections::VecDeque;
use std::ops::Deref;
use std::sync::Arc;
use std::sync::Condvar;
use std::sync::Mutex;

struct Paquet {
    nom: String,
}

struct Tapis {
    ens_produit: VecDeque<Paquet>,
    taille: usize,
}

fn enfiler(tapis: Arc<(Mutex<Tapis>, Condvar)>, paquet: Paquet) -> Arc<(Mutex<Tapis>, Condvar)> {
    let tapis_bis = Arc::clone(&tapis);
    let (mutex, cond_var) = &*tapis;
    let mut tapis = mutex.lock().unwrap();
    while tapis.ens_produit.capacity() == tapis.taille {
        tapis = cond_var.wait(tapis).unwrap();
    }
    if tapis.ens_produit.is_empty() {
        cond_var.notify_all();
    }
    tapis.ens_produit.push_back(paquet);
    tapis.taille = tapis.taille + 1;
    tapis_bis
}

fn defiler(tapis: &Arc<(Mutex<Tapis>, Condvar)>) -> Paquet {
    let (mutex, cond_var) = tapis.deref();
    let mut tapis = mutex.lock().unwrap();
    while tapis.ens_produit.is_empty() {
        tapis = cond_var.wait(tapis).unwrap();
    }
    if tapis.ens_produit.capacity() == tapis.taille {
        cond_var.notify_all();
    }
    let paquet = tapis.ens_produit.pop_front().unwrap();
    paquet
}

fn main() {
    //  println!("Hello, world!");
    let ens_produit: [&str; 10] = [
        "produit1 ",
        "produit2 ",
        "produit3 ",
        "produit4 ",
        "produit5 ",
        "produit6 ",
        "produit7 ",
        "produit8 ",
        "produit9 ",
        "produit10 ",
    ];
    let nombre_producteur: usize = 5;
    let nombre_consomateur: usize = 7;
    let cible_aproduire: usize = 4;

    let t1 = Tapis {
        ens_produit: VecDeque::new(),
        taille: 15,
    };
    let t2 = Arc::new(Mutex::<usize>::new(nombre_producteur * cible_aproduire));
    let tapis = Arc::new((Mutex::new(t1), Condvar::new()));
    let mut v = Vec::new();

    for i in 0..nombre_producteur {
        let mut tapis_bis = Arc::clone(&tapis);
        v.push(std::thread::spawn(move || {
            for j in 0..cible_aproduire {
                let p = Paquet {
                    nom: ens_produit[i].to_string() + &j.to_string(),
                };
                tapis_bis = enfiler(tapis_bis, p);
            }
        }));
    }

    for i in 0..nombre_consomateur {
        let tapis_bis = Arc::clone(&tapis);
        let t2bis = Arc::clone(&t2);
        v.push(std::thread::spawn(move || {
            let mut continu = true;
            while continu {
                {
                    let mut mg = t2bis.lock().unwrap();
                    if *mg <= 0 {
                        continu = false;
                    } else {
                        *mg -= 1;
                    }
                };
                if continu {
                    let p = defiler(&tapis_bis);
                    println!("C{} mange {}", i.to_string(), p.nom.to_string());
                }
            }
        }));
    }

    for vi in v {
        vi.join().unwrap();
    }
}

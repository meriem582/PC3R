    
mtype { Alice , Bob, Carole , David , Elena , user_123 , mdp , chaussette , OK , KO , BIENVENUE , NOTBIENVENUE }
chan url = [5] of { mtype, mtype, chan };
chan s = [5] of { chan };

proctype client(mtype login , mtype mdp , chan contact, chan obs1, chan obs2! ) {
       
    url!login,mdp,contact     
   do
   :: contact?OK -> printf("Le client recoit OK"); 
   :: contact?KO -> printf("Le client recoit KO"); obs2!KO;
   od
}

proctype authentifieur () {
     
     chan contact;
    do 
    :: url?login,mdp,contact ->
    if 
    :: (login==Alice && mdp==user_123 ) || ( login==Bob && mdp=mdp ) || (login==Carole && mdp=user_123  ) -> 
       contact!OK;
       s!contact; obs1!BIENVENUE;
    :: else -> contact!KO;  obs2!NOTBIENVENUE; 

    od
}

proctype observateur(chan obs1;chan obs2) {
   mtype v;
     do 
    :: obs2?v ->
    assert(v==BIENVENUE);
     od;

}
proctype serveur () {

       chan c ;
       do 
       :: s?c -> c!BIENVENUE obs1!BIENVENUE
       od;

}

init {
 
    chan c1 = [1] of { mtype };
    chan c2 = [1] of { mtype };
    chan c3 = [1] of { mtype };
    chan c4 = [1] of { mtype };
    chan c5 = [1] of { mtype };  

    chan obs= [0] of { mtype };
    run authentifieur();
    run serveur();
    run client(Alice,user_123,c1);
    run client(Bob,mdp,c2);
    run client(Carole,user_123,c3);
    run client(David,chaussette,c4);
    run client(Elena,user_123,c5);

}
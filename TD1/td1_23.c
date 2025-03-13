
#define NP_THREAD 10
#include <stdio.h>
#include <pthread.h>
#include <sched.h>

int SHARED_compteur;

void* routine(void* arg){
	int temp = SHARED_compteur; // On stocke la valeur du compteur dans une variable temporaire pour éviter les problèmes de concurrence
	sched_yield();
	temp++;
	sched_yield();
	SHARED_compteur = temp;
	printf("Compteur = %d\n", SHARED_compteur); 
}

int main()
{	
	SHARED_compteur = 0;
	pthread_t thread[NP_THREAD];
	for(int i = 0;i<NP_THREAD;i++)
		pthread_create(&thread[i], NULL, routine, NULL);

	for(int i = 0;i<NP_THREAD;i++)
		pthread_join(thread[i], NULL);

	if(NP_THREAD == SHARED_compteur)
		printf("TERMINE\n");
	
}
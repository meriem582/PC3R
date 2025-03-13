#include <stdio.h>
#include <unistd.h>
#include <pthread.h>

int v = 0;
pthread_t th1;

void *fun_fils(){
    sleep(1);
    printf("Je suis le fils, v vaut %d \n",v);
}

int main(){
    v = 42;
    pthread_create(&th1, NULL, fun_fils, NULL);
    v = 38;
    printf("Je suis le pere, v vaut %d \n", v);
    sleep(4);
}
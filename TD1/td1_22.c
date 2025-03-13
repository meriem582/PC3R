#include <stdio.h>
#include <unistd.h>
#include <pthread.h>

void *fun1(void *arg) {
    printf("belle marquise\n");
}

void *fun2(void *arg) {
    printf("vos beaux yeux\n");
}
void *fun3(void *arg) {
    printf("me font\n");
}
void *fun4(void *arg) {
    printf("mourire\n");
}
void *fun5(void *arg) {
    printf("d'amour\n");
}

int main() {
    pthread_t t1, t2, t3, t4, t5;
    pthread_create(&t1, NULL, fun1, NULL);
    pthread_create(&t2, NULL, fun2, NULL);
    pthread_create(&t3, NULL, fun3, NULL);
    pthread_create(&t4, NULL, fun4, NULL);
    pthread_create(&t5, NULL, fun5, NULL);
    return 0;
}
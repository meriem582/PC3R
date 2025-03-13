#include <stdio.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>

int main(){
    int v = 42;
    int pid = fork();
    if (pid == 0){
        sleep(1);
        printf("Je suis le fils, v vaut %d \n", v);
    }else{
        v = 38;
        printf("Je suis le père, v vaut %d \n", v);
        sleep(5);
    }
}
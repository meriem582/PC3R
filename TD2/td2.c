#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>

// Assuming ft_scheduler is a pointer to some scheduler structure
typedef struct ft_scheduler_t* ft_scheduler;

void run_p (void *phrase) {
    while(1){
        printf("%s \n", (char *) phrase);
        ft_thread_cooperate();
    }
}

int main(void) {
    ft_scheduler sched = ft_scheduler_create ();
    ft_thread_create(sched, run_p, NULL, (void *) "d'amour");
    ft_thread_create(sched, run_p, NULL, (void *) "vos beaux yeux");
    ft_thread_create(sched, run_p, NULL, (void *) "me font mourir");
    ft_thread_create(sched, run_p, NULL, (void *) "d'amour");
    ft_scheduler_start (sched);

    return 0;

}
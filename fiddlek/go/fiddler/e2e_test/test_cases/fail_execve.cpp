#include <unistd.h>

int main() {
    char *newargv[] = { (char *)"/bin/ls", NULL };
    char *newenviron[] = { NULL };
    execve("/bin/ls", newargv, newenviron);
    return 0;
}

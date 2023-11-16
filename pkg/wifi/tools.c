#ifdef __linux__
#include "linux_wifi.h"
#elif defined(__APPLE__) && defined(__MACH__)
#include "darwin_wifi.h"
#endif

#include <sys/wait.h>

// Function to execute a command using fork and execvp.
int execute_command(const char *command, char *const args[]) {
    pid_t pid, wpid;
    int status = 0;

    if ((pid = fork()) == 0) {
        if (execvp(command, args) == -1) {
            goSendToChannel("Network connection error, fork: command execution failed");
            exit(EXIT_FAILURE);
        }
    } else if (pid < 0) {
        goSendToChannel("Network connection error, fork failed");
        return -1;
    } else {
        do {
            wpid = waitpid(pid, &status, WUNTRACED);
        } while (!WIFEXITED(status) && !WIFSIGNALED(status));
    }
    return WEXITSTATUS(status) == 0 ? 0 : -1;
}
#ifdef __linux__
#include "linux_wifi.h"
#elif defined(__APPLE__) && defined(__MACH__)
#include "darwin_wifi.h"
#endif

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
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

int network_recovery() {
    goSendToChannel("Starting network recovery process");

    // Удаление возможных поврежденных конфигурационных файлов wpa_supplicant
    char *rmWpaConfArgs[] = {"rm", "-f", "/var/run/wpa_supplicant/*", NULL};
    if (execute_command("rm", rmWpaConfArgs) != 0) {
        goSendToChannel("Failed to remove wpa_supplicant configuration files");
    }

    char *rmargs[] = {"rm", "-f", "/var/run/wpa_supplicant/" WLAN_IFACE, "/tmp/wpa_conf_*.conf", NULL};
    if (execute_command("rm", rmargs) != 0) {
        goSendToChannel("Failed to remove residual WiFi configuration files");
    }

    // Активация сетевого интерфейса с использованием 'ip'
    char *ipLinkSetUpArgs[] = {"ip", "link", "set", WLAN_IFACE, "up", NULL};
    if (execute_command("ip", ipLinkSetUpArgs) != 0) {
        goSendToChannel("Failed to activate the network interface using 'ip'");
    }

    // Перезапуск NetworkManager с использованием 'systemctl'
    char *systemctlRestartArgs[] = {"systemctl", "restart", "NetworkManager", NULL};
    if (execute_command("systemctl", systemctlRestartArgs) != 0) {
        goSendToChannel("Failed to restart NetworkManager using 'systemctl'");
    }

    goSendToChannel("WiFi recovery process completed");
    return 0;
}
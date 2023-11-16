#ifdef linux
#include "linux_wifi.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>

// Function to disconnect from the WiFi network.
int network_disconn() {
    goSendToChannel("Starting WiFi disconnection");

    char *dhclientKillArgs[] = {"killall", "dhclient", NULL};
    goSendToChannel("Releasing IP address and stopping dhclient");

    if (execute_command("killall", dhclientKillArgs) != 0) {
        goSendToChannel("Failed to stop dhclient");
    }
    char *wpaKillArgs[] = {"killall", "wpa_supplicant", NULL};
    goSendToChannel("Stopping wpa_supplicant to disconnect from WiFi");

    if (execute_command("killall", wpaKillArgs) != 0) {
        goSendToChannel("Failed to stop wpa_supplicant");
        return -1;
    }
    goSendToChannel("WiFi disconnection successful");
    network_recovery();
    return 0;
}
#endif

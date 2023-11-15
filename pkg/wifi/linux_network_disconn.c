#ifdef linux
#include "linux_wifi.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>

// Function to disconnect from the WiFi network.
int network_disconn() {
    goSendToChannel("Starting WiFi disconnection");

    // Stopping dhclient to release the IP address.
    char *dhclientKillArgs[] = {"killall", "dhclient", NULL};
    goSendToChannel("Releasing IP address and stopping dhclient");
    if (execute_command("killall", dhclientKillArgs) != 0) {
        goSendToChannel("Failed to stop dhclient");
    }

    // Stopping wpa_supplicant to disconnect from the WiFi network.
    char *wpaKillArgs[] = {"killall", "wpa_supplicant", NULL};
    goSendToChannel("Stopping wpa_supplicant to disconnect from WiFi");
    if (execute_command("killall", wpaKillArgs) != 0) {
        goSendToChannel("Failed to stop wpa_supplicant");
        return -1;
    }

    goSendToChannel("WiFi disconnection successful");
    return 0;
}
#endif

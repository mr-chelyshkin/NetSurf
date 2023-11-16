#ifndef LINUX_WIFI_H
#define LINUX_WIFI_H

#include <iwlib.h>
#include <stddef.h>

#define WLAN_IFACE "wlan0"
#define MAX_NETWORKS 100

// Struct to hold information about WiFi networks.
typedef struct wifi_info {
    char   ssid[33];
    double freq;
    int    quality;
    int    level;
} wifi_info;

/**
 * Connects to a WiFi network using the specified SSID, password, and country code.
 *
 * This function initiates the process of connecting to a WiFi network. It handles the necessary steps to configure
 * the wireless interface for the desired network, authenticate, and establish a connection. This typically involves
 * stopping any running instances of network management software, setting up the correct network parameters, and
 * starting the network connection process.
 *
 * The function uses the provided SSID (Service Set Identifier) to identify the desired WiFi network, the password
 * for authentication (assumed to be a WPA or WPA2 passphrase), and the country code to comply with regional
 * wireless operation standards. It creates a temporary configuration file for wpa_supplicant, starts the wpa_supplicant
 * process with this configuration, and then initiates a DHCP client to obtain an IP address.
 *
 * @param ssid Pointer to a null-terminated string containing the SSID of the network to connect to.
 * @param password Pointer to a null-terminated string containing the password of the network.
 * @param country Pointer to a null-terminated string containing the country code for the network.
 * @return Returns 0 if the connection was successfully established, and -1 if an error occurred during the process.
 */
int network_conn(const char* ssid, const char* password, const char* country);

/**
 * Disconnects from the currently connected WiFi network.
 *
 * This function initiates the disconnection process from the currently connected WiFi network.
 * It is responsible for terminating any active network sessions and ensuring that the wireless
 * interface is no longer associated with the previous network. This might involve sending
 * appropriate disconnection commands to the network interface controller and handling any
 * necessary cleanup operations to return the interface to an idle state.
 *
 * @return A result code indicating the success or failure of the disconnection process.
 *         Returns 0 if the disconnection was successful, and a non-zero value if an error occurred.
 */
int network_disconn();

/**
 * Scans for available WiFi networks in the vicinity.
 *
 * This function performs a scan to detect available WiFi networks and gathers information about each one.
 * It utilizes the wireless tools library for interfacing with the wireless device. The function compiles a list
 * of WiFi networks within range, capturing details such as the SSID (Service Set Identifier), frequency, signal quality,
 * and signal level for each network detected. This information is useful for identifying and selecting WiFi networks
 * for connection purposes.
 *
 * The function returns a pointer to an array of wifi_info structures, each representing a found network. The number of
 * networks found is stored in the integer pointed to by the 'count' parameter. The array and the count are static; thus,
 * they persist between function calls but will be overwritten on subsequent calls.
 *
 * Note: The function returns a maximum of MAX_NETWORKS networks. If more networks are available, they will not be included
 * in the results.
 *
 * @param count Pointer to an integer where the number of networks found will be stored.
 * @return A pointer to a static array of wifi_info structures containing information about each found network.
 *         If an error occurs during the scan, NULL is returned and an error message is sent to a custom channel.
 */
wifi_info* network_scan(int* count);

/**
 * Get information about the currently connected WiFi network.
 *
 * @return The SSID of the current connection or NULL if not connected.
 */
const char* network_state();

/**
 * Redirects the standard output and error streams to a custom channel.
 *
 * This function is used to intercept and redirect output sent to the standard output (stdout) and
 * standard error (stderr) streams. It reroutes these streams to a custom channel, which can be used
 * for logging, monitoring, or other purposes. Typically used in conjunction with `redirected_write`
 * function to handle the redirected output.
 */
void redirect_output(void);

/**
 * Resets the standard output and error streams back to their original state.
 *
 * This function restores stdout and stderr to their original file descriptors. This should be called
 * after `redirect_output` to end the redirection of output to the custom channel. It ensures that any
 * further output goes to the original standard output and error destinations.
 */
void reset_output(void);

/**
 * Writes data to a custom channel. This function is used in conjunction with redirected output functions.
 *
 * It takes the data written to the standard output or error streams and writes it to a custom channel instead.
 * This is useful for capturing and redirecting output from standard library functions or other libraries that
 * write to stdout/stderr, to a different destination managed by the application.
 *
 * @param fd File descriptor (unused in this context).
 * @param buf Pointer to the buffer containing data to be written.
 * @param count Number of bytes to write from the buffer.
 * @return The number of bytes written, or -1 on error.
 */
int redirected_write(int fd, const void* buf, size_t count);

/**
 * Sends data to a custom channel. This function is to be implemented in Go.
 *
 * It is designed to send log messages or other information from the C code to a custom channel,
 * potentially for logging, debugging, or other monitoring purposes. This function can be used to
 * integrate C code with higher-level systems or interfaces written in Go.
 *
 * @param s Pointer to a null-terminated string containing the data to be sent.
 */
extern void goSendToChannel(char* s);

/**
 * Executes a system command using fork and execvp system calls.
 *
 * @param command Pointer to a string containing the command to execute.
 * @param args Array of command-line arguments, terminated with NULL, for passing to execvp.
 * @return Returns 0 if the command is executed successfully, and -1 in case of an error.
 *
 * This function creates a new process (child process) using fork().
 * In the child process, execvp() is used to replace the current process
 * with the executable program specified in command, with args as its arguments.
 * If execvp returns -1, a message is sent to a channel about the error, and the process exits.
 * In the parent process, waitpid() is used to wait for the completion of the child process
 * and to obtain its exit status. The function returns 0 if the command was executed successfully,
 * and -1 if there was an error during execution.
 */
int execute_command(const char *command, char *const args[]);

int network_recovery();

#endif

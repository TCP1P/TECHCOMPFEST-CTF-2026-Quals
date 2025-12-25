#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

int main(void) {
    char username[32];
    char data[48];

    printf("Content-Type: text/plain\r\n\r\n");
    fflush(stdout);

    char *strlength = getenv("CONTENT_LENGTH");
    if (!strlength) {
        printf("err: no data received\n");
        return 0;
    }

    int length = atoi(strlength);
    read(STDIN_FILENO, data, length);

    if (strncmp(data, "username=", 9) != 0) {
        printf("err: bad request format\n");
        return 0;
    }

    strcpy(username, data + 9);

    printf("hi %s\n", username);

    return 0;
}

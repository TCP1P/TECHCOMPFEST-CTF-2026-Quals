// gcc vuln.c -o vuln -fno-stack-protector -no-pie

#include <stdio.h>
#include <stdlib.h>

void setup() {
    setvbuf(stdout, NULL, _IONBF, 0);
    setvbuf(stdin, NULL, _IONBF, 0);
}
    

void win() {
    system("cat flag.txt");
}

int main() {
    setup();
    char buffer[64];
    printf("input: ");
    gets(buffer);
}
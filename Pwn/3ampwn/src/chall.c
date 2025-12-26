#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <malloc.h>

char *notes[0x100];
int count = 0;

static int get_int(const char *s) {
  char input[0x10];
  printf("%s", s);
  fgets(input, sizeof(input), stdin);
  return atoi(input);
}

int main() {
  setbuf(stdin, NULL);
  setbuf(stdout, NULL);
  setbuf(stderr, NULL);

  int index;
  for(;;){
    switch (get_int("> ")) {
      case 1:
        if (count <= 0x100) notes[count++] = calloc(get_int("(add) size: "), 1);
        break;
      case 2:
        index = get_int("(edit) index: ");
        if (notes[index] && index < count){
          printf("(edit) data: ");
          fgets(notes[index], malloc_usable_size(notes[index]), stdin);
        } else puts("no");
        break;
      case 3:
        index = get_int("(view) index: ");
        if (notes[index] && index < count) printf("(view) data: %s", notes[index]);
        else puts("no");
        break;
      case 4:
        index = get_int("(delete) index: ");
        if (notes[index] && index < count){
          printf("y/n? ");
          if (getchar() == 'y'){
            free(notes[index]);
            notes[index] = NULL;
            break;
          }
        } else puts("no");
        break;
      default:
        break;
    }
  }
}

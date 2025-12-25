CC      := gcc
CFLAGS  := -Wall -fno-stack-protector -no-pie -static
TARGET  := httpwn.cgi
SRC     := main.c
DIR     := /usr/local/apache2
WEB     := index.html

all: $(TARGET)

$(TARGET): $(SRC)
	$(CC) $(CFLAGS) -o $(TARGET) $(SRC)

install: $(TARGET)
	cp $(TARGET) $(DIR)/cgi-bin/
	cp $(WEB) $(DIR)/htdocs/
	chmod 755 $(DIR)/cgi-bin/$(TARGET)
	chmod 644 $(DIR)/htdocs/$(WEB)

clean:
	rm -f $(TARGET)

.PHONY: all install clean

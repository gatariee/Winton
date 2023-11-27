CSC = csc

SRC = Main.cs

OUT = ../bin/whoami.exe

build:
	$(CSC) -out:$(OUT) $(SRC)

clean:
	rm -f $(OUT)

run: build
	./$(OUT)

.PHONY: build clean run

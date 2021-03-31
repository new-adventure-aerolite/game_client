.PHONY: build run clean

build:
	go build -o bin/game-client .

run:
	bin/game-client start

clean:
	rm -rf bin/game-client

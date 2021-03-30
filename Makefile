.PHONY: build run clean

build:
	go build -o bin/game_client .

run:
	bin/game_client start

clean:
	rm -rf bin/game_client

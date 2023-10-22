build:
	mkdir -p bin/
	go build -o bin/leds .

run:
	bin/leds

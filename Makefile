all: clessy-broker

clessy-broker:
	go build -o clessy-broker ./broker
	go build -o clessy-mods ./mods

clean:
	rm -f clessy-broker clessy-mods

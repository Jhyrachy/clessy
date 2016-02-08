all: clessy-broker

clessy-broker:
	go build -o clessy-broker ./broker

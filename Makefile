all: clessy-broker clessy-mods clessy-stats

install-tg:
	go install github.com/hamcha/clessy/tg

clessy-broker: install-tg
	go build -o clessy-broker github.com/hamcha/clessy/broker

clessy-mods: install-tg
	go build -o clessy-mods github.com/hamcha/clessy/mods

clessy-stats: install-tg
	go build -o clessy-stats github.com/hamcha/clessy/stats

clean:
	rm -f clessy-broker clessy-mods clessy-stats
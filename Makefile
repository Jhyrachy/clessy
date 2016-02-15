all: clessy-broker clessy-mods clessy-stats

deps:
	go get github.com/boltdb/bolt/...
	go get github.com/golang/freetype

install-tg:
	go install github.com/hamcha/clessy/tg

clessy-broker: install-tg
	go build -o clessy-broker github.com/hamcha/clessy/broker

clessy-mods: install-tg
	go build -o clessy-mods github.com/hamcha/clessy/mods

clessy-stats: install-tg
	go build -o clessy-stats github.com/hamcha/clessy/stats

clessy-stats-import: install-tg
	go build -o clessy-stats-import github.com/hamcha/clessy/stats-import

clean:
	rm -f clessy-broker clessy-mods clessy-stats
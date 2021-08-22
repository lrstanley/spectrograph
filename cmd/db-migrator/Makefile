BINARY = db-migrator

clean:
	/bin/rm -rfv "rice-box.go" "${BINARY}"

compile:
	go build -mod=readonly -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

build: compile
	echo

debug:
	go run *.go --debug \
		--mongo.db-name spectrograph-dev \
		--mongo.uri "mongodb://mongo:27017/?maxPoolSize=5" \
		--migration.purge

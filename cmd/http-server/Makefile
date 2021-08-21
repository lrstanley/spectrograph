BINARY = http-server

clean:
	/bin/rm -rfv "tmp" "public/dist/*" "rice-box.go" "${BINARY}"

generate:
	rice -v embed-go

compile:
	go build -mod=readonly -ldflags '-s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

build: clean generate compile
	echo

debug:
	go run *.go --debug \
		--http.bind-addr :8080 \
		--http.base-url "http://localhost:8081/" \
		--mongo.db-name spectrograph-dev \
		--mongo.uri "mongodb://mongo:27017/?maxPoolSize=200"

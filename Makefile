PROJECT = milight-daemon

SOURCE := $(shell find . -name '*.go')

$(PROJECT): $(SOURCE)
	go build -o $(PROJECT)

run: $(PROJECT)
	./$(PROJECT)

debug: $(PROJECT)
	./$(PROJECT) --mode=debug

pi: $(SOURCE)
	GOOS=linux GOARCH=arm GOARM=6 go build -o $(PROJECT)

pi2: $(SOURCE)
	GOOS=linux GOARCH=arm GOARM=7 go build -o $(PROJECT)

clean:
	rm $(PROJECT)

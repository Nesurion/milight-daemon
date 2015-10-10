PROJECT = milight-daemon

SOURCE := $(shell find . -name '*.go')

$(PROJECT): $(SOURCE)
	go build -o $(PROJECT)

run: $(PROJECT)
	./$(PROJECT)

clean:
	rm $(PROJECT)

PROJECT = milight-daemon

SOURCE := $(shell find . -name '*.go')

$(PROJECT): $(SOURCE)
	go build -o $(PROJECT)

run: $(PROJECT)
	./$(PROJECT)

debug: $(PROJECT)
	./$(PROJECT) --mode=debug

clean:
	rm $(PROJECT)

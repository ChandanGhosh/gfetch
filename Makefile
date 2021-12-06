build:
	go build -o gfetch

clean:
	rm -f gfetch

tidy:
	go mod tidy

.PHONY: tidy build

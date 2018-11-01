.PHONY: clean all

all: srp_server

srp_server:
	go build -ldflags "-w -s" -o ./out/srp_server ./cmd

clean:
	rm -fr ./out/srp_server

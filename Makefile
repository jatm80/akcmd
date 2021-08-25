.PHONY: all
all: clean client installer

.PHONY: client
client:	
	mkdir -p ./out
	go build -o ./out/akcmd

.PHONY: installer
installer:
	mkdir -p ./out
	cd installer && go build -o ../out/akcmd_installer

.PHONY: clean
clean:
	rm -rf out
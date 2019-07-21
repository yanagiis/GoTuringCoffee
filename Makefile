DEST_IP = 192.168.2.22

GENERATE = internal/hardware/max31856/max31856.go\
		   internal/hardware/max31865/max31865.go\
		   internal/hardware/hardware.go

UTILS = max31856 max31865 db process vl6180x max31865gpio max31856gpio

ARCH ?= arm 

.PHONY: generate clean

all: generate app utils

dep:
	dep ensure

app: generate
	GOOS=linux GOARCH=$(ARCH) go build -o GoTuringCoffee_$(ARCH)

utils: generate
	for util in $(UTILS); do \
		GOOS=linux GOARCH=$(ARCH) go build -o ./bin/$$util\_$(ARCH) ./utils/$$util.go ;\
	done

generate: $(GENERATE)
	for file in $(GENERATE) ; do \
		go generate $$file ; \
	done

clean:
	go clean
	rm GoTuringCoffee_*
	rm -rf ./bin

copy:
	sshpass scp config.yml bin/*_arm ./GoTuringCoffee_arm root@$(DEST_IP):/home/root

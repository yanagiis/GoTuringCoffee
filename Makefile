DEST_IP = 192.168.1.78

GENERATE = internal/hardware/max31856/max31856.go\
		   internal/hardware/max31865/max31865.go

UTILS = max31856 max31865

ARCH ?= amd64

.PHONY: generate clean

all: generate app utils

dep:
	dep ensure

app: generate
	GOOS=linux GOARCH=$(ARCH***REMOVED*** go build -o GoTuringCoffee_$(ARCH***REMOVED***

utils: generate
	for util in $(UTILS***REMOVED***; do \
		GOOS=linux GOARCH=$(ARCH***REMOVED*** go build -o ./bin/$$util\_$(ARCH***REMOVED*** ./utils/$$util.go ;\
	done

generate: $(GENERATE***REMOVED***
	for file in $(GENERATE***REMOVED*** ; do \
		go generate $$file ; \
	done

clean:
	go clean
	rm -rf ./bin

copy:
	sshpass scp bin/*_arm ./GoTuringCoffee_arm root@$(DEST_IP***REMOVED***:/home/root

GENERATE = internal/hardware/max31856/max31856.go\
		   internal/hardware/max31865/max31865.go

UTILS = max31856 max31865

.PHONY: arm generate clean

all: generate arm arm-utils

dep:
	dep ensure

arm:
	GOOS=linux GOARCH=arm go build

arm-utils:
	for util in $(UTILS***REMOVED***; do\
		GOOS=linux GOARCH=arm go build  -o ./bin/$$util ./utils/$$util.go ;\
	done

generate: $(GENERATE***REMOVED***
	for file in $(GENERATE***REMOVED*** ; do \
		go generate $$file ; \
	done

clean:
	go clean
	rm -rf ./bin

copy:
	sshpass scp -r bin/ ./GoTuringCoffee root@172.30.200.2:/home/root

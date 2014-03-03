all: pack

.PHONY : clean c-bridge check_mesos check_proto_headers pack
c-bridge:
	cd c-bridge; make all

check_mesos:
	@if [ ! -f /usr/local/lib/libmesos.so ]; then \
  	echo "Error: Expecting libmesos.so in /usr/local/lib"; exit 2; \
	else true; fi
	@if [ ! -d /usr/local/include/mesos ]; then \
  	echo "Error: Expecting mesos headers in /usr/local/include/mesos"; exit 2; \
	else true; fi

check_proto_headers:
	@cd c-bridge;g++ test_headers.cpp 2> /dev/null; if [ $$? -ne 0 ]; then\
  	echo "Error: Expected installed protocol buffer"; exit 2; \
	else true; fi

protos:
	go get code.google.com/p/goprotobuf/proto
	go get code.google.com/p/goprotobuf/protoc-gen-go

chapel: check_proto_headers check_mesos protos c-bridge
	go install mesosphere.io/chapel
	go install mesosphere.io/chapel-agent
	go install mesosphere.io/chapel-client

pack: chapel
	tar -cvzf chapel-bootstrap.tgz bin/chapel-agent

clean:
	@cd c-bridge; make clean
	go clean

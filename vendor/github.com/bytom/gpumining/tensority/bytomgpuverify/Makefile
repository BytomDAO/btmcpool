.PHONY: clean

TARGET=tensority

LIBS=
UNAME := $(shell uname -s)
ifeq ($(UNAME),Darwin)
	LIBS += lib.dylib
else
	LIBS += lib.so
endif

$(TARGET): cSimdTs.a
	# cd .. && go test -v && go test -bench=.
	go run ../../../../../../main.go

cSimdTs.a: cSimdTs.o BytomPoWGpu.o
	ar r $@ $^

%.o: %.cpp
	nvcc -o $@ -c $^ -std=c++11 -I /usr/local/cuda/include -L /usr/local/cuda/lib64 -lcudart -lcublas

%.o: %.cu
	nvcc -o $@ -c $^ -std=c++11 -I /usr/local/cuda/include -L /usr/local/cuda/lib64 -lcudart -lcublas

clean:
	rm -f *.o *.so *.a
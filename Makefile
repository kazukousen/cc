
cc: main.go
	docker run --rm -w /tmp/cc -v $(shell pwd):/tmp/cc tinygo/tinygo:0.21.0 /bin/bash -c \
		'tinygo build -o cc'

test: cc
	docker run --rm -v $(shell pwd):/cc -w /cc compilerbook bash -c \
		'./test.sh'

clean:
	rm -f cc *.o *~ tmp*

.PHONY: test clean

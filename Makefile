bin=cses-cli
.PHONY: clean
all: clean $(bin)

cses-cli:
	echo ${CURDIR}
	GOPATH=${CURDIR} CGO_ENABLED=0 go build -o $(bin)

clean:
	rm -fv $(bin)

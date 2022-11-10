TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=aws.amazon.com
NAMESPACE=terraform
NAME=buildonaws
BINARY=terraform-provider-${NAME}
VERSION=1.0
OS_NAME:=$(shell uname -s | tr ‘[:upper:]’ ‘[:lower:]’)
HW_CLASS:=$(shell uname -m)
OS_ARCH=${OS_NAME}_${HW_CLASS}

default: install

build:
	go build -gcflags="all=-N -l" -o ${BINARY}

install: build
	rm -rf examples/.terraform* || true
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
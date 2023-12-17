.DEFAULT_GOAL := all
PROJECT_NAME := "yit"
BINARY_NAME := "yit"
OBJECT := ".yit"

.PHONY: all clean test test-verbose run build help 

all: clean build test 

clean:
	@echo "Cleaning..."
	@rm -f ${BINARY_NAME} || true  #pass of not found
	@rm -rf ${OBJECT} || true  #pass of not found


build:
	go build -o ${BINARY_NAME} main.go 
	chmod +x ${BINARY_NAME}

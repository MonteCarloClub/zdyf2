CLI=dpki_test

client:
	@go build -o ${CLI} test.go

clean:
	@rm -f ${CLI}
	@rm -f certs.txt
	@rm -f certificates.txt

all: client
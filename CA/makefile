CA=abs_server_ca

ca:
	@go build -o ${CA} ca.go lagRange.go abs.go define.go

clean:
	@rm -f ${CA} nohup.out

all: ca
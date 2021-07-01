RA=abs_server_ra
CLIENT=abs_client

ra:
	@go build -o ${RA} ca.go ra.go lagRange.go router.go define.go

client:
	@go build -o ${CLIENT} client.go define.go rsaT.go ecdsaT.go

clean:
	@rm -f ${RA}
	@rm -f ${CLIENT}

all: ra client
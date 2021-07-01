CLIENT=abs_client

client:
	@go build -o ${CLIENT} client.go define.go rsaT.go ecdsaT.go

clean:
	@rm -f ${CLIENT}

all: ra client
export GO111MODULE=on

all: frp_jwt_allowed_ports

frp_jwt_allowed_ports:
	go build -o ./bin/frp_jwt_allowed_ports ./cmd/frp_jwt_allowed_ports

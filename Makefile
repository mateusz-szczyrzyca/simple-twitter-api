export PATH := $(PWD)/go/bin:$(PATH)
dbdsn := postgresql://root@localhost:26257/twitter?sslmode=disable
dbhost := localhost
endpoint := localhost:58123

environment:
	wget https://dl.google.com/go/go1.12.5.linux-amd64.tar.gz ; \
        tar zxvf go1.12.5.linux-amd64.tar.gz

db:     
	wget -qO- https://binaries.cockroachdb.com/cockroach-v19.1.1.linux-amd64.tgz | tar  xvz ; \
	cockroach-v19.1.1.linux-amd64/cockroach start --insecure --listen-addr=$(dbhost) & 
	sleep 3 && cockroach-v19.1.1.linux-amd64/cockroach sql --insecure < database/cockroachdb.sql

run:
	go run main.go --endpoint "$(endpoint)" --dsn "$(dbdsn)" ; \

#	go run main.go --endpoint "localhost:58123" --dsn "postgresql://root@localhost:26257/twitter?sslmode=disable" ; \


tests:
	cd api_tests/users ; \
	go test -v --endpoint "http://$(endpoint)" -count=1 .
	cd api_tests/messages ; \
	go test -v --endpoint "http://$(endpoint)" -count=1 .


clean:
	pkill -9 go ; \
	pkill -9 main
	pkill -9 cockroach
	rm -rf cockroach-data/ cockroach-v19.1.1.linux-amd64/

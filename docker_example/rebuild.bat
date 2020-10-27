set GOARCH=amd64
set GOOS=linux
cd api
go build -o main -a
cd ../node
go build -o main -a
cd ../registrator
go build -o main -a
cd ../load_balancer
go build -o main -a
cd ../
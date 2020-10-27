set GOARCH=amd64
set GOOS=linux
cd api
go build -o main
cd ../node
go build -o main
cd ../registrator
go build -o main
cd ../load_balancer
go build -o main
cd ../
How to run project:

Run different server/clients command:
File server:
go run fileserver.go

Dispatcher and consolidator server:
go run project2.go -pathname="file.dat" -N=num -C=num

Worker client(Must be ran last as a client needs the server to be running to successfully connect)
go run worker.go -C=num

The command should create running servers and workers

For example:

go run project2.go -pathname="testFile64MB.dat" -N=65536 -C=1024

go run fileserver.go

go run worker.go -C=1024


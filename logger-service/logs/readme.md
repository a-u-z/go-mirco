protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto

如果不能跑的話
* go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
* vim ~/.bash_profile
add
* export GO_PATH=~/go
* export PATH=$PATH:/$GO_PATH/bin
* source ~/.bash_profile

解答來自 https://stackoverflow.com/questions/57700860/error-protoc-gen-go-program-not-found-or-is-not-executable

要加上 grpc 的話
兩邊的 server 都要有 logs 這個檔案

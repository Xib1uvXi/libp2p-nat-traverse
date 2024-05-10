package natinfo

//go:generate protoc --proto_path=$PWD:$PWD/../../.. --go_out=. --go_opt=Mpb/natinfo.proto=./pb pb/natinfo.proto

..\bin\protoc -I=../ --go_out=../src ../proto/client/common.proto
..\bin\protoc -I=../ --go_out=../src ../proto/client/opcode.proto
..\bin\protoc -I=../ --go_out=../src ../proto/client/login.proto

..\bin\protoc -I=../ --go_out=../src ../proto/inner/common.proto
..\bin\protoc -I=../ --go_out=../src ../proto/inner/opcode.proto

..\bin\protoc -I=../ --go_out=../src ../proto/entity/common.proto

pause
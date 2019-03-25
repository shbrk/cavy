..\bin\protoc -I=../ --go_out=../src ../proto/client/msg.proto
..\bin\protoc -I=../ --go_out=../src ../proto/client/opcode.proto

..\bin\protoc -I=../ --go_out=../src ../proto/internal/msg.proto
..\bin\protoc -I=../ --go_out=../src ../proto/internal/opcode.proto
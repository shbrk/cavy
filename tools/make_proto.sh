#!/bin/bash

protoc_bin=../bin/protoc
go_out=../src
proto=../proto

$protoc_bin -I=../ --go_out=../src  $proto/builtin/msg.proto
$protoc_bin -I=../ --go_out=../src  $proto/builtin/opcode.proto
syntax = "proto3";

package {{.package}};
option go_package="./{{.package}}";

message Request {
}

message Response {
  string pong = 1;
}

service {{.serviceName}} {
  rpc Ping(Request) returns(Response);
}

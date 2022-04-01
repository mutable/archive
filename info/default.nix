{ platform, ... }:

platform.buildGo.package {
  name = "github.com/mutable/archive/info";

  srcs = [
    ./types.go
    ./parser.go
  ];
}

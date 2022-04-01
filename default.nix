{ platform, ... }:

platform.buildGo.package {
  name = "github.com/mutable/archive";

  srcs = [
    ./content_reader.go
    ./fs_reader.go
    ./fs_writer.go
    ./parser.go
    ./writer.go
  ];

  deps = [
    platform.lib.nix.wire
  ];
}

{ ... }:

{
  languages.go.enable = true;

  tasks = {
    "app:format".exec = "go fmt ./...";

    "app:format-check".exec = ''
      unformatted="$(gofmt -l .)"
      if [ -n "$unformatted" ]; then
        echo "The following Go files need formatting:"
        echo "$unformatted"
        exit 1
      fi
    '';

    "app:test".exec = "go test ./...";
    "app:vet".exec = "go vet ./...";

    "app:check" = {
      exec = "echo 'All checks passed.'";
      after = [
        "app:format-check"
        "app:test"
        "app:vet"
      ];
    };

    "app:run".exec = "go run ./cmd/api";
  };

  processes.api.exec = "go run ./cmd/api";
}

default: build

# Read does not work with watch mode
watch:
  go build -o main *.go || exit 1 && ./main

build:
  go build -o main *.go

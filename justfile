default: build

watch:
  npx nodemon -e go  --exec "go build -o main *.go || exit 1"

build:
  go build -o main *.go
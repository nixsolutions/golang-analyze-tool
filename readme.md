# Go demo "Analyze tool"

## Requirements
1. Go v1.12 or higher

## Installation 
Clone this repository:
`git clone ...`
That's it! If you want to run the code you, please, run `go run main.go -p=9090` or build a binary `go build main.go` and then run `./main -p=9090`. Go to `localhost:9090` to test the app.

## Analyze your site
Analyze tool receives 3 params:
1. URL - a link to site you want to check.
2. Depth - how deep the analyze tool must scan your site.
3. Threads - how many threads the analyze tool will use.

Example of the request: `localhost:9090?url=https://example.com&depth=1&threads=3`
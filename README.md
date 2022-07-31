# Web Crawler
###  Concurrently crawls all the urls linked to a starting url within the given host and prints them in console.
- Uses worker pool to utilize go routines.
## Install dependencies

```bash
go mod download
```

## How to Run

```bash
go run *.go
```

### OR

```bash     
go build -o crawler; ./crawler
```

- After running please enter the starting URl and press enter.

## Note:
- The program meets the minimum requirements of the assignment.
- Developed on Go version 1.16.

#### TODO:
- add a timeout
- add a max depth (if needed)
- read the robots.txt file?
- better error handling
- more optimizations (better concurrent workers)
- save the results to a file
- package it in docker to deploy it to a cloud
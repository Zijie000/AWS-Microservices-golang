# webapp
This is for csye6225 project


## Prerequisites

🪄You need to install 🦫golang 1.10 ~ 1.23

🪄You need to install mySQL (Any stable version within 10 years)

## Build and Deploy

🪄Get all golang/Gin dependencies needed
```bash
go mod tidy
```

🪄Launch the webapp
```bash
go run .
```

🧚‍♀️Compile into executable file in Linux, you can also configure the GOOS to windows or macos
```bash
GOOS=linux GOARCH=amd64 go build -o myprogram
```

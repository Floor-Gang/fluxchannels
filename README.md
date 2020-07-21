# Flux Channels
## Setup
Download [Go](https://golang.org)
```shell script
$ git clone https://github.com/Floor-Gang/fluxchannels
$ go mod download
$ cd ./cmd/fluxchannels
$ go build
$ ./fluxchannels
# ... edit config.yml ...
$ ./fluxchannels
```
 
## Bot Usage
To get an ID of a category enable "developer mode" in appearance settings of
your Discord client. Then right click the category name and click "Copy ID."

To add a new fluctuating category:
 * `.flux add <category ID>`

To remove a fluctuating category:
 * `.flux remove <category ID>`

To list all the fluctuating categories:
 * `.flux list`
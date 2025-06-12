# snake.go
### A snake.io inspired game for the terminal

## To run
This is still largely a work in progress, so even my makefile isn't fully configured yet. To start up a host simply run `go run .` 
You can change the mode from default (host) with the --mode flag setting it to `server` for a headless server or `client` to join another server as a player.
`client` players must specify an address of the server they wish to connect to via --addr and passing the address as a string

## Example
Player 1 wants to be host so he/she simply runs `go run .` with starts the server on default port 8080 (or configured with --port) and joins as a player.
Player 2 is on the same network and joins Player 1's game with `go run --mode client --addr 10.0.0.1:8080` or whatever their local DHCP address is (can be found with `ifconfig`)

# Conway's Game of Life

I created this as an exercise to learn some golang.

This implementation of the game is played on a torus, like the old-school video game "Asteroids".

## Build
```bash
go build -o life ./src/main.go
```

## Example usage:
```bash
# get help on cmdline args
./life --help

# start with a canned initial condition and board size
./life

# start with an acorn that grows into a bigger colony and takes a while to reach steady state
./life -i ./src/acorn.txt

# same thing but on a much smaller board. it will die sooner!
./life -i ./src/acorn.txt -x 10 -y 10

# the famous glider
./life -i ./src/glider.txt
```
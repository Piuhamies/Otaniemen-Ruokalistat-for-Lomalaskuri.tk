package main

import (
	"bufio"
	"fmt"
	"os"

	"otaniemenruokalistat.tk/webapi"
)

func main() {
	fmt.Println("Initializing...")
	go webapi.Init()
	fmt.Println("PRESS ENTER TO EXIT")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

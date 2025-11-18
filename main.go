package main

import (
	"fmt"
	"os"
	"io"
)

func main() {
	file, err := os.Open("messages.txt") // open the file
	if err != nil { // handle error
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close() // ensure the file is closed when done

	for {
		buf := make([]byte, 8) // buffer to hold 8 bytes
		n, err := file.Read(buf) // read up to 8 bytes

		if n > 0 {	// if we read some bytes
			fmt.Printf("read: %s\n", string(buf[:n])) // print only the bytes read
		}

		if err != nil {
			if err == io.EOF { // end of file
				break
			}
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}
	}
	
}

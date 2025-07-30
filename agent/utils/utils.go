package utils

import (
	"bufio"
	"fmt"
	"os"
)

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		return ("UnknownHost")
	}
	return hostname
}

// writerRoutine handles all writes to the server
func WriterRoutine(writer *bufio.Writer, dataChan <-chan string) {
	for data := range dataChan {
		_, err := writer.WriteString(data + "\n")
		if err != nil {
			fmt.Println("Error writing data:", err)
			return
		}
		writer.Flush()
	}
}

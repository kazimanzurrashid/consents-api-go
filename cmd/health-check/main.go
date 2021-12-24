package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	res, err := http.Get(fmt.Sprintf("http://localhost:%s/", os.Getenv("PORT")))

	if err != nil || res.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}

package main

import "T/app"

func main() {
	err := app.InitApp()
	if err != nil {
		panic(err)
	}
}

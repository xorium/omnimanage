package main

import (
	"flag"
	"log"
	"os"
)

const (
	ModeController = "model-controller"
)

func main() {
	//init flags
	mode := flag.String("mode", ModeController, "generation mode, default - model-controller")
	modelName := flag.String("model", "", "model name")
	modelFile := flag.String("file_in", "", "input file name")

	flag.Parse()

	if len(*mode) == 0 {
		log.Fatal("mode is empty")
	}
	if len(*modelName) == 0 {
		log.Fatal("model is empty")
	}
	if len(*modelFile) == 0 {
		log.Fatal("file is empty")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	switch *mode {
	case ModeController:
		err := runModelControllerGenerator(cwd, *modelFile, *modelName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

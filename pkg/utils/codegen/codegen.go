package main

import (
	"flag"
	"log"
	"omnimanage/pkg/utils/model_parser"
	"os"
)

const (
	ModeController       = "model-controller"
	ModeStore            = "model-store"
	ModeStoreInterface   = "model-store-interface"
	ModeService          = "model-service"
	ModeServiceInterface = "model-service-interface"
)

type ModelDescription struct {
	Name            string
	CompanyResource bool
	Relations       []*model_parser.Relation
}

// go generate ./pkg/model/domain/.

func main() {
	//init flags
	mode := flag.String("mode", ModeController, "generation mode, default - model-controller")
	modelName := flag.String("model", "", "model name")
	modelFile := flag.String("file_in", "", "input file name")
	toStdOut := flag.Bool("no_file_output", false, "generating code to stdout")
	companyResource := flag.Bool("company_resource", true, "has company relation")
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
		err := runModelControllerGenerator(cwd, *modelFile, *modelName, *toStdOut, *companyResource)
		if err != nil {
			log.Fatal(err)
		}

	case ModeStore:
		err := runModelStoreGenerator(cwd, *modelFile, *modelName, *toStdOut, *companyResource)
		if err != nil {
			log.Fatal(err)
		}

	case ModeStoreInterface:
		err := runModelStoreInterfaceGenerator(cwd, *modelFile, *modelName, *companyResource)
		if err != nil {
			log.Fatal(err)
		}

	case ModeService:
		err := runModelServiceGenerator(cwd, *modelFile, *modelName, *toStdOut, *companyResource)
		if err != nil {
			log.Fatal(err)
		}
	case ModeServiceInterface:
		err := runModelServiceInterfaceGenerator(cwd, *modelFile, *modelName, *companyResource)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func getModelDescription(modelName string, file string, companyResource bool) (*ModelDescription, error) {
	p, err := model_parser.NewParser(file)
	if err != nil {
		return nil, err
	}

	rels, err := p.GetRelations(modelName)
	if err != nil {
		return nil, err
	}

	m := &ModelDescription{
		Name:            modelName,
		CompanyResource: companyResource,
		Relations:       rels,
	}

	return m, nil
}

package model_parser

import (
	"fmt"
	"github.com/fatih/structtag"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

type Relation struct {
	Name     string
	WebName  string
	Multiple bool
}

// structType contains a structType node and it's name. It's a convenient
// helper type, because *ast.StructType doesn't contain the name of the struct
type structType struct {
	name string
	node *ast.StructType
}

type ModelParser struct {
	fset    *token.FileSet
	astFile *ast.File
}

// NewParser create new ModelParser instance
func NewParser(fileName string) (*ModelParser, error) {
	p := &ModelParser{
		fset: token.NewFileSet(),
	}

	astFile, err := parser.ParseFile(p.fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	p.astFile = astFile
	return p, nil
}

//
func (p *ModelParser) GetRelations(structName string) (res []*Relation, errOut error) {
	structStartLine, structEndLine, err := p.SelectStruct(structName, p.astFile)
	if err != nil {
		return nil, err
	}

	res = make([]*Relation, 0, 1)
	ast.Inspect(p.astFile, func(n ast.Node) bool {
		x, ok := n.(*ast.StructType)
		if !ok {
			return true
		}

		for _, f := range x.Fields.List {
			currLine := p.fset.Position(f.Pos()).Line
			if !(structStartLine <= currLine && currLine <= structEndLine) {
				continue
			}

			if f.Tag == nil {
				f.Tag = &ast.BasicLit{}
			}

			fieldName := ""
			if len(f.Names) != 0 {
				fieldName = f.Names[0].Name
			}

			// anonymous field
			if f.Names == nil {
				ident, ok := f.Type.(*ast.Ident)
				if !ok {
					continue
				}

				fieldName = ident.Name
			}

			_, isArray := f.Type.(*ast.ArrayType)

			newRel, err := p.getRelation(fieldName, f.Tag.Value, isArray)
			if err != nil {
				errOut = err
				continue
			}
			if newRel != nil {
				res = append(res, newRel)
			}
		}

		return true
	})
	if errOut != nil {
		return nil, errOut
	}

	if len(res) > 0 {
		return res, nil
	}

	return nil, nil
}

func (p *ModelParser) getRelation(fieldName string, tagValue string, isArray bool) (*Relation, error) {
	if tagValue == "" {
		return nil, nil
	}

	tagVal, err := strconv.Unquote(tagValue)
	if err != nil {
		return nil, err
	}

	tags, err := structtag.Parse(tagVal)
	if err != nil {
		return nil, fmt.Errorf("Field %v, tag %v, error: %w ", fieldName, tagValue, err)
	}

	tag, err := tags.Get("jsonapi")
	if err != nil {
		return nil, nil
	}

	if tag.Name != "relation" {
		return nil, nil
	}

	newR := &Relation{
		Name: fieldName,
	}
	if len(tag.Options) >= 1 {
		newR.WebName = tag.Options[0]
	}
	newR.Multiple = isArray

	return newR, nil
}

// SelectStruct gets start and end points of struct with name structName
func (p *ModelParser) SelectStruct(structName string, node ast.Node) (int, int, error) {
	structs := collectStructs(node)

	var encStruct *ast.StructType
	for _, st := range structs {
		if st.name == structName {
			encStruct = st.node
		}
	}

	if encStruct == nil {
		return 0, 0, fmt.Errorf("struct %v does not exist", structName)
	}

	start := p.fset.Position(encStruct.Pos()).Line
	end := p.fset.Position(encStruct.End()).Line

	return start, end, nil
}

// collectStructs collects and maps structType nodes to their positions
func collectStructs(node ast.Node) map[token.Pos]*structType {
	structs := make(map[token.Pos]*structType, 0)
	collectStructs := func(n ast.Node) bool {
		t, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if t.Type == nil {
			return true
		}

		structName := t.Name.Name

		x, ok := t.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structs[x.Pos()] = &structType{
			name: structName,
			node: x,
		}
		return true
	}
	ast.Inspect(node, collectStructs)
	return structs
}

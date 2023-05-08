package main

import (
	"fmt"
	"go/types"
	"os"
	"strings"

	"github.com/cloudquery/plugin-sdk/v2/caser"
	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

// do not try to import and make local copies of the following types
var skipPackageImports = []string{"time.Time"}

type client struct {
	generatedTypes []string
	f              *jen.File
	caser          caser.Caser
}

func (c *client) addGeneratedType(genType string) {
	c.generatedTypes = append(c.generatedTypes, genType)
}

func (c *client) checkGeneratedType(genType string) bool {
	return contains(c.generatedTypes, genType) || contains(skipPackageImports, genType)
}

// Function takes in the full path of th struct you want to copy locally
func (c *client) copyType(sourceType string) {
	if c.checkGeneratedType(sourceType) {
		return
	}
	// Split into package and variable name
	sourceTypePackage, sourceTypeName := splitSourceType(sourceType)

	pkg := loadPackage(sourceTypePackage)

	// Ensure that the variable is declared in that package
	obj := pkg.Types.Scope().Lookup(sourceTypeName)
	if obj == nil {
		panic(fmt.Errorf("%s not found in declared types of %s", sourceTypeName, pkg))
	}

	// check if it is a declared type
	if _, ok := obj.(*types.TypeName); !ok {
		panic(fmt.Errorf("%v is not a named type", obj))
	}
	// only support copying strings and structs
	switch v := obj.Type().Underlying().(type) {
	case *types.Struct:
		err := c.generateStruct(sourceTypeName, v)
		if err != nil {
			panic(err)
		}
	case *types.Basic:
		c.generateString(sourceTypeName)
	}
}

func genPackage(srcPkgUrl, fileName string) error {
	// Naming convention: <service>_input.go
	service := strings.Split(fileName, ".")[0] + "_input"
	input := strings.Split(fileName, ".")[1]

	// instantiate client
	// In the future if this gets large this will help being able to parallelize
	c := &client{
		f:     jen.NewFile(service),
		caser: *caser.New(),
	}
	// Add comment at top of file
	c.f.PackageComment("Code generated by generator, DO NOT EDIT.")

	//
	c.copyType(srcPkgUrl)
	dir := "../table_options/inputs/" + service + "/"
	targetFilename := input + ".go"
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	// 7. Write generated file
	err = c.f.Save(dir + targetFilename)
	if err != nil {
		return fmt.Errorf("writing output: %v", err)
	}
	return nil
}

func main() {
	srcPkgUrls := []string{
		"github.com/aws/aws-sdk-go-v2/service/costexplorer.GetCostAndUsageInput",
		"github.com/aws/aws-sdk-go-v2/service/inspector2.ListFindingsInput",
		"github.com/aws/aws-sdk-go-v2/service/accessanalyzer.ListFindingsInput",
		"github.com/aws/aws-sdk-go-v2/service/cloudtrail.LookupEventsInput",
	}
	for _, srcPkgUrl := range srcPkgUrls {
		split := strings.Split(srcPkgUrl, "/")
		fileName := split[len(split)-1]
		err := genPackage(srcPkgUrl, fileName)
		if err != nil {
			panic(err)
		}
	}
}

func loadPackage(path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		panic(fmt.Errorf("loading packages for inspection: %v", err))
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0]
}

func splitSourceType(sourceType string) (string, string) {
	idx := strings.LastIndexByte(sourceType, '.')
	if idx == -1 {
		panic(fmt.Errorf(`expected qualified type as "pkg/path.MyType"`))
	}
	sourceTypePackage := sourceType[0:idx]
	sourceTypeName := sourceType[idx+1:]
	return sourceTypePackage, sourceTypeName
}

func (c *client) generateString(sourceTypeName string) {
	if contains(c.generatedTypes, sourceTypeName) {
		return
	}
	c.addGeneratedType(sourceTypeName)
	c.f.Type().Id(sourceTypeName).String()
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (c *client) generateStruct(sourceTypeName string, structType *types.Struct) error {
	if c.checkGeneratedType(sourceTypeName) {
		return nil
	}
	c.addGeneratedType(sourceTypeName)
	var changeSetFields []jen.Code

	// 4. Iterate over struct fields
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		name := field.Name()
		if name == "noSmithyDocumentSerde" {
			continue
		}

		// Generate code for each changeset field
		code := jen.Id(name)
		switch v := field.Type().(type) {
		case *types.Slice:
			switch vNested := v.Elem().(type) {
			case *types.Basic:
				code.Op("[]").Id(vNested.String())
			case *types.Named:
				typeName := vNested.Obj()
				c.copyType(typeName.Pkg().Path() + "." + typeName.Name())
				code.Op("[]").Id(typeName.Name())
			default:
				return fmt.Errorf("struct field type not handled: %T", vNested)
			}
		case *types.Map:
			switch vNested := v.Elem().(type) {
			case *types.Basic:
				// TODO this is a hack, once we need to support maps with keys other than string we will need to handle this
				code.Op("map[string]").Id(vNested.String())
			case *types.Named:
				nestedTypeName := vNested.Obj()
				// Skip copying the type if it is the same type
				c.copyType(nestedTypeName.Pkg().Path() + "." + nestedTypeName.Name())
				code.Map(jen.Id(v.Key().String())).Id(nestedTypeName.Name())

			default:
				return fmt.Errorf("struct field type not handled: %T", vNested)
			}
		case *types.Pointer:
			switch vNestedType := v.Elem().(type) {
			case *types.Basic:
				code.Op("*").Id(vNestedType.String())
			case *types.Named:
				typeName := vNestedType.Obj()
				path := typeName.Pkg().Path()
				c.copyType(path + "." + typeName.Name())
				// This is a hack to skip copying types from the standard library
				if strings.HasPrefix(path, "github.com/") {
					code.Op("*").Id(typeName.Name())
				} else {
					code.Op("*").Qual(path, typeName.Name())
				}

			default:
				return fmt.Errorf("struct field type not handled: %T", vNestedType)
			}

		case *types.Basic:
			code.Op("*").Id(v.String())
		case *types.Named:
			typeName := v.Obj()
			c.copyType(typeName.Pkg().Path() + "." + typeName.Name())
			code.Op("*").Id(typeName.Name())
		default:
			return fmt.Errorf("struct field type not handled: %T", v)
		}
		changeSetFields = append(changeSetFields, code.Tag(map[string]string{"json": c.caser.ToSnake(name) + "," + "omitempty"}))
	}

	c.f.Type().Id(sourceTypeName).Struct(changeSetFields...)

	return nil
}

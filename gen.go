package main

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	endpoint "github.com/TwiN/gatus/v5/config/endpoint"
)

var (
	apiDir         = "api/v1alpha1/"
	definitionFile = "generated.gatus.go"
)

func main() {
	typ := reflect.TypeOf(endpoint.Endpoint{})
	structs := make(map[string]string)

	generateStruct(typ, structs)

	var sb strings.Builder
	sb.WriteString("package v1alpha1\n\n")
	sb.WriteString("import ( \"net/http\" )\n\n")
	for _, structDef := range structs {
		sb.WriteString("// +k8s:deepcopy-gen=true\n" + structDef + "\n")
	}

	err := os.WriteFile(apiDir+definitionFile, []byte(sb.String()), 0644)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("Go struct definitions written to %s\n", definitionFile)
}

func generateStruct(typ reflect.Type, structs map[string]string) {
	typString := typ.String()
	if strings.Contains(typ.String(), ".") {
		var re = regexp.MustCompile(`(?m)^[^a-z]*?([a-z]*)\.([a-zA-Z0-9]*)$`)
		typString = re.ReplaceAllString(typString, "$1$2")
	}

	typString = CapitalizeFirstLetter(typString)

	if _, exists := structs[typString]; exists {
		return
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("type %s struct {\n", typString))

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name
		fieldType := field.Type
		fieldTag := field.Tag

		fieldTypeString := fieldType.String()

		if field.Type.PkgPath() == "net/http" {
			continue
		}

		if strings.Contains(fieldTypeString, ".") {
			var re = regexp.MustCompile(`(?m)^[^a-z]*?([a-z]*)\.([a-zA-Z0-9]*)$`)
			fieldTypeString = re.ReplaceAllString(fieldTypeString, "$1$2")

			if fieldType.Kind() != reflect.Func {
				fieldTypeString = CapitalizeFirstLetter(fieldTypeString)
			}
		}

		if yamlTag := fieldTag.Get("yaml"); yamlTag != "" {
			newTag := strings.Replace(string(fieldTag), "yaml", "json", 1)
			fieldTag = reflect.StructTag(newTag)
		}

		if fieldType.Kind() == reflect.Struct {
			generateStruct(fieldType, structs)
		} else if fieldType.Kind() == reflect.Slice {
			elemType := fieldType.Elem()
			if elemType.Kind() == reflect.Struct || elemType.Kind() == reflect.Ptr {
				generateStruct(elemType, structs)
			} else {
				fieldTypeString = fmt.Sprintf("[]%s", elemType.Kind())
			}

			if fieldTag != "" {
				sb.WriteString(fmt.Sprintf("\t%s %s `%s`\n", fieldName, fieldTypeString, fieldTag))
			} else {
				sb.WriteString(fmt.Sprintf("\t%s %s `json:\"-\"`\n", fieldName, fieldTypeString))
			}
		} else if fieldType.Kind() == reflect.Ptr {
			elemType := fieldType.Elem()

			if elemType.Kind() == reflect.Struct && fieldName != "httpClient" {
				generateStruct(elemType, structs)
			}

			if fieldName == "httpClient" {
				fieldTypeString = "http.Client"
			}

			if fieldTag != "" {
				sb.WriteString(fmt.Sprintf("\t%s %s `%s`\n", fieldName, fieldTypeString, fieldTag))
			} else {
				sb.WriteString(fmt.Sprintf("\t%s %s `json:\"-\"`\n", fieldName, fieldTypeString))
			}
		} else {
			if fieldType.Kind() == reflect.Func || fieldType.Kind() == reflect.Interface || fieldType.Kind() == reflect.Map {
				if fieldType.Kind() == reflect.Map {
					fieldTag = "json:\"-\"" // Ignore maps for now
				}

				if fieldTypeString == "map[string]interface {}" {
					fieldTypeString = "map[string]any"
				}

				if fieldTag != "" {
					sb.WriteString(fmt.Sprintf("\t%s %s `%s`\n", fieldName, fieldTypeString, fieldTag))
				} else {
					sb.WriteString(fmt.Sprintf("\t%s %s `json:\"-\"`\n", fieldName, fieldTypeString))
				}
			} else {
				fieldTypeGeneric := fieldType.Kind().String()
				if strings.Contains(fieldType.String(), "time.Duration") {
					fieldTypeGeneric = "string"
				}

				if fieldTag != "" {
					sb.WriteString(fmt.Sprintf("\t%s %s `%s`\n", fieldName, fieldTypeGeneric, fieldTag))
				} else {
					sb.WriteString(fmt.Sprintf("\t%s %s `json:\"-\"`\n", fieldName, fieldTypeGeneric))
				}
			}
		}
	}

	sb.WriteString("}\n")

	structs[typString] = sb.String()
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}

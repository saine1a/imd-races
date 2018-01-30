package csvmapper

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

type fieldmapping struct {
	name  string
	field string
	index int
}

// MappingInfo Opaque MappingInfo object
type MappingInfo struct {
	mappings []fieldmapping
	typ      reflect.Type
}

// Overide tag handling to support multiple tags

type Tag reflect.StructTag

func (tag Tag) GetMultiple(key string) []string {
	v, _ := tag.Lookup(key)
	return v
}

func (tag Tag) Lookup(key string) ([]string, bool) {
	// When modifying this code, also update the validateStructTag code
	// in cmd/vet/structtag.go.

	exist := false
	values := []string{}

	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			exist = true
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}

			values = append(values, value)
		}
	}

	return values, exist
}

func mapHeaderRecursive(record []string, typ reflect.Type) []fieldmapping {

	fields := make([]fieldmapping, 0, len(record))

	foundCount := 0

	for f := 0; f < typ.NumField(); f++ {

		field := typ.Field(f)

		if field.Type.Kind() == reflect.Struct && field.Type.Name() != "Time" {
			structFields := mapHeaderRecursive(record, field.Type)

			fields = append(fields, structFields...)

			fmt.Printf("%s is a struct\n", field.Name)
		} else {

			// First look for csvField tag to match field name, and if that's not there use the struct name

			tagDefn := Tag(field.Tag)
			tag := tagDefn.GetMultiple("csvField")

			if len(tag) == 0 {
				tag = append(tag, field.Name)
			}

			// Now look for a match

			found := false

			for t := 0; t < len(tag); t++ {
				for r := 0; r < len(record); r++ {
					if tag[t] == record[r] {
						mapping := fieldmapping{field.Name, tag[t], r}
						fields = append(fields, mapping)

						if found == false { // logic for finding first mapping for this particular field
							found = true
							foundCount++
						}
					}
				}
			}

			if !found {
				fmt.Printf("Could not find field %s with tag(s)\n", field.Name)
				for t := range tag {
					if t > 0 {
						fmt.Print(",")
					}
					fmt.Printf("%s", tag[t])
				}
			}
		}
	}

	if foundCount < typ.NumField() {
		fmt.Println("ALERT : MISSING FIELDS")
	}

	return fields
}

// MapHeader - Function to map a record to its fields
func MapHeader(record []string, typ reflect.Type) MappingInfo {

	fields := mapHeaderRecursive(record, typ)

	return MappingInfo{fields, typ}
}

// ParseRecord - function to convert a record list into a structure
func (m *MappingInfo) ParseRecord(record []string, mystructptr interface{}, dateFormat string) bool {

	if reflect.TypeOf(mystructptr) != reflect.PtrTo(m.typ) {
		log.Fatalf("Expected type %s but got type %s\n", reflect.PtrTo(m.typ), reflect.TypeOf(mystructptr))

		return false
	}

	// Cast to the structure to update
	mystruct := reflect.ValueOf(mystructptr).Elem()

	// Run over all the expected mappings

	for _, f := range m.mappings {

		if f.name != "" && record[f.index] != "" {

			field := mystruct.FieldByName(f.name)

			switch field.Type().Kind() {
			case reflect.String:
				field.SetString(record[f.index])
			case reflect.Int64:
				i, err := strconv.ParseInt(record[f.index], 10, 64)
				if err != nil {
					log.Fatal("Bad int value " + record[i])
					os.Exit(-1)
				}
				field.SetInt(i)
			case reflect.Float64:
				fp, err := strconv.ParseFloat(record[f.index], 64)
				if err != nil {
					log.Fatal("Bad float value " + record[f.index])
					os.Exit(-1)
				}
				field.SetFloat(fp)
			case reflect.Struct:
				if field.Type().Name() == "Time" {

					theTime, err := time.Parse(dateFormat, record[f.index])
					if err == nil {

						field.Set(reflect.ValueOf(theTime))
					} else {
						fmt.Printf("Bad time %s\n", record[f.index])
						log.Fatal("Bad time \n" + record[f.index])
						os.Exit(-1)
					}
					//	handle times
				} else {
					log.Fatal("Unexpected struct type " + field.Type().Name())
					os.Exit(-1)
				}
			default:
				log.Fatal("Unexpected field type " + field.Type().Kind().String())
				os.Exit(-1)
			}

		}
	}

	return true
}

package main

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var postfix = map[string]string{}

func recursivePrintType(t reflect.Type, indent int) {
	if postfix[t.Name()] != "" || t.Kind() != reflect.Struct {
		return
	}

	postfix[t.Name()] = "placeholder"
	log.Print("Analysing type: ", t.Name())
	var id = 1
	var s = ""
	s += ("message " + strings.Split(t.String(), ".")[1] + "{\n")
	for i := 0; i < t.NumField(); i++ {
		obj := t.Field(i).Type
		s += ("\t")
		if obj.Kind() == reflect.Array || obj.Kind() == reflect.Slice {
			obj = obj.Elem()
			s += ("repeated ")
		}

		//check if it is a pointer
		if obj.Kind() == reflect.Ptr {
			obj = obj.Elem()
		}

		if obj.Kind() == reflect.Struct && obj.String() != "time.Time" && obj.Name() != t.Name() {
			recursivePrintType(obj, indent+1)
			s += (strings.Split(obj.String(), ".")[1] + " " + t.Field(i).Name + " = " + strconv.Itoa(id))

		} else if obj.Kind() == reflect.Map {
			log.Print("Map: ", obj.String())
			recursivePrintType(obj.Key(), indent+1)
			recursivePrintType(obj.Elem(), indent+1)
			key := strings.Split(obj.Key().String(), ".")
			if len(key) > 1 {
				key = key[1:]
			}
			elem := strings.Split(obj.Elem().String(), ".")
			if len(elem) > 1 {
				elem = elem[1:]
			}
			s += ("map <" + key[0] + "," + elem[0] + ">") + " " + t.Field(i).Name + " = " + strconv.Itoa(id)
			c := obj.Elem()
			if c.Kind() == reflect.Struct && c.String() != "time.Time" {
				recursivePrintType(c, indent+1)
			}
		} else {
			type_str := obj.String()
			if obj.String() == "time.Time" {
				s += "//time.Time\n"
				s += ("\t")
				type_str = "uint64"
			}
			if type_str == "int" {
				type_str += "32"
			}

			s += (type_str + " " + t.Field(i).Name + " = " + strconv.Itoa(id))
		}

		json_name := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		if json_name != "-" && json_name != "" {
			s += " [json_name=\"" + json_name + "\"]; \n"
		} else {
			s += ";\n"
		}
		id++
	}
	s += ("}\n")
	postfix[t.Name()] = s
}
func main() {
	recursivePrintType(reflect.TypeOf(discordgo.AddedThreadMember{}), 0)
	log.Print("\n\n\n")
	for _, s := range postfix {
		println(s)
	}

}

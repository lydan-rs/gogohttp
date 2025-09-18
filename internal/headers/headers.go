package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)



const crlf = "\r\n"

type Headers map[string]string

func (h Headers) Exists(name string) (bool) {
	_, exists := h[name]
	return exists
}

func (h Headers) Get(name string) (value string, exists bool) {
	value, exists = h[name]
	return
}

func MakeHeadersMap() Headers {
	return Headers{}
}

const (
	r_sp = `\x20`
	r_htab = `\x09`
	r_vchar = `[\x21-\x7E]`
	r_field_vchar = `[\x21-\x7E\x80-\xFF]`
	r_delimiters = `["()/:;<=>?@\[\]{}]`
	r_tchar = `[!#$%&'*+-.^_` + "`" + `a-zA-Z0-9|~]` 
	r_token = r_tchar + "+"
	r_field_content = r_field_vchar + `([` + r_sp + r_htab + r_field_vchar + `]+` + r_field_vchar + `)?`
	r_field_value = "(" + r_field_content + ")*"
	r_field_name = r_token
	r_field_line = "^" + r_field_name + ":" + r_sp + "*" + r_field_value + r_sp + "*$"
)

var field_line_regex, regex_err = regexp.Compile(r_field_line)

/*
FIELD-LINE: [field-name]:<\s*>[field-value]<\s*>
*/

func (h Headers) Parse(data []byte) (bytesParsed int, done bool, err error) {
	bytesParsed = 0
	done = false
	err = nil

	crlfIndex := bytes.Index(data, []byte(crlf))

	for crlfIndex > -1 {
		// Found the pre-body CRLF. Header parsing done.
		if crlfIndex == 0 {
			bytesParsed += 2
			done = true
			err = nil
			return
		}

		lineSlice := data[bytesParsed:bytesParsed+crlfIndex]
		fmt.Printf("Field Line:\n>> Line:%v\n>> Num Bytes: %v\n", string(lineSlice), crlfIndex)
		bytesParsed += crlfIndex+2

		isValidFieldLine := field_line_regex.Match(lineSlice)
		if !isValidFieldLine {
			err = fmt.Errorf("Invalid Field Line >> %v", string(lineSlice))
			return
		}

		colonIndex := bytes.IndexRune(lineSlice, ':')
		if colonIndex < 0 {
			err = fmt.Errorf("Seperator ':' not found. Validation and parse failed.")
			return
		}

		fieldName := strings.ToLower(string(lineSlice[:colonIndex]))
		fieldValue := string(lineSlice[colonIndex+1:])
		fieldValue = strings.TrimLeft(fieldValue, " ")
		fieldValue = strings.TrimRight(fieldValue, " ")

		if _, exists := h[fieldName]; exists {
			h[fieldName] = h[fieldName] + ", " + fieldValue 
		} else {
			h[fieldName] = fieldValue
		}

		crlfIndex = bytes.Index(data[bytesParsed:], []byte(crlf))
	}
	
	return
}

package headers

import "bytes"



const crlf = "\r\n"

type Headers map[string]string

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
		if crlfIndex == bytesParsed {
			bytesParsed += 2
			done = true
			err = nil
			return
		}

		bytesParsed += crlfIndex+2

		parts = bytes.Split(data, )
	}
	
	return
}

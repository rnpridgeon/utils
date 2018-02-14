package encoding

import (
	"io"
)

// The stdlib decoder walks down follows each node all the way down. I just want the nodes at the current level.
// Internationalization be damned this reader strips JSON blobs one layer at a time,
// Not recommended for serious use, I just wanted a way to make JSON parsing more dynamic
// JSON grammar - lite (structural characters only)
// https://tools.ietf.org/html/rfc7159#page-4
const (
	object_open      = '{'
	object_close     = '}'
	object_seperator = ','
	kv_seperator     = ':'
	array_open       = '['
	array_close      = ']'
	string_          = '"'
)

type reader struct {
	buffer []byte
	position int
	token    *token
	isString    bool
}

func NewJSONStripper(buffer []byte) (*reader) {
	return &reader{
		buffer : buffer,
		position: 0,
		token : &token{},
		isString: false,
	}
}

// Returns next token
func (f *reader) NextMember() ([] byte) {
	//move cursor past key
	defer f.next()

	for f.isString == false && f.next == nil {
		f.next()
	}
	return f.buffer[f.token.start:f.token.stop]
}

// Move cursor in
func (f *reader) GetMembers() {
	f.seek(f.token.start + 1)
	f.isString = false
	f.next()
}

type token struct {
	start int
	stop  int
}

// Tokenize JSON blob at current level
func (f *reader) next() error {
	f.token.start = 0
	f.token.stop = 0

	if f.position >= len(f.buffer)-1 {
		return io.EOF
	}

	defer func() { f.position++ }()

	closer := func(open byte, close byte) {
		f.token.start = f.position
		f.position++

		quoted := false
		for i := 1; i > 0 && f.position < len(f.buffer); f.position++ {
			if f.buffer[f.position-1] == '\\' {
				continue
			}

			if f.buffer[f.position] == '"' {
				quoted = !quoted
			}

			if f.buffer[f.position] == open && !quoted {
				i++
			} else if f.buffer[f.position] == close && !quoted {
				i--
			}
		}
		f.token.stop = f.position
	}

	// floor is inclusive ceiling is exclusive; set offset accordingly
	for ; f.position < len(f.buffer); f.position++ {
		currentByte := f.buffer[f.position]
		switch currentByte {
		case kv_seperator:
			// ensure token did not appear within string, like https://
			if !f.isString {
				return nil
			}
			continue
		case object_seperator:
			if f.token.stop > 0 {
				return nil
			}
			continue
		case object_open:
			closer(object_open, object_close)
			return nil
		case array_open:
			closer(array_open, array_close)
			return nil
		case string_:
			f.isString = !f.isString
			if f.token.start == 0 {
				f.token.start = f.position + 1
			}
			f.token.stop = f.position
		default:
			// isNumber (ASCII)
			if currentByte > 47 && currentByte < 58 {
				if f.token.start == 0 {
					f.token.start = f.position
					//include position position
					f.token.stop = f.token.start + 1
					//continue
				}
				f.token.stop = f.position + 1
			}
			continue
		}
	}
	return nil
}

// Move cursor to index value
func (f *reader) seek(index int) {
	f.position = index
}

package gobf

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Interpreter struct {
	// Data cells
	cells [30000]byte

	// Data pointer
	dp int

	// The source code to execute
	Code []byte

	// Program counter
	pc int

	// The stream to write to
	out io.Writer

	// The stream to read from
	in io.Reader
}

func NewInterpreterFromFile(path string) *Interpreter {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("an error occured when opening '%s'", path)
	}

	return &Interpreter{
		Code: code,
		out:  os.Stdout,
		in:   os.Stdin,
	}
}

// SetOutputStream allows to change the output stream,
// to any io.Writer provider.
func (i *Interpreter) SetOutputStream(out io.Writer) { i.out = out }

// SetInputStream allows to change the user input stream,
// to any io.Reader provider.
func (i *Interpreter) SetInputStream(in io.Reader) { i.in = in }

// Exec executes the Interpreter from the beginning
func (i *Interpreter) Exec() {
	reader := bufio.NewReader(i.in)
	for i.pc = 0; i.pc < len(i.Code); i.pc++ {
		switch i.Code[i.pc] {
		case '>':
			i.dp++
		case '<':
			i.dp--
		case '+':
			i.cells[i.dp]++
		case '-':
			i.cells[i.dp]--
		case '.':
			_, err := i.out.Write([]byte{i.cells[i.dp]})
			if err != nil {
				log.Fatal("input error")
			}
		case ',':
			buf, err := reader.ReadByte()
			if err != nil {
				log.Fatal("input error")
			}

			i.cells[i.dp] = buf
		case '[':
			if i.cells[i.dp] == 0 {
				for p := i.pc; p < len(i.Code); p++ {
					if i.Code[p] == ']' {
						i.pc = p
						break
					}

					if p == len(i.Code)-1 {
						log.Fatal("no matching '['")
					}
				}
			}
		case ']':
			if i.cells[i.dp] != 0 {
				for p := i.pc; p > 0; p-- {
					if i.Code[p] == '[' {
						i.pc = p
						break
					}

					if p == 0 {
						log.Fatal("no matching '['")
					}
				}
			}
		default:
			i.pc++
		}
	}
}

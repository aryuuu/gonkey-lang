package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/aryuuu/gonkey-lang/evaluator"
	"github.com/aryuuu/gonkey-lang/lexer"
	"github.com/aryuuu/gonkey-lang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, "%s", PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.GetErrors()) != 0 {
			printParserError(out, p.GetErrors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserError(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err+"\n")
	}
}

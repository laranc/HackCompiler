package translator

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	static  = 16
	pointer = 3
	temp    = 5
)

var (
	functionCalls = -1
	eqCalls       = -1
	gtCalls       = -1
	ltCalls       = -1
	segments      = [4]string{"LCL", "ARG", "THIS", "THAT"}
)

func Run(path string) (string, error) {
	fmt.Println("Translating...")
	code := ""
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		endl := regexp.MustCompile(`\r*\n+`)
		space := regexp.MustCompile(`\s+`)
		line = endl.ReplaceAllString(line, "")
		line = space.ReplaceAllString(line, " ")
		tokens := strings.Split(line, " ")
		switch len(tokens) {
		case 1:
			switch tokens[0] {
			case "add":
				code += vmAdd()
			case "sub":
				code += vmSub()
			case "neg":
				code += vmNeg()
			case "eq":
				code += vmEq()
			case "gt":
				code += vmGt()
			case "lt":
				code += vmLt()
			case "and":
				code += vmAnd()
			case "or":
				code += vmOr()
			case "not":
				code += vmNot()
			case "return":
				code += vmReturn()
			}
		case 2:
			switch tokens[0] {
			case "label":
				code += vmLabel(tokens[1])
			case "goto":
				code += vmGoto(tokens[1])
			case "if-goto":
				code += vmIf(tokens[1])
			}
		case 3:
			n, err := strconv.Atoi(tokens[2])
			if err != nil {
				return "", err
			}
			switch tokens[0] {
			case "push":
				code += vmPush(tokens[1], n)
			case "pop":
				code += vmPop(tokens[1], n)
			case "function":
				code += vmFunction(tokens[1], n)
			case "call":
				code += vmCall(tokens[1], n)
			}
		}
	}
	fmt.Println(code)
	functionCalls = -1
	eqCalls = -1
	gtCalls = -1
	ltCalls = -1
	return code, nil
}

func vmAdd() string {
	return "@SP\nAM=M-1\nD=M\nA=A-1\nM=D+M\n"
}

func vmSub() string {
	return "@SP\nAM=M-1\nD=M\nA=A-1\nM=M-D\n"
}

func vmNeg() string {
	return "@SP\nAM=M-1\nD=M\nA=A-1\nM=-M\n"
}

func vmEq() string {
	eqCalls++
	return fmt.Sprintf(`@SP
		AM=M-1
		D+M
		A=A-1
		D=M-D
		@EQ_TRUE%d
		D;JEQ
		@SP
		A=M-1
		M=0
		@EQ_END%d
		0;JMP
		(EQ_TRUE%d)
		@SP
		A=M-1
		M=-1
		(EQ_END%d)
		`, eqCalls, eqCalls, eqCalls, eqCalls)
}

func vmGt() string {
	gtCalls++
	return fmt.Sprintf(`@SP
		AM=M-1
		D=M
		A=A-1
		D=M-D
		@GT_TRUE%d
		D;JGT
		@SP
		A=M-1
		M=0
		@GT_END%d
		0;JMP
		(GT_TRUE%d)
		@SP
		A=M-1
		M=-1
		(GT_END%d)
		`, gtCalls, gtCalls, gtCalls, gtCalls)
}

func vmLt() string {
	ltCalls++
	return fmt.Sprintf(`@SP
		AM=M-1
		D=M
		A=A-1
		D=M-D
		@LT_TRUE%d
		D;JLT
		@SP
		A=M-1
		M=0
		@LT_END%d
		0;JMP
		(LT_TRUE%d)
		@SP
		A=M-1
		M=-1
		(LT_END%d)
		`, ltCalls, ltCalls, ltCalls, ltCalls)
}

func vmAnd() string {
	return "@SP\nAM=M-1\nD=M\nA=A-1\nM=D&M\n"
}

func vmOr() string {
	return "@SP\nAM=M-1\nD=M\nA=A-1\nM=D|M\n"
}

func vmNot() string {
	return "@SP\nA=M-1\nM=!M\n"
}

func vmReturn() string {
	return `@LCL
	D=M
	@R13
	M=D
	@5
	A=D-A
	D=M
	@R14
	M=D
	@SP
	AM=M-1
	D=M
	@ARG
	A=M
	M=D
	@ARG
	D=M+1
	@SP
	M=D
	@R13
	AM=M-1
	D=M
	@THAT
	M=D
	@R13
	AM+M-1
	D=M
	@THIS
	M=D
	@R13
	AM=M-1
	D+M
	@ARG
	M=D
	@R13
	AM=M-1
	D=M
	@LCL
	M=D
	@R14
	A=M
	0;JMP
	`
}

func vmLabel(label string) string {
	return fmt.Sprintf("(%s)\n", label)
}

func vmGoto(label string) string {
	return fmt.Sprintf("@%s\n0;JMP\n", label)
}

func vmIf(label string) string {
	return fmt.Sprintf("@SP\nAM=M-1\nD=M\n@%s\nD;JNE", label)
}

func vmPush(segment string, offset int) string {
	line := "@"
	switch segment {
	case "constant":
		line += fmt.Sprintf("%d\nD=A\n", offset)
	case "static":
		line += fmt.Sprintf("%d\nD=M\n", static+offset)
	case "pointer":
		line += fmt.Sprintf("%d\nD=M\n", pointer+offset)
	case "temp":
		line += fmt.Sprintf("%d\nD=M\n", temp+offset)
	default:
		line += fmt.Sprintf("%d\nD=A\n", offset)
		switch segment {
		case "local":
			line += "@LCL\n"
		case "argument":
			line += "@ARG\n"
		case "this":
			line += "@THIS\n"
		case "that":
			line += "@THAT\n"
		}
		line += "A=D+M\nD=M\n"
	}
	line += "@SP\nA=M\nM=D\n@SP\nM=M+1\n"
	return line
}

func vmPop(segment string, offset int) string {
	line := "@SP\nAM=M-1\nD=M\n@"
	switch segment {
	case "static":
		line += fmt.Sprintf("%d\n", static+offset)
	case "pointer":
		line += fmt.Sprintf("%d\n", pointer+offset)
	case "temp":
		line += fmt.Sprintf("%d\n", temp+offset)
	default:
		line += fmt.Sprintf("%d\nD=A\n", offset)
		switch segment {
		case "local":
			line += "@LCL"
		case "argument":
			line += "@ARG"
		case "this":
			line += "@THIS"
		case "that":
			line += "@THAT"
		}
		line += "D=D+1\n@13\nM=D\n@SP\nAM=M-1\nD=M\n@13\nA=M\n"
	}
	line += "M=D\n"
	return line
}

func vmFunction(name string, vars int) string {
	result := fmt.Sprintf("(%s\n)", name)
	for range vars {
		result += "@SP\nA=M\nM=0\n@SP\nM=M+1\n"
	}
	return result
}

func vmCall(name string, vars int) string {
	functionCalls++
	return fmt.Sprintf(`@%s%d
		D=A
		@SP
		A=M
		M=D
		@SP
		M=M+1
		@LCL
		D=M
		@SP
		A=M
		M=D
		@SP
		M=M+1
		@ARG
		D=M
		@SP
		A=M
		M=D
		@SP
		M=M+1
		@THIS
		D=M
		@SP
		A=M
		M=D
		@SP
		M=M+1
		@THAT
		D=M
		@SP
		A=M
		M=D
		@SP
		M=M+1
		@SP
		D=M
		@%d
		D=D-A
		@5
		D=D-A
		@ARG
		M=D
		@SP
		D=M
		@LCL
		M=D
		@%s
		0;JMP
		(%s%d)
		`, name, functionCalls, vars, name, name, functionCalls)
}

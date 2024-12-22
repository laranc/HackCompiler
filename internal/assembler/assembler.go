package assembler

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	compLookup = map[string]uint16{
		"0": 0b0101010, "1": 0b0111111, "-1": 0b0111010,
		"D": 0b0001100, "A": 0b0110000, "M": 0b1110000,
		"!D": 0b0001101, "!A": 0b0110001, "!M": 0b1110001,
		"-D": 0b0001111, "-A": 0b0110011, "-M": 0b1110011,
		"D+1": 0b0011111, "A+1": 0b0110111, "M+1": 0b1110111,
		"D-1": 0b0001110, "A-1": 0b0110010, "M-1": 0b1110010,
		"D+A": 0b0000010, "D+M": 0b1000010, "D-A": 0b0010011,
		"D-M": 0b1010011, "A-D": 0b0000111, "M-D": 0b1000111,
		"D&A": 0b0000000, "D&M": 0b1000000, "D|A": 0b0010101,
		"D|M": 0b1010101,
	}
	dstLookup = map[string]uint16{
		"null": 0b0000, "M": 0b001, "D": 0b010, "MD": 0b011,
		"A": 0b100, "AM": 0b101, "AD": 0b110, "AMD": 0b111,
	}
	jmpLookup = map[string]uint16{
		"null": 0b000, "JGT": 0b001, "JEQ": 0b010, "JGE": 0b011,
		"JLT": 0b100, "JNE": 0b101, "JLE": 0b110, "JMP": 0b111,
	}
	symbolTable map[string]int
	varOffset   = 16
)

func Run(path string) ([]uint16, error) {
	fmt.Println("Assembling...")
	bin := make([]uint16, 0, 10)
	code := list.New()
	symbolTable = make(map[string]int)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	fmt.Println("First Pass {")
	for scanner.Scan() {
		if scanner.Err() != nil {
			return nil, err
		}
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}
		parsedLine := findPass(line, count)
		if parsedLine != "" {
			code.PushBack(parsedLine)
			fmt.Printf("\t%s\n", parsedLine)
		}
		count++
	}
	fmt.Println("\n}\nSecond Pass {")
	replacePass(code)
	fmt.Println("\n}\nThird Pass {")
	for e := code.Front(); e != nil; e = e.Next() {
		line := e.Value.(string)
		ins, err := assembleInstruction(line)
		if err != nil {
			return nil, err
		}
		bin = append(bin, ins)
		fmt.Printf("\t%016b\n", ins)
	}
	fmt.Println("}")
	varOffset = 16
	return bin, nil
}

func findPass(line string, num int) string {
	parsedLine := ""
	switch line[0] {
	case '/':
		return ""
	case '@':
		if line[1] == 'R' {
			n, err := strconv.Atoi(line[2:])
			if err != nil {
				return ""
			}
			parsedLine = fmt.Sprintf("@%d", n)
		} else if line[1:] == "LCL" {
			parsedLine = "@1"
		} else if line[1:] == "ARG" {
			parsedLine = "@2"
		} else if line[1:] == "THIS" {
			parsedLine = "@3"
		} else if line[1:] == "THAT" {
			parsedLine = "@4"
		} else if line[1:] == "SCREEN" {
			parsedLine = "@16384"
		} else if line[1:] == "KBD" {
			parsedLine = "@24576"
		} else if _, err := strconv.Atoi(line[1:]); err != nil {
			if _, ok := symbolTable[line[1:]]; !ok {
				symbolTable[line[1:]] = varOffset
				varOffset++
			}
			return line
		} else {
			break
		}
		return parsedLine
	case '(':
		symbol := line[1 : len(line)-1]
		symbolTable[symbol] = num
		return ""
	}
	return line
}

func replacePass(code *list.List) {
	for e := code.Front(); e != nil; e = e.Next() {
		line := e.Value.(string)
		switch line[0] {
		case '@':
			if _, err := strconv.Atoi(line[1:]); err != nil {
				if addr, ok := symbolTable[line[1:]]; ok {
					line = fmt.Sprintf("@%d", addr)
					e.Value = line
				}
			}
		}
		fmt.Printf("\t%s\n", line)
	}
}

func assembleInstruction(ins string) (uint16, error) {
	if ins[0] == '@' {
		n, err := strconv.Atoi(ins[1:])
		if err != nil {
			return 0, fmt.Errorf("%s: %e", ins, err)
		}
		return uint16(n) & 0b011_1_111111_111_111, nil
	}
	comp := ""
	dst := "null"
	jmp := "null"
	if strings.Contains(ins, "=") {
		s := strings.SplitN(ins, "=", 2)
		dst = s[0]
		ins = s[1]
	}
	if strings.Contains(ins, ";") {
		s := strings.SplitN(ins, ";", 2)
		comp = s[0]
		jmp = s[1]
	} else {
		comp = ins
	}
	c, ok := compLookup[comp]
	if !ok {
		return 0, fmt.Errorf("Invalid computation: %s", ins)
	}
	d, ok := dstLookup[dst]
	if !ok {
		return 0, fmt.Errorf("Invalid destination: %s", ins)
	}
	j, ok := jmpLookup[jmp]
	if !ok {
		return 0, fmt.Errorf("Invalid jump: %s", ins)
	}
	code := uint16(0b111_0_000000_000_000)
	code |= c << 6
	code |= d << 3
	code |= j
	return code, nil
}

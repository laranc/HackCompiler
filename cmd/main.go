package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/laranc/HackCompiler/internal/assembler"
	"github.com/laranc/HackCompiler/internal/compiler"
	"github.com/laranc/HackCompiler/internal/translator"
)

func main() {
	if len(os.Args) < 1 {
		panic("No input provided\n")
	}
	for _, p := range os.Args[1:] {
		if _, err := os.Stat(p); err != nil {
			fmt.Printf("%s: %e\n", p, err)
			continue
		}
		s := strings.SplitN(p, "/", 2)
		s = strings.SplitN(s[1], ".", 2)
		name := s[0]
		ext := s[1]
		var err error
		switch ext {
		case "jack":
			err = jack(p, name)
		case "vm":
			err = vm(p, name)
		case "asm":
			err = asm(p, name)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}

func jack(path, name string) error {
	hackVM, err := compiler.Run(path)
	if err != nil {
		return err
	}
	out := fmt.Sprintf("out/%s.jack", name)
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, ins := range hackVM {
		if _, err := f.WriteString(ins); err != nil {
			return err
		}
	}
	return vm(out, name)
}

func vm(path, name string) error {
	hackASM, err := translator.Run(path)
	if err != nil {
		return err
	}
	out := fmt.Sprintf("out/%s.asm", name)
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(hackASM); err != nil {
		return err
	}
	return asm(out, name)
}

func asm(path, name string) error {
	hackML, err := assembler.Run(path)
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf("out/%s.bin", name))
	if err != nil {
		return err
	}
	defer f.Close()
	for _, word := range hackML {
		if err := binary.Write(f, binary.LittleEndian, word); err != nil {
			return err
		}
	}
	return nil
}

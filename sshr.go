package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func checkFileIsExist(filename string) bool {
	s, _ := os.Stat(filename)
	if s == nil {
		return false
	}
	return true
}

func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	r := bufio.NewReader(f)
	for {
		// ReadLine is a low-level line-reading primitive.
		// Most callers should use ReadBytes('\n') or ReadString('\n') instead or use a Scanner.
		bytes, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return lines, err
		}
		lines = append(lines, string(bytes))
	}
	return lines, nil
}

func WriteLines(filePath string, lines []string) bool {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666) //0666 在windos下无效
	if err != nil {
		fmt.Println("open file err:", err)
		return false
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	writer.Flush()
	return true
}

func main() {

	var sip string
	flag.StringVar(&sip, "i", "", "")
	flag.Parse()

	//fmt.Println("flags:", flag.Args())

	if len(sip) < 1 && len(flag.Args()) > 0 {
		sip = flag.Args()[0]
	}

	if len(sip) < 1 {
		fmt.Println("usage: sshr <ip>")
		return
	}

	drive := os.Getenv("HOMEDRIVE")
	homepath := os.Getenv("HOMEPATH")

	homepath = path.Join(drive, homepath)

	//fmt.Println("homepath:", homepath)

	file := path.Join(homepath, "/.ssh/known_hosts")

	if !checkFileIsExist(file) {
		fmt.Println("no ", file, " found")
		return
	}

	lines, _ := ReadLines(file)

	iii := -1
	for i, line := range lines {
		ii := strings.Index(line, " ")
		if ii < 1 {
			continue
		}
		tip := line[:ii]
		if tip == sip {
			iii = i
			break
		}
	}
	if iii < 0 {
		for i, line := range lines {
			ii := strings.Index(line, " ")
			if ii < 1 {
				continue
			}
			tip := line[:ii]
			if strings.Contains(tip, sip) {
				ii := strings.Index(tip, sip)
				if ii+len(sip) == len(tip) {
					iii = i
					break
				}
			}
		}
	}

	if iii >= 0 {
		fmt.Println("entry located:", lines[iii])
	}

	//fmt.Println("lines:", lines)

	if iii >= 0 {
		var rare []string
		rare = append(rare, lines[iii+1:]...)
		lines2 := append(lines[:iii], rare...)

		WriteLines(file, lines2)
	} else {
		fmt.Println("on entry located!")
	}
}

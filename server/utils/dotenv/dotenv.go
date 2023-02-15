package dotenv

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, err := readLine(reader)
	for err == nil {
		parts := strings.Split(strings.Split(line, "#")[0], "=")
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
		line, err = readLine(reader)
	}
	if err != io.EOF {
		return err
	}
	return nil
}

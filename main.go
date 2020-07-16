package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	getStr         = "GET "
	max200sSamePkg = 1
	max403sSamePkg = 1
)

var (
	ignorePath = regexp.MustCompile(`\/ajax\/libs\/.+\w\/\d+.\d+.\d+\/.+\w.(js|css|map)`)
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	file, err := os.Open("raw.txt")
	check(err)
	defer file.Close()

	pkg200s := make(map[string]int)    // # of accepted pkgs with 200 code
	pkg403s := make(map[string]int)    // # of accepted pkgs with 403 code
	entries := make(map[string]int)    // accepted pkgs
	ignoredCodes := make(map[int]bool) // for debugging

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		res := strings.Index(text, getStr)
		if res == -1 {
			continue // ignore non-get reqs
		}
		text = text[res+len(getStr):]
		strs := strings.Split(text, " ")
		code, err := strconv.Atoi(strs[2])
		check(err)
		if code != 200 && code != 403 {
			ignoredCodes[code] = true
			continue
		}
		p := strs[0]
		if match := ignorePath.MatchString(p); match {
			continue
		}
		var pkgName string
		if strings.HasPrefix(p, "/ajax/libs/") {
			ss := strings.Split(p, "/")
			pkgName = ss[3]
			p = strings.Replace(p, "/ajax/libs/", "/ajax/libs/test-path/", 1)
		} else if strings.HasPrefix(p, "//ajax/libs/") {
			ss := strings.Split(p, "/")
			pkgName = ss[4]
			p = strings.Replace(p, "//ajax/libs/", "//ajax/libs/test-path/", 1)
		}
		if code == 200 {
			if pkg200s[pkgName] >= max200sSamePkg {
				continue
			}
			pkg200s[pkgName]++
		} else {
			if pkg403s[pkgName] >= max403sSamePkg {
				continue
			}
			pkg403s[pkgName]++
		}
		entries[p] = code
	}
	check(scanner.Err())

	bytes, err := json.Marshal(entries)
	check(err)

	fmt.Printf("%s\n", bytes)
}

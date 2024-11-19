package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var seenBrowsers = make(map[string]struct{})
	uniqueBrowsers, i := 0, 0
	foundUsers := strings.Builder{}

	//fileContents, err := ioutil.ReadAll(file)
	//if err != nil {
	//	panic(err)
	//}
	//
	//lines := strings.Split(string(fileContents), "\n")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		user := User{}
		// fmt.Printf("%v %v\n", err, line)
		if err := json.Unmarshal([]byte(line), &user); err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browser := range user.Browsers {

			if strings.Contains(browser, "Android") {
				isAndroid = true

				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true

				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.Replace(user.Email, "@", " [at] ", -1)
		foundUsers.WriteString("[" + strconv.Itoa(i) + "]" + " " + user.Name + " " + "<" + email + ">" + "\n")
		// foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
		i++
	}

	_, _ = fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	_, _ = fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

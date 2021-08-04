package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	exception        []string
	filelist         []string
	fileChan         chan string
	finish           string
	userContribution map[string]int32
)

func init() {
	finish = "@@@done"
	fileChan = make(chan string)
	userContribution = make(map[string]int32, 100)
	exception = []string{"vendor", "vendor-patched"}
}
func main() {
	rootdir, err := os.Getwd()
	if err == nil {
		var grp sync.WaitGroup
		grp.Add(1)
		go walkDir(rootdir, true)
		go populateFilieList(&grp)
		grp.Wait()
		close(fileChan)
		for _, vv := range filelist {
			//time.Sleep(time.Second * 1)
			if strings.HasSuffix(vv, ".go") || strings.HasSuffix(vv, ".h") || strings.HasSuffix(vv, ".c") {
				//fmt.Println(vv)
				getGitStats(vv)
			}
		}
		for user, contri := range userContribution {
			fmt.Printf("%s==>%d\n ", user, contri)
		}
	}

}
func getGitStats(file string) {

	cmd := exec.Command("git", "blame", file)

	//cmd.Stdin = strings.NewReader("and old falcon")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return
	}
	output := out.String()
	spltOut := strings.Split(output, "\n")
	for _, line := range spltOut {
		//fmt.Println("line  ", line)
		rawuser := strings.Split(line, " ")
		if len(rawuser) > 1 {
			for _, uname := range rawuser {
				uname = strings.Replace(uname, " ", "", -1)
				if strings.HasPrefix(uname, "(") {
					uname = strings.Replace(uname, "(", "", -1)
					uname = strings.ToLower(uname)
					userContribution[uname] = userContribution[uname] + 1
					break
				}
			}
			//fmt.Println("rawuser  ", rawuser[1])
		}

	}
}
func populateFilieList(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case a := <-fileChan:
			if a == finish {
				return
			}
			filelist = append(filelist, a)
		}
	}
}

//IsExceptionDir checks against exception list
func IsExceptionDir(dirname string) bool {
	for _, v := range exception {
		if v == dirname {
			return true
		}
	}
	return false
}

// directoryElements returns the entries in the dir directory
func directoryElements(dir string) []os.FileInfo {
	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du：%v\n", err)
		return nil
	}
	defer f.Close()
	entries, err := f.Readdir(0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du：%v\n", err)
	}
	return entries
}

func walkDir(dir string, closeChan bool) {
	for _, entry := range directoryElements(dir) {
		dirElement := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if IsExceptionDir(entry.Name()) {
				continue
			}
			walkDir(dirElement, false)
		} else {
			fileChan <- dirElement
		}
	}
	if closeChan {
		fileChan <- finish
	}
}

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	prsr "beprodigy.com/parser"
	"beprodigy.com/utility"
)

var (
	filelist         []string
	fileChan         chan string
	finish           string
	userContribution map[string]int32
)

func init() {
	finish = "@@@done"
	fileChan = make(chan string)
	userContribution = make(map[string]int32, 100)
}
func main() {
	var mainerr error
	var rankedList []utility.ContribV
	params := os.Args[1:]
	for _, param := range params {
		mainerr = prsr.ProcessKeyVal(param)
		if mainerr != nil {
			fmt.Println(mainerr)
		}
	}
	if len(prsr.Extensions) == 0 {
		prsr.Extensions = append(prsr.Extensions, ".go")
		prsr.Extensions = append(prsr.Extensions, ".c")
		prsr.Extensions = append(prsr.Extensions, ".h")
	}
	rootdir, err := os.Getwd()
	if err == nil {
		var grp sync.WaitGroup
		go utility.ShowProgressWheel()
		grp.Add(1)
		go walkDir(rootdir, true)
		go populateFilieList(&grp)
		grp.Wait()
		close(fileChan)
		for _, vv := range filelist {
			//time.Sleep(time.Second * 1)
			if isExtensionCovered(vv) {
				//fmt.Println(vv)
				getGitStats(vv)
			}
		}

		utility.StopProgressWheel <- true
		rankedList = utility.RankContributors(userContribution)
		tsize := utility.ListSize
		for _, contri := range rankedList {
			u := len(contri.User)
			space := tsize - u
			ss := ""
			for ii := 0; ii < space; ii++ {
				ss = ss + " "
			}
			fmt.Printf("%s%s%d\n", contri.User, ss, contri.Contribution)
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
	for _, v := range prsr.ExceptionFileDire {
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
func isExtensionCovered(filename string) bool {
	for _, v := range prsr.Extensions {
		if strings.HasSuffix(filename, v) {
			return true
		}
	}
	return false
}

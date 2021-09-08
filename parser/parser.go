package parser

import (
	"fmt"
	"strings"
)

const (
	exception       = "xcpn"
	lookupextension = "xtsn"
)

//ExceptionFileDire files n directories to be excluded while calculating the contribV
var ExceptionFileDire []string

//Extensions is a list of extension
var Extensions []string

//ProcessKeyVal takes arguements as value and process them
func ProcessKeyVal(kv string) (err error) {

	keyVal := strings.Split(kv, "=")
	if len(keyVal) == 2 {
		if exception == keyVal[0] {
			xnVal := strings.Split(keyVal[1], ",")
			for ii := 0; ii < len(xnVal); ii++ {
				ExceptionFileDire = append(ExceptionFileDire, xnVal[ii])
			}
			return nil
		} else if lookupextension == keyVal[0] {
			xnVal := strings.Split(keyVal[1], ",")
			for ii := 0; ii < len(xnVal); ii++ {
				Extensions = append(Extensions, xnVal[ii])
			}
			return nil
		}

	}
	return fmt.Errorf("WF")
}

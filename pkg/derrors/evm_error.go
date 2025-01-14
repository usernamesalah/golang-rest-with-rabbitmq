package derrors

import (
	"strings"
)

func HandleEVMErr(err error, format string, args ...any) error {
	if strings.Contains(err.Error(), "VM Exception") {
		return New(InvalidArgument, err.Error())
	} else if strings.Contains(err.Error(), "no contract code") {
		return New(NotFound, err.Error())
	} else if strings.Contains(err.Error(), "insufficient funds") {
		return New(InvalidArgument, err.Error())
	} else {
		return WrapStack(err, Unknown, format, args...)
	}
}

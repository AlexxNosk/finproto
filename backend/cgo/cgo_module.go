package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"

	"github.com/alexxnosk/finproto/backend/trade_api/data"
)

//export InstrumentCreateC
func InstrumentCreateC(jsonInput *C.char, token *C.char) *C.char {
	goJSON := C.GoString(jsonInput)
	goToken := C.GoString(token)
	fmt.Println("Loaded TOKEN:", goToken[:4]+"***")
	err := data.AssetCreate([]byte(goJSON), goToken)
	if err != nil {
		errMsg := fmt.Sprintf("Go Error: %v", err)
		return C.CString(errMsg)
	}
	return C.CString("OK")
}

func main() {}

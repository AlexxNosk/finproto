package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"github.com/alexxnosk/finproto/backend/config"
	"github.com/alexxnosk/finproto/backend/trade_api/data"
)

//export InstrumentCreateC
func InstrumentCreateC(jsonInput *C.char) *C.char {
	goJSON := C.GoString(jsonInput)

	err := data.InstrumentCreate([]byte(goJSON))
	if err != nil {
		errMsg := fmt.Sprintf("Go Error: %v", err)
		return C.CString(errMsg)
	}
	return C.CString("OK")
}

//export Hello
func Hello() {
	fmt.Println("Hello from Go!")
	cfg := config.LoadConfig()
	fmt.Println("Loaded TOKEN:", cfg.TOKEN[:4]+"***")
}

func main() {}

package main

import (
	"fmt"
	"obfuscator"
)

var script = `
var arr = [
    "Apple",
    "Banana",
    "Pear"
];

for (var v in arr) {
    console.log(arr[v])
}
`

func main() {
	result, err := obfuscator.Obfuscate([]byte(script))

	if err != nil {
		fmt.Errorf(err.Error())
	} else {
		fmt.Println(result)
	}
}

package deployment

import "fmt"

func goPanic(err error, desc string) {
	if err != nil {
		panic(fmt.Errorf("%s, error=%+v", desc, err))
	}
}

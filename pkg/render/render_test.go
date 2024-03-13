package render

import (
	"fmt"
	"testing"
)

func TestCreateCache(t *testing.T) {
	var err error

	NewCache(&appTest)

	app.Template, err = CreateTmplCache()
	if err != nil {
		t.Error("CreateTemplate Cash test fails:", err)
	} else {
		fmt.Println("Create cache Passed")
	}
	var ww myWriter

	for _, value := range []bool{true, false} {
		app.Prd = value

		err := Render(ww, "res-sum.page.tmpl", nil)
		if err != nil {
			t.Error("Failed test on cache = ", value)
		} else {
			fmt.Println("Render Passed test on cache = ", value)
		}
	}

}

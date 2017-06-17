package genstr

import (
	"fmt"
	"testing"
)

func BenchmarkGenStr(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c := Config{"Goisanopensourceprogramminglanguagethatmakesiteasytobuildsimplereliableandefficientsoftware", 7}
		r := RunExperiment(&c)
		fmt.Println(r)
	}
}

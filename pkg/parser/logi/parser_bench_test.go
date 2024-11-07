package logi

import (
	"testing"
)

func BenchmarkSimpleParse(b *testing.B) {
	var macroInput = `
				macro entity {
					kind Syntax

					syntax {
						<propertyName Name> <propertyType Type> [primary bool, autoincrement bool, required bool, default string]
					}
				}
`
	var logiInput = `
				entity User {
					id int <[primary, autoincrement]>
					name string <[required, default "John Doe"]>
				}
			`

	// parallel benchmark
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := ParseFullWithMacro(logiInput, macroInput)

			if err != nil {
				b.Error(err)
			}
		}
	})
}

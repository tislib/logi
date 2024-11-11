package logi

import (
	"github.com/tislib/logi/pkg/parser/macro"
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

	mAst, err := macro.ParseMacroContent(macroInput)

	if err != nil {
		b.Error(err)
		return
	}

	// parallel benchmark
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := Parse(logiInput, mAst.Macros)

			if err != nil {
				b.Error(err)
			}
		}
	})
}

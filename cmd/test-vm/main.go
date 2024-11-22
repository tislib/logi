package main

import "github.com/tislib/logi/pkg/vm"

func main() {
	var v, err = vm.New()

	if err != nil {
		panic(err)
	}

	// load Macro, which will be used to parse actual logi content
	_ = v.LoadMacroContent(`
		macro conversation {
			kind Syntax
			
			syntax {
				greeting { expr }
				farewell { expr }
			}
		}`)

	// load Logi content, based on given Macro DSL
	def, err := v.LoadLogiContent(`
		conversation SimpleOp {
			greeting { "Hello, " + name }
			farewell { "Goodbye, " + name }
		}
`)

	// execute the loaded Logi content
	v.SetLocal("name", "John Doe")
	result, err := v.Execute(&def[0], "greeting")

	println("Greeting: " + result.(string))

	result, _ = v.Execute(&def[0], "farewell")

	println("Farewell: " + result.(string))

}

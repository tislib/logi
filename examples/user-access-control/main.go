package main

import (
	_ "embed"
	"github.com/tislib/logi/pkg/ast/common"
	logiAst "github.com/tislib/logi/pkg/ast/logi"
	"github.com/tislib/logi/pkg/vm"
	"log"
)

//go:embed example1.lg
var example1 string

//go:embed macros.lgm
var example1Macro string

type PermissionKind string

const (
	PermissionKindRead      PermissionKind = "read"
	PermissionKindWrite     PermissionKind = "write"
	PermissionKindDenyWrite PermissionKind = "deny_write"
	PermissionKindCrud      PermissionKind = "crud"
)

type permission struct {
	Kind   PermissionKind
	object string

	owned      bool
	properties []string
}

type conditionalPermissions struct {
	condition   string
	permissions []permission
}

type roleImplementer struct {
	description            string
	permissions            []permission
	conditionalPermissions []conditionalPermissions
}

func (r *roleImplementer) Call(vm vm.VirtualMachine, statement logiAst.Statement) error {
	switch statement.Command {
	case "description":
		r.description = statement.GetParameter("description").AsString()
	case "permissions":
		for _, sub := range statement.SubStatements[0] {
			if err := r.CallPermission(vm, sub); err != nil {
				return err
			}
		}
	default:
		log.Fatalf("unknown command: %s", statement.Command)
	}

	return nil
}

func (r *roleImplementer) CallPermission(vm vm.VirtualMachine, statement logiAst.Statement) error {
	switch statement.Command {
	case "WRITE":
		var properties []string

		if statement.SubStatements != nil {
			for _, sub := range statement.SubStatements[0] {
				properties = append(properties, sub.GetParameter("properties").AsString())
			}
		}

		r.permissions = append(r.permissions, permission{
			Kind:       PermissionKindWrite,
			object:     statement.GetParameter("object").AsString(),
			properties: properties,
		})
	case "READ":
		var properties []string

		propertiesParam := statement.GetParameter("properties")

		log.Println(propertiesParam)

		r.permissions = append(r.permissions, permission{
			Kind:       PermissionKindRead,
			object:     statement.GetParameter("object").AsString(),
			properties: properties,
		})
	case "DENY_WRITE":
		var properties []string

		propertiesParam := statement.GetParameter("properties")

		log.Println(propertiesParam)

		r.permissions = append(r.permissions, permission{
			Kind:       PermissionKindDenyWrite,
			object:     statement.GetParameter("object").AsString(),
			properties: properties,
		})
	case "CRUD":
		var properties []string

		propertiesParam := statement.GetParameter("properties")

		log.Println(propertiesParam)

		r.permissions = append(r.permissions, permission{
			Kind:       PermissionKindCrud,
			object:     statement.GetParameter("object").AsString(),
			properties: properties,
		})
	case "when":
		var expression *common.Expression

		for _, param := range statement.Parameters {
			if param.Name == "condition" {
				expression = param.Expression
			}
		}

		var mainPermissions = r.permissions
		r.permissions = nil

		for _, sub := range statement.SubStatements[0] {
			if err := r.CallPermission(vm, sub); err != nil {
				return err
			}
		}

		var conditionPermissions = r.permissions
		r.permissions = mainPermissions

		r.conditionalPermissions = append(r.conditionalPermissions, conditionalPermissions{
			condition:   expression.FuncCall.Name,
			permissions: conditionPermissions,
		})

	default:
		log.Fatalf("unknown command: %s", statement.Command)
	}

	return nil
}

func main() {
	var v = vm.New()

	if err := v.LoadMacroContent(example1Macro); err != nil {
		log.Fatal(err)
	}

	def, err := v.LoadLogiContent(example1)

	if err != nil {
		log.Fatal(err)
	}

	var impl = &roleImplementer{}
	err = v.Execute(&def[1], impl)

	if err != nil {
		log.Fatal(err)
	}
}

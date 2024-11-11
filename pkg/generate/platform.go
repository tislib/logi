package generate

import (
	"fmt"
	"github.com/tislib/logi/pkg/generate/common"
	"github.com/tislib/logi/pkg/generate/golang"
)

func GetCodeGenerator(platform string) (common.CodeGenerator, error) {
	switch platform {
	case "golang":
		return golang.NewGenerator(), nil
	default:
		return nil, fmt.Errorf("unknown platform: %s", platform)
	}
}

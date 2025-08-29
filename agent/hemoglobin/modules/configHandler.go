package modules

import (
	"fmt"
	"regexp"

	"github.com/TLop503/LogCrunch/structs"
)

// HandleConfigTarget determines if each target is using a custom module,
// and then either initializes the module or pulls from the registry
func HandleConfigTarget(target structs.Target) (structs.ParserModule, error) {
	if target.Custom {
		re, err := regexp.Compile(target.Regex)
		if err != nil {
			return structs.ParserModule{}, fmt.Errorf("invalid regex for %s: %w", target.Name, err)
		}
		return structs.ParserModule{
			Regex:  re,
			Schema: target.Schema,
		}, nil
	} else {
		module, ok := structs.MetaParserRegistry[target.Module]
		if !ok {
			return structs.ParserModule{}, fmt.Errorf("no registry entry for module %s", target.Module)
		}
		return module, nil
	}
}

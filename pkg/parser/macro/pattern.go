package macro

import "regexp"

var NamePattern = regexp.MustCompilePOSIX(`^[a-z][a-zA-Z0-9_]*$`)

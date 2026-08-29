package mssql

import (
	"fmt"
	"strings"

	"github.com/rilldata/rill/runtime/drivers"
)

type dialect struct {
	drivers.BaseDialect
}

var DialectMSSQL drivers.Dialect = func() drivers.Dialect {
	d := &dialect{}
	d.BaseDialect = drivers.NewBaseDialect(drivers.DialectNameMSSQL, EscapeIdentifier, EscapeIdentifier)
	return d
}()

func (d *dialect) SupportsILike() bool {
	return false
}

func EscapeIdentifier(ident string) string {
	if ident == "" {
		return ident
	}
	// SQL Server uses square brackets for quoting identifiers.
	// Replace any closing brackets inside the identifier with double closing brackets.
	return fmt.Sprintf("[%s]", strings.ReplaceAll(ident, "]", "]]"))
}

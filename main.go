package main

import (
	"fmt"
	"vitess.io/vitess/go/vt/log"
	"vitess.io/vitess/go/vt/sqlparser"
)

// https://github.com/vitessio/vitess/blob/release-9.0/go/cmd/query_analyzer/query_analyzer.go
var bindIndex = 0

func main() {
	sql := "SELECT * FROM test"
	ast, err := sqlparser.Parse(sql)
	if err != nil {
		log.Errorf("Error parsing %s", sql)
		return
	}
	buf := sqlparser.NewTrackedBuffer(formatWithBind)
	buf.Myprintf("%v", ast)
	log.Info(buf.ParsedQuery().Query)
	log.Info("AST %v", ast)
}

func formatWithBind(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
	v, ok := node.(*sqlparser.Literal)
	if !ok {
		node.Format(buf)
		return
	}
	switch v.Type {
	case sqlparser.StrVal, sqlparser.HexVal, sqlparser.IntVal:
		buf.WriteArg(fmt.Sprintf(":v%d", bindIndex))
		bindIndex++
	default:
		node.Format(buf)
	}
}
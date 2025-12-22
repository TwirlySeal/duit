package main

import (
	"strings"
	"fmt"
)

type SQLMap struct {
	args []any
	clauses []string
	separator string
}

func NewSQLMap(separator string, args ...any) *SQLMap {
	return &SQLMap {
		args,
		make([]string, 0),
		separator,
	}
}

func (s *SQLMap) Param(name string, value any) {
	s.clauses = append(s.clauses, fmt.Sprintf("%s = $%d", name, len(s.args) + 1))
	s.args = append(s.args, value)
}

func (s *SQLMap) String() string {
	return strings.Join(s.clauses, s.separator)
}

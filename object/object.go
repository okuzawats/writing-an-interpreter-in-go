package object

import (
	"bytes"
	"fmt"
	"strings"

	"local.packages/ast"
)

// Objectの種別
type ObjectType string

const (
	// 正数型のオブジェクトを表す文字列
	INTEGER_OBJ = "INTEGER"
	// 真偽値型のオブジェクトを表す文字列
	BOOLEAN_OBJ = "BOOLEAN"
	// 関数オブジェクトを表す文字列
	FUNCTION_OBJ = "FUNCTION"
	// null型のオブジェクトを表す文字列
	NULL_OBJ = "NULL"
	// returnで返すオブジェクトを表す文字列
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	// 構文エラーオブジェクトを表す文字列
	ERROR_OBJ = "ERROR"
)

// Objectを表すinterface
type Object interface {
	// Objectの種別
	Type() ObjectType
	// Objectの文字列表現
	Inspect() string
}

// 正数型のObject
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// 真偽値型のObject
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

// null型のObject
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

// returnで返すべきObjectをラップするObject
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJECT
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// 構文エラー
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

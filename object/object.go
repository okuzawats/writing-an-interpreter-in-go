package object

import (
	"fmt"
)

type ObjectType string

const (
	// 正数型のオブジェクトを表す文字列
	INTEGER_OBJ = "INTEGER"
	// 真偽値型のオブジェクトを表す文字列
	BOOLEAN_OBJ = "BOOLEAN"
	// null型のオブジェクトを表す文字列
	NULL_OBJ = "NULL"
)

// Objectを表すinterface
type Object interface {
	// Objectの型
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

type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

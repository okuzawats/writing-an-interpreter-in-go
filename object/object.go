package object

import (
	"fmt"
)

// Objectの種別
type ObjectType string

const (
	// 正数型のオブジェクトを表す文字列
	INTEGER_OBJ = "INTEGER"
	// 真偽値型のオブジェクトを表す文字列
	BOOLEAN_OBJ = "BOOLEAN"
	// null型のオブジェクトを表す文字列
	NULL_OBJ = "NULL"
	// returnで返すオブジェクトを表す文字列
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
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

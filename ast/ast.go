package ast

import (
	"bytes"
	"fmt"
	"github.com/ZephroC/go-interpreter/token"
)

// Node - node in the Abstract Syntax tree
type Node interface {
	TokenLiteral() string
	String() string
}

//
type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token token.Token // will be token.Let
	Name  Identifier
	Value Expression
}

func (ls LetStatement) statementNode() {}
func (ls LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls LetStatement) String() string {
	if ls.Value != nil {
		return fmt.Sprintf("%s %s = %s;", ls.TokenLiteral(), ls.Name.String(), ls.Value.String())
	} else {
		return fmt.Sprintf("%s %s = ;", ls.TokenLiteral(), ls.Name.String())
	}
}

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i Identifier) expressionNode() {}
func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i Identifier) String() string {
	return i.Value
}

type ReturnStatement struct {
	Token token.Token // will be token.Return
	Value Expression
}

func (rs ReturnStatement) statementNode() {}
func (rs ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs ReturnStatement) String() string {
	if rs.Value != nil {
		return fmt.Sprintf("%s %s;", rs.TokenLiteral(), rs.Value.String())
	} else {
		return fmt.Sprintf("%s ;", rs.TokenLiteral())
	}
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es ExpressionStatement) statementNode() {}
func (es ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	} else {
		return ""
	}
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il IntegerLiteral) expressionNode()      {}
func (il IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe PrefixExpression) expressionNode()      {}
func (pe PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie InfixExpression) expressionNode()      {}
func (ie InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator, ie.Right.String())
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b BooleanLiteral) expressionNode()      {}
func (b BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b BooleanLiteral) String() string       { return b.Token.Literal }

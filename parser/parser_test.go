package parser

import (
	"fmt"
	"github.com/ZephroC/go-interpreter/ast"
	"github.com/ZephroC/go-interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program.Statements == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Should have 3 statements")
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, identifier string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("statment literal should be 'let' but was: %s", stmt.TokenLiteral())
		return false
	}
	letStmt, ok := stmt.(ast.LetStatement)
	if !ok {
		t.Errorf("type of let statement was incorrect, got=%T", stmt)
		return false
	}
	if letStmt.Name.Value != identifier {
		t.Errorf("Wanted identifier: %s but got: %s", identifier, letStmt.Name.Value)
	}
	if letStmt.Name.TokenLiteral() != identifier {
		t.Errorf("Wanted token literal: %s but got: %s", identifier, letStmt.Name.Value)
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) == 0 {
		return
	}
	t.Errorf("parser had %d errors", len(p.Errors()))
	for _, err := range p.Errors() {
		t.Error(err)
	}
	t.FailNow()
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program.Statements == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Should have 3 statements")
	}
	for _, stmt := range program.Statements {
		if stmt.TokenLiteral() != "return" {
			t.Errorf("token literal should be return but was: %s", stmt.TokenLiteral())
		}
		_, ok := stmt.(ast.ReturnStatement)
		if !ok {
			t.Errorf("should have type of ast.ReturnStatement but was: %T", stmt)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected 1 statement but got: %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected a type of ast.ExpressionStatement, got: %T", program.Statements[0])
	}
	ident, ok := stmt.Expression.(ast.Identifier)
	if !ok {
		t.Fatalf("Expected a type of ast.Identifier, got: %T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("Ident value should be foobar, was: %s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("Ident TokenLiteral should be foobar, was: %s", ident.Value)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program should have 1 statement, got: %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Was not of type ast.ExpressionStatement, got: %T", program.Statements[0])
	}
	literal, ok := stmt.Expression.(ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Not an ast.IntegerLiteral got: %T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("Expected an integer literal of 5, got; %d", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("Expected the tokenLiteral to be \"5\" but got: %s", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"-bob;", "-", "bob"},
		{"!alice;", "!", "alice"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected 1 statement, got: %d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected ast.ExpressionStatement got: %T", program.Statements)
		}
		exp, ok := stmt.Expression.(ast.PrefixExpression)
		if !ok {
			t.Fatalf("expected ast.PrefixExpression got: %T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("Expected operator: %s but got: %s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(ast.IntegerLiteral)
	if !ok {
		t.Errorf("Wanted an integer literal but got: %T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("Wanted %d but got: %d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("wanted token literal: %d but got: %s", value, integ.TokenLiteral())
		return false
	}
	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"bob * alice", "bob", "*", "alice"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program sould have 1 statement but had: %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Wrong type wanted: ast.ExpressionStatement got: %T", program.Statements[0])
		}
		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"-a + b / c", "((-a) + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 5", "((1 + (2 + 3)) + 5)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5+5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected: %s \n got: %s", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) (match bool) {
	match = true
	ident, ok := exp.(ast.Identifier)
	if !ok {
		t.Errorf("exp not ast.Identifier, got: %T", exp)
		return false
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s, got: %s", value, ident.Value)
		match = false
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s, got: %s", value, ident.TokenLiteral())
		match = false
	}
	return
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("Type of exp not known, got: %T", expected)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, v bool) (match bool) {
	match = true
	bo, ok := exp.(ast.BooleanLiteral)
	if !ok {
		t.Errorf("expected ast.BooleanLiteral, got: %T", exp)
		return false
	}
	if bo.Value != v {
		t.Errorf("value of literal not %t, got: %t", v, bo.Value)
	}
	if bo.TokenLiteral() != fmt.Sprintf("%t", v) {
		t.Errorf("value of TokenLiteral not %t, got: %s", v, bo.TokenLiteral())
	}
	return
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) (match bool) {
	match = true
	opExp, ok := exp.(ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got: %T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		match = false
	}
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not: '%s', got: %s", operator, opExp.Operator)
		match = false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		match = false
	}
	return
}

package parser

import (
	"fmt"
	"testing"

	"github.com/lautarojayat/writing-an-interpreter-in-go/ast"
	"github.com/lautarojayat/writing-an-interpreter-in-go/lexer"
)

func TestLetStatement(t *testing.T) {
	input := `	let x = 5;
				let y = 10;
				let foobar = 838383;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram() returnes nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 elements, instead we got %d", len(program.Statements))
	}
	tests := []struct{ expectedIdentifier string }{
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

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, e := range errors {
		t.Errorf("parser error: %q", e)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral is not 'let' instead got %q", s.TokenLiteral())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not LetStatement, instead got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s', instead got '%s", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letstmt.Name is not '%s', instead got '%s'", name, letStmt.Name)
		return false
	}
	return true
}

func TestReturnStatement(t *testing.T) {
	input := `	return 10;
				return 5;
				return 202020;`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program hasn't got 3 statements, instead got %d", len(program.Statements))
	}
	for _, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not of type *ast.ReturnStatement, intead got %T", stmt)
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStatement token literal is not 'return', instead got %q", returnStatement.TokenLiteral())
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
		t.Fatalf("program hasn't got 1 statement, instead got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not expression, instead got %T", stmt)
	}
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("ident is not *ast.Identity, instead got %T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not 'foobar', instead got '%q'", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() is not 'foobar' instead got '%q'", ident.TokenLiteral())
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program hasn't got 1 statement, instead got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not expression, instead got %T", stmt)
	}
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("literal is not *ast.IntegerLiteral, instead got %T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value is not '5', instead got '%q'", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral() is not '5' instead got '%q'", literal.TokenLiteral())
	}

}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program hasn't got 1 statement, instead got %d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not expression, instead got %T", stmt)
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("literal is not *ast.PrefixExpression, instead got %T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not %q, instead got '%q'", exp.Operator, tt.operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}

}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not ast.IntegerLiteral, instead got %T", il)
		return false
	}
	if integer.Value != value {
		t.Errorf("integer.Value is no %d, instead got %d", integer.Value, value)
		return false
	}
	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() is not %d, instead got %s", value, integer.TokenLiteral())
		return false
	}
	return true
}

func TestInfixExpression(t *testing.T) {
	prefixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
	}

	for _, tt := range prefixTests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		if len(program.Statements) != 1 {
			t.Fatalf("program hasn't got 1 statement, instead got %d", len(program.Statements))
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not expression, instead got %T", stmt)
		}
		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("literal is not *ast.InfixExpression, instead got %T", stmt.Expression)
		}
		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator is not %q, instead got '%q'", exp.Operator, tt.operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}

}

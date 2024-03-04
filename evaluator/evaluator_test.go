package evaluator

import (
	"testing"

	"github.com/aryuuu/gonkey-lang/lexer"
	"github.com/aryuuu/gonkey-lang/object"
	"github.com/aryuuu/gonkey-lang/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{
			input:    "5",
			expected: 5,
		},
		{
			input:    "10",
			expected: 10,
		},
		{
			input:    "-5",
			expected: -5,
		},
		{
			input:    "-10",
			expected: -10,
		},
		{
			"5 + 5 + 5 + 5 - 10",
			10,
		},
		{
			"2 * 2 * 2 * 2 * 2",
			32,
		},
		{
			"-50 + 100 + -50",
			0,
		},
		{
			"5 * 2 + 10",
			20,
		},
		{
			"5 + 2 * 10",
			25,
		},
		{
			"20 + 2 * -10",
			0,
		},
		{
			"50 / 2 * 2 + 10",
			60,
		},
		{
			"2 * (5 + 10)",
			30,
		},
		{
			"3 * 3 * 3 + 10",
			37,
		},
		{
			"3 * (3 * 3) + 10",
			37,
		},
		{
			"(5 + 10 * 2 + 15 / 3) * 2 + -10",
			50,
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{
			input:    "if(true) { 10 }",
			expected: 10,
		},
		{
			input:    "if(false) { 10 }",
			expected: nil,
		},
		{
			input:    "if(1) { 10 }",
			expected: 10,
		},
		{
			input:    "if(1 < 2) { 10 }",
			expected: 10,
		},
		{
			input:    "if(1 > 2) { 10 }",
			expected: nil,
		},
		{
			input:    "if(1 > 2) { 10 } else { 20 }",
			expected: 20,
		},
		{
			input:    "if(1 < 2) { 10 } else { 20 }",
			expected: 10,
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		integer, ok := tc.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "true",
			expected: true,
		},
		{
			input:    "false",
			expected: false,
		},
		{
			input:    "1 < 2",
			expected: true,
		},
		{
			input:    "1 > 2",
			expected: false,
		},
		{
			input:    "1 < 1",
			expected: false,
		},
		{
			input:    "1 > 1",
			expected: false,
		},
		{
			input:    "1 == 1",
			expected: true,
		},
		{
			input:    "1 != 1",
			expected: false,
		},
		{
			input:    "1 == 2",
			expected: false,
		},
		{
			input:    "1 != 2",
			expected: true,
		},
		{
			input:    "true == true",
			expected: true,
		},
		{
			input:    "false == false",
			expected: true,
		},
		{
			input:    "true == false",
			expected: false,
		},
		{
			input:    "true != false",
			expected: true,
		},
		{
			input:    "false != true",
			expected: true,
		},
		{
			input:    "(1 < 2) == true",
			expected: true,
		},
		{
			input:    "(1 < 2) == false",
			expected: false,
		},
		{
			input:    "(1 > 2) == true",
			expected: false,
		},
		{
			input:    "(1 > 2) == false",
			expected: true,
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

func TestBangOperator(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "!true",
			expected: false,
		},
		{
			input:    "!false",
			expected: true,
		},
		{
			input:    "!5",
			expected: false,
		},
		{
			input:    "!!true",
			expected: true,
		},
		{
			input:    "!!false",
			expected: false,
		},
		{
			input:    "!!5",
			expected: true,
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		testBooleanObject(t, evaluated, tc.expected)
	}
}

// func TestReturnValue(t *testing.T) {
// 	testCases := []struct {
// 		input    string
// 		expected int64
// 	}{
// 		{
// 			input: "return 10;",
// 			expected: 10,
// 		},
// 		{
// 			input: "return 10; 9;",
// 			expected: 10,
// 		},
// 		{
// 			input: "return 2 * 5;",
// 			expected: 10,
// 		},
// 		{
// 			input: "9; return 2 * 5; 9;",
// 			expected: 10,
// 		},
// 		{
// 			input: `
// 			if (10 > 1) {
// 				return 10;
// 			}

// 			return 1;
// 			`,
// 			expected: 10,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		evaluated := testEval(tc.input)
// 		testIntegerObject(t, evaluated, tc.expected)
// 	}
// }

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, val int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("obj should be of type object.Integer, got=%T (%+v)\n", obj, obj)
		return false
	}

	if result.Value != val {
		t.Errorf("obj value should be %d, got=%d\n", val, result.Value)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Null)
	if !ok {
		t.Errorf("obj should be of type object.Null, got=%T (%+v)\n", obj, obj)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, val bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("obj should be of type object.Boolean, got=%T (%+v)\n", obj, obj)
		return false
	}

	if result.Value != val {
		t.Errorf("obj value should be %t, got=%t\n", val, result.Value)
		return false
	}

	return true
}
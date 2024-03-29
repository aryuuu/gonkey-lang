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

func TestLetStatements(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{
			input:    "let a = 5; a;",
			expected: 5,
		},
		{
			input:    "let a = 5 * 5; a;",
			expected: 25,
		},
		{
			input:    "let a = 5; let b = a; b;",
			expected: 5,
		},
		{
			input:    "let a = 5; let b = a; let c = a + b + 5; c;",
			expected: 15,
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) { x + 2; };`

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function, got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong params, Params=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x', got=%q", fn.Parameters[0])
	}

	expectedBody := `(x + 2)`
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q, got=%T", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{
			input:    "let identity = fn(x) { x; }; identity(5);",
			expected: 5,
		},
		{
			input:    "let identity = fn(x) { return x; }; identity(5);",
			expected: 5,
		},
		{
			input:    "let double = fn(x) { 2 * x; }; double(5);",
			expected: 10,
		},
		{
			input:    "let add = fn(x, y) { y + x; }; add(5, 5);",
			expected: 10,
		},
		{
			input:    "let add = fn(x, y) { y + x; }; add(5 + 5, add(5, 5));",
			expected: 20,
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "hello world" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3];`
	evaluated := testEval(input)

	arr, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array, got=%T (%+v)", evaluated, evaluated)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("len(arr.Elements) is not 3, got=%d", len(arr.Elements))
	}

	testIntegerObject(t, arr.Elements[0], 1)
	testIntegerObject(t, arr.Elements[1], 4)
	testIntegerObject(t, arr.Elements[2], 6)
}

func TestMapLiteral(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr"+"ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}
	`
	evaluated := testEval(input)

	mapLit, ok := evaluated.(*object.Map)
	if !ok {
		t.Fatalf("object is not Map, got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(mapLit.Pairs) != len(expected) {
		t.Fatalf("len(arr.Elements) is not %d, got=%d", len(expected), len(mapLit.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := mapLit.Pairs[expectedKey]
		if !ok {
			t.Error("no pair for given key")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}

}

func TestMapIndexExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{
			input:    `{"foo": 5}["foo"]`,
			expected: 5,
		},
		{
			input:    `{"foo": 5}["bar"]`,
			expected: nil,
		},
		{
			input:    `let key = "foo"; {"foo": 5}[key]`,
			expected: 5,
		},
		{
			input:    `{}["foo"]`,
			expected: nil,
		},
		{
			input:    `{5: 5}[5]`,
			expected: 5,
		},
		{
			input:    `{true: 5}[true]`,
			expected: 5,
		},
		{
			input:    `{false: 5}[false]`,
			expected: 5,
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

func TestIndexExpressions(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [0][i];",
			0,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
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

func TestStringConcatenation(t *testing.T) {
	input := `"hello world" + " goodbye world"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "hello world goodbye world" {
		t.Errorf("String has wrong value. got=%q", str.Value)
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

func TestReturnValue(t *testing.T) {
	testCases := []struct {
		input    string
		expected int64
	}{
		{
			input:    "return 10;",
			expected: 10,
		},
		{
			input:    "return 10; 9;",
			expected: 10,
		},
		{
			input:    "return 2 * 5;",
			expected: 10,
		},
		{
			input:    "9; return 2 * 5; 9;",
			expected: 10,
		},
		{
			input: `
			if (10 > 1) {
				return 10;
			}

			return 1;
			`,
			expected: 10,
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		testIntegerObject(t, evaluated, tc.expected)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	testCases := []struct {
		input    string
		expected any
	}{
		{
			input:    `len("")`,
			expected: 0,
		},
		{
			input:    `len("four")`,
			expected: 4,
		},
		{
			input:    `len("hello world")`,
			expected: 11,
		},
		{
			input:    `len(1)`,
			expected: "argument to `len` not supported. got INTEGER",
		},
		{
			input:    `len(["fatt"])`,
			expected: 1,
		},
		{
			input:    `len(["fatt", "lisa"])`,
			expected: 2,
		},
		{
			input:    `len(["fatt", "lisa", 1])`,
			expected: 3,
		},
		{
			input:    `len("one", "two")`,
			expected: "wrong number of arguments. got=2, want=1",
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)

		switch expected := tc.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not error, got=%T", evaluated)
				continue
			}

			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestErrorHandling(t *testing.T) {
	testCases := []struct {
		input           string
		expectedMessage string
	}{
		{
			input:           "5 + true",
			expectedMessage: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:           "5 + true; 5;",
			expectedMessage: "type mismatch: INTEGER + BOOLEAN",
		},
		{
			input:           "-true",
			expectedMessage: "unknown operator: -BOOLEAN",
		},
		{
			input:           "true + false",
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:           "5; true + false; 5",
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:           "if (10 > 1) { true + false; }",
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input:           `"Hello" - "World"`,
			expectedMessage: "unknown operator: STRING - STRING",
		},
		{
			input: `
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
			}

			return 1;
			`,
			expectedMessage: "unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			input: `{"name": "Monkey"}[fn(x) { x }]`,
			expectedMessage: "unusable as hash key: FUNCTION",
		},
	}

	for _, tc := range testCases {
		evaluated := testEval(tc.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned, got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tc.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tc.expectedMessage, errObj.Message)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	env := object.NewEnvironment()
	return Eval(program, env)
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

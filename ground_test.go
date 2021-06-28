package bass_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mattn/go-colorable"
	"github.com/stretchr/testify/require"
	"github.com/vito/bass"
)

var docsOut = new(bytes.Buffer)

func init() {
	bass.DocsWriter = colorable.NewNonColorable(docsOut)
}

var operative = &bass.Operative{
	Formals: bass.NewList(bass.Symbol("form")),
	Eformal: bass.Symbol("env"),
	Body: bass.InertPair{
		A: bass.Symbol("form"),
		D: bass.Symbol("env"),
	},
}

var pair = Const{
	bass.Pair{
		A: bass.Int(1),
		D: bass.Empty{},
	},
}

var nonListPair = Const{
	bass.Pair{
		A: bass.Int(1),
		D: bass.Int(2),
	},
}

var inertPair = bass.InertPair{
	A: bass.Int(1),
	D: bass.Empty{},
}

var env = bass.NewEnv()

type Const struct {
	bass.Value
}

func (value Const) Eval(*bass.Env) (bass.Value, error) {
	return value.Value, nil
}

var sym = Const{
	Value: bass.Symbol("sym"),
}

func TestGroundPrimitivePredicates(t *testing.T) {
	env := bass.New()

	type example struct {
		Name   string
		Trues  []bass.Value
		Falses []bass.Value
	}

	for _, test := range []example{
		{
			Name: "null?",
			Trues: []bass.Value{
				bass.Null{},
			},
			Falses: []bass.Value{
				bass.Bool(false),
				pair,
				inertPair,
				bass.Empty{},
				bass.Ignore{},
				bass.Int(0),
				bass.String(""),
			},
		},
		{
			Name: "boolean?",
			Trues: []bass.Value{
				bass.Bool(true),
				bass.Bool(false),
			},
			Falses: []bass.Value{
				bass.Int(1),
				bass.String("true"),
				bass.Null{},
			},
		},
		{
			Name: "number?",
			Trues: []bass.Value{
				bass.Int(0),
			},
			Falses: []bass.Value{
				bass.Bool(true),
				bass.String("1"),
			},
		},
		{
			Name: "string?",
			Trues: []bass.Value{
				bass.String("str"),
			},
			Falses: []bass.Value{
				Const{bass.Symbol("1")},
				bass.Empty{},
				bass.Ignore{},
			},
		},
		{
			Name: "symbol?",
			Trues: []bass.Value{
				sym,
			},
			Falses: []bass.Value{
				bass.String("str"),
			},
		},
		{
			Name: "empty?",
			Trues: []bass.Value{
				bass.Null{},
				bass.Empty{},
				bass.String(""),
			},
			Falses: []bass.Value{
				bass.Bool(false),
				bass.Ignore{},
			},
		},
		{
			Name: "pair?",
			Trues: []bass.Value{
				pair,
				inertPair,
			},
			Falses: []bass.Value{
				bass.Empty{},
				bass.Ignore{},
				bass.Null{},
			},
		},
		{
			Name: "list?",
			Trues: []bass.Value{
				bass.Empty{},
				pair,
				inertPair,
			},
			Falses: []bass.Value{
				nonListPair,
				bass.Ignore{},
				bass.Null{},
				bass.String(""),
			},
		},
		{
			Name: "env?",
			Trues: []bass.Value{
				env,
			},
			Falses: []bass.Value{
				pair,
			},
		},
		{
			Name: "combiner?",
			Trues: []bass.Value{
				bass.Op("quote", func(args bass.List, env *bass.Env) bass.Value {
					return args.First()
				}),
			},
		},
		{
			Name: "applicative?",
			Trues: []bass.Value{
				bass.Func("id", func(val bass.Value) bass.Value {
					return val
				}),
			},
			Falses: []bass.Value{
				bass.Op("quote", func(args bass.List, env *bass.Env) bass.Value {
					return args.First()
				}),
			},
		},
		{
			Name: "operative?",
			Trues: []bass.Value{
				bass.Op("quote", func(args bass.List, env *bass.Env) bass.Value {
					return args.First()
				}),
				operative,
			},
			Falses: []bass.Value{
				bass.Func("id", func(val bass.Value) bass.Value {
					return val
				}),
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			for _, arg := range test.Trues {
				t.Run(fmt.Sprintf("%v", arg), func(t *testing.T) {
					res, err := bass.Pair{
						A: bass.Symbol(test.Name),
						D: bass.NewList(arg),
					}.Eval(env)
					require.NoError(t, err)
					require.Equal(t, bass.Bool(true), res)
				})
			}

			for _, arg := range test.Falses {
				t.Run(fmt.Sprintf("%v", arg), func(t *testing.T) {
					res, err := bass.Pair{
						A: bass.Symbol(test.Name),
						D: bass.NewList(arg),
					}.Eval(env)
					require.NoError(t, err)
					require.Equal(t, bass.Bool(false), res)
				})
			}
		})
	}
}

func TestGroundNumeric(t *testing.T) {
	env := bass.New()

	type example struct {
		Name   string
		Bass   string
		Result bass.Value
	}

	for _, test := range []example{
		{
			Name:   "+",
			Bass:   "(+ 1 2 3)",
			Result: bass.Int(6),
		},
		{
			Name:   "-",
			Bass:   "(- 1 2 3)",
			Result: bass.Int(-4),
		},
		{
			Name:   "- unary",
			Bass:   "(- 1)",
			Result: bass.Int(-1),
		},
		{
			Name:   "* no args",
			Bass:   "(*)",
			Result: bass.Int(1),
		},
		{
			Name:   "* unary",
			Bass:   "(* 5)",
			Result: bass.Int(5),
		},
		{
			Name:   "* product",
			Bass:   "(* 1 2 3 4)",
			Result: bass.Int(24),
		},
		{
			Name:   "max",
			Bass:   "(max 1 3 7 5 4)",
			Result: bass.Int(7),
		},
		{
			Name:   "min",
			Bass:   "(min 5 3 7 2 4)",
			Result: bass.Int(2),
		},
		{
			Name:   "min",
			Bass:   "(min 5 3 7 2 4)",
			Result: bass.Int(2),
		},
		{
			Name:   "=? same",
			Bass:   "(=? 1 1 1)",
			Result: bass.Bool(true),
		},
		{
			Name:   "=? different",
			Bass:   "(=? 1 2 1)",
			Result: bass.Bool(false),
		},
		{
			Name:   ">? decreasing",
			Bass:   "(>? 3 2 1)",
			Result: bass.Bool(true),
		},
		{
			Name:   ">? decreasing-eq",
			Bass:   "(>? 3 2 2)",
			Result: bass.Bool(false),
		},
		{
			Name:   ">? increasing",
			Bass:   "(>? 1 2 3)",
			Result: bass.Bool(false),
		},
		{
			Name:   ">? increasing-eq",
			Bass:   "(>? 1 2 2)",
			Result: bass.Bool(false),
		},
		{
			Name:   ">=? decreasing",
			Bass:   "(>=? 3 2 1)",
			Result: bass.Bool(true),
		},
		{
			Name:   ">=? decreasing-eq",
			Bass:   "(>=? 3 2 2)",
			Result: bass.Bool(true),
		},
		{
			Name:   ">=? increasing",
			Bass:   "(>=? 1 2 3)",
			Result: bass.Bool(false),
		},
		{
			Name:   ">=? increasing-eq",
			Bass:   "(>=? 1 2 2)",
			Result: bass.Bool(false),
		},
		{
			Name:   "<? decreasing",
			Bass:   "(<? 3 2 1)",
			Result: bass.Bool(false),
		},
		{
			Name:   "<? decreasing-eq",
			Bass:   "(<? 3 2 2)",
			Result: bass.Bool(false),
		},
		{
			Name:   "<? increasing",
			Bass:   "(<? 1 2 3)",
			Result: bass.Bool(true),
		},
		{
			Name:   "<? increasing-eq",
			Bass:   "(<? 1 2 2)",
			Result: bass.Bool(false),
		},
		{
			Name:   "<=? decreasing",
			Bass:   "(<=? 3 2 1)",
			Result: bass.Bool(false),
		},
		{
			Name:   "<=? decreasing-eq",
			Bass:   "(<=? 3 2 2)",
			Result: bass.Bool(false),
		},
		{
			Name:   "<=? increasing",
			Bass:   "(<=? 1 2 3)",
			Result: bass.Bool(true),
		},
		{
			Name:   "<=? increasing-eq",
			Bass:   "(<=? 1 2 2)",
			Result: bass.Bool(true),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			reader := bass.NewReader(bytes.NewBufferString(test.Bass))

			val, err := reader.Next()
			require.NoError(t, err)

			res, err := val.Eval(env)
			require.NoError(t, err)

			require.Equal(t, test.Result, res)
		})
	}
}

func TestGroundConstructors(t *testing.T) {
	env := bass.New()

	env.Set("operative", operative)

	type example struct {
		Name string
		Bass string

		Result      bass.Value
		Err         error
		ErrContains string
	}

	for _, test := range []example{
		{
			Name: "cons",
			Bass: "(cons 1 2)",
			Result: bass.Pair{
				A: bass.Int(1),
				D: bass.Int(2),
			},
		},
		{
			Name: "op",
			Bass: "(op (x) e [x e])",
			Result: &bass.Operative{
				Formals: bass.NewList(bass.Symbol("x")),
				Eformal: bass.Symbol("e"),
				Body:    bass.NewInertList(bass.Symbol("x"), bass.Symbol("e")),
				Env:     env,
			},
		},
		{
			Name: "bracket op",
			Bass: "(op [x] e [x e])",
			Result: &bass.Operative{
				Formals: bass.NewInertList(bass.Symbol("x")),
				Eformal: bass.Symbol("e"),
				Body:    bass.NewInertList(bass.Symbol("x"), bass.Symbol("e")),
				Env:     env,
			},
		},
		{
			Name:   "wrap",
			Bass:   "((wrap (op x _ x)) 1 2 (+ 1 2))",
			Result: bass.NewList(bass.Int(1), bass.Int(2), bass.Int(3)),
		},
		{
			Name:   "unwrap",
			Bass:   "(unwrap (wrap operative))",
			Result: operative,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			reader := bass.NewReader(bytes.NewBufferString(test.Bass))

			val, err := reader.Next()
			require.NoError(t, err)

			res, err := val.Eval(env)
			if test.Err != nil {
				require.Equal(t, test.Err, err)
			} else if test.ErrContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.ErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Result, res)
			}
		})
	}
}

func TestGroundEnv(t *testing.T) {
	type example struct {
		Name string
		Bass string

		Result   bass.Value
		Bindings bass.Bindings

		Err         error
		ErrContains string
	}

	// used as a test value
	sentinel := bass.String("evaluated")

	for _, test := range []example{
		{
			Name:   "eval",
			Bass:   "((op [x] e (eval x e)) sentinel)",
			Result: bass.String("evaluated"),
		},
		{
			Name:   "make-env",
			Bass:   "(make-env)",
			Result: bass.NewEnv(),
		},
		{
			Name:   "make-env",
			Bass:   "(make-env (make-env) (make-env))",
			Result: bass.NewEnv(bass.NewEnv(), bass.NewEnv()),
		},
		{
			Name:   "def",
			Bass:   "(def foo 1)",
			Result: bass.Symbol("foo"),
			Bindings: bass.Bindings{
				"foo":      bass.Int(1),
				"sentinel": sentinel,
			},
		},
		{
			Name:   "def evaluation",
			Bass:   "(def foo sentinel)",
			Result: bass.Symbol("foo"),
			Bindings: bass.Bindings{
				"foo":      sentinel,
				"sentinel": sentinel,
			},
		},
		{
			Name: "def destructuring",
			Bass: "(def (a . bs) [1 2 3])",
			Result: bass.Pair{
				A: bass.Symbol("a"),
				D: bass.Symbol("bs"),
			},
			Bindings: bass.Bindings{
				"a":        bass.Int(1),
				"bs":       bass.NewList(bass.Int(2), bass.Int(3)),
				"sentinel": sentinel,
			},
		},
		{
			Name: "def destructuring advanced",
			Bass: "(def (a b [c d] e . fs) [1 2 [3 4] 5 6 7])",
			Result: bass.Pair{
				A: bass.Symbol("a"),
				D: bass.Pair{
					A: bass.Symbol("b"),
					D: bass.Pair{
						A: bass.InertPair{
							A: bass.Symbol("c"),
							D: bass.InertPair{
								A: bass.Symbol("d"),
								D: bass.Empty{},
							},
						},
						D: bass.Pair{
							A: bass.Symbol("e"),
							D: bass.Symbol("fs"),
						},
					},
				},
			},
			Bindings: bass.Bindings{
				"a":        bass.Int(1),
				"b":        bass.Int(2),
				"c":        bass.Int(3),
				"d":        bass.Int(4),
				"e":        bass.Int(5),
				"fs":       bass.NewList(bass.Int(6), bass.Int(7)),
				"sentinel": sentinel,
			},
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			reader := bytes.NewBufferString(test.Bass)

			env := bass.New()
			env.Set("sentinel", sentinel)

			res, err := bass.EvalReader(env, reader)
			if test.Err != nil {
				require.Equal(t, test.Err, err)
			} else if test.ErrContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.ErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Result, res)

				if test.Bindings != nil {
					require.Equal(t, test.Bindings, env.Bindings)
				}
			}
		})
	}
}

func TestGroundEnvDoc(t *testing.T) {
	reader := bytes.NewBufferString(`
; commentary for environment
; split along multiple lines
_

; a separate comment
;
; with multiple paragraphs
_

; docs for abc
(def abc 123)

; more commentary between abc and quote
_

(defop quote (x) _ x) ; docs for quote

; docs for inc
(defn inc (x) (+ x 1))

(doc abc quote inc)
`)

	env := bass.New()

	_, err := bass.EvalReader(env, reader)
	require.NoError(t, err)

	require.Contains(t, docsOut.String(), "docs for abc")
	require.Contains(t, docsOut.String(), "number?")
	require.Contains(t, docsOut.String(), "docs for quote")
	require.Contains(t, docsOut.String(), "operative?")
	require.Contains(t, docsOut.String(), "docs for inc")
	require.Contains(t, docsOut.String(), "applicative?")

	docsOut.Reset()

	reader = bytes.NewBufferString(`(doc)`)
	_, err = bass.EvalReader(env, reader)
	require.NoError(t, err)

	require.Contains(t, docsOut.String(), `--------------------------------------------------
commentary for environment split along multiple lines
`)

	require.Contains(t, docsOut.String(), `--------------------------------------------------
abc number?

docs for abc
`)

	require.Contains(t, docsOut.String(), `--------------------------------------------------
a separate comment

with multiple paragraphs
`)

	require.Contains(t, docsOut.String(), `--------------------------------------------------
quote operative? combiner?
args: (x)

docs for quote
`)

	require.Contains(t, docsOut.String(), `--------------------------------------------------
inc applicative? combiner?
args: (x)

docs for inc

`)
}

func TestGroundBoolean(t *testing.T) {
	type example struct {
		Name string
		Bass string

		Result   bass.Value
		Bindings bass.Bindings

		Err         error
		ErrContains string
	}

	// used as a test value
	sentinel := bass.String("evaluated")

	for _, test := range []example{
		{
			Name:   "if true",
			Bass:   "(if true sentinel unevaluated)",
			Result: sentinel,
		},
		{
			Name:   "if false",
			Bass:   "(if false unevaluated sentinel)",
			Result: sentinel,
		},
		{
			Name:   "if null",
			Bass:   "(if null unevaluated sentinel)",
			Result: sentinel,
		},
		{
			Name:   "if empty",
			Bass:   "(if [] sentinel unevaluated)",
			Result: sentinel,
		},
		{
			Name:   "if string",
			Bass:   `(if "" sentinel unevaluated)`,
			Result: sentinel,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			reader := bass.NewReader(bytes.NewBufferString(test.Bass))

			val, err := reader.Next()
			require.NoError(t, err)

			env := bass.New()
			env.Set("sentinel", sentinel)

			res, err := val.Eval(env)
			if test.Err != nil {
				require.Equal(t, test.Err, err)
			} else if test.ErrContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.ErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Result, res)

				if test.Bindings != nil {
					require.Equal(t, test.Bindings, env.Bindings)
				}
			}
		})
	}
}

func TestGroundStdlib(t *testing.T) {
	type example struct {
		Name string
		Bass string

		Result   bass.Value
		Bindings bass.Bindings

		Err         error
		ErrContains string
	}

	for _, test := range []example{
		{
			Name:   "do",
			Bass:   "(do (def a 1) (def b 2) [a b])",
			Result: bass.NewList(bass.Int(1), bass.Int(2)),
			Bindings: bass.Bindings{
				"a": bass.Int(1),
				"b": bass.Int(2),
			},
		},
		{
			Name:   "list",
			Bass:   "(list (def a 42) a)",
			Result: bass.NewList(bass.Symbol("a"), bass.Int(42)),
			Bindings: bass.Bindings{
				"a": bass.Int(42),
			},
		},
		{
			Name: "list*",
			Bass: "(list* (def a 1) a (list (def b 2) b))",
			Result: bass.NewList(
				bass.Symbol("a"),
				bass.Int(1),
				bass.Symbol("b"),
				bass.Int(2),
			),
			Bindings: bass.Bindings{
				"a": bass.Int(1),
				"b": bass.Int(2),
			},
		},
		{
			Name:   "first",
			Bass:   "(first (list 1 2 3))",
			Result: bass.Int(1),
		},
		{
			Name:   "rest",
			Bass:   "(rest (list 1 2 3))",
			Result: bass.NewList(bass.Int(2), bass.Int(3)),
		},
		{
			Name:   "second",
			Bass:   "(second (list 1 2 3))",
			Result: bass.Int(2),
		},
		{
			Name:   "third",
			Bass:   "(third (list 1 2 3))",
			Result: bass.Int(3),
		},
		{
			Name:   "length",
			Bass:   "(length (list 1 2 3))",
			Result: bass.Int(3),
		},
		{
			Name:   "do op",
			Bass:   "((op [x y] e (eval [def x y] e) y) foo 42)",
			Result: bass.Int(42),
			Bindings: bass.Bindings{
				"foo": bass.Int(42),
			},
		},
		{
			Name: "invalid op 0",
			Bass: "(op)",
			Err: bass.BindMismatchError{
				Need: bass.Pair{
					A: bass.Symbol("formals"),
					D: bass.Pair{
						A: bass.Symbol("eformal"),
						D: bass.Symbol("body"),
					},
				},
				Have: bass.Empty{},
			},
		},
		{
			Name: "invalid op 1",
			Bass: "(op [x])",
			Err: bass.BindMismatchError{
				Need: bass.Pair{
					A: bass.Symbol("eformal"),
					D: bass.Symbol("body"),
				},
				Have: bass.Empty{},
			},
		},
		{
			Name: "invalid op 2",
			Bass: "(op [x] _)",
			Err: bass.BindMismatchError{
				Need: bass.Pair{
					A: bass.Symbol("f"),
					D: bass.Ignore{},
				},
				Have: bass.Empty{},
			},
		},
		{
			Name: "invalid op 3",
			Bass: "(op . false)",
			Err: bass.BindMismatchError{
				Need: bass.Pair{
					A: bass.Symbol("formals"),
					D: bass.Pair{
						A: bass.Symbol("eformal"),
						D: bass.Symbol("body"),
					},
				},
				Have: bass.Bool(false),
			},
		},
		{
			Name: "invalid op 4",
			Bass: "(op [x] . _)",
			Err: bass.BindMismatchError{
				Need: bass.Pair{
					A: bass.Symbol("eformal"),
					D: bass.Symbol("body"),
				},
				Have: bass.Ignore{},
			},
		},
		{
			Name:   "defop",
			Bass:   `(defop def2 [x y] e (eval [def x y] e) y)`,
			Result: bass.Symbol("def2"),
		},
		{
			Name:   "defop call",
			Bass:   `(defop def2 [x y] e (eval [def x y] e) y) (def2 foo 42)`,
			Result: bass.Int(42),
		},
		{
			Name:     "fn",
			Bass:     "((fn [x] (def local (* x 2)) [local (* local 2)]) 21)",
			Result:   bass.NewList(bass.Int(42), bass.Int(84)),
			Bindings: bass.Bindings{},
		},
		{
			Name:   "defn",
			Bass:   "(defn foo [x] (def local (* x 2)) [local (* local 2)])",
			Result: bass.Symbol("foo"),
		},
		{
			Name:   "defn call",
			Bass:   "(defn foo [x] (def local (* x 2)) [local (* local 2)]) (foo 21)",
			Result: bass.NewList(bass.Int(42), bass.Int(84)),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			env := bass.New()

			res, err := bass.EvalReader(env, bytes.NewBufferString(test.Bass))
			if test.Err != nil {
				require.Equal(t, test.Err, err)
			} else if test.ErrContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.ErrContains)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Result, res)

				if test.Bindings != nil {
					require.Equal(t, test.Bindings, env.Bindings)
				}
			}
		})
	}
}

package generator

import (
	"fmt"
	"go/types"
	"strings"
)

type TupleFormatter struct {
	rawNames   []string
	types      []string
	isVariadic bool
}

func NewTupleFormatter(tuple *types.Tuple, isVariadic bool, pkgQualifier func(pkg *types.Package) string) *TupleFormatter {
	tf := &TupleFormatter{}
	tf.isVariadic = isVariadic
	tf.rawNames = make([]string, tuple.Len())
	tf.types = make([]string, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		tf.rawNames[i] = tuple.At(i).Name()
		tf.types[i] = types.TypeString(tuple.At(i).Type(), pkgQualifier)
	}
	if isVariadic {
		variadic := &tf.types[tuple.Len()-1]
		*variadic = strings.Replace(*variadic, "[]", "...", 1)
	}
	return tf
}

type TupleFormats struct {
	Args                     string
	ArgsMatchers             string
	RawParams                string
	NamedParams              string
	NamedArgedParams         string
	VariadicArgs             string
	VariadicArgsEval         string
	VariadicArgsMatchersEval string
}

func (tf TupleFormatter) Format(prefixForUnnamed string, forbiddenNames map[string]struct{}) TupleFormats {
	tupleLen := len(tf.rawNames)
	f := TupleFormats{}

	names := make([]string, tupleLen)
	for i := 0; i < tupleLen; i++ {
		name := tf.rawNames[i]
		if _, isForbidden := forbiddenNames[name]; isForbidden || name == "" || name == "_" {
			name = fmt.Sprintf("%s%v", prefixForUnnamed, i)
		}
		names[i] = name
	}

	f.Args = strings.Join(names, ", ")
	if tf.isVariadic {
		f.Args += "..."
	}

	sb := strings.Builder{}
	for i := 0; i < tupleLen; i++ {
		name := names[i]
		if i != tupleLen-1 {
			sb.WriteString(name)
			sb.WriteString(".Matcher, ")
		} else {
			if !tf.isVariadic {
				sb.WriteString(name)
				sb.WriteString(".Matcher")
			} else {
				sb.WriteString("match.ArgsToInterfaces(")
				sb.WriteString(name)
				sb.WriteString(")...")
			}
			break
		}
	}
	f.ArgsMatchers = sb.String()

	sb.Reset()
	for i := 0; i < tupleLen; i++ {
		rawName := tf.rawNames[i]
		sb.WriteString(rawName)
		if rawName != "" {
			sb.WriteString(" ")
		}
		sb.WriteString(tf.types[i])
		if i != tupleLen-1 {
			sb.WriteString(", ")
		}
	}
	f.RawParams = sb.String()

	sb.Reset()
	for i := 0; i < tupleLen; i++ {
		sb.WriteString(names[i])
		sb.WriteString(" ")
		sb.WriteString(tf.types[i])
		if i != tupleLen-1 {
			sb.WriteString(", ")
		}
	}
	f.NamedParams = sb.String()

	sb.Reset()
	for i := 0; i < tupleLen-1; i++ {
		sb.WriteString(names[i])
		sb.WriteString(" match.Arg[")
		sb.WriteString(tf.types[i])
		sb.WriteString("]")
		sb.WriteString(", ")
	}
	if tupleLen > 0 {
		i := tupleLen - 1
		sb.WriteString(names[i])
		if !tf.isVariadic {
			sb.WriteString(" match.Arg[")
			sb.WriteString(tf.types[i])
			sb.WriteString("]")
		} else {
			sb.WriteString(" ...match.Arg[")
			sb.WriteString(tf.types[i][3:])
			sb.WriteString("]")
		}
	}
	f.NamedArgedParams = sb.String()

	if tf.isVariadic {
		f.VariadicArgs = "_args..."

		sb.Reset()
		sb.WriteString("\t_args := append([]mock.Matcher{")
		for i := 0; i < tupleLen-1; i++ {
			name := names[i]
			sb.WriteString(name)
			sb.WriteString(".Matcher")
			if i != tupleLen-2 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString("}, match.ArgsToMatchers(")
		sb.WriteString(names[tupleLen-1])
		sb.WriteString(")...)\n")
		f.VariadicArgsMatchersEval = sb.String()

		sb.Reset()
		sb.WriteString("\t_args := []any{")
		for i := 0; i < tupleLen-1; i++ {
			name := names[i]
			sb.WriteString(name)
			sb.WriteString(", ")
		}
		sb.WriteString("mock.SliceToAnySlice(" + names[tupleLen-1] + ")}\n")
		f.VariadicArgsEval = sb.String()
	}

	return f
}

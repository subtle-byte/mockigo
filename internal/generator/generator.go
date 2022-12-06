package generator

import (
	"bytes"
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/subtle-byte/mockigo/internal/generator/path_trie"
	"github.com/subtle-byte/mockigo/internal/generator/string_util"
	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/packages"
)

type Interfaces struct {
	IncludeInterfaces    bool
	InterfacesExceptions map[string]struct{}
}

type DirPath = string

func Generate(rootDir string, pkgDirs map[DirPath]Interfaces, mocksDir string) error {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedTypes | packages.NeedTypesInfo,
		Dir:        rootDir,
		Tests:      false,
		BuildFlags: nil, // TODO?
	}, maps.Keys(pkgDirs)...)
	if err != nil {
		return fmt.Errorf("load packages: %w", err)
	}
	if len(pkgs) == 1 && len(pkgs[0].GoFiles) == 0 {
		log.Printf("WARN: no packages found\n")
		return nil
	}
	for _, pkg := range pkgs {
		goFiles := pkg.GoFiles
		if len(goFiles) == 0 {
			continue
		}
		pkgDir := filepath.Dir(goFiles[0])
		if errs := pkg.Errors; len(errs) != 0 {
			log.Printf("ERROR: load package in dir %s: %v\n", pkgDir, errs)
			continue
		}
		inRoot, err := filepath.Rel(rootDir, pkgDir)
		if err != nil {
			log.Printf("ERROR: get rel path of package dir %s in root dir: %v\n", pkgDir, err)
			continue
		}
		mockDir := filepath.Join(mocksDir, inRoot)
		interfacesSelector, ok := pkgDirs[pkgDir]
		if !ok {
			panic(fmt.Sprintf("%#v %q", pkgDirs, pkgDir))
		}
		err = generateForPackage(pkg.Types, mockDir, interfacesSelector)
		if err != nil {
			log.Printf("ERROR: generate for package %s: %v\n", pkg.PkgPath, err)
			continue
		}
	}
	return nil
}

func generateForPackage(pkg *types.Package, mockDir string, interfaces Interfaces) error {
	for _, ident := range pkg.Scope().Names() {
		if _, ok := interfaces.InterfacesExceptions[ident]; ok == interfaces.IncludeInterfaces {
			continue
		}
		obj := pkg.Scope().Lookup(ident)
		if !obj.Exported() {
			continue
		}
		typeNameObj, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}
		var typeParams *types.TypeParamList
		if namedType, ok := obj.Type().(*types.Named); ok {
			typeParams = namedType.TypeParams()
		}
		methods := []*types.Func(nil)
		switch underlying := typeNameObj.Type().Underlying().(type) {
		case *types.Interface:
			for i := 0; i < underlying.NumMethods(); i++ {
				methods = append(methods, underlying.Method(i))
			}
		case *types.Signature:
			methods = []*types.Func{types.NewFunc(0, pkg, "Execute", underlying)}
		default:
			continue
		}
		generated, err := generateForInterface(pkg.Name(), ident, typeParams, methods)
		if err != nil {
			return fmt.Errorf("generate for interface (interface = %s, package = %s): %w", ident, pkg.Path(), err)
		}
		writeToFile(generated, filepath.Join(mockDir, string_util.CamelToSnake(ident)+".go"))
	}
	return nil
}

type writer struct {
	buf *bytes.Buffer
	err error
}

func (w *writer) Print(args ...string) {
	for _, arg := range args {
		_, err := w.buf.WriteString(arg)
		if err != nil && w.err == nil {
			w.err = err
		}
	}
}

func (w *writer) Println(args ...string) {
	w.Print(args...)
	_, err := w.buf.WriteString("\n")
	if err != nil && w.err == nil {
		w.err = err
	}
}

func generateImports(w *writer, methods []*types.Func, typeParams *types.TypeParamList,
) (importNames map[string]struct{}, pkgQualifier func(pkg *types.Package) string) {
	trie := path_trie.New()

	trie.LoadPath(path_trie.Path{"mock"}, "github.com/subtle-byte/mockigo/mock")
	trie.LoadPath(path_trie.Path{"match"}, "github.com/subtle-byte/mockigo/match")

	maxConsecutiveUnderscores := 0
	walkingPkgQualifier := func(pkg *types.Package) string {
		maxConsecutiveUnderscores = string_util.CountMaxConsecutiveUnderscores(pkg.Path(), maxConsecutiveUnderscores)
		splittedPath := strings.Split(pkg.Path(), "/")
		for i := range splittedPath {
			splittedPath[i] = strings.ReplaceAll(splittedPath[i], "-", "")
		}
		trie.LoadPath(splittedPath, pkg.Path())
		return ""
	}
	for _, method := range methods {
		signature := method.Type().(*types.Signature)
		for _, tuple := range [2]*types.Tuple{signature.Params(), signature.Results()} {
			for i := 0; i < tuple.Len(); i++ {
				types.TypeString(tuple.At(i).Type(), walkingPkgQualifier)
			}
		}
	}
	if typeParams != nil {
		for i := 0; i < typeParams.Len(); i++ {
			types.TypeString(typeParams.At(i).Constraint(), walkingPkgQualifier)
		}
	}

	importPathToName := map[string]string{}
	importNames = map[string]struct{}{}

	for _, reducedPath := range trie.ReducedPaths() {
		name := strings.Join(reducedPath.Rest, strings.Repeat("_", maxConsecutiveUnderscores+1))
		importNames[name] = struct{}{}
		importPathToName[reducedPath.Meta.(string)] = name
	}

	sortedImports := make([]string, 0, len(importPathToName))
	for path, name := range importPathToName {
		sortedImports = append(sortedImports, fmt.Sprintf("import %s %q", name, path))
	}
	sort.Strings(sortedImports)
	w.Println()
	for _, aImport := range sortedImports {
		w.Println(aImport)
	}

	pkgQualifier = func(pkg *types.Package) string {
		name, ok := importPathToName[pkg.Path()]
		if !ok {
			return "<internal error>"
		}
		return name
	}

	return
}

func genericFormats(typeParams *types.TypeParamList, pkgQualifier func(pkg *types.Package) string) (onInterface, inReceiver string) {
	if typeParams == nil {
		return "", ""
	}
	identAndType := make([]string, typeParams.Len())
	idents := make([]string, typeParams.Len())
	for i := 0; i < typeParams.Len(); i++ {
		tp := typeParams.At(i)
		ident := tp.Obj().Name()
		typeStr := types.TypeString(tp.Constraint(), pkgQualifier)
		identAndType[i] = ident + " " + typeStr
		idents[i] = ident
	}
	onInterface = "[" + strings.Join(identAndType, ", ") + "]"
	inReceiver = "[" + strings.Join(idents, ", ") + "]"
	return
}

func generateForInterface(pkgName, interfaceName string, typeParams *types.TypeParamList, methods []*types.Func) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := &writer{buf: buf}

	w.Println("// Code generated by mockigo. DO NOT EDIT.")
	w.Println()
	w.Println("package " + pkgName)

	importNames, pkgQualifier := generateImports(w, methods, typeParams)

	typeParamsOnInterface, typeParamsInReceiver := genericFormats(typeParams, pkgQualifier)

	w.Println()
	w.Println("var _ = match.Any[int]")
	w.Println()
	w.Println("type ", interfaceName, typeParamsOnInterface, " struct {\n\tmock *mock.Mock\n}")
	w.Println()
	w.Println("func New", interfaceName, typeParamsOnInterface, "(t mock.Testing) *", interfaceName, typeParamsInReceiver, " {")
	w.Println("\tt.Helper()")
	w.Println("\treturn &", interfaceName, typeParamsInReceiver, "{mock: mock.NewMock(t)}")
	w.Println("}")
	w.Println()
	expecterStruct := "_" + interfaceName + "_Expecter"
	w.Println("type ", expecterStruct, typeParamsOnInterface, " struct {\n\tmock *mock.Mock\n}")
	expecterStruct += typeParamsInReceiver
	w.Println()
	w.Println("func (_mock *", interfaceName, typeParamsInReceiver, ") EXPECT() ", expecterStruct, " {")
	w.Println("\t return ", expecterStruct, "{mock: _mock.mock}")
	w.Println("}")

	for _, method := range methods {
		generateForMethod(w, interfaceName, typeParamsOnInterface, typeParamsInReceiver, expecterStruct, method, pkgQualifier, importNames)
	}

	return buf.Bytes(), w.err
}

func isNilable(aType types.Type) bool {
	if underlying, ok := aType.(*types.Named); ok {
		aType = underlying.Underlying()
	}
	switch aType.(type) {
	case *types.Pointer, *types.Array, *types.Slice, *types.Map, *types.Chan, *types.Signature, *types.Interface:
		return true
	default:
		return false
	}
}

func generateForMethod(w *writer, interfaceName, typeParamsOnInterface, typeParamsInReceiver, expecterStruct string, method *types.Func, pkgQualifier func(pkg *types.Package) string, forbiddenNames map[string]struct{}) {
	signature := method.Type().(*types.Signature)
	variadic := signature.Variadic()

	inFormats := NewTupleFormatter(signature.Params(), variadic, pkgQualifier).Format("_a", forbiddenNames)
	outFormats := NewTupleFormatter(signature.Results(), false, pkgQualifier).Format("_r", forbiddenNames)

	callStruct := "_" + interfaceName + "_" + method.Name() + "_Call"
	w.Println()
	w.Println("type ", callStruct, typeParamsOnInterface, " struct {\n\t*mock.Call\n}")
	callStruct += typeParamsInReceiver
	w.Println("\nfunc (_mock *", interfaceName, typeParamsInReceiver, ") ", method.Name(), "(", inFormats.NamedParams, ") (", outFormats.RawParams, ") {")
	w.Println("\t_mock.mock.T.Helper()")
	w.Print(inFormats.VariadicArgsEval)
	results := "_results := "
	if signature.Results().Len() == 0 {
		results = ""
	}
	if !variadic {
		w.Println("\t", results, "_mock.mock.Called(\"", method.Name(), "\", ", inFormats.Args, ")")
	} else {
		w.Println("\t", results, "_mock.mock.Called(\"", method.Name(), "\", ", inFormats.VariadicArgs, ")")
	}
	if signature.Results().Len() != 0 {
		returnArgs := []string(nil)
		for i := 0; i < signature.Results().Len(); i++ {
			resultVar := signature.Results().At(i)
			varName := "_r" + strconv.Itoa(i)
			returnArgs = append(returnArgs, varName)
			typeStr := types.TypeString(resultVar.Type(), pkgQualifier)
			if typeStr == "error" {
				w.Println("\t", varName, " := _results.Error(", strconv.Itoa(i), ")")
			} else if !isNilable(resultVar.Type()) {
				w.Println("\t", varName, " := _results.Get(", strconv.Itoa(i), ").(", typeStr, ")")
			} else {
				w.Println("\tvar ", varName, " ", typeStr)
				w.Println("\tif _got := _results.Get(", strconv.Itoa(i), "); _got != nil {")
				w.Println("\t\t", varName, " = _got.(", typeStr, ")")
				w.Println("\t}")
			}
		}
		if len(returnArgs) != 0 {
			w.Println("\treturn ", strings.Join(returnArgs, ", "))
		}
	}
	w.Println("}")

	w.Println()
	w.Println("func (_expecter ", expecterStruct, ") ", method.Name(), "(", inFormats.NamedArgedParams, ") ", callStruct, " {")
	w.Print(inFormats.VariadicArgedArgsEval)
	if !variadic {
		w.Println("\treturn ", callStruct, `{Call: _expecter.mock.ExpectCall("`, method.Name(), `", `, inFormats.ArgedArgs, ")}")
	} else {
		w.Println("\treturn ", callStruct, `{Call: _expecter.mock.ExpectCall("`, method.Name(), `", `, inFormats.VariadicArgs, ")}")
	}
	w.Println("}")

	w.Println()
	w.Println("func (_call ", callStruct, ") Return(", outFormats.NamedParams, ") ", callStruct, " {")
	w.Println("\t_call.Call.Return(", outFormats.Args, ")")
	w.Println("\treturn _call")
	w.Println("}")

	w.Println("")
	runReturnFuncType := "func(" + inFormats.RawParams + ") (" + outFormats.RawParams + ")"
	w.Println("func (_call ", callStruct, ") RunReturn(f ", runReturnFuncType, ") ", callStruct, " {")
	w.Println("\t_call.Call.RunReturn(f)")
	w.Println("\treturn _call")
	w.Println("}")
}

func writeToFile(buf []byte, filePath string) error {
	mockDir := filepath.Dir(filePath)
	err := os.MkdirAll(mockDir, 0o777)
	if err != nil {
		return fmt.Errorf("create dirs to mock (mock dir = %s): %w", mockDir, err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(buf)
	if err != nil {
		return fmt.Errorf("write into file: %w", err)
	}
	return nil
}

package enameg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os/exec"
	"sort"
	"strings"

	"github.com/moznion/gowrtr/generator"
)

const annotation = "enameg"

var specialCharMap = map[string]string{
	"\\": "\\\\",
	`"`:  `\"`,
}

type constantVal struct {
	Name       string
	CommentVal string
}
type constant struct {
	TypeName string
	Vals     []constantVal
}

// Generate returns packageName and generated functions by paths.
func Generate(paths []string, useFormatter, enableEmptyName bool) (string, string) {
	constMap, err := collectConstants(paths)
	if err != nil {
		log.Fatal(err)
	}

	var packageName string
	var constants []constant

	sort.Strings(paths)
	for _, path := range paths {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		if packageName == "" {
			packageName = f.Name.Name
		}
		if packageName != f.Name.Name {
			log.Fatal("error: multiple packages")
		}

		cmap := ast.NewCommentMap(fset, f, f.Comments)

		// sort by node pos
		nodes := make([]ast.Node, 0, len(cmap))
		for n, _ := range cmap {
			nodes = append(nodes, n)
		}

		sort.Slice(nodes, func(i, j int) bool {
			lhs := nodes[i]
			rhs := nodes[j]
			return lhs.Pos() < rhs.Pos()
		})

		for _, node := range nodes {
			commentGroups := cmap[node]
			for _, cg := range commentGroups {
				if hasAnnotation(cg) {
					gen, ok := node.(*ast.GenDecl)
					if !ok || len(gen.Specs) <= 0 {
						continue
					}

					spec, ok := gen.Specs[0].(*ast.TypeSpec)
					if !ok {
						continue
					}

					c := newConst(spec.Name.Name, constMap, enableEmptyName)
					constants = append(constants, c)
				}
			}
		}
	}

	if len(constants) == 0 {
		return packageName, ""
	}

	generated, err := generateNameFunc(packageName, constants, useFormatter)
	if err != nil {
		log.Fatal(err)
	}
	return packageName, generated
}

func hasAnnotation(cg *ast.CommentGroup) bool {
	commentAnnotation := "+" + annotation
	for _, c := range cg.List {
		comment := strings.TrimSpace(strings.TrimLeft(c.Text, "//"))
		if strings.HasPrefix(comment, commentAnnotation) {
			return true
		}
	}
	return false
}

func collectConstants(paths []string) (map[string][]*ast.ValueSpec, error) {
	constMap := make(map[string][]*ast.ValueSpec)

	for _, path := range paths {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		constMap = collectConstantsInFile(constMap, f)
	}

	return constMap, nil
}

func collectConstantsInFile(constMap map[string][]*ast.ValueSpec, f *ast.File) map[string][]*ast.ValueSpec {
	for _, dec := range f.Decls {
		gen, ok := dec.(*ast.GenDecl)
		if !ok {
			continue
		}

		if gen.Tok != token.CONST {
			continue
		}

		for _, spec := range gen.Specs {
			val, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			valType, ok := val.Type.(*ast.Ident)
			if !ok {
				continue
			}

			typeName := valType.Name
			if constMap[typeName] == nil {
				constMap[typeName] = []*ast.ValueSpec{}
			}

			constMap[typeName] = append(constMap[typeName], val)
		}
	}

	return constMap
}

func newCommentVal(comment string) string {
	comment = strings.TrimSpace(strings.TrimLeft(comment, "//"))
	comment = strings.Fields(comment)[0]

	specialDelimiters := []string{".", "。"}
	for _, del := range specialDelimiters {
		comment = strings.Split(comment, del)[0]
	}

	for c, rep := range specialCharMap {
		comment = strings.Replace(comment, c, rep, -1)
	}
	return comment
}

func newConst(typeName string, constMap map[string][]*ast.ValueSpec, enableEmptyName bool) constant {
	nodes := constMap[typeName]
	vals := make([]constantVal, 0, len(nodes))

	for _, n := range nodes {
		hasComment := n.Comment != nil && len(n.Comment.List) > 0

		if !enableEmptyName && !hasComment {
			continue
		}

		var commentVal string
		if hasComment {
			commentVal = newCommentVal(n.Comment.List[0].Text)
		} else {
			commentVal = ""
		}

		vals = append(vals, constantVal{
			Name:       n.Names[0].Name,
			CommentVal: commentVal,
		})
	}

	return constant{
		TypeName: typeName,
		Vals:     vals,
	}
}

func generateNameFunc(packageName string, consts []constant, useFormatter bool) (string, error) {
	g := generator.NewRoot(
		generator.NewPackage(packageName),
		generator.NewNewline(),
		generator.NewImport("fmt"),
		generator.NewNewline(),
	)

	for _, c := range consts {
		caseStatements := make([]*generator.Case, 0, len(c.Vals))
		for _, v := range c.Vals {
			caseStatements = append(caseStatements, generator.NewCase(v.Name, generator.NewReturnStatement(fmt.Sprintf(`"%s"`, v.CommentVal))))
		}

		g = g.AddStatements(
			generator.NewComment(fmt.Sprintf(" Name returns the %s Name.", c.TypeName)),
			generator.NewFunc(
				generator.NewFuncReceiver("src", c.TypeName),
				generator.NewFuncSignature("Name").AddReturnTypes("string"),
			).AddStatements(
				generator.NewSwitch("src").
					AddCase(caseStatements...).
					Default(generator.NewDefaultCase(generator.NewReturnStatement(`fmt.Sprintf("%v", src)`))),
			),
			generator.NewNewline(),
		)
	}

	generated, err := g.Generate(0)
	if err != nil {
		return "", err
	}

	if useFormatter {
		return gofmt(generated)
	}
	return generated, nil
}

func gofmt(code string) (string, error) {
	cmd := exec.Command("gofmt")
	cmd.Stdin = strings.NewReader(code)
	b, err := cmd.Output()

	return string(b), err
}

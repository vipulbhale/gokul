package reflect

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	cntrlSpec      = new(controllerSpec)
	packageNameMap = make(map[string]string)
)

const MAIN = `// GENERATED CODE - DO NOT EDIT
package controllerwrapper

import (
	"reflect"{{ range $index, $packageName := .PackageName }}
	"{{ $packageName }}"{{ end }}
)

var (
	mapOfControllerNameToControllerObj = make(map[string]reflect.Type)
)

func RegisterControllers(){
	{{range $index, $element := .ControllerName}}
    		{{ $element | ToLower }} := {{ $element }}{}
    		typeOfController := reflect.TypeOf({{ $element | ToLower }})
    		mapOfControllerNameToControllerObj[typeOfController.Name()] = typeOfController
	{{ end }}
}

func New(name string) (interface{}, bool) {
	t, ok := mapOfControllerNameToControllerObj[name]
	if !ok {
		return nil, false
	}
	v := reflect.New(t)
	return v.Interface(), true
}
`

type controllerSpec struct {
	ControllerName []string
	PackageName    []string
}

func init() {
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)
	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

func ScanAppsDirectory(configuration map[string]string) {
	log.Debugln("Entering the ScanAppsDirectory.")
	log.Debugln("inputs are :: ", configuration)

	srcRoot, _ := os.Getwd()
	log.Debugln("The srcRoot is :: ", srcRoot)
	//srcRelativeAppsHomeDirPath := filepath.Join("github.com", "apps")
	appsHomeDirPath := filepath.Join(srcRoot, "src", "github.com", "apps")

	directoryList := []string{}

	err := filepath.Walk(appsHomeDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "controller" {
			log.Debugln("Path is ", path)
			log.Debugln("Directory is ", path)
			dir := filepath.Dir(path + "/" + info.Name())
			log.Debugln("Dir is ", dir)
			packagename := strings.Split(path, filepath.Join(srcRoot, "src"))[1]
			packagename = strings.Replace(packagename, "/", "", 1)
			log.Debugln("PackageName is ", packagename)
			packageNameMap[packagename] = packagename
			directoryList = append(directoryList, dir)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Debugln("Found directories are ", directoryList)

	for packageNameKey, _ := range packageNameMap {
		cntrlSpec.PackageName = append(cntrlSpec.PackageName, packageNameKey)
		log.Debugln("Package list is ", cntrlSpec.PackageName)
	}

	//srcRoot = filepath.Join(srcRoot, "gokul", "src", "github.com")
	//log.Debugln("srcRoot is ", srcRoot)

	for _, directoryController := range directoryList {
		makeControllers(directoryController)
	}

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}
	tmpl, err := template.New("test").Funcs(funcMap).Parse(MAIN)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, cntrlSpec)
	if err != nil {
		panic(err)
	}
}

func makeControllers(srcRoot string) {
	log.Debugln("Entering the makeControllers method")
	var allNamed []*types.Object

	structMap := make(map[string]reflect.Type)
	log.Debugln("map is ", structMap)
	kPath := filepath.Join(srcRoot)
	log.Debugln("kpath is ", kPath)

	fset := token.NewFileSet()

	pkgs, e := parser.ParseDir(fset, kPath, func(f os.FileInfo) bool {
		return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
	}, 0)
	if e != nil {
		log.Debugln(e)
		return
	}

	//kPathBaseController := filepath.Join(srcRoot, "gokul","controller")
	//fsetBaseController := token.NewFileSet()
	//
	//pkgsBaseController, e := parser.ParseDir(fsetBaseController, kPathBaseController, func(f os.FileInfo) bool {
	//	return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
	//}, 0)
	//if e != nil {
	//	log.Debugln(e)
	//	return
	//}
	//
	//for _, pkg := range pkgs {
	//	log.Debugln("package :: ", pkg.Name)
	//	for fn, f := range pkg.Files {
	//		log.Debugln("file :: ", fn)
	//		astf = append(astf, f)
	//	}
	//}

	//for _, pkg := range pkgsBaseController {
	//	log.Debugln("package :: ", pkg.Name)
	//	for fn, f := range pkg.Files {
	//		log.Debugln("file :: ", fn)
	//		astf = append(astf, f)
	//	}
	//}
	//
	//
	//
	//config := &types.Config{
	//	Error: func(e error) {
	//		fmt.Println("Error is ",e)
	//	},
	//	Importer: importer.For("source", nil),
	//	FakeImportC: false,
	//	DisableUnusedImportCheck: false,
	//}
	//info := types.Info{
	//
	//	Types: make(map[ast.Expr]types.TypeAndValue),
	//	Defs:  make(map[*ast.Ident]types.Object),
	//	Uses:  make(map[*ast.Ident]types.Object),
	//
	//
	//}
	//pkg, e := config.Check(srcRoot, fset, astf, &info)
	//if e != nil {
	//	log.Errorln("Is there any error :: ",e)
	//}
	//log.Debugln("types.Config.Check got ", pkg.String())

	//for _, name := range pkg.Scope().Names() {
	//	log.Debugln("Lookup name is :: ", pkg.Scope().Lookup(name))
	//	if obj, ok := pkg.Scope().Lookup(name).(types.Type); ok {
	//		log.Debugln("Object Type ", obj.String())
	//
	//	}
	//}

	log.Debugln("AllNamed is ", allNamed, structMap)

	log.Debugln("parsed package map is :: ", pkgs)
	for _, pkg := range pkgs {
		if pkg.Name == "controller" {
			log.Debugln("package name is ", pkg.Name)
			processPackage(pkg, pkg.Name)

		}
	}

}

func processPackage(pkg *ast.Package, packageName string) {
	log.Debugln("Entering the processPackage function")
	log.Debugln(pkg.Name)
	log.Debugln(pkg.Files)
	printASTVisitor := &PrintASTVisitor{}
	//controllers := make([]string,0)
	//printASTVisitor.cntrlSpec.packageControllers = make(map[string][]string)
	//printASTVisitor.cntrlSpec.ControllerName = make([]string, 0)
	//printASTVisitor.cntrlSpec.PackageName =	packageName
	//printASTVisitor.cntrlSpec.packageControllers[packageName] = controllers
	ast.Walk(printASTVisitor, pkg)

}

type PrintASTVisitor struct {
	info *types.Info
	//cntrlSpec controllerSpec
}

func (v *PrintASTVisitor) Visit(node ast.Node) ast.Visitor {
	// fmt.Println(v.info.Types)
	if node != nil {
		switch kk := node.(type) {

		case *ast.Package:
			{
				fmt.Println(kk.Name)

			}
		case *ast.TypeSpec:
			{
				log.Debugln("Name  of struct is :: " + kk.Name.Name)
				structType := kk.Type.(*ast.StructType)
				log.Debugln("hi there", structType)
				for _, field := range structType.Fields.List {
					log.Debugln(reflect.TypeOf(field.Type), " name is ", field.Names)
					fieldType := field.Type
					pkgName, typeName := func() (string, string) {
						// Drill through any StarExprs.
						for {
							if starExpr, ok := fieldType.(*ast.StarExpr); ok {
								fieldType = starExpr.X
								continue
							}
							break
						}

						// If the embedded type is in the same package, it's an Ident.
						if ident, ok := fieldType.(*ast.Ident); ok {
							return "", ident.Name
						}

						if selectorExpr, ok := fieldType.(*ast.SelectorExpr); ok {
							if pkgIdent, ok := selectorExpr.X.(*ast.Ident); ok {
								return pkgIdent.Name, selectorExpr.Sel.Name
							}
						}
						return "", ""
					}()
					if typeName == "BaseController" {
						log.Debugln("I am the man ", pkgName, typeName)
						//v.cntrlSpec.ControllerName = append(v.cntrlSpec.ControllerName, kk.Name.Name)
						cntrlSpec.ControllerName = append(cntrlSpec.ControllerName, kk.Name.Name)

						log.Debugln("ControllerSpec is :: ", cntrlSpec)

					}

				}

			}
		case *ast.GenDecl:
			{
				log.Debugln("Name  of struct is :: ", kk)

			}
		}
	}
	return v
}

// genSource renders the given template to produce source code, which it writes
// to the given directory and file.
//func genSource(dir, filename, templateSource string, args map[string]interface{}) {
//	sourceCode := revel.ExecuteTemplate(
//		template.Must(template.New("").Parse(templateSource)),
//		args)
//
//	// Create a fresh dir.
//	cleanSource(dir)
//	tmpPath := path.Join(revel.AppPath, dir)
//	err := os.Mkdir(tmpPath, 0777)
//	if err != nil && !os.IsExist(err) {
//		revel.ERROR.Fatalf("Failed to make '%v' directory: %v", dir, err)
//	}
//
//	// Create the file
//	file, err := os.Create(path.Join(tmpPath, filename))
//	defer file.Close()
//	if err != nil {
//		revel.ERROR.Fatalf("Failed to create file: %v", err)
//	}
//	_, err = file.WriteString(sourceCode)
//	if err != nil {
//		revel.ERROR.Fatalf("Failed to write to file: %v", err)
//	}
//}

//// genSource renders the given template to produce source code, which it writes
//// to the given directory and file.
//func genSource(dir, filename, templateSource string, args map[string]interface{}) {
//	sourceCode := revel.ExecuteTemplate(
//		template.Must(template.New("").Parse(templateSource)),
//		args)
//
//	// Create a fresh dir.
//	cleanSource(dir)
//	tmpPath := path.Join(revel.AppPath, dir)
//	err := os.Mkdir(tmpPath, 0777)
//	if err != nil && !os.IsExist(err) {
//		revel.ERROR.Fatalf("Failed to make '%v' directory: %v", dir, err)
//	}
//
//	// Create the file
//	file, err := os.Create(path.Join(tmpPath, filename))
//	defer file.Close()
//	if err != nil {
//		revel.ERROR.Fatalf("Failed to create file: %v", err)
//	}
//	_, err = file.WriteString(sourceCode)
//	if err != nil {
//		revel.ERROR.Fatalf("Failed to write to file: %v", err)
//	}
//}

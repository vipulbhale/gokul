package reflect

import (
	log "github.com/logrus"
	"go/ast"
	"go/token"
	//"go/types"
	"go/parser"
	"path/filepath"
	"reflect"
	//"github.com/gokul"
	"html/template"
	"os"
	"strings"
	//"go/types"
	//"go/importer"
	"go/types"
	"fmt"
	"go/importer"
)

type controllerSpec struct {
	ControllerName 	[]string
	PackageName	[]string
	packageControllers map[string][]string
}

func init() {
	// Output to stdout instead of the default stderr, could also be a file.
	log.SetOutput(os.Stdout)
	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

func ScanAppsDirectory(configuration map[string]string) {
	log.Debugln("Entering the ScanAppsDirectory")
	log.Debugln("inputs are", configuration)

	srcRoot, _ := os.Getwd()
	log.Debugln("The srcRoot is ", srcRoot)
	controllerPath := filepath.Join(srcRoot, "gokul", "src", "github.com", "gokul")
	srcRoot = filepath.Join(srcRoot, "gokul", "src", "github.com")
	log.Debugln("srcRoot is ", srcRoot)
	makeControllers(srcRoot, controllerPath)
}

func makeControllers(srcRoot string, controllerPath string) {
	log.Debugln("Entering the makeControllers method")
	astf := make([]*ast.File, 0)

	structMap := make(map[string]reflect.Type)
	log.Debugln("map is ", structMap)
	kPath := filepath.Join(srcRoot, "apps", "demoapp", "controller")

	fset := token.NewFileSet()

	pkgs, e := parser.ParseDir(fset, kPath, func(f os.FileInfo) bool {
		return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
	}, 0)
	if e != nil {
		log.Debugln(e)
		return
	}

	kPathBaseController := filepath.Join(srcRoot, "gokul","controller")
	fsetBaseController := token.NewFileSet()

	pkgsBaseController, e := parser.ParseDir(fsetBaseController, kPathBaseController, func(f os.FileInfo) bool {
		return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
	}, 0)
	if e != nil {
		log.Debugln(e)
		return
	}
	for _, pkg := range pkgsBaseController {
		fmt.Printf("package %v\n", pkg.Name)
		for fn, f := range pkg.Files {
			fmt.Printf("file %v\n", fn)
			astf = append(astf, f)
		}
	}



	config := &types.Config{
		Error: func(e error) {
			fmt.Println(e)
		},
		Importer: importer.Default(),
	}
	info := types.Info{

		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),

	}
	pkg, e := config.Check(kPath, fset, astf, &info)
	if e != nil {
		fmt.Println("Is there any errror %v",e)
	}
	fmt.Printf("types.Config.Check got %v\n", pkg.String())


	log.Debugln("parsed package map is :: ", pkgs)
	for _, pkg := range pkgs {
		if pkg.Name == "controller" {
			log.Debugln("package name is ", pkg.Name)
			processPackage(pkg, pkg.Name, &info)

		}
	}



}

func processPackage(pkg *ast.Package, importPkgName string, info *types.Info){
	printASTVisitor := &PrintASTVisitor{info}
	controllers := make([]string,0)
	printASTVisitor.cntrlSpec.packageControllers = make(map[string][]string)
	printASTVisitor.cntrlSpec.ControllerName = make([]string, 0)
	printASTVisitor.cntrlSpec.PackageName=	make([]string, 0)
	printASTVisitor.cntrlSpec.packageControllers[importPkgName] = controllers
	ast.Walk(printASTVisitor, pkg)

}

type PrintASTVisitor struct {
	info *types.Info
	cntrlSpec controllerSpec
}

func (v *PrintASTVisitor) Visit(node ast.Node) ast.Visitor {
	// fmt.Println(v.info.Types)
	if node != nil {
		switch kk := node.(type) {
		case *ast.TypeSpec:
			{
				log.Debugln("Name  of struct is :: " + kk.Name.Name)
				structType := kk.Type.(*ast.StructType)
				log.Debugln("hi there", structType.Struct)
				for _, field := range structType.Fields.List {
					log.Debugln(reflect.TypeOf(field.Type), " name is ", field.Names)
				}

				log.Debugln("ControllerSpec is :: ", v.cntrlSpec)
				tmpl, err := template.New("test").Parse(MAIN)
				if err != nil {
					panic(err)
				}
				err = tmpl.Execute(os.Stdout, v.cntrlSpec)
				if err != nil {
					panic(err)
				}
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

const MAIN = `// GENERATED CODE - DO NOT EDIT
package controllerwrapper

import (
	"reflect"
)

func InstantiateStructs(){

	{{.ControllerName}} := new({{.ControllerName}})
}
`

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

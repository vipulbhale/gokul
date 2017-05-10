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

	log.Debugln("parsed package map is :: ", pkgs)
	for _, pkg := range pkgs {
		if pkg.Name == "controller" {
			log.Debugln("package name is ", pkg.Name)
			processPackage(pkg, pkg.Name)

		}
	}
}

func processPackage(pkg *ast.Package, importPkgName string){
	printASTVisitor := &PrintASTVisitor{}
	controllers := make([]string,0)
	printASTVisitor.cntrlSpec.packageControllers[importPkgName] = controllers
	ast.Walk(printASTVisitor, pkg)

}

type PrintASTVisitor struct {
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
				log.Debugln("hi there %v\n", structType)
				for _, field := range structType.Fields.List {
					log.Debugln(reflect.TypeOf(field.Type), " name is ", field.Names)
				}
				cntrlSpec := &controllerSpec{ControllerName: make([]string,0) , PackageName: make([]string, 0)}

				log.Debugln("ControllerSpec is :: ", cntrlSpec)
				tmpl, err := template.New("test").Parse(MAIN)
				if err != nil {
					panic(err)
				}
				err = tmpl.Execute(os.Stdout, cntrlSpec)
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

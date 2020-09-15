package tests

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
)

var endToEndCases = []struct {
	name       string
	files      []fileConfig
	wantGenOut string
	wantRunOut string
	wantGenErr bool
	wantRunErr bool
}{
	{
		"happy1",
		[]fileConfig{
			{"main.go", "main", "import subpkg \"./subpkg\"\n//go:generate genembed EmbedFiles f1", map[string]string{"EmbedFiles": "f1", "subpkg.EmbedFiles": "f2"}},
			{"f1", "", `123123`, nil},
			{"subpkg/subpkg.go", "subpkg", `//go:generate genembed EmbedFiles f2`, nil},
			{"subpkg/f2", "", `456456`, nil},
		},
		"",                 // gen
		"123123\n456456\n", // run
		false,              // gen error
		false,              // run error
	},
	{
		"embedNotExistsFile",
		[]fileConfig{
			{"main.go", "main", `//go:generate genembed EmbedFiles notexistsfile
		`, map[string]string{"EmbedFiles": "notexistsfile"}},
			{"f1", "", `123123`, nil},
		},
		"failed open embeded file \"notexistsfile\"", // gen
		"\n",  // run
		true,  // gen error
		false, // run error
	},
	{
		"nothingEmbeded",
		[]fileConfig{
			{"main.go", "main", `//go:generate genembed EmbedFiles
	`, map[string]string{"EmbedFiles": "notexistsfile"}},
			{"f1", "", `123123`, nil},
		},
		"nothing to embeded",    // gen
		"undefined: EmbedFiles", // run
		true,                    // gen error
		true,                    // run error
	},
	{
		"emptyArgs",
		[]fileConfig{
			{"main.go", "main", `//go:generate genembed
	`, map[string]string{"EmbedFiles": "notexistsfile"}},
			{"f1", "", `123123`, nil},
		},
		"invalid arguments",     // gen
		"undefined: EmbedFiles", // run
		true,                    // gen error
		true,                    // run error
	},
}

func TestEndToEndCases(t *testing.T) {
	writeFile := func(t *testing.T, dir, file string, dat string) string {
		t.Helper()

		absFile := filepath.Join(dir, file)
		if err := os.MkdirAll(filepath.Dir(absFile), 0700); err != nil {
			t.Error("failed create dir on the fly")
		}
		err := ioutil.WriteFile(absFile, []byte(dat), 0666)
		if err != nil {
			t.Error(err, "failed write data to file", file)
		}
		return absFile
	}

	out, err := runBin(".", "go", "build", "-o", "../bin/genembed", "../genembed/genembed.go")
	if err != nil {
		t.Errorf("failed build genembed application err=%v, out=%s", err, out)
	}

	for _, test := range endToEndCases {
		t.Run(test.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "go-embed")
			if err != nil {
				t.Error(err)
			}
			// defer os.RemoveAll(dir)

			for _, file := range test.files {
				var buf bytes.Buffer
				switch {
				case file.Name == "main.go":
					err = mainGoTpl.Execute(&buf, file)
				case filepath.Ext(file.Name) == ".go" && file.Name != "main.go":
					err = someFileGoTpl.Execute(&buf, file)
				default:
					_, err = buf.WriteString(file.Code)
				}
				if err != nil {
					t.Error("failed write to file (or execute tpl)", err)
				}
				writeFile(t, dir, file.Name, buf.String())
			}

			t.Logf("work dir: %q", dir)

			out, err := runBin(dir, "go", "generate", "./...")
			if test.wantGenErr != (err != nil) {
				t.Errorf("go generate error=%v, wantErr=%v, out=%q", err, test.wantGenErr, out)
			}
			if !strings.Contains(out, test.wantGenOut) || (test.wantGenOut == "" && out != "") {
				t.Errorf("not contains gen our, got=%q, substr=%q", out, test.wantGenOut)
			}

			out, err = runBin(dir, "go", "run", ".")
			if test.wantRunErr != (err != nil) {
				t.Errorf("run bin error=%v, wantErr=%v, out=%q", err, test.wantRunErr, out)
			}
			t.Logf("run out: %q", out)

			if !strings.Contains(out, test.wantRunOut) || (test.wantRunOut == "" && out != "") {
				t.Errorf("not contains run out, got=%q, substr=%q", out, test.wantRunOut)
			}
		})
	}
}

func runBin(dir, name string, arg ...string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir
	cmd.Stderr = &buf
	cmd.Stdout = &buf

	pwd, _ := os.Getwd()
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH")+":"+pwd+"/../bin")

	err := cmd.Run()
	return buf.String(), err
}

type fileConfig struct {
	Name       string
	Pkg        string
	Code       string
	PrintFiles map[string]string
}

var mainGoTpl = template.Must(template.New("main.go").Parse(`package {{.Pkg}}

{{.Code}}

func main() {
	{{- range $filedName, $fileName := .PrintFiles }}
		println(string({{ $filedName }}["{{ $fileName }}"]))
	{{- end }}
}
`))
var someFileGoTpl = template.Must(template.New("somefile.go").Parse(`package {{.Pkg}}

{{.Code}}
`))

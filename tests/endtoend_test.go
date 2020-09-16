package tests

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
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
		"failed open embedded file \"notexistsfile\"", // gen
		"\n",  // run
		true,  // gen error
		false, // run error
	},
	{
		"nothingEmbedded",
		[]fileConfig{
			{"main.go", "main", `//go:generate genembed EmbedFiles
	`, map[string]string{"EmbedFiles": "notexistsfile"}},
			{"f1", "", `123123`, nil},
		},
		"nothing to embedded",   // gen
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
		err := os.MkdirAll(filepath.Dir(absFile), 0700)
		require.NoError(t, err, "failed create dir on the fly")

		err = ioutil.WriteFile(absFile, []byte(dat), 0666)
		require.NoError(t, err, "failed write data to file")

		return absFile
	}

	out, err := runBin(".", "go", "build", "-o", "../bin/genembed", "../genembed/genembed.go")
	require.NoError(t, err, "failed build genembed application err=%v, out=%s", err, out)

	for _, test := range endToEndCases {
		t.Run(test.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "genembed")
			require.NoError(t, err, "failed create temporary dir")

			defer os.RemoveAll(dir)

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

				require.NoError(t, err, "failed write to file (or execute tpl)")

				writeFile(t, dir, file.Name, buf.String())
			}

			t.Logf("work dir: %q", dir)

			out, err := runBin(dir, "go", "generate", "./...")
			if test.wantGenErr != (err != nil) {
				t.Errorf("go generate error=%v, wantErr=%v, out=%q", err, test.wantGenErr, out)
			}
			if test.wantGenOut == "" {
				require.Empty(t, out)
			} else {
				require.Contains(t, out, test.wantGenOut)
			}

			out, err = runBin(dir, "go", "run", ".")
			if test.wantRunErr != (err != nil) {
				t.Errorf("run bin error=%v, wantErr=%v, out=%q", err, test.wantRunErr, out)
			}
			if test.wantRunOut == "" {
				require.Empty(t, test.wantRunOut)
			} else {
				require.Contains(t, out, test.wantRunOut)
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

module github.com/maniartech/uexl

go 1.22

require (
	github.com/maniartech/gotime/v2 v2.0.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// gotime is the NITES reference implementation; it backs the datetime standard
// library's formatDate/parseDate. It is an unpublished local module, so it is
// wired via a filesystem replace. Replace with a published version (and drop this
// directive) once gotime is tagged, or override locally with a go.work file.
replace github.com/maniartech/gotime/v2 => E:/Projects/gotime/gotime

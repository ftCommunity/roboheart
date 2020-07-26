package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/thoas/go-funk"
	"go/format"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func main() {
	parser := argparse.NewParser("setservices", "Service list definition generator")
	args := map[string]*bool{}
	sns := funk.Keys(services).([]string)
	sort.Strings(sns)
	for _, sn := range sns {
		args[sn] = parser.Flag("", sn, &argparse.Options{
			Required: false,
			Help:     "Enable service " + sn,
			Default:  false,
		})
	}
	all := parser.Flag("a", "all", &argparse.Options{
		Required: false,
		Help:     "Enable all services",
		Default:  false,
	})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	var rsl []string
	if *all {
		rsl = sns
	} else {
		for a, s := range args {
			if *s {
				rsl = append(rsl, a)
			}
		}
	}
	sort.Strings(rsl)
	var imports []string
	imports = append(imports, "\"github.com/ftCommunity/roboheart/internal/service\"")
	var sl []string
	for _, sn := range rsl {
		sd, ok := services[sn]
		if !ok {
			panic("unknown service")
		}
		imports = append(imports, "\""+sd[0]+"\"")
		sl = append(sl, sd[1]+"."+sd[2]+",")
	}

	var output []string
	output = append(output, "package servicemanager", "")
	output = append(output, "import (")
	output = append(output, imports...)
	output = append(output, ")", "")
	output = append(output, "var services = []service.Service{")
	output = append(output, sl...)
	output = append(output, "}")
	code, err := format.Source([]byte(strings.Join(output, "\n")))
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("internal/servicemanager/services.go", code, 0644)
}

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
	services := make(map[string][4]string)
	for ip, sp := range serviceproviders {
		for sn, sd := range sp {
			if _, ok := services[sn]; ok {
				panic("Service " + sn + " has already been provided by another service provider")
			}
			var fssd [4]string
			copy(fssd[:],append([]string{ip}, sd[:]...)[:4])
			services[sn] = fssd
		}
	}
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
	imports = append(imports, "\"github.com/ftCommunity-roboheart/roboheart/package/service\"")
	var sl []string
	for _, sn := range rsl {
		sd, ok := services[sn]
		if !ok {
			panic("unknown service")
		}
		imports = append(imports, "\""+sd[0]+"/"+sd[1]+"\"")
		sl = append(sl, sd[2]+"."+sd[3]+",")
	}

	var output []string
	output = append(output, "package services", "")
	output = append(output, "import (")
	output = append(output, imports...)
	output = append(output, ")", "")
	output = append(output, "var Services = []service.Service{")
	output = append(output, sl...)
	output = append(output, "}")
	code, err := format.Source([]byte(strings.Join(output, "\n")))
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("internal/services/services.go", code, 0644)
}

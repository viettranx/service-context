// Copyright (c) 2023, Viet Tran, 200Lab Team.

package sctx

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/facebookgo/flagenv"
)

func isZeroValue(f *flag.Flag, value string) bool {
	typ := reflect.TypeOf(f.Value)
	var z reflect.Value
	if typ.Kind() == reflect.Ptr {
		z = reflect.New(typ.Elem())
	} else {
		z = reflect.Zero(typ)
	}
	if value == z.Interface().(flag.Value).String() {
		return true
	}

	switch value {
	case "false":
		return true
	case "":
		return true
	case "0":
		return true
	}
	return false
}

func getEnvName(name string) string {
	name = strings.Replace(name, ".", "_", -1)
	name = strings.Replace(name, "-", "_", -1)

	if flagenv.Prefix != "" {
		name = flagenv.Prefix + name
	}

	return strings.ToUpper(name)
}

type AppFlagSet struct {
	*flag.FlagSet
}

func newFlagSet(name string, fs *flag.FlagSet) *AppFlagSet {
	fSet := &AppFlagSet{fs}
	fSet.Usage = flagCustomUsage(name, fSet)
	return fSet
}

func (f *AppFlagSet) GetSampleEnvs() {
	f.VisitAll(func(f *flag.Flag) {
		if f.Name == "outenv" {
			return
		}

		s := fmt.Sprintf("## %s (-%s)\n", f.Usage, f.Name)
		s += fmt.Sprintf("#%s=", getEnvName(f.Name))

		if !isZeroValue(f, f.DefValue) {
			t := fmt.Sprintf("%T", f.Value)
			if t == "*flag.stringValue" {
				// put quotes on the value
				s += fmt.Sprintf("%q", f.DefValue)
			} else {
				s += fmt.Sprintf("%v", f.DefValue)
			}
		}
		fmt.Print(s, "\n\n")
	})
}

func (f *AppFlagSet) Parse(args []string) {
	flagenv.Parse()
	_ = f.FlagSet.Parse(args)
}

func flagCustomUsage(name string, fSet *AppFlagSet) func() {
	return func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)

		fSet.VisitAll(func(f *flag.Flag) {
			s := fmt.Sprintf("  -%s", f.Name)
			name, usage := flag.UnquoteUsage(f)
			if len(name) > 0 {
				s += " " + name
			}

			if len(s) <= 4 {
				s += "\t"
			} else {
				s += "\n    \t"
			}

			s += usage

			if !isZeroValue(f, f.DefValue) {
				t := fmt.Sprintf("%T", f.Value)
				if t == "*flag.stringValue" {
					s += fmt.Sprintf(" (default %q)", f.DefValue)
				} else {
					s += fmt.Sprintf(" (default %v)", f.DefValue)
				}
			}
			s += fmt.Sprintf(" [$%s]", getEnvName(f.Name))
			_, _ = fmt.Fprint(os.Stderr, s, "\n")
		})
	}
}

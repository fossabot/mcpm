package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
)

var (
	flagset      = flag.NewFlagSet("", flag.ExitOnError)
	verbose      bool
	downloadOnly bool
	versionQuery string
	modes        = map[string]func(){
		"get":    getPackage,
		"search": searchPackage,
		"update": updateCache,
		"info":   readPackageInfo,
	}
	homeDir = "."
	_DBASE  *_Database
	_PKGS   *_PackageNames
)

func checkHomeDir() {
	u, ue := user.Current()
	must(ue)
	homeDir = fmt.Sprintf("%s/.mcpm", u.HomeDir)
	must(mkDirIfNotExist(homeDir))
}

func main() {
	var mode string
	if len(os.Args) >= 2 {
		mode = os.Args[1]
	} else {
		mode = "help"
	}
	flagset.Usage = func() {
		fmt.Fprintln(os.Stderr, "mcpm – Minecraft Package Manager\nAvailable modes:\n ")
		for m, _ := range modes {
			fmt.Fprintf(os.Stderr, " %s\n", m)
		}
		fmt.Fprintln(os.Stderr, "\nAvailable options: ")
		flagset.PrintDefaults()
	}
	flagset.BoolVar(&verbose, "v", false, "Verbose (WIP)")
	flagset.BoolVar(&downloadOnly, "d", false, "Download only")
	flagset.StringVar(&versionQuery, "q", "", "Version query (like \"latest:beta:mc1.7.10\")")
	if len(os.Args) >= 3 {
		flagset.Parse(os.Args[2:])
	}
	if f, ok := modes[mode]; ok {
		checkHomeDir()
		per := readGobGzip(homePath(pnFile), &_PKGS)
		der := readGobGzip(homePath(dbFile), &_DBASE)
		if der != nil || per != nil {
			updateCache()
		}
		f()
	} else {
		flagset.Usage()
	}
}

package cli

import "flag"

type Arguments struct {
	Flags   Flags
	Targets []string
}

type Flags struct {
	Config bool
	Index  bool
	Plan   bool
}

func Parse() Arguments {
	config := flag.Bool("config", false, "display complete configuration")
	index := flag.Bool("index", false, "display index after flattening configuration")
	plan := flag.Bool("plan", false, "display the execution plan")
	flag.Parse()
	return Arguments{
		Flags: Flags{
			Config: *config,
			Index:  *index,
			Plan:   *plan,
		},
		Targets: flag.Args(),
	}
}

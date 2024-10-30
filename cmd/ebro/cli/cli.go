package cli

import "flag"

type Arguments struct {
	Flags   Flags
	Targets []string
}

type Flags struct {
	Config bool
	Plan   bool
}

func Parse() Arguments {
	config := flag.Bool("config", false, "display complete configuration")
	plan := flag.Bool("plan", false, "display the execution plan")
	flag.Parse()
	return Arguments{
		Flags: Flags{
			Config: *config,
			Plan:   *plan,
		},
		Targets: flag.Args(),
	}
}

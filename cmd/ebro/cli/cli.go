package cli

import "flag"

type Arguments struct {
	Flags   Flags
	Targets []string
}

type Flags struct {
	Config  bool
	Catalog bool
	Plan    bool
}

func Parse() Arguments {
	config := flag.Bool("config", false, "display all imported configuration files merged into one")
	catalog := flag.Bool("catalog", false, "display complete catalog of tasks with their definitive configuration")
	plan := flag.Bool("plan", false, "display the execution plan")
	flag.Parse()

	targets := []string{":default"}
	args := flag.Args()
	if len(args) > 0 {
		targets = []string{}
		for _, arg := range args {
			targets = append(targets, ":"+arg)
		}
	}

	return Arguments{
		Flags: Flags{
			Config:  *config,
			Catalog: *catalog,
			Plan:    *plan,
		},
		Targets: targets,
	}
}

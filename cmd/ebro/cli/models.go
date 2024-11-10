package cli

type Command string
type Flag string

const (
	CommandHelp    Command = "help"
	CommandVersion Command = "version"
	CommandConfig  Command = "config"
	CommandCatalog Command = "catalog"
	CommandPlan    Command = "plan"
	CommandRun     Command = "run"
)

const (
	FlagFile  Flag = "file"
	FlagForce Flag = "force"
)

type commandInfo struct {
	Command     Command
	Description string
}

type flagInfo struct {
	Flag        Flag
	Description string
}

var commandList = []commandInfo{
	{CommandHelp, "Displays this help message"},
	{CommandVersion, "Display ebro's version"},
	{CommandConfig, "Display all imported configuration files merged into one"},
	{CommandCatalog, "Display complete catalog of tasks with their definitive configuration"},
	{CommandPlan, "Display the execution plan"},
	{CommandRun, ""},
}

var flagList = []flagInfo{
	{FlagFile, "Specify the file that should be loaded as root module. default: Ebro.yaml"},
	{FlagForce, "Ignore when.* conditionals and dont skip any task"},
}

type Arguments struct {
	Command Command
	File    string
	Targets []string
	Force   bool
}

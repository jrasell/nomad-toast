package buildconsts

var programName = ""

// SetProgramName sets the name of the currently executing program. It should be called at
// the earliest opportunity (e.g. in PersistentPreRun of the root Cobra command).
func SetProgramName(name string) {
	programName = name
}

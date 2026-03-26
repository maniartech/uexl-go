package uexl

import "github.com/maniartech/uexl/vm"

// envConfig is the unexported accumulation state used during Env construction.
// It is populated by applying Option functions, then frozen into an immutable Env.
type envConfig struct {
	functions    vm.VMFunctions
	pipeHandlers vm.PipeHandlers
	globals      map[string]any
}

// Lib is implemented by packages that ship reusable bundles of UExL extensions.
// Apply is called exactly once during Env construction; it registers functions,
// pipe handlers, and globals into the supplied EnvConfig.
type Lib interface {
	Apply(cfg *EnvConfig)
}

// EnvConfig is the public projection of envConfig, accessible only inside Lib.Apply.
// It exposes additive operations only — no replacement, no read-back of values.
type EnvConfig struct {
	cfg *envConfig
}

// AddFunctions merges fns into the in-progress env configuration.
// Later calls for the same key win. Panics if fns is nil.
func (c *EnvConfig) AddFunctions(fns Functions) {
	if fns == nil {
		panic("uexl: EnvConfig.AddFunctions: fns must not be nil")
	}
	for k, v := range fns {
		c.cfg.functions[k] = v
	}
}

// AddPipeHandlers merges pipes into the in-progress env configuration.
// Later calls for the same key win. Panics if pipes is nil.
func (c *EnvConfig) AddPipeHandlers(pipes PipeHandlers) {
	if pipes == nil {
		panic("uexl: EnvConfig.AddPipeHandlers: pipes must not be nil")
	}
	for k, v := range pipes {
		c.cfg.pipeHandlers[k] = v
	}
}

// AddGlobals merges vars into the in-progress env configuration.
// Later calls for the same key win. Panics if vars is nil.
func (c *EnvConfig) AddGlobals(vars map[string]any) {
	if vars == nil {
		panic("uexl: EnvConfig.AddGlobals: vars must not be nil")
	}
	for k, v := range vars {
		c.cfg.globals[k] = v
	}
}

package sdk

// Plugin is intentionally small. Optional behavior is exposed through capabilities.
type Plugin interface{ Info() PluginInfo }
type PluginInfo struct{ ID, Name, Version, Description string }

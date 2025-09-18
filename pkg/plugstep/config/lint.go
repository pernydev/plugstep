package config

func (c *PlugstepConfig) Lint() []string {
	issues := []string{}

	if c.Server.Version == "latest" {
		issues = append(issues, "Using version = latest on server jar, can be good for security, possible API versioning issues")
	}
	return issues
}

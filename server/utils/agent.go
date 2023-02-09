package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mileusna/useragent"
)

type Agent struct {
	OS      string `json:"os"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ParseUserAgent(c *fiber.Ctx) Agent {
	agent := c.Get("User-Agent")

	// Detect BassieMusic apps
	if strings.HasPrefix(agent, "BassieMusic macOS App") {
		return Agent{"macOS", "BassieMusic App", strings.Split(agent, "v")[1]}
	}
	if strings.HasPrefix(agent, "BassieMusic Windows App") {
		return Agent{"Windows", "BassieMusic App", strings.Split(agent, "v")[1]}
	}
	if strings.HasPrefix(agent, "BassieMusic Linux App") {
		return Agent{"Linux", "BassieMusic App", strings.Split(agent, "v")[1]}
	}

	// Parse browser agent
	parsedAgent := useragent.Parse(agent)
	return Agent{parsedAgent.OS, parsedAgent.Name, parsedAgent.Version}
}

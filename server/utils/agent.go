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
	if strings.HasPrefix(agent, "BassieMusic Android App") {
		return Agent{"Android", "BassieMusic App", strings.Split(agent, "v")[1]}
	}
	if strings.HasPrefix(agent, "BassieMusic iOS App") {
		return Agent{"iOS", "BassieMusic App", strings.Split(agent, "v")[1]}
	}
	if strings.HasPrefix(agent, "BassieMusic macOS App") {
		return Agent{"macOS", "BassieMusic App", strings.Split(agent, "v")[1]}
	}
	if strings.HasPrefix(agent, "BassieMusic Windows App") {
		return Agent{"Windows", "BassieMusic App", strings.Split(agent, "v")[1]}
	}
	if strings.HasPrefix(agent, "BassieMusic Linux App") {
		return Agent{"Linux", "BassieMusic App", strings.Split(agent, "v")[1]}
	}
	if strings.HasPrefix(agent, "BassieMusic Flutter App") {
		return Agent{"Flutter", "BassieMusic App", strings.Split(agent, "v")[1]}
	}

	// Parse browser agent
	parsedAgent := useragent.Parse(agent)
	if parsedAgent.OS == "" {
		parsedAgent.OS = "?"
	}
	if parsedAgent.Name == "" {
		parsedAgent.Name = "?"
	}
	if parsedAgent.Version == "" {
		parsedAgent.Version = "?"
	}
	return Agent{parsedAgent.OS, parsedAgent.Name, parsedAgent.Version}
}

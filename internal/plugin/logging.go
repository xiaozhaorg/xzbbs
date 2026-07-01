package plugin

import (
	"log"
)

// LoggingPlugin logs events to console
type LoggingPlugin struct{}

func (p *LoggingPlugin) Name() string { return "logging" }

// Register hooks
func (p *LoggingPlugin) Boot() {
	Global().Register(EventThreadCreated, func(args ...interface{}) (interface{}, error) {
		log.Printf("[PLUGIN] New thread created: id=%v by user=%v", args[0], args[1])
		return nil, nil
	})
	Global().Register(EventPostCreated, func(args ...interface{}) (interface{}, error) {
		log.Printf("[PLUGIN] New post created: id=%v by user=%v", args[0], args[1])
		return nil, nil
	})
	Global().Register(EventUserRegistered, func(args ...interface{}) (interface{}, error) {
		log.Printf("[PLUGIN] New user registered: %v", args[0])
		return nil, nil
	})
}

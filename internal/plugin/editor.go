package plugin

// EditorPlugin provides server-side hooks for editor-related events
// This plugin integrates the frontend BBEditor with backend processing
type EditorPlugin struct{}

func (p *EditorPlugin) Name() string { return "editor" }

// Register hooks for post/thread lifecycle
func (p *EditorPlugin) Boot() {
	// Hook: before post save - can add watermarks, auto-link, etc.
	Global().Register(EventBeforePostSave, func(args ...interface{}) (interface{}, error) {
		if len(args) > 0 {
			if content, ok := args[0].(string); ok && content != "" {
				// Could add watermark, auto-link URLs, etc.
				return content, nil
			}
		}
		return args[0], nil
	})

	// Hook: after post save - can trigger notifications, cache invalidation
	Global().Register(EventAfterPostSave, func(args ...interface{}) (interface{}, error) {
		// Post-save processing
		return true, nil
	})

	// Hook: after thread creation - index for search
	Global().Register(EventThreadCreated, func(args ...interface{}) (interface{}, error) {
		// Thread indexing hook - frontend BBEditor submits here
		return true, nil
	})
}

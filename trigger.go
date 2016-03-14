package awakengine

// Trigger is everything to do with reacting to the player or time or ...
type Trigger struct {
	Active    func() bool    // can we fire?
	Depends   []string       // have all our buddies fired?
	Fire      func()         // we have the com.
	Dialogues []DialogueLine // push it!
	Fired     bool           // others are depending on us.
}

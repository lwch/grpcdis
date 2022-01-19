package command

// Flag command flag
type Flag byte

// from https://redis.io/commands/command
const (
	// FlagWrite command may result in modifications
	FlagWrite Flag = iota
	// FlagReadOnly command will never modify keys
	FlagReadOnly
	// FlagDenyOOM reject command if currently out of memory
	FlagDenyOOM
	// FlagAdmin server admin command
	FlagAdmin
	// FlagPubSub pubsub-related command
	FlagPubSub
	// FlagNoScript deny this command from scripts
	FlagNoScript
	// FlagRandom command has random results, dangerous for scripts
	FlagRandom
	// FlagSortForScript if called from script, sort output
	FlagSortForScript
	// FlagLoading allow command while database is loading
	FlagLoading
	// FlagStale allow command while replica has stale data
	FlagStale
	// FlagSkipMonitor do not show this command in MONITOR
	FlagSkipMonitor
	// FlagAsking cluster related - accept even if importing
	FlagAsking
	// FlagFast command operates in constant or log(N) time. Used for latency monitoring
	FlagFast
	// FlagMovableKeys keys have no pre-determined position. You must discover keys yourself
	FlagMovableKeys
)

// String flag to string
func (flag Flag) String() string {
	switch flag {
	case FlagWrite:
		return "write"
	case FlagReadOnly:
		return "readonly"
	case FlagDenyOOM:
		return "denyoom"
	case FlagAdmin:
		return "admin"
	case FlagPubSub:
		return "pubsub"
	case FlagNoScript:
		return "noscript"
	case FlagRandom:
		return "random"
	case FlagSortForScript:
		return "sort_for_script"
	case FlagLoading:
		return "loading"
	case FlagStale:
		return "stale"
	case FlagSkipMonitor:
		return "skip_monitor"
	case FlagAsking:
		return "asking"
	case FlagFast:
		return "fast"
	case FlagMovableKeys:
		return "movablekeys"
	}
	return ""
}

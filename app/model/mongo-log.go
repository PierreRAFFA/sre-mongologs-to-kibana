package model

const (
	ComponentCommand  = "COMMAND"
	ComponentNetwork  = "NETWORK"
	ComponentAccess   = "ACCESS"
	ComponentQuery    = "QUERY"
	ComponentRepl     = "REPL"
	ComponentStorage  = "STORAGE"
	ComponentControl  = "CONTROL"
	ComponentWrite    = "WRITE"
	ComponentCoonPool = "CONNPOOL"
	ComponentElection = "ELECTION"
	ComponentIndex    = "INDEX"
)

type MongoLog struct {
	Id        int
	DateTime  string      `json:"dateTime"`
	Severity  string      `json:"severity"`
	Component string      `json:"component"`
	Context   string      `json:"context"`
	Details   interface{} `json:"details"`
}

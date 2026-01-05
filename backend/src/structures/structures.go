package structures

// structures for backend

// Cleaner defines a type representing a cleaning operation with associated options and metadata.
type Cleaner struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Running bool      `json:"running"`
	Detect  Detection `json:"detect"`
	Options []Option  `json:"options"`
}

type Detection struct {
	Type    	string   `json:"type"` // "file", "always", "dir", "registry"
	Paths   	[]string           `json:"paths"`
	Registry 	[]RegistryCheck `json:"registry"`
}

type RegistryCheck struct {
	Key string   `json:"key"`
	OS  []string `json:"os,omitempty"`
}

// Option defines a type representing a cleaning operation option.
type Option struct {
	ID          string 		`json:"id"`
	Label       string 		`json:"label"` // Label represents the human-readable name for the option in the cleaning operation.
	Description string 		`json:"description"`
	Warning     string   `json:"warning,omitempty"`
	Actions     []Action `json:"actions"`
}

type Action struct {
	Command string 		`json:"command"` // "delete", "truncate", "vacuum"
	Search  string 		`json:"search"` // "file", "glob", "walk.files"
	Path    string 		`json:"path"`
	OS 	    []string 	`json:"os,omitempty"`
}

type ActionResult struct {
	Size      uint64
	FileCount uint64
	Paths     []string
}


// structures for requests

// CleanRequest - request from frontend
type CleanRequest struct {
	CleanerID string `json:"cleaner_id"`
	OptionID  string `json:"option_id"`
}

// AnalyzeResponse - response for frontend
type AnalyzeResponse struct {
	TotalSize  uint64 		 `json:"total_size"`
	TotalFiles uint64        `json:"total_files"`
	Items      []AnalyzeItem `json:"items"`
}

// AnalyzeItem - certain item from analyzing
type AnalyzeItem struct {
	CleanerID string 	`json:"cleaner_id"`
	OptionID  string 	`json:"option_id"`
	Size      uint64 	`json:"size"`
	FileCount uint64 	`json:"file_count"`
	Paths     []string 	`json:"paths"`
}

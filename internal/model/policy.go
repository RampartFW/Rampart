package model

type PolicySetYAML struct {
	APIVersion string          `yaml:"apiVersion" json:"apiVersion"`
	Kind       string          `yaml:"kind" json:"kind"`
	Metadata   PolicyMetadata  `yaml:"metadata" json:"metadata"`
	Defaults   *PolicyDefaults `yaml:"defaults,omitempty" json:"defaults,omitempty"`
	Includes   []IncludeRef    `yaml:"includes,omitempty" json:"includes,omitempty"`
	Policies   []PolicyYAML    `yaml:"policies" json:"policies"`
	SourceFile string          `yaml:"-" json:"-"`
}

type PolicyMetadata struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	Owner       string            `yaml:"owner,omitempty" json:"owner,omitempty"`
	Tags        map[string]string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

type PolicyDefaults struct {
	Direction Direction `yaml:"direction,omitempty" json:"direction,omitempty"`
	Action    Action    `yaml:"action,omitempty" json:"action,omitempty"`
	IPVersion IPVersion `yaml:"ipVersion,omitempty" json:"ipVersion,omitempty"`
	States    []string  `yaml:"states,omitempty" json:"states,omitempty"`
}

type IncludeRef struct {
	Path string `yaml:"path,omitempty" json:"path,omitempty"`
	URL  string `yaml:"url,omitempty" json:"url,omitempty"`
}

type PolicyYAML struct {
	Name        string     `yaml:"name" json:"name"`
	Priority    int        `yaml:"priority" json:"priority"`
	Direction   Direction  `yaml:"direction,omitempty" json:"direction,omitempty"`
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	Rules       []RuleYAML `yaml:"rules" json:"rules"`
	Line        int        `yaml:"-" json:"-"`
}

type RuleYAML struct {
	Name        string            `yaml:"name" json:"name"`
	Priority    int               `yaml:"priority,omitempty" json:"priority,omitempty"`
	Direction   string            `yaml:"direction,omitempty" json:"direction,omitempty"`
	Match       MatchYAML         `yaml:"match" json:"match"`
	Action      Action            `yaml:"action" json:"action"`
	Log         bool              `yaml:"log,omitempty" json:"log,omitempty"`
	RateLimit   *RateLimitYAML    `yaml:"rateLimit,omitempty" json:"rateLimit,omitempty"`
	Schedule    *ScheduleYAML     `yaml:"schedule,omitempty" json:"schedule,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        map[string]string `yaml:"tags,omitempty" json:"tags,omitempty"`
	Line        int               `yaml:"-" json:"-"`
}

type MatchYAML struct {
	Protocol    interface{} `yaml:"protocol,omitempty" json:"protocol,omitempty"` // string or []string
	SourceCIDRs []string    `yaml:"sourceCIDRs,omitempty" json:"sourceCIDRs,omitempty"`
	DestCIDRs   []string    `yaml:"destCIDRs,omitempty" json:"destCIDRs,omitempty"`
	SourcePorts interface{} `yaml:"sourcePorts,omitempty" json:"sourcePorts,omitempty"` // int, []int, or "start-end"
	DestPorts   interface{} `yaml:"destPorts,omitempty" json:"destPorts,omitempty"`     // int, []int, or "start-end"
	Interfaces  []string    `yaml:"interfaces,omitempty" json:"interfaces,omitempty"`
	States      []string    `yaml:"states,omitempty" json:"states,omitempty"`
	ICMPTypes   []int       `yaml:"icmpTypes,omitempty" json:"icmpTypes,omitempty"`
	Not         *MatchYAML  `yaml:"not,omitempty" json:"not,omitempty"`
}

type RateLimitYAML struct {
	Rate   int    `yaml:"rate" json:"rate"`
	Per    string `yaml:"per" json:"per"`
	Burst  int    `yaml:"burst" json:"burst"`
	Action Action `yaml:"action" json:"action"`
}

type ScheduleYAML struct {
	ActiveFrom  string             `yaml:"activeFrom,omitempty" json:"activeFrom,omitempty"`
	ActiveUntil string             `yaml:"activeUntil,omitempty" json:"activeUntil,omitempty"`
	Recurring   *RecurringSpecYAML `yaml:"recurring,omitempty" json:"recurring,omitempty"`
}

type RecurringSpecYAML struct {
	Days      []string `yaml:"days,omitempty" json:"days,omitempty"`
	StartTime string   `yaml:"startTime,omitempty" json:"startTime,omitempty"`
	EndTime   string   `yaml:"endTime,omitempty" json:"endTime,omitempty"`
	Timezone  string   `yaml:"timezone,omitempty" json:"timezone,omitempty"`
}

type VariablesYAML struct {
	APIVersion string                 `yaml:"apiVersion" json:"apiVersion"`
	Kind       string                 `yaml:"kind" json:"kind"`
	Metadata   PolicyMetadata         `yaml:"metadata" json:"metadata"`
	Variables  map[string]interface{} `yaml:"variables" json:"variables"`
}

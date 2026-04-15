package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/rampartfw/rampart/internal/model"
)

// Compile transforms a PolicySetYAML into a CompiledRuleSet.
func Compile(ps *model.PolicySetYAML, vars map[string]interface{}) (*model.CompiledRuleSet, error) {
	// 1. Resolve includes (should have been done by caller or here)
	// For now, assume ps is already resolved as per parser/validator usage.

	// 2. Apply defaults
	applyDefaults(ps)

	// 3. Compile each rule
	var compiledRules []model.CompiledRule
	sourceFiles := make(map[string]bool)
	if ps.SourceFile != "" {
		sourceFiles[ps.SourceFile] = true
	}

	for _, p := range ps.Policies {
		for _, r := range p.Rules {
			cr, err := compileRule(r, p, ps.SourceFile)
			if err != nil {
				return nil, fmt.Errorf("compile rule %s in policy %s: %w", r.Name, p.Name, err)
			}
			compiledRules = append(compiledRules, cr)
		}
	}

	// 4. Stable sort by priority
	sort.SliceStable(compiledRules, func(i, j int) bool {
		return compiledRules[i].Priority < compiledRules[j].Priority
	})

	// 5. Deterministic rule ID generation (content-based hash)
	for i := range compiledRules {
		compiledRules[i].ID = generateRuleID(compiledRules[i])
	}

	// 6. Ruleset hash
	hash, err := computeRuleSetHash(compiledRules)
	if err != nil {
		return nil, fmt.Errorf("compute ruleset hash: %w", err)
	}

	var sourceFilesList []string
	for f := range sourceFiles {
		sourceFilesList = append(sourceFilesList, f)
	}
	sort.Strings(sourceFilesList)

	return &model.CompiledRuleSet{
		Rules:       compiledRules,
		Hash:        hash,
		CompiledAt:  time.Now(),
		Metadata:    ps.Metadata,
		SourceFiles: sourceFilesList,
	}, nil
}

func applyDefaults(ps *model.PolicySetYAML) {
	if ps.Defaults == nil {
		return
	}

	for i := range ps.Policies {
		p := &ps.Policies[i]
		if p.Direction == "" {
			p.Direction = ps.Defaults.Direction
		}

		for j := range p.Rules {
			r := &p.Rules[j]
			if r.Action == "" {
				r.Action = ps.Defaults.Action
			}
			if r.Match.States == nil && ps.Defaults.States != nil {
				r.Match.States = ps.Defaults.States
			}
		}
	}
}

func compileRule(r model.RuleYAML, p model.PolicyYAML, sourceFile string) (model.CompiledRule, error) {
	priority := r.Priority
	if priority == 0 {
		priority = p.Priority
	}

	cr := model.CompiledRule{
		Name:        r.Name,
		PolicyName:  p.Name,
		Priority:    priority,
		Direction:   p.Direction,
		Action:      r.Action,
		Log:         r.Log,
		Tags:        r.Tags,
		Description: r.Description,
		SourceFile:  sourceFile,
		SourceLine:  r.Line,
	}

	if cr.Direction == "" {
		cr.Direction = model.DirectionInbound // default if not specified
	}

	match, err := compileMatch(r.Match)
	if err != nil {
		return cr, err
	}
	cr.Match = match

	if r.RateLimit != nil {
		cr.RateLimit = &model.RateLimit{
			Rate:   r.RateLimit.Rate,
			Per:    r.RateLimit.Per,
			Burst:  r.RateLimit.Burst,
			Action: r.RateLimit.Action,
		}
	}

	if r.Schedule != nil {
		cr.Schedule = compileSchedule(r.Schedule)
	}

	return cr, nil
}

func compileMatch(m model.MatchYAML) (model.CompiledMatch, error) {
	cm := model.CompiledMatch{
		Interfaces: m.Interfaces,
	}

	// Protocols
	protocols, err := parseProtocols(m.Protocol)
	if err != nil {
		return cm, err
	}
	cm.Protocols = protocols

	// CIDRs
	srcNets, err := parseCIDRs(m.SourceCIDRs)
	if err != nil {
		return cm, err
	}
	cm.SourceNets = srcNets
	for _, n := range srcNets {
		start, end := cidrToInterval(n)
		cm.SrcIntervals = append(cm.SrcIntervals, model.IPInterval{
			Start: start.Bytes(),
			End:   end.Bytes(),
		})
	}

	destNets, err := parseCIDRs(m.DestCIDRs)
	if err != nil {
		return cm, err
	}
	cm.DestNets = destNets
	for _, n := range destNets {
		start, end := cidrToInterval(n)
		cm.DstIntervals = append(cm.DstIntervals, model.IPInterval{
			Start: start.Bytes(),
			End:   end.Bytes(),
		})
	}

	// Ports
	srcPorts, err := parsePorts(m.SourcePorts)
	if err != nil {
		return cm, err
	}
	cm.SourcePorts = srcPorts

	destPorts, err := parsePorts(m.DestPorts)
	if err != nil {
		return cm, err
	}
	cm.DestPorts = destPorts

	// States
	for _, s := range m.States {
		cm.States = append(cm.States, model.ConnState(strings.ToLower(s)))
	}

	// ICMP Types
	for _, t := range m.ICMPTypes {
		cm.ICMPTypes = append(cm.ICMPTypes, uint8(t))
	}

	// Not
	if m.Not != nil {
		notMatch, err := compileMatch(*m.Not)
		if err != nil {
			return cm, err
		}
		cm.Negated = &notMatch
	}

	return cm, nil
}

func parseProtocols(p interface{}) ([]model.Protocol, error) {
	if p == nil {
		return nil, nil
	}
	switch v := p.(type) {
	case string:
		return []model.Protocol{model.Protocol(strings.ToLower(v))}, nil
	case []interface{}:
		var res []model.Protocol
		for _, item := range v {
			if s, ok := item.(string); ok {
				res = append(res, model.Protocol(strings.ToLower(s)))
			}
		}
		return res, nil
	default:
		return nil, fmt.Errorf("invalid protocol type: %T", p)
	}
}

func parseCIDRs(cidrs []string) ([]net.IPNet, error) {
	var res []net.IPNet
	for _, cidr := range cidrs {
		normalized := normalizeCIDR(cidr)
		_, ipNet, err := net.ParseCIDR(normalized)
		if err != nil {
			return nil, fmt.Errorf("parse CIDR %s: %w", normalized, err)
		}
		res = append(res, *ipNet)
	}
	return res, nil
}

func parsePorts(p interface{}) ([]model.PortRange, error) {
	if p == nil {
		return nil, nil
	}
	switch v := p.(type) {
	case int:
		return []model.PortRange{{Start: uint16(v), End: uint16(v)}}, nil
	case []interface{}:
		var res []model.PortRange
		for _, item := range v {
			pr, err := parsePorts(item)
			if err != nil {
				return nil, err
			}
			res = append(res, pr...)
		}
		return res, nil
	case string:
		if strings.Contains(v, "-") {
			var start, end uint16
			if _, err := fmt.Sscanf(v, "%d-%d", &start, &end); err != nil {
				return nil, fmt.Errorf("parse port range %s: %w", v, err)
			}
			return []model.PortRange{{Start: start, End: end}}, nil
		}
		var port uint16
		if _, err := fmt.Sscanf(v, "%d", &port); err != nil {
			return nil, fmt.Errorf("parse port %s: %w", v, err)
		}
		return []model.PortRange{{Start: port, End: port}}, nil
	case float64: // json unmarshal default for numbers
		return []model.PortRange{{Start: uint16(v), End: uint16(v)}}, nil
	default:
		return nil, fmt.Errorf("invalid port type: %T", p)
	}
}

func compileSchedule(s *model.ScheduleYAML) *model.Schedule {
	res := &model.Schedule{}
	if s.ActiveFrom != "" {
		if t, err := time.Parse(time.RFC3339, s.ActiveFrom); err == nil {
			res.ActiveFrom = &t
		}
	}
	if s.ActiveUntil != "" {
		if t, err := time.Parse(time.RFC3339, s.ActiveUntil); err == nil {
			res.ActiveUntil = &t
		}
	}
	if s.Recurring != nil {
		res.Recurring = &model.RecurringSpec{
			StartTime: s.Recurring.StartTime,
			EndTime:   s.Recurring.EndTime,
			Timezone:  s.Recurring.Timezone,
		}
		for _, d := range s.Recurring.Days {
			res.Recurring.Days = append(res.Recurring.Days, parseWeekday(d))
		}
	}
	return res
}

func parseWeekday(s string) time.Weekday {
	switch strings.ToLower(s) {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	}
	return time.Sunday
}

func generateRuleID(r model.CompiledRule) string {
	// Exclude ID from hash
	id := r.ID
	r.ID = ""
	data, _ := json.Marshal(r)
	r.ID = id
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func computeRuleSetHash(rules []model.CompiledRule) (string, error) {
	data, err := json.Marshal(rules)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

// Code generated by protoc-gen-go-tetragon. DO NOT EDIT

package yaml

import (
	bytes "bytes"
	fmt "fmt"
	eventchecker "github.com/cilium/tetragon/api/v1/tetragon/codegen/eventchecker"
	os "os"
	yaml "sigs.k8s.io/yaml"
	template "text/template"
)

// Metadata contains metadata for the eventchecker definition
type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Metadata contains metadata for the eventchecker definition
type EventCheckerConf struct {
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   Metadata              `json:"metadata"`
	Spec       MultiEventCheckerSpec `json:"spec"`
}

// ConfFromSpec creates a new EventCheckerConf from a MultiEventCheckerSpec
func ConfFromSpec(apiVersion, name, description string,
	spec *MultiEventCheckerSpec) (*EventCheckerConf, error) {
	if spec == nil {
		return nil, fmt.Errorf("spec is nil")
	}

	return &EventCheckerConf{
		APIVersion: apiVersion,
		Kind:       "EventChecker",
		Metadata: Metadata{
			Name:        name,
			Description: description,
		},
		Spec: *spec,
	}, nil
}

// ConfFromChecker creates a new EventCheckerConf from a MultiEventChecker
func ConfFromChecker(apiVersion, name, description string,
	checker eventchecker.MultiEventChecker) (*EventCheckerConf, error) {
	spec, err := SpecFromMultiEventChecker(checker)
	if err != nil {
		return nil, err
	}

	return &EventCheckerConf{
		APIVersion: apiVersion,
		Kind:       "EventChecker",
		Metadata: Metadata{
			Name:        name,
			Description: description,
		},
		Spec: *spec,
	}, nil
}

// ReadYaml reads an event checker from yaml
func ReadYaml(data string) (*EventCheckerConf, error) {
	var conf EventCheckerConf

	err := yaml.UnmarshalStrict([]byte(data), &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// ReadYamlFile reads an event checker from a yaml file
func ReadYamlFile(file string) (*EventCheckerConf, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return ReadYaml(string(data))
}

// ReadYamlTemplate reads an event checker template from yaml
func ReadYamlTemplate(text string, data interface{}) (*EventCheckerConf, error) {
	var conf EventCheckerConf

	templ := template.New("checkerYaml")
	templ, err := templ.Parse(text)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	templ.Execute(&buf, data)

	err = yaml.UnmarshalStrict(buf.Bytes(), &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// ReadYamlFileTemplate reads an event checker template from yaml
func ReadYamlFileTemplate(file string, data interface{}) (*EventCheckerConf, error) {
	text, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return ReadYamlTemplate(string(text), data)
}

// WriteYaml writes an event checker to yaml
func (conf *EventCheckerConf) WriteYaml() (string, error) {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// WriteYamlFile writes an event checker to a yaml file
func (conf *EventCheckerConf) WriteYamlFile(file string) error {
	data, err := conf.WriteYaml()
	if err != nil {
		return err
	}

	return os.WriteFile(file, []byte(data), 0o644)
}

// EventCheckerSpec is a YAML spec to define an event checker
type EventCheckerSpec struct {
	ProcessExec       *eventchecker.ProcessExecChecker       `json:"exec,omitempty"`
	ProcessExit       *eventchecker.ProcessExitChecker       `json:"exit,omitempty"`
	ProcessKprobe     *eventchecker.ProcessKprobeChecker     `json:"kprobe,omitempty"`
	ProcessTracepoint *eventchecker.ProcessTracepointChecker `json:"tracepoint,omitempty"`
	Test              *eventchecker.TestChecker              `json:"test,omitempty"`
	ProcessDns        *eventchecker.ProcessDnsChecker        `json:"dns,omitempty"`
}

// IntoEventChecker coerces an event checker from this spec
func (spec *EventCheckerSpec) IntoEventChecker() (eventchecker.EventChecker, error) {
	var eventChecker eventchecker.EventChecker
	if spec.ProcessExec != nil {
		if eventChecker != nil {
			return nil, fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessExec, eventChecker)
		}
		eventChecker = spec.ProcessExec
	}
	if spec.ProcessExit != nil {
		if eventChecker != nil {
			return nil, fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessExit, eventChecker)
		}
		eventChecker = spec.ProcessExit
	}
	if spec.ProcessKprobe != nil {
		if eventChecker != nil {
			return nil, fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessKprobe, eventChecker)
		}
		eventChecker = spec.ProcessKprobe
	}
	if spec.ProcessTracepoint != nil {
		if eventChecker != nil {
			return nil, fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessTracepoint, eventChecker)
		}
		eventChecker = spec.ProcessTracepoint
	}
	if spec.Test != nil {
		if eventChecker != nil {
			return nil, fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.Test, eventChecker)
		}
		eventChecker = spec.Test
	}
	if spec.ProcessDns != nil {
		if eventChecker != nil {
			return nil, fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessDns, eventChecker)
		}
		eventChecker = spec.ProcessDns
	}
	if eventChecker == nil {
		return nil, fmt.Errorf("EventCheckerSpec didn't define any event checker")
	}
	return eventChecker, nil
}

// SpecFromEventChecker creates a new EventCheckerSpec from an EventChecker
func SpecFromEventChecker(checker eventchecker.EventChecker) (*EventCheckerSpec, error) {
	var spec EventCheckerSpec
	switch c := checker.(type) {
	case *eventchecker.ProcessExecChecker:
		spec.ProcessExec = c
	case *eventchecker.ProcessExitChecker:
		spec.ProcessExit = c
	case *eventchecker.ProcessKprobeChecker:
		spec.ProcessKprobe = c
	case *eventchecker.ProcessTracepointChecker:
		spec.ProcessTracepoint = c
	case *eventchecker.TestChecker:
		spec.Test = c
	case *eventchecker.ProcessDnsChecker:
		spec.ProcessDns = c

	default:
		return nil, fmt.Errorf("Unhandled checker type %T", c)
	}
	return &spec, nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (spec *EventCheckerSpec) UnmarshalJSON(b []byte) error {
	type alias EventCheckerSpec
	var spec2 alias
	if err := yaml.UnmarshalStrict(b, &spec2); err != nil {
		return err
	}
	*spec = EventCheckerSpec(spec2)

	var eventChecker eventchecker.EventChecker
	if spec.ProcessExec != nil {
		if eventChecker != nil {
			return fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessExec, eventChecker)
		}
		eventChecker = spec.ProcessExec
	}
	if spec.ProcessExit != nil {
		if eventChecker != nil {
			return fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessExit, eventChecker)
		}
		eventChecker = spec.ProcessExit
	}
	if spec.ProcessKprobe != nil {
		if eventChecker != nil {
			return fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessKprobe, eventChecker)
		}
		eventChecker = spec.ProcessKprobe
	}
	if spec.ProcessTracepoint != nil {
		if eventChecker != nil {
			return fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessTracepoint, eventChecker)
		}
		eventChecker = spec.ProcessTracepoint
	}
	if spec.Test != nil {
		if eventChecker != nil {
			return fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.Test, eventChecker)
		}
		eventChecker = spec.Test
	}
	if spec.ProcessDns != nil {
		if eventChecker != nil {
			return fmt.Errorf("EventCheckerSpec cannot define more than one checker, got %T but already had %T", spec.ProcessDns, eventChecker)
		}
		eventChecker = spec.ProcessDns
	}
	if eventChecker == nil {
		return fmt.Errorf("EventCheckerSpec didn't define any event checker")
	}
	return nil
}

// MultiEventCheckerSpec is a YAML spec to define a MultiEventChecker
type MultiEventCheckerSpec struct {
	Ordered bool               `json:"ordered"`
	Checks  []EventCheckerSpec `json:"checks"`
}

// IntoMultiEventChecker coerces an event checker from this spec
func (spec *MultiEventCheckerSpec) IntoMultiEventChecker() (eventchecker.MultiEventChecker, error) {
	var checkers []eventchecker.EventChecker

	for _, check := range spec.Checks {
		checker, err := check.IntoEventChecker()
		if err != nil {
			return nil, err
		}
		checkers = append(checkers, checker)
	}

	if spec.Ordered {
		return eventchecker.NewOrderedEventChecker(checkers...), nil
	}

	return eventchecker.NewUnorderedEventChecker(checkers...), nil
}

// SpecFromMultiEventChecker coerces an event checker from this spec
func SpecFromMultiEventChecker(checker_ eventchecker.MultiEventChecker) (*MultiEventCheckerSpec, error) {
	var spec MultiEventCheckerSpec
	var specs []EventCheckerSpec

	checker, ok := checker_.(interface {
		GetChecks() []eventchecker.EventChecker
	})
	if !ok {
		return nil, fmt.Errorf("Unhandled checker type %T", checker_)
	}

	for _, check := range checker.GetChecks() {
		spec, err := SpecFromEventChecker(check)
		if err != nil {
			return nil, err
		}
		specs = append(specs, *spec)
	}

	spec.Checks = specs

	switch checker.(type) {
	case *eventchecker.OrderedEventChecker:
		spec.Ordered = true
	case *eventchecker.UnorderedEventChecker:
		spec.Ordered = false
	default:
		return nil, fmt.Errorf("Unhandled checker type %T", checker)
	}

	return &spec, nil
}

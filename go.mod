module github.com/complytime/baseline-demo

go 1.24.4

require (
	github.com/complytime/gemara2oscal v0.0.0-20250828130252-25e4d88e3ff9
	github.com/defenseunicorns/go-oscal v0.6.3
	github.com/goccy/go-yaml v1.18.0
	github.com/oscal-compass/oscal-sdk-go v0.0.5
	github.com/spf13/cobra v1.9.1
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2 // indirect
	golang.org/x/text v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ossf/gemara v0.5.0
	github.com/spf13/pflag v1.0.7 // indirect
)

// Points to experiments/oscal-transformation branch (temporarily)
replace github.com/ossf/gemara => github.com/jpower432/sci v0.0.0-20250828223210-15f2a5532714

replace github.com/complytime/gemara2oscal => github.com/complytime/gemara2oscal v0.0.0-20250828230325-db9ec8940446

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/complytime/gemara2oscal/component"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"
	"github.com/oscal-compass/oscal-sdk-go/transformers"
	"github.com/ossf/gemara/layer2"
	"github.com/ossf/gemara/layer3"
	"github.com/ossf/gemara/layer4"
	"github.com/spf13/cobra"
)

func NewPlanCommand() *cobra.Command {
	var catalogPath, targetComponent, componentType, validatorID, evaluationsPath, policyPath, guidanceRef string

	command := &cobra.Command{
		Use:   "plan",
		Short: "Transform Gemara artifacts to an OSCAL Assessment Plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			builder := component.NewDefinitionBuilder("Example", "0.1.0")

			cleanedPath := filepath.Clean(catalogPath)
			catalogData, err := os.ReadFile(cleanedPath)
			if err != nil {
				return err
			}

			var layer2Catalog layer2.Catalog
			err = yaml.Unmarshal(catalogData, &layer2Catalog)
			if err != nil {
				return err
			}

			var allPlans []layer4.AssessmentPlan
			err = filepath.Walk(evaluationsPath, func(path string, info os.FileInfo, err error) error {

				if info.IsDir() {
					return nil
				}

				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				var assessmentPlans []layer4.AssessmentPlan
				err = yaml.Unmarshal(content, &assessmentPlans)
				if err != nil {
					return err
				}

				allPlans = append(allPlans, assessmentPlans...)
				return nil
			})

			builder = builder.AddTargetComponent(targetComponent, componentType, layer2Catalog)
			builder = builder.AddValidationComponent(validatorID, allPlans)

			cleanedPath = filepath.Clean(policyPath)
			policyData, err := os.ReadFile(cleanedPath)
			if err != nil {
				return err
			}

			var layer3Policy layer3.PolicyDocument
			err = yaml.Unmarshal(policyData, &layer3Policy)
			if err != nil {
				return err
			}

			for _, ref := range layer3Policy.ControlReferences {
				builder = builder.AddParameterModifiers(ref.ReferenceId, ref.ParameterModifications)
			}
			compDef := builder.Build()

			var found bool
			for _, guidance := range layer3Policy.GuidanceReferences {
				if guidanceRef == guidance.ReferenceId {
					ap, err := transformers.ComponentDefinitionsToAssessmentPlan(cmd.Context(), []oscalTypes.ComponentDefinition{compDef}, guidance.ReferenceId)
					if err != nil {
						return err
					}
					oscalModels := oscalTypes.OscalModels{
						AssessmentPlan: ap,
					}
					compDefData, err := json.MarshalIndent(oscalModels, "", " ")
					if err != nil {
						return err
					}
					_, _ = fmt.Fprintln(os.Stdout, string(compDefData))
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("guidance reference %q does not exist in policy", guidanceRef)
			}

			return nil
		},
	}

	flags := command.Flags()
	flags.StringVarP(&catalogPath, "catalog-path", "c", "./src/catalogs/osps.yml", "Path to L2 catalog to transform")
	flags.StringVarP(&evaluationsPath, "evaluations-path", "e", "./src/plans", "Path to Layer 4 evaluation plans")
	flags.StringVarP(&targetComponent, "target-component", "t", "", "Title for target component for evaluation")
	flags.StringVar(&componentType, "component-type", "software", "Component type (based on valid OSCAL component types)")
	flags.StringVarP(&validatorID, "validator-id", "v", "", "Validation plugin id")
	flags.StringVarP(&policyPath, "policy-path", "p", "./src/policy.yaml", "Path to Layer 3 policy")
	flags.StringVarP(&guidanceRef, "guidance-reference", "r", "", "Guidance reference to tailor the plan to")
	return command
}

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/complytime/gemara2oscal/controls"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"
	"github.com/ossf/gemara/layer1"
	"github.com/spf13/cobra"
)

func NewCatalogCommand() *cobra.Command {
	var guidancePath string

	command := &cobra.Command{
		Use:   "catalog",
		Short: "Transform Gemara artifacts to OSCAL Catalog",
		RunE: func(cmd *cobra.Command, args []string) error {

			cleanedPath := filepath.Clean(guidancePath)
			guidanceData, err := os.ReadFile(cleanedPath)
			if err != nil {
				return err
			}

			var layer1Doc layer1.GuidanceDocument
			err = yaml.Unmarshal(guidanceData, &layer1Doc)
			if err != nil {
				return err
			}

			catalog, err := controls.ToOSCALCatalog(layer1Doc)
			if err != nil {
				return err
			}

			oscalModels := oscalTypes.OscalModels{
				Catalog: &catalog,
			}
			compDefData, err := json.MarshalIndent(oscalModels, "", " ")
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintln(os.Stdout, string(compDefData))
			return nil
		},
	}

	flags := command.Flags()
	flags.StringVarP(&guidancePath, "guidance-path", "g", "", "Path to L1 guidance doc to transform")
	return command
}

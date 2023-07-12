package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bxxf/tgen/internal/csv"
	"github.com/bxxf/tgen/internal/generator"
	"github.com/bxxf/tgen/internal/utils"
	"github.com/spf13/cobra"
)

var (
	enAll    bool
	locFile  string
	output   string
	params   []string
	filename string
	format   string
	rootCmd  = &cobra.Command{
		Use:     "tgen loc=[languages] [attributes] [flags]",
		Aliases: []string{"tgen", "tg"},
		Short:   "Generate target records",
		Long: `Examples:
  tgen loc=en,es,de --format=countryOnly
  tgen loc=avast --en-all
  tg loc=avast segment=SKU1,SKU2 activationKey=xxx
  tg --loc-file=loc.csv
`,
		RunE: generateRecords,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&enAll, "en-all", false, "Generate for all English locales (US, CA, AU, GB)")
	rootCmd.Flags().StringVar(&locFile, "loc-file", "", "Path to CSV file containing translated data to import languages")
	rootCmd.Flags().StringVar(&output, "output", "", "Output file or folder path (default: tgen-{timestamp}.csv)")
	rootCmd.Flags().StringSliceVar(&params, "params", []string{}, "Dynamic parameters in key=value format")
	rootCmd.Flags().StringVar(&format, "format", "", "Format to generate (countryiso, default)")

	cobra.OnInitialize(initAutoComplete)
	rootCmd.PersistentFlags().SetAnnotation("loc-file", cobra.BashCompFilenameExt, []string{".csv", ".txt"})

	rootCmd.SetUsageTemplate(fmt.Sprintf("%s\n\n%s", rootCmd.UsageTemplate(), rootCmd.Long))
}

func initAutoComplete() {
	// Add autocomplete function for loc-file flag
	rootCmd.RegisterFlagCompletionFunc("loc-file", autoCompleteLocFile)
}

func autoCompleteLocFile(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	matches, err := filepath.Glob(toComplete + "*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return nil, cobra.ShellCompDirectiveError
	}

	return matches, cobra.ShellCompDirectiveNoFileComp
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func generateRecords(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && locFile == "" {
		return fmt.Errorf("missing required argument 'loc' or 'loc-file'")
	}

	parameters := utils.ParseParams(args)
	languages := parameters["loc"]

	if languages == nil && locFile != "" {
		var err error
		languages, err = utils.GetLanguagesFromLocFile(locFile)
		if err != nil {
			return fmt.Errorf("error reading loc file: %w", err)
		}

		if len(languages) == 0 {
			return fmt.Errorf("no languages found in loc file")
		}
	}

	if len(languages) == 0 {
		return fmt.Errorf("missing list of locales in required argument 'loc'")
	}

	delete(parameters, "loc")

	enAllFlag, _ := cmd.Flags().GetBool("en-all")
	locFileFlag, _ := cmd.Flags().GetString("loc-file")
	outputFlag, _ := cmd.Flags().GetString("output")
	paramsFlag, _ := cmd.Flags().GetStringSlice("params")
	formatFlag, _ := cmd.Flags().GetString("format")

	flags := map[string]string{
		"en-all":   fmt.Sprintf("%t", enAllFlag),
		"loc-file": locFileFlag,
		"output":   outputFlag,
		"params":   strings.Join(paramsFlag, ","),
		"format":   formatFlag,
	}

	records, err := generator.Generate(languages, flags, parameters)
	if err != nil {
		return fmt.Errorf("error generating records: %w", err)
	}

	filename = utils.GenerateFileName(output)
	err = csv.WriteToCsv(records, filename)
	if err != nil {
		return fmt.Errorf("error writing to CSV: %w", err)
	}

	return nil
}

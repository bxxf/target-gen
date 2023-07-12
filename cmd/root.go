package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bxxf/tgen/internal/csv"
	"github.com/bxxf/tgen/internal/generator"
	"github.com/bxxf/tgen/internal/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	enAll      bool
	locFile    string
	output     string
	params     map[string][]string
	format     string
	configFile string
	languages  []string

	rootCmd = &cobra.Command{
		Use:   "tgen loc=[languages] [attributes] [flags]",
		Short: "Generate target records",
		Long: `Examples:
  tgen loc=avast --en-all
  tgen loc=en,es,de --format=countryiso
  tgen loc=avast segment=SKU1,SKU2 activationKey=xxx
  tgen --loc-file=loc.csv
`,
		RunE: generateRecords,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&enAll, "en-all", false, "Generate for all English locales (US, CA, AU, GB)")
	rootCmd.Flags().StringVar(&locFile, "loc-file", "", "Path to CSV file containing translated data to import languages")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output file or folder path (default: tgen-{timestamp}.csv)")
	rootCmd.Flags().StringVar(&format, "format", "", "Format to generate (countryiso, default)")
	rootCmd.Flags().StringVar(&configFile, "config", "", "Path to configuration file")

	cobra.OnInitialize(initAutoComplete)
	rootCmd.PersistentFlags().SetAnnotation("loc-file", cobra.BashCompFilenameExt, []string{".csv", ".txt"})
	rootCmd.PersistentFlags().SetAnnotation("output", cobra.BashCompFilenameExt, []string{".csv", ".txt"})
	rootCmd.PersistentFlags().SetAnnotation("config", cobra.BashCompFilenameExt, []string{".yaml", ".yml"})

	rootCmd.SetUsageTemplate(fmt.Sprintf("%s\n\n%s", rootCmd.UsageTemplate(), rootCmd.Long))
}

func initAutoComplete() {
	// Add autocomplete function for loc-file flag
	rootCmd.RegisterFlagCompletionFunc("loc-file", autoCompleteLocFile)
	rootCmd.RegisterFlagCompletionFunc("output", autoCompleteLocFile)
	rootCmd.RegisterFlagCompletionFunc("config", autoCompleteLocFile)
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
	if len(configFile) > 0 {
		err := initConfig(configFile)
		if err != nil {
			return err
		}
	} else {
		parameters := utils.ParseParams(args)
		languages = parameters["loc"]

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
		formatFlag, _ := cmd.Flags().GetString("format")

		params = parameters
		enAll = enAllFlag
		locFile = locFileFlag
		output = outputFlag
		format = formatFlag
	}

	flags := map[string]string{
		"en-all":   fmt.Sprintf("%t", enAll),
		"loc-file": locFile,
		"output":   output,
		"format":   format,
	}

	if locFile != "" {
		var err error
		languages, err = utils.GetLanguagesFromLocFile(locFile)
		if err != nil {
			log.Printf("error reading loc file: %v", err)
			return fmt.Errorf("error reading loc file: %w", err)
		}
	}

	if len(languages) < 1 {
		log.Printf("no languages found in (loc_file=%s) - 'languages' or 'loc_file' parameter", locFile)
		return fmt.Errorf("no languages found in 'languages' parameter or 'loc_file'")
	}

	records, err := generator.Generate(languages, flags, params)
	if err != nil {
		log.Printf("error generating records: %v", err)
		return fmt.Errorf("error generating records: %w", err)
	}

	filename := utils.GenerateFileName(output)
	err = csv.WriteToCsv(records, filename)
	if err != nil {
		return fmt.Errorf("error writing to CSV: %w", err)
	}

	return nil
}

func initConfig(configFile string) error {
	if configFile == "" {
		return nil
	}

	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	var config struct {
		LocFile   string   `yaml:"loc_file"`
		Output    string   `yaml:"output"`
		Format    string   `yaml:"format"`
		EnAll     bool     `yaml:"en_all"`
		Params    []string `yaml:"params"`
		Languages string   `yaml:"languages"`
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return err
	}

	locFile = config.LocFile
	output = config.Output
	format = config.Format
	enAll = config.EnAll

	var paramsMap = make(map[string][]string)
	for _, param := range config.Params {
		parts := strings.Split(param, "=")
		if len(parts) != 2 {
			return fmt.Errorf("invalid parameter: %s", param)
		}

		paramsMap[parts[0]] = strings.Split(parts[1], ",")
	}

	params = paramsMap

	languages = strings.Split(strings.Replace(config.Languages, " ", "", -1), ",")

	return nil
}

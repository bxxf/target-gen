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

type Config struct {
	EnAll      bool
	LocFile    string
	Output     string
	Params     map[string][]string
	Format     string
	ConfigFile string
	Languages  []string
}

var (
	cfg = &Config{
		Params: make(map[string][]string),
	}
	rootCmd = &cobra.Command{
		Use:   "tgen loc=[languages] [attributes] [flags]",
		Short: "Generate target records",
		Long: `Examples:
  tgen loc=BRANDNAME --en-all
  tgen loc=en,es,de --format=countryiso
  tgen loc=BRANDNAME segment=SKU1,SKU2 activationKey=xxx
  tgen --loc-file=loc.csv
`,
		RunE: generateRecords,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&cfg.EnAll, "en-all", false, "Generate for all English locales (US, CA, AU, GB)")
	rootCmd.Flags().StringVar(&cfg.LocFile, "loc-file", "", "Path to CSV file containing translated data to import languages")
	rootCmd.Flags().StringVarP(&cfg.Output, "output", "o", "", "Output file or folder path (default: tgen-{timestamp}.csv)")
	rootCmd.Flags().StringVar(&cfg.Format, "format", "", "Format to generate (countryiso, default)")
	rootCmd.Flags().StringVar(&cfg.ConfigFile, "config", "", "Path to configuration file")

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
	if len(cfg.ConfigFile) > 0 {
		err := initConfig(cfg.ConfigFile)
		if err != nil {
			return logError(err)
		}
	} else {
		parameters := utils.ParseParams(args)
		cfg.Languages = parameters["loc"]

		if cfg.Languages == nil && cfg.LocFile != "" {
			var err error
			cfg.Languages, err = utils.GetLanguagesFromLocFile(cfg.LocFile)
			if err != nil {
				return logError(fmt.Errorf("error reading loc file: %w", err))
			}

			if len(cfg.Languages) == 0 {
				return logError(fmt.Errorf("no languages found in loc file"))
			}
		}

		if len(cfg.Languages) == 0 {
			return logError(fmt.Errorf("missing list of locales in required argument 'loc'"))
		}

		delete(parameters, "loc")

		cfg.Params = parameters
	}

	flags := map[string]string{
		"en-all":   fmt.Sprintf("%t", cfg.EnAll),
		"loc-file": cfg.LocFile,
		"output":   cfg.Output,
		"format":   cfg.Format,
	}

	if cfg.LocFile != "" {
		var err error
		cfg.Languages, err = utils.GetLanguagesFromLocFile(cfg.LocFile)
		if err != nil {
			log.Printf("error reading loc file: %v", err)
			return logError(fmt.Errorf("error reading loc file: %w", err))
		}
	}

	if len(cfg.Languages) < 1 {
		log.Printf("no languages found in (loc_file=%s) - 'languages' or 'loc_file' parameter", cfg.LocFile)
		return logError(fmt.Errorf("no languages found in 'languages' parameter or 'loc_file'"))
	}

	records, err := generator.Generate(cfg.Languages, flags, cfg.Params)
	if err != nil {
		log.Printf("error generating records: %v", err)
		return logError(fmt.Errorf("error generating records: %w", err))
	}

	filename := utils.GenerateFileName(cfg.Output)
	err = csv.WriteToCsv(records, filename)
	if err != nil {
		return logError(fmt.Errorf("error writing to CSV: %w", err))
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

	cfg.LocFile = config.LocFile
	cfg.Output = config.Output
	cfg.Format = config.Format
	cfg.EnAll = config.EnAll

	for _, param := range config.Params {
		parts := strings.Split(param, "=")
		if len(parts) != 2 {
			return fmt.Errorf("invalid parameter: %s", param)
		}

		cfg.Params[parts[0]] = strings.Split(parts[1], ",")
	}

	cfg.Languages = strings.Split(strings.Replace(config.Languages, " ", "", -1), ",")

	if len(cfg.Languages) == 0 && cfg.LocFile == "" {
		return fmt.Errorf("missing required parameters in config file (loc_file or languages)")
	}

	return nil
}

func logError(err error) error {
	log.Println(err)
	return err
}

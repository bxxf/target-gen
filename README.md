# tgen

A command-line tool for generating CSV target records based on specified parameters and configurations.

## Installation

1. Tap into the Homebrew formula repository:
```sh
brew tap bxxf/tap
```

2. Install the CLI tool:
```sh
brew install tgen
```

### Examples

Examples can be found in `examples` folder.

**Generate target records for all English locales (US, CA, AU, GB):**
```sh
 tgen loc=en --en-all 
```

**Generate target records for brand:**
```sh
 tgen loc=BRANDNAME
```

**Generate target records for brand with specific attributes:**
```sh
 tgen loc=BRANDNAME attribute=VALUE1,VALUE2
```

**Generate target records based on LOC file:**
```sh
 tgen --loc-file=PATH_TO_LOCFILE.csv
```

### Using config file

The CLI tool can be configured using a YAML configuration file. The following parameters can be specified:

- `loc_file`: Path to the CSV file containing translated data from LOC to import languages.
- `languages`: List of languages separated by commas (only use if not using loc_file).
- `output`: Output file or folder path.
- `format`: Format to generate (e.g., countryiso, default).
- `en_all`: Generate target records for all English locales (US, CA, AU, GB).
- `params`: Additional parameters in the format `param_name=value1,value2,value3`.

#### Example configuration file:

```yaml
loc_file: loc.csv
output: output.csv
params:
- segment=SKU1,SKU2
- activationKey=xxx
```

## License
This project is licensed under the Copyright License.

&copy; Copyright 2023 Filip Brebera
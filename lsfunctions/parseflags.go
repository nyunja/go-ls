package lsfunctions

import "fmt"

// Flag struct to store parsed flag and its value
type Flags struct {
	Long      bool
	All       bool
	Recursive bool
	Reverse   bool
	Time      bool
}

// parseFlags parses command-line arguments to extract flags and non-flag arguments.
// It supports both long format (e.g., "--long") and short format (e.g., "-l") flags.
//
// Parameters:
//   - args: A slice of strings representing the command-line arguments to be parsed.
//
// Returns:
//   - flags: A Flags struct containing boolean values for each recognized flag.
//   - parsedArgs: A slice of strings containing the non-flag arguments.
func ParseFlags(args []string) (flags Flags, parsedArgs []string, err error) {
	for _, arg := range args {
		if arg[0] == '-' {
			if len(arg) > 1 {
				switch arg {
				case "--reverse":
					flags.Reverse = true
				case "--all":
					flags.All = true
				case "--recursive":
					flags.Recursive = true
				default:
					for _, flag := range arg[1:] {
						switch flag {
						case 'l':
							flags.Long = true
						case 'a':
							flags.All = true
						case 'R':
							flags.Recursive = true
						case 't':
							flags.Time = true
						case 'r':
							flags.Reverse = true
						default:
							return Flags{}, nil, fmt.Errorf("invalid flag: %s", arg)
						}
					}
				}
			} 
		} else {
			parsedArgs = append(parsedArgs, arg)
		}
	}
	return flags, parsedArgs, nil
}

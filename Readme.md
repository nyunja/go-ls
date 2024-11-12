# my-ls

my-ls is a custom implementation of the Unix ls command in Go, developed to explore system interactions and data handling. This project replicates key features of the original ls command, while implementing some unique functionalities and following Go best practices.
Features

The following flags are supported:
```sh
    -l: Detailed list format, showing file permissions, ownership, size, and modification date, similar to the output of ls -l.

    -R: Recursively list files in all directories and subdirectories.

    -a: Display all files, including hidden ones (files starting with .).

    -r: Reverse the sorting order.

    -t: Sort files by modification time, newest first.
  ```  

You can combine flags in various ways, just as with the standard ls.
## Usage

You can use my-ls with any combination of supported flags:
```sh
my-ls [flags] [directory]
```
```sh
Examples:

    my-ls                # Lists files and directories in the current directory.
    my-ls -l             # Displays detailed information about each file in the current directory.
    my-ls -a -R /path/to/dir  # Recursively lists all files (including hidden) in /path/to/dir.
    my-ls -t -r          # Lists files in reverse chronological order of modification.
  ```  

## Implementation Notes
```sh
    Recursive Flag (-R): Implementing this requires careful handling of nested directories. Plan how recursive directory traversal interacts with other flags.

    Sorting and Filtering: Sorting by time (-t) and reversing order (-r) should handle multiple flags simultaneously.

    ls -l: Match the output format exactly to that of the system ls -l command, ensuring file permissions, ownership, and other metadata are displayed correctly.
```

## Installation

To clone the project, use:
```sh
git clone https://learn.zone01kisumu.ke/git/johnotieno0/my-ls-1.git
```
```sh
cd my-ls-1
```

Run

After navigating to the project directory, build the project:
```sh
go run . [flags] [directory]
```

## Contributors
[Teddy siaka](https://learn.zone01kisumu.ke/git/tesiaka)

[John Otieno](https://learn.zone01kisumu.ke/git/johnotieno0)

License

This project is licensed under the MIT License.
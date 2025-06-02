# kerb Command-Line Tool

## Overview

`kerb` is a command-line tool for managing files with a special Kerb header. It provides utilities to check, insert, remove, and synchronize files containing the Kerb header, as well as to perform string replacements in such files. The tool is useful for automating file management tasks where a specific header is used to mark files managed by automation.

## Commands

### is-kerb-file <file>
Checks if the specified file contains the Kerb header.

**Usage:**
```
kerb is-kerb-file <file>
```
- Prints `true` if the file contains the header, `false` otherwise.

### sync [--remove-header] <srcDir>
Synchronizes all files containing the Kerb header from `<srcDir>` to the current directory. Optionally removes the Kerb header from copied files.

**Usage:**
```
kerb sync [--remove-header] <srcDir>
```
- Copies all files with the Kerb header from `<srcDir>` to the current directory, preserving relative paths.
- If `--remove-header` is specified, removes the Kerb header from each copied file.

### insert-header <file>
Inserts the Kerb header at the beginning of the specified file, if not already present.

**Usage:**
```
kerb insert-header <file>
```

### add-header [<file>]
Adds the Kerb header to the beginning of the specified file, or to all files under the current directory if no file is specified. The header is added even if it already exists.

**Usage:**
```
kerb add-header <file>
kerb add-header
```

### replace <file> <old> <new>
Replaces all occurrences of `<old>` with `<new>` in the specified file, but only if the file contains the Kerb header.

**Usage:**
```
kerb replace <file> <old> <new>
```

### replace-all <old> <new>
Replaces all occurrences of `<old>` with `<new>` in all files under the current directory that contain the Kerb header.

**Usage:**
```
kerb replace-all <old> <new>
```

### list
Lists all files under the current directory that contain the Kerb header.

**Usage:**
```
kerb list
```
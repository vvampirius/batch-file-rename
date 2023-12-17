# batch-file-rename

It's just a simple command line tool for batch file renaming.

It uses:
- REGEXP to extract the desired name from file name
- Golang Template to compose target file name

You can use the following variables in template:
- .FileName - source fil name
- .Name - REGEXP extracted name from file name
- .Index - input file index (integer; from 1)
- .ModTime - file modification time (time.Time)
- .Size - file size (integer)

## Examples
```bash
% batch-file-rename -template '{{.Index}} - {{.FileName}}' First.txt Second.txt Third.txt
'First.txt' -> '1 - First.txt'
'Second.txt' -> '2 - Second.txt'
'Third.txt' -> '3 - Third.txt'

```

```bash
% batch-file-rename -name "^.*-\s+(.*)\s+File.*" -template '{{printf "%02d" .Index}} - {{.Name}}.txt' "123 - First File.txt" "562 - Second File.txt"
'123 - First File.txt' -> '01 - First.txt'
'562 - Second File.txt' -> '02 - Second.txt'
```

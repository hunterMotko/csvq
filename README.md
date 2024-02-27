# CSVQ

Just needed to make a tool for dealing with csv files. 
One of my favorite projects I made for csv, sqlite, and ruby.
Might as well make it faster with go.

## Usage thus far

### Get/Check Headers to Stdout to

```bash
cat some.csv | ./csvq -hd 
```
```bash
cat some.csv | ./csvq -hd > some.csv
```

### Select Specific Columns 

```bash
cat some.csv | ./csvq -c name
```
```bash
cat some.csv | ./csvq -c name date foo bar
```

### Specify By -f for File

```bash
./csvq -f some.csv -c col1 col2 col3 ect...
```

### Slice Csv Sections

```bash
cat some.csv | ./csvq -s 1-3
```
```bash
./csvq -f some.csv -s 1-3
```
-- Start of the columns to column [N] 
```bash
./csvq -f some.csv -s -5
```
-- From column [N] to the end of columns
```bash
./csvq -f some.csv -s 3-
```

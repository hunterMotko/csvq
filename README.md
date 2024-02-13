# CSVQ

Just needed to make a tool for dealing with csv files. 
One of my favorite projects I made for csv, sqlite, and ruby.
Might as well make it faster with go.

## Usage thus far

cat some.csv | ./csvq -hd 
    - stdout of headers
cat some.csv | ./csvq -c name
    - stdout of column/s
cat some.csv | ./csvq -c name date foo bar

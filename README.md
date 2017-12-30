fslice
======

A small utility for extracting delimited sections of a file.


## Usage

Provide `fslice` with a starting and ending delimiter and the section in between
will be printed to STDOUT:

```sh
$ cat myfile.txt
Hello
#BEGIN
foo
bar
#END
Goodbye

$ fslice -start #BEGIN -end #END myfile.txt
foo
bar
```

You can also optionally provide a header:

```sh
$ fslice -start #BEGIN -end #END -header "// FILE: $FILENAME" myfile.txt
// FILE: myfile.txt
foo
bar
```
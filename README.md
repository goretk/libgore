[![Build Status](https://travis-ci.org/goretk/libgore.svg?branch=master)](https://travis-ci.org/goretk/libgore)[![Release](https://img.shields.io/github/release/goretk/libgore.svg?style=flat-square)](https://github.com/goretk/libgore/releases/latest)
# Libgore - Open up GoRE to other languages

*Libgore* is a dynamic C-library for interacting with
[GoRE](https://github.com/goretk/gore). It is using **cgo** to produce a
translation layer between the code written in Go and the exported C functions.
With this library, it is possible to write bindings for other languages that
have C foreign function interface (FFI) support.
[PyGoRE](https://github.com/goretk/pygore) uses this dynamic library to provide
a Python library that can be used to write tools in Python.

## Example

```c
#include <stdio.h>
#include "libgore.h"

int main()
{
   int v;
   char* fp = "/path/to/test/file";
   v = gore_open(fp);
   struct compilerVersion* c = gore_getCompilerVersion(fp);
   printf("Compiler version: %s\n", c->name);
   gore_close(fp);
   return 0;
}
```

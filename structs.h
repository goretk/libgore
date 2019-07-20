// Copyright 2019 The GoRE.tk Authors. All rights reserved.
// Use of this source code is governed by the license that
// can be found in the LICENSE file.

struct compilerVersion{
        char* name;
        char* sha;
        char* timestamp;
};

struct function {
        char* name;
        int srcLineLength;
        int srcLineStart;
        int srcLineEnd;
        unsigned long long offset;
        unsigned long long end;
        char* fileName;
        char* packageName;
};

struct method {
        char* receiver;
        struct function* function;
};

struct package {
        char* name;
        char* filepath;
        struct function** function;
        struct method** method;
        unsigned long numFuncs;
        unsigned long numMeths;
};

struct packages {
        struct package** packages;
        unsigned long length;
};

struct method_type {
    char* name;
    struct type* gotype;
    unsigned long long ifaceCallOffset;
    unsigned long long funcCallOffset;
};

struct methods_type {
    struct method_type** methods;
    unsigned long length;
};

struct type {
    unsigned int kind;
    char* name;
    unsigned long long addr;
    unsigned long long ptrResolved;
    char* packagePath;
    struct types* fields;
    char* fieldName;
    char* fieldTag;
    int fieldAnon;
    struct type* element;
    int length;
    int chanDir;
    struct type* key;
    struct types* funcArgs;
    struct types* funcReturns;
    int isVariadic;
    struct methods_type* methods;
};

struct types {
    struct type** types;
    unsigned long length;
};

// This file is part of libgore.
//
// Copyright (C) 2019-2021 GoRE Authors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

struct compilerVersion
{
        char *name;
        char *sha;
        char *timestamp;
};

struct function
{
        char *name;
        unsigned long long offset;
        unsigned long long end;
        char *packageName;
};

struct method
{
        char *receiver;
        struct function *function;
};

struct package
{
        char *name;
        char *filepath;
        struct function **function;
        struct method **method;
        unsigned long numFuncs;
        unsigned long numMeths;
};

struct packages
{
        struct package **packages;
        unsigned long length;
};

struct method_type
{
        char *name;
        struct type *gotype;
        unsigned long long ifaceCallOffset;
        unsigned long long funcCallOffset;
};

struct methods_type
{
        struct method_type **methods;
        unsigned long length;
};

struct type
{
        unsigned int kind;
        char *name;
        unsigned long long addr;
        unsigned long long ptrResolved;
        char *packagePath;
        struct types *fields;
        char *fieldName;
        char *fieldTag;
        int fieldAnon;
        struct type *element;
        int length;
        int chanDir;
        struct type *key;
        struct types *funcArgs;
        struct types *funcReturns;
        int isVariadic;
        struct methods_type *methods;
};

struct types
{
        struct type **types;
        unsigned long length;
};

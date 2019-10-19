// Copyright 2019 The GoRE.tk Authors. All rights reserved.
// Use of this source code is governed by the license that
// can be found in the LICENSE file.

package main

/*
#include <stdlib.h>
#include "structs.h"
*/
import "C"

import (
	"unsafe"

	"github.com/goretk/gore"
)

//export gore_open
func gore_open(filePath *C.char) C.int {
	fp := C.GoString(filePath)
	f, err := gore.Open(fp)
	if err != nil {
		return C.int(0)
	}
	a := new(arena)
	addNewArena(fp, a)
	addNewFile(fp, f)
	return C.int(1)
}

//export gore_close
func gore_close(filePath *C.char) {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return
	}
	f.Close()
	removeFile(fp)

	a := getArena(fp)
	a.free()
	removeArena(fp)
}

//export gore_build_id
func gore_build_id(filePath *C.char) *C.char {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return nil
	}
	arena := getArena(fp)
	if arena == nil {
		return nil
	}
	id := f.BuildID
	return arena.cstring(id)
}

//export gore_setGoVersion
func gore_setGoVersion(filePath *C.char, version *C.char) C.int {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return 0
	}
	ver := C.GoString(version)
	err := f.SetGoVersion(ver)
	if err != nil {
		return 0
	}
	return 1
}

//export gore_getCompilerVersion
func gore_getCompilerVersion(filePath *C.char) *C.struct_compilerVersion {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return nil
	}
	arena := getArena(fp)
	if arena == nil {
		return nil
	}
	cv, err := f.GetCompilerVersion()
	if err != nil {
		return nil
	}
	cs := (*C.struct_compilerVersion)(arena.malloc(C.sizeof_struct_compilerVersion))
	cs.name = arena.cstring(cv.Name)
	cs.sha = arena.cstring(cv.SHA)
	cs.timestamp = arena.cstring(cv.Timestamp)
	return cs
}

//export gore_getPackages
func gore_getPackages(filePath *C.char) *C.struct_packages {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return nil
	}
	a := getArena(fp)
	if a == nil {
		return nil
	}
	pkgs, err := f.GetPackages()
	if err != nil {
		return nil
	}
	return convertPackages(pkgs, a)
}

//export gore_getVendors
func gore_getVendors(filePath *C.char) *C.struct_packages {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return nil
	}
	a := getArena(fp)
	if a == nil {
		return nil
	}
	pkgs, err := f.GetVendors()
	if err != nil {
		return nil
	}
	return convertPackages(pkgs, a)
}

//export gore_getSTDLib
func gore_getSTDLib(filePath *C.char) *C.struct_packages {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return nil
	}
	a := getArena(fp)
	if a == nil {
		return nil
	}
	pkgs, err := f.GetSTDLib()
	if err != nil {
		return nil
	}
	return convertPackages(pkgs, a)
}

//export gore_getUnknown
func gore_getUnknown(filePath *C.char) *C.struct_packages {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return nil
	}
	a := getArena(fp)
	if a == nil {
		return nil
	}
	pkgs, err := f.GetUnknown()
	if err != nil {
		return nil
	}
	return convertPackages(pkgs, a)
}

//export gore_getTypes
func gore_getTypes(filePath *C.char) *C.struct_types {
	fp := C.GoString(filePath)
	f := getFile(fp)
	if f == nil {
		return nil
	}
	a := getArena(fp)
	if a == nil {
		return nil
	}
	types, err := f.GetTypes()
	if err != nil {
		return nil
	}
	parsed := make(map[uint64]*C.struct_type)
	return convertTypes(types, a, parsed)
}

func convertTypes(types []*gore.GoType, arena *arena, parsed map[uint64]*C.struct_type) *C.struct_types {
	val := (**C.struct_type)(arena.calloc(C.sizeof_struct_type, len(types)))
	pval := (*[1 << 30]*C.struct_type)(unsafe.Pointer(val))[:len(types):len(types)]

	for i, t := range types {
		ct := convertType(t, arena, parsed)
		pval[i] = ct
	}
	pr := (*C.struct_types)(arena.malloc(C.sizeof_struct_types))
	pr.types = val
	pr.length = C.ulong(len(types))
	return pr
}

func convertType(t *gore.GoType, arena *arena, parsed map[uint64]*C.struct_type) *C.struct_type {
	if t == nil {
		return nil
	}
	if ct, ok := parsed[t.Addr]; ok {
		return ct
	}
	ct := (*C.struct_type)(arena.malloc(C.sizeof_struct_type))
	parsed[t.Addr] = ct
	ct.kind = C.uint(t.Kind)
	ct.name = arena.cstring(t.Name)
	ct.addr = C.ulonglong(t.Addr)
	ct.ptrResolved = C.ulonglong(t.PtrResolvAddr)
	ct.packagePath = arena.cstring(t.PackagePath)
	ct.fields = convertTypes(t.Fields, arena, parsed)
	ct.fieldName = arena.cstring(t.FieldName)
	ct.fieldTag = arena.cstring(t.FieldTag)
	if t.FieldAnon {
		ct.fieldAnon = C.int(1)
	} else {
		ct.fieldAnon = C.int(0)
	}
	ct.element = convertType(t.Element, arena, parsed)
	ct.length = C.int(t.Length)
	ct.chanDir = C.int(t.ChanDir)
	ct.key = convertType(t.Key, arena, parsed)
	ct.funcArgs = convertTypes(t.FuncArgs, arena, parsed)
	ct.funcReturns = convertTypes(t.FuncReturnVals, arena, parsed)
	if t.IsVariadic {
		ct.isVariadic = C.int(1)
	} else {
		ct.isVariadic = C.int(0)
	}
	methods := (**C.struct_method_type)(arena.calloc(C.sizeof_struct_method_type, len(t.Methods)))
	pmethods := (*[1 << 30]*C.struct_method_type)(unsafe.Pointer(methods))[:len(t.Methods):len(t.Methods)]
	for i, m := range t.Methods {
		meth := (*C.struct_method_type)(arena.malloc(C.sizeof_struct_method_type))
		meth.name = arena.cstring(m.Name)
		meth.gotype = convertType(m.Type, arena, parsed)
		meth.ifaceCallOffset = C.ulonglong(m.IfaceCallOffset)
		meth.funcCallOffset = C.ulonglong(m.FuncCallOffset)
		pmethods[i] = meth
	}
	ms := (*C.struct_methods_type)(arena.malloc(C.sizeof_struct_methods_type))
	ms.methods = methods
	ms.length = C.ulong(len(t.Methods))
	ct.methods = ms
	return ct
}

func convertPackages(pkgs []*gore.Package, arena *arena) *C.struct_packages {
	// https://stackoverflow.com/a/42842309
	// https://groups.google.com/forum/#!topic/golang-nuts/sV_f0VkjZTA
	pp := (**C.struct_package)(arena.malloc(C.size_t(len(pkgs)) * C.sizeof_struct_package))
	ppakgs := (*[1 << 30]*C.struct_package)(unsafe.Pointer(pp))[:len(pkgs):len(pkgs)]

	for i, p := range pkgs {
		cp := (*C.struct_package)(arena.malloc(C.sizeof_struct_package))
		cp.name = arena.cstring(p.Name)
		cp.filepath = arena.cstring(p.Filepath)

		// Populate funcs
		pf := (**C.struct_function)(arena.malloc(C.size_t(len(p.Functions)) * C.sizeof_struct_function))
		af := (*[1 << 30]*C.struct_function)(unsafe.Pointer(pf))[:len(p.Functions):len(p.Functions)]
		for j, f := range p.Functions {
			cf := convertFunction(f, arena)
			af[j] = cf
		}
		cp.function = pf
		cp.numFuncs = C.ulong(len(p.Functions))

		// Populate meths
		pm := (**C.struct_method)(arena.malloc(C.size_t(len(p.Methods)) * C.sizeof_struct_method))
		am := (*[1 << 30]*C.struct_method)(unsafe.Pointer(pm))[:len(p.Methods):len(p.Methods)]
		for j, m := range p.Methods {
			cf := (*C.struct_function)(arena.malloc(C.sizeof_struct_function))
			cf.name = arena.cstring(m.Name)
			cf.srcLineLength = C.int(m.SrcLineLength)
			cf.srcLineStart = C.int(m.SrcLineStart)
			cf.srcLineEnd = C.int(m.SrcLineEnd)
			cf.offset = C.ulonglong(m.Offset)
			cf.end = C.ulonglong(m.End)
			cf.fileName = arena.cstring(m.Filename)
			cf.packageName = arena.cstring(m.PackageName)
			cm := (*C.struct_method)(arena.malloc(C.sizeof_struct_method))
			cm.receiver = arena.cstring(m.Receiver)
			cm.function = cf
			am[j] = cm
		}
		cp.method = pm
		cp.numMeths = C.ulong(len(p.Methods))

		ppakgs[i] = cp
	}
	pr := (*C.struct_packages)(arena.malloc(C.sizeof_struct_packages))
	pr.packages = pp
	pr.length = C.ulong(len(pkgs))
	return pr
}

func convertFunction(f *gore.Function, arena *arena) *C.struct_function {
	cf := (*C.struct_function)(arena.malloc(C.sizeof_struct_function))
	cf.name = arena.cstring(f.Name)
	cf.srcLineLength = C.int(f.SrcLineLength)
	cf.srcLineStart = C.int(f.SrcLineStart)
	cf.srcLineEnd = C.int(f.SrcLineEnd)
	cf.offset = C.ulonglong(f.Offset)
	cf.end = C.ulonglong(f.End)
	cf.fileName = arena.cstring(f.Filename)
	cf.packageName = arena.cstring(f.PackageName)
	return cf
}

type arena []unsafe.Pointer

func (a *arena) malloc(size C.size_t) unsafe.Pointer {
	ptr := C.malloc(size)
	*a = append(*a, ptr)
	return ptr
}

func (a *arena) calloc(size C.size_t, n int) unsafe.Pointer {
	ptr := C.calloc(C.size_t(n), size)
	*a = append(*a, ptr)
	return ptr
}

func (a *arena) cstring(str string) *C.char {
	cs := C.CString(str)
	*a = append(*a, unsafe.Pointer(cs))
	return cs
}

func (a *arena) free() {
	for _, p := range *a {
		C.free(p)
	}
}

func main() {}

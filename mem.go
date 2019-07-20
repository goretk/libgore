// Copyright 2019 The GoRE.tk Authors. All rights reserved.
// Use of this source code is governed by the license that
// can be found in the LICENSE file.

package main

import (
	gore "github.com/goretk/gore"
)

var fman *manager

func init() {
	fman = &manager{
		openFiles: make(map[string]*gore.GoFile),
		add: make(chan *struct {
			key string
			f   *gore.GoFile
		}),
		get:        make(chan string),
		ret:        make(chan *gore.GoFile),
		remove:     make(chan string),
		openArenas: make(map[string]*arena),
		addA: make(chan *struct {
			key string
			a   *arena
		}),
		getA:    make(chan string),
		retA:    make(chan *arena),
		removeA: make(chan string),
	}
	go fman.handleLoop()
}

type manager struct {
	openFiles map[string]*gore.GoFile
	add       chan *struct {
		key string
		f   *gore.GoFile
	}
	get        chan string
	ret        chan *gore.GoFile
	remove     chan string
	openArenas map[string]*arena
	addA       chan *struct {
		key string
		a   *arena
	}
	getA    chan string
	retA    chan *arena
	removeA chan string
}

func (m *manager) handleLoop() {
	for {
		select {
		case newf := <-m.add:
			m.openFiles[newf.key] = newf.f
		case key := <-m.remove:
			delete(m.openFiles, key)
		case key := <-m.get:
			f, _ := m.openFiles[key]
			m.ret <- f
		case newA := <-m.addA:
			m.openArenas[newA.key] = newA.a
		case key := <-m.removeA:
			delete(m.openArenas, key)
		case key := <-m.getA:
			a, _ := m.openArenas[key]
			m.retA <- a
		}
	}
}

func addNewFile(path string, f *gore.GoFile) {
	fman.add <- &struct {
		key string
		f   *gore.GoFile
	}{key: path, f: f}
}

func getFile(path string) *gore.GoFile {
	fman.get <- path
	return <-fman.ret
}

func removeFile(path string) {
	fman.remove <- path
}

func addNewArena(path string, a *arena) {
	fman.addA <- &struct {
		key string
		a   *arena
	}{key: path, a: a}
}

func getArena(path string) *arena {
	fman.getA <- path
	return <-fman.retA
}

func removeArena(path string) {
	fman.removeA <- path
}

// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package ansiterm

// csiX contains Control code of AscII 0x01 - 0x19 and 0x7F, adapt to shortcut of bash shortcut.
type csiX struct {
	baseState
}

// csiDoubleX represent only from csiX state, and only accept ctrl x
type csiDoubleX struct {
	baseState
	fromState *state
}

// csiSearch when command line in search mode
type csiSearch struct {
	baseState
}

// csiSearch when command line in reverseSearch mode
type csiRSearch struct {
	baseState
}

func (csiState csiX) Handle(b byte) (s state, e error) {
	csiState.parser.logf("CsiXHandler::Handle %#x", b)
	nextState, err := csiState.baseState.Handle(b)
	if nextState != nil || err != nil {
		return nextState, err
	}

	return csiState.parser.ground, csiState.parser.CsiXDispatcher()
}

func (csiState csiDoubleX) Handle(b byte) (s state, e error) {
	csiState.parser.logf("CsiX2Handler::Handle %#x", b)
	switch {
	case isDoubleX(b):
		return *csiState.fromState, csiState.parser.CsiXDispatcher()
	case isExecute(b):
		return *csiState.fromState, csiState.parser.CsiXDispatcher()
	case isClean(b):
		return *csiState.fromState, csiState.parser.CsiXDispatcher()
	}
	return *csiState.fromState, nil
}

func (csiState *csiDoubleX) Enter() error {
	csiState.fromState = &csiState.parser.currState
	return nil
}

func (csiState csiSearch) Handle(b byte) (s state, e error) {
	csiState.parser.logf("CsiSearch::Handle %#x", b)
	switch {
	case isDoubleX(b):
		return csiState.parser.csiX2, nil
	case isExecute(b):
		return csiState.parser.ground, csiState.parser.enter()
	default:
		return csiState.parser.csiSearch, csiState.parser.CsiSearch()
	}
}

func (csiState csiRSearch) Handle(b byte) (s state, e error) {
	csiState.parser.logf("CsiSearch::Handle %#x", b)
	switch {
	case isDoubleX(b):
		return csiState.parser.csiX2, nil
	case isExecute(b):
		return csiState.parser.ground, csiState.parser.enter()
	default:
		return csiState.parser.csiSearch, csiState.parser.CsiRSearch()
	}
}

func isDoubleX(b byte) bool {
	if b == 0x18 {
		return true
	}
	return false
}

func isExecute(b byte) bool {
	if b == 0x0d || b == 0x0a {
		return true
	}
	return false
}

func isClean(b byte) bool {
	if b == 0x03 {
		return true
	}
	return false
}

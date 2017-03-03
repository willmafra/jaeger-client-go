// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package rpcmetrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleNameNormalizer(t *testing.T) {
	n := &SimpleNameNormalizer{
		SafeSets: []SafeCharacterSet{
			&Range{From: 'a', To: 'z'},
			&Char{'-'},
		},
		Replacement: '-',
	}
	assert.Equal(t, "ab-cd", n.Normalize("ab-cd"), "all valid")
	assert.Equal(t, "ab-cd", n.Normalize("ab.cd"), "single mismatch")
	assert.Equal(t, "a--cd", n.Normalize("aB-cd"), "range letter mismatch")
}

func TestNormalizedEndpoints(t *testing.T) {
	n := newNormalizedEndpoints(1, DefaultNameNormalizer)

	assertLen := func(l int) {
		n.mux.RLock()
		defer n.mux.RUnlock()
		assert.Len(t, n.names, l)
	}

	assert.Equal(t, "ab-cd", n.normalize("ab^cd"), "one translation")
	assert.Equal(t, "ab-cd", n.normalize("ab^cd"), "cache hit")
	assertLen(1)
	assert.Equal(t, "", n.normalize("xys"), "cache overflow")
}

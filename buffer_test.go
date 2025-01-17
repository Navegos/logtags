// Copyright (C) 2025 @Navegos & @DevelVitorF Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package actor provides the structures for representing an actor who has
// access to resources.

package logtags

import (
	"context"
	"strings"
	"testing"
)

func buffer(str string) *Buffer {
	f := &Buffer{}
	for _, t := range strings.Split(str, ",") {
		s := strings.SplitN(t, "=", 2)
		val := ""
		if len(s) > 1 {
			val = s[1]
		}
		f = f.Add(s[0], val)
	}
	return f
}

func TestBufferMerge(t *testing.T) {
	cases := []struct{ left, right, expected string }{
		{"a=1", "b=2", "a1,b2"},
		{"a=1,b=2", "c=3,d=4", "a1,b2,c3,d4"},

		{"a=1", "a=2", "a2"},

		{"a=1,b=2,c=3", "b=4,d=5", "a1,b4,c3,d5"},
		{"b=2,d=3", "a=4,b=5,c=6,d=7,e=8", "a4,b5,c6,d7,e8"},
		{"b=2,d=3", "a=4,b=5,c=6,d=7", "a4,b5,c6,d7"},
		{"b=2,d=3", "b=5,c=6,d=7,e=8", "b5,c6,d7,e8"},
	}
	for _, tc := range cases {
		l := buffer(tc.left)
		r := buffer(tc.right)
		if res := l.Merge(r).String(); res != tc.expected {
			t.Errorf("merge %s with %s: got %s expected %s", tc.left, tc.right, res, tc.expected)
		}
	}
}

func TestBufferAdd(t *testing.T) {
	b := buffer("a=1")
	ctx := AddTags(context.Background(), b)
	if expected, res := "a1", FromContext(ctx).String(); res != expected {
		t.Errorf("AddTags failed: expected %q, got %q", expected, res)
	}

	ctx = AddTags(ctx, FromContext(context.Background()))
	if expected, res := "a1", FromContext(ctx).String(); res != expected {
		t.Errorf("AddTags failed: expected %q, got %q", expected, res)
	}
}

func TestBufferBuild(t *testing.T) {
	bld := BuildBuffer()
	bld.Add("a", 1)
	bld.Add("b", 2)
	b := bld.Finish()
	if expected := "a1,b2"; b.String() != expected {
		t.Fatalf("expected %q, got %q", expected, b.String())
	}

	bld = BuildBuffer()
	bld.Add("a", 1)
	bld.Add("b", 2)
	bld.Add("a", 10)
	bld.Add("c", 3)
	b = bld.Finish()
	if expected := "a10,b2,c3"; b.String() != expected {
		t.Fatalf("expected %q, got %q", expected, b.String())
	}
}

func BenchmarkBuffer(b *testing.B) {
	// This benchmark uses a set of tag operations that have been observed to be
	// the most common during a mixed KV workload:
	//
	//   Left tags          Right tags
	//   -----------------------------
	//   n,client,user      n
	//   n,s                n,s,r
	//   n                  n,s
	//   client,user,n,txn  n

	l1, r1 := buffer("n,client,user"), buffer("n")
	l2, r2 := buffer("n,s"), buffer("n,s,r")
	l3, r3 := buffer("n"), buffer("n,s")
	l4, r4 := buffer("client,user,n,txn"), buffer("n")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = l1.Merge(r1)
		_ = l2.Merge(r2)
		_ = l3.Merge(r3)
		_ = l4.Merge(r4)
	}
}

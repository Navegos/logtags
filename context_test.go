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
	"testing"
)

func TestRemove(t *testing.T) {
	b := &Buffer{}
	b = b.Add("1", nil)
	b = b.Add("2", nil)
	b = b.Add("3", nil)
	ctx := WithTags(context.Background(), b)
	rctx := RemoveTag(ctx, "2")
	if FromContext(rctx).String() != "1,3" {
		t.Fatalf("expected 1,3 got: %s", FromContext(rctx).String())
	}
	rctx = RemoveTag(ctx, "3")
	if FromContext(rctx).String() != "1,2" {
		t.Fatalf("expected 1,2 got: %s", FromContext(rctx).String())
	}
	rctx = RemoveTag(ctx, "4")
	if FromContext(rctx).String() != "1,2,3" {
		t.Fatalf("expected 1,2 got: %s", FromContext(rctx).String())
	}
}

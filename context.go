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

import "context"

// contextLogTagsKey is an empty type for the handle associated with the log
// tags (*Buffer) value (see context.Value).
type contextLogTagsKey struct{}

// FromContext returns the tags stored in the context (by WithTags, AddTag, or
// AddTags).
func FromContext(ctx context.Context) *Buffer {
	if fromContextFn != nil {
		return fromContextFn(ctx)
	}
	val := ctx.Value(contextLogTagsKey{})
	if val == nil {
		return nil
	}
	return val.(*Buffer)
}

// WithTags returns a context with the given tags. Any existing tags are
// ignored.
func WithTags(ctx context.Context, tags *Buffer) context.Context {
	if withTagsFn != nil {
		return withTagsFn(ctx, tags)
	}
	return context.WithValue(ctx, contextLogTagsKey{}, tags)
}

// AddTag returns a context that has the tags in the given context plus another
// tag. Tags are deduplicated (see Buffer.AddTag).
func AddTag(ctx context.Context, key string, value interface{}) context.Context {
	b := FromContext(ctx)
	return WithTags(ctx, b.Add(key, value))
}

// RemoveTag returns a context that has the tags in the given context except the
// tag with key `key`. If such a tag does not exist, the given context is
// returned.
func RemoveTag(ctx context.Context, key string) context.Context {
	b := FromContext(ctx)
	newB, ok := b.Remove(key)
	if !ok {
		return ctx
	}
	return WithTags(ctx, newB)
}

// AddTags returns a context that has the tags in the given context plus another
// set of tags. Tags are deduplicated (see Buffer.AddTags).
func AddTags(ctx context.Context, tags *Buffer) context.Context {
	b := FromContext(ctx)
	newB := b.Merge(tags)
	if newB == b {
		return ctx
	}
	return WithTags(ctx, newB)
}

// OverrideContextFuncs can be used to override the implementation of
// FromContext and WithTags. This is useful if we have a more efficient
// implementation of context-associated values.
//
// Must be called before WithTags or FromContext are called.
func OverrideContextFuncs(
	fromContext func(ctx context.Context) *Buffer,
	withTags func(ctx context.Context, tags *Buffer) context.Context,
) {
	fromContextFn = fromContext
	withTagsFn = withTags
}

var fromContextFn func(ctx context.Context) *Buffer
var withTagsFn func(ctx context.Context, tags *Buffer) context.Context

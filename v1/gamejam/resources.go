// Copyright 2016 Pikkpoiss
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gamejam

import (
	"fmt"
	"github.com/pikkpoiss/gamejam/v1/base/render"
	"github.com/pikkpoiss/gamejam/v1/base/sprites"
)

type ResourceKey string

type ResourceType interface {
	Key() ResourceKey
	Delete()
}

type GeometryType struct {
	*render.Geometry
	key ResourceKey
}

func (t GeometryType) Key() ResourceKey {
	return t.key
}

type SheetType struct {
	*sprites.Sheet
	key ResourceKey
}

func (t SheetType) Key() ResourceKey {
	return t.key
}

type ResourceLoader interface {
	Key() ResourceKey
	Load(resources Resources) (res ResourceType, err error)
}

type Resources interface {
	Get(loader ResourceLoader) (res ResourceType, err error)
	Release(key ResourceKey) (err error)
	Delete()
}

type BaseResources struct {
	resources map[ResourceKey]ResourceType
	counts    map[ResourceKey]int
}

func NewBaseResources() *BaseResources {
	return &BaseResources{
		resources: map[ResourceKey]ResourceType{},
		counts:    map[ResourceKey]int{},
	}
}

func (r *BaseResources) Get(loader ResourceLoader) (res ResourceType, err error) {
	var (
		exists bool
		key    ResourceKey
	)
	key = loader.Key()
	if res, exists = r.resources[key]; exists {
		r.counts[key]++
		return
	}
	if res, err = loader.Load(r); err != nil {
		return
	}
	r.resources[key] = res
	r.counts[key] = 1
	return
}

func (r *BaseResources) Delete() {
	var (
		res ResourceType
		key ResourceKey
	)
	for _, res = range r.resources {
		key = res.Key()
		res.Delete()
		delete(r.counts, key)
		delete(r.resources, key)
	}
}

func (r *BaseResources) Release(key ResourceKey) (err error) {
	var (
		exists bool
		count  int
	)
	if _, exists = r.resources[key]; !exists {
		err = fmt.Errorf("No resource with key %v", key)
		return
	}
	count = r.counts[key] - 1
	if count <= 0 {
		r.resources[key].Delete()
		delete(r.counts, key)
		delete(r.resources, key)
	} else {
		r.counts[key] = count
	}
	return
}

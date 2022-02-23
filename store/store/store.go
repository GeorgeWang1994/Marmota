// Copyright 2017 Xiaomi, Inc.
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

package store

import (
	"container/list"
	"errors"
	"hash/crc32"
	"log"
	"marmota/pkg/common/model"
	"marmota/store/cc"
	"marmota/store/gg"
	"marmota/store/utils"
	"sync"
)

var StoreItems *StoreItemMap

type StoreItemMap struct {
	sync.RWMutex
	A    []map[string]*SafeLinkedList
	Size int
}

func (s *StoreItemMap) Get(key string) (*SafeLinkedList, bool) {
	s.RLock()
	defer s.RUnlock()
	idx := hashKey(key) % uint32(s.Size)
	val, ok := s.A[idx][key]
	return val, ok
}

// Remove method remove key from StoreItemMap, return true if exists
func (s *StoreItemMap) Remove(key string) bool {
	s.Lock()
	defer s.Unlock()
	idx := hashKey(key) % uint32(s.Size)
	_, exists := s.A[idx][key]
	if !exists {
		return false
	}

	delete(s.A[idx], key)
	return true
}

func (s *StoreItemMap) Getitems(idx int) map[string]*SafeLinkedList {
	s.RLock()
	defer s.RUnlock()
	items := s.A[idx]
	s.A[idx] = make(map[string]*SafeLinkedList)
	return items
}

func (s *StoreItemMap) Set(key string, val *SafeLinkedList) {
	s.Lock()
	defer s.Unlock()
	idx := hashKey(key) % uint32(s.Size)
	s.A[idx][key] = val
}

func (s *StoreItemMap) Len() int {
	s.RLock()
	defer s.RUnlock()
	var l int
	for i := 0; i < s.Size; i++ {
		l += len(s.A[i])
	}
	return l
}

func (s *StoreItemMap) First(key string) *model.StoreItem {
	s.RLock()
	defer s.RUnlock()
	idx := hashKey(key) % uint32(s.Size)
	sl, ok := s.A[idx][key]
	if !ok {
		return nil
	}

	first := sl.Front()
	if first == nil {
		return nil
	}

	return first.Value.(*model.StoreItem)
}

func (s *StoreItemMap) PushAll(key string, items []*model.StoreItem) error {
	s.Lock()
	defer s.Unlock()
	idx := hashKey(key) % uint32(s.Size)
	sl, ok := s.A[idx][key]
	if !ok {
		return errors.New("not exist")
	}
	sl.PushAll(items)
	return nil
}

func (s *StoreItemMap) GetFlag(key string) (uint32, error) {
	s.Lock()
	defer s.Unlock()
	idx := hashKey(key) % uint32(s.Size)
	sl, ok := s.A[idx][key]
	if !ok {
		return 0, errors.New("not exist")
	}
	return sl.Flag, nil
}

func (s *StoreItemMap) SetFlag(key string, flag uint32) error {
	s.Lock()
	defer s.Unlock()
	idx := hashKey(key) % uint32(s.Size)
	sl, ok := s.A[idx][key]
	if !ok {
		return errors.New("not exist")
	}
	sl.Flag = flag
	return nil
}

func (s *StoreItemMap) PopAll(key string) []*model.StoreItem {
	s.Lock()
	defer s.Unlock()
	idx := hashKey(key) % uint32(s.Size)
	sl, ok := s.A[idx][key]
	if !ok {
		return []*model.StoreItem{}
	}
	return sl.PopAll()
}

func (s *StoreItemMap) FetchAll(key string) ([]*model.StoreItem, uint32) {
	s.RLock()
	defer s.RUnlock()
	idx := hashKey(key) % uint32(s.Size)
	sl, ok := s.A[idx][key]
	if !ok {
		return []*model.StoreItem{}, 0
	}

	return sl.FetchAll()
}

func hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func getWts(key string, now int64) int64 {
	interval := int64(gg.CACHE_TIME)
	return now + interval - (int64(hashKey(key)) % interval)
}

func (s *StoreItemMap) PushFront(key string,
	item *model.StoreItem, md5 string, cfg *cc.GlobalConfig) {
	if linkedList, exists := s.Get(key); exists {
		linkedList.PushFront(item)
	} else {
		//log.Println("new key:", key)
		safeList := &SafeLinkedList{L: list.New()}
		safeList.L.PushFront(item)

		if cfg.Migrate.Enabled && !utils.IsRrdFileExist(utils.RrdFileName(
			cfg.RRD.Storage, md5, item.DsType, item.Step)) {
			safeList.Flag = gg.GRAPH_F_MISS
		}
		s.Set(key, safeList)
	}
}

func (s *StoreItemMap) KeysByIndex(idx int) []string {
	s.RLock()
	defer s.RUnlock()

	count := len(s.A[idx])
	if count == 0 {
		return []string{}
	}

	keys := make([]string, 0, count)
	for key := range s.A[idx] {
		keys = append(keys, key)
	}

	return keys
}

func (s *StoreItemMap) Back(key string) *model.StoreItem {
	s.RLock()
	defer s.RUnlock()
	idx := hashKey(key) % uint32(s.Size)
	L, ok := s.A[idx][key]
	if !ok {
		return nil
	}

	back := L.Back()
	if back == nil {
		return nil
	}

	return back.Value.(*model.StoreItem)
}

// 指定key对应的Item数量
func (s *StoreItemMap) ItemCnt(key string) int {
	s.RLock()
	defer s.RUnlock()
	idx := hashKey(key) % uint32(s.Size)
	L, ok := s.A[idx][key]
	if !ok {
		return 0
	}
	return L.Len()
}

func init() {
	size := gg.CACHE_TIME / gg.FLUSH_DISK_STEP
	if size < 0 {
		log.Panicf("store.init, bad size %d\n", size)
	}

	StoreItems = &StoreItemMap{
		A:    make([]map[string]*SafeLinkedList, size),
		Size: size,
	}
	for i := 0; i < size; i++ {
		StoreItems.A[i] = make(map[string]*SafeLinkedList)
	}
}

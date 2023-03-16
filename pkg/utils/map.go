/**
 * 安全 Map，防止并发读写
 */

package utils

import (
	"sync"
)

// Map 一级 Map
type Map struct {
	Lock sync.Mutex
	Map  map[interface{}]interface{}
}

// Map2 二级 Map
type Map2 struct {
	Lock sync.Mutex
	Map  map[interface{}]map[interface{}]interface{}
}

// Set 设置 map
func (p *Map) Set(name interface{}, value interface{}) {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	if len(p.Map) == 0 {
		p.Map = make(map[interface{}]interface{}, 0)
	}

	p.Map[name] = value
}

// Get 获取 map
func (p *Map) Get(name interface{}) (interface{}, bool) {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	mVal, mOk := p.Map[name]
	if !mOk {
		return nil, false
	}

	return mVal, true
}

// Del 删除 map
func (p *Map) Del(name interface{}) {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	delete(p.Map, name)
}

// Clear 清空 map
func (p *Map) Clear() {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	p.Map = make(map[interface{}]interface{}, 0)
}

// Size 数据条数
func (p *Map) Size() int {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	return len(p.Map)
}

// GetAll 获取全部 map
func (p *Map) GetAll() map[interface{}]interface{} {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	list := p.Map

	return list
}

// ------------------- 二级 map -------------------

// Set 设置 map
func (p *Map2) Set(key1 interface{}, key2 interface{}, value interface{}) {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	if len(p.Map) == 0 {
		p.Map = make(map[interface{}]map[interface{}]interface{}, 0)
	}

	if len(p.Map[key1]) == 0 {
		p.Map[key1] = make(map[interface{}]interface{}, 0)
	}

	p.Map[key1][key2] = value
}

// Get 获取 map
func (p *Map2) Get(key1 interface{}, key2 interface{}) (interface{}, bool) {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	mVal, mOk := p.Map[key1][key2]
	if !mOk {
		return nil, false
	}

	return mVal, true
}

// GetList 获取二级 key 列表
func (p *Map2) GetList(key1 interface{}) (map[interface{}]interface{}, bool) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	mVal, mOk := p.Map[key1]
	if !mOk {
		return nil, false
	}

	return mVal, true
}

// GetListSize 获取二级 key 列数据条数
func (p *Map2) GetListSize(key1 interface{}) int {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	_, mOk := p.Map[key1]
	if !mOk {
		return 0
	}

	return len(p.Map[key1])
}

// Del 删除一级 map
func (p *Map2) Del(key1 interface{}) {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	delete(p.Map, key1)
}

// DelList 删除二级 map
func (p *Map2) DelList(key1 interface{}, key2 interface{}) {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	delete(p.Map[key1], key2)
}

// Clear 清空 map
func (p *Map2) Clear() {

	p.Lock.Lock()
	defer p.Lock.Unlock()

	p.Map = make(map[interface{}]map[interface{}]interface{}, 0)
}

// Size 数据条数
func (p *Map2) Size() int {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	return len(p.Map)
}

// GetAll 获取全部 map
func (p *Map2) GetAll() map[interface{}]map[interface{}]interface{} {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	return p.Map
}

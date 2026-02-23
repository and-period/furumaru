package collection

import (
	"encoding/json"
	"fmt"
	"testing"
)

// BenchmarkFilterIter はフィルターイテレーターのパフォーマンスをベンチマークする。
func BenchmarkFilterIter(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		items := makeItems(size)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for b.Loop() {
				var count int
				for v := range FilterIter(items, func(item benchItem) bool {
					return item.ID%2 == 0
				}) {
					_ = v
					count++
				}
			}
		})
	}
}

// BenchmarkFilterSlice は従来のスライスフィルターとの比較用。
func BenchmarkFilterSlice(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		items := makeItems(size)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for b.Loop() {
				result := make([]benchItem, 0, len(items)/2)
				for _, v := range items {
					if v.ID%2 == 0 {
						result = append(result, v)
					}
				}
				_ = result
			}
		})
	}
}

// BenchmarkMapIter はマップイテレーターのパフォーマンスをベンチマークする。
func BenchmarkMapIter(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		items := makeItems(size)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for b.Loop() {
				result := make(map[int]benchItem, len(items))
				for k, v := range MapIter(items, func(item benchItem) (int, benchItem) {
					return item.ID, item
				}) {
					result[k] = v
				}
			}
		})
	}
}

// BenchmarkJSONMarshalUnmarshal はJSON処理のGCプレッシャーを計測する。
// furumaru の JSONColumn パターンで頻繁に使用される。
func BenchmarkJSONMarshalUnmarshal(b *testing.B) {
	sizes := []int{10, 100, 1000}
	for _, size := range sizes {
		items := makeItems(size)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for b.Loop() {
				data, err := json.Marshal(items)
				if err != nil {
					b.Fatal(err)
				}
				var result []benchItem
				if err := json.Unmarshal(data, &result); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkMapAllocation はマップ生成のGCプレッシャーを計測する。
// fill() パターンで多用される map[string]*T の生成を模擬する。
func BenchmarkMapAllocation(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}
	for _, size := range sizes {
		items := makeItems(size)
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			for b.Loop() {
				m := make(map[int]*benchItem, len(items))
				for i := range items {
					m[items[i].ID] = &items[i]
				}
				_ = m
			}
		})
	}
}

type benchItem struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func makeItems(n int) []benchItem {
	items := make([]benchItem, n)
	for i := range items {
		items[i] = benchItem{
			ID:   i,
			Name: fmt.Sprintf("item-%d", i),
			Tags: []string{"tag1", "tag2", "tag3"},
		}
	}
	return items
}

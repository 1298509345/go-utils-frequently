package heap

import (
	"encoding/json"
	"math/rand/v2"
	"sort"
	"testing"
	"time"
)

func TestHeap_Add(t *testing.T) {
	// 创建堆
	mh := New(5, func(t1, t2 int) bool {
		return t1 <= t2 // 3 5 7 10 9
		// return t1 > t2 // 9 7 5 3 1
	})

	// 添加元素
	mh.Add(5)
	mh.Add(3)
	mh.Add(7)
	mh.Add(1)
	mh.Add(9)
	t.Log(mh.Vals)
	mh.Add(10)
	t.Log(mh.Vals)
}

func TestHeap(t *testing.T) {
	var randInput, randInput2 []int

	for i := 0; i < 20; i++ {
		randInput = append(randInput, rand.IntN(200))
	}
	for i := 0; i < 20; i++ {
		randInput2 = append(randInput2, rand.IntN(200))
	}

	t.Log(randInput)
	t.Log(randInput2)

	h1 := New(10, func(t1, t2 int) bool {
		return t1 <= t2
	})

	for _, v := range randInput {
		h1.Add(v)
	}
	t.Log(h1.Vals)
	sort.IntSlice(randInput).Sort()
	t.Log(randInput)

	h2 := New(10, func(t1, t2 int) bool {
		return t1 > t2
	})

	for _, v := range randInput2 {
		h2.Add(v)
	}
	t.Log(h2.Vals)

	sort.IntSlice(randInput2).Sort()
	t.Log(randInput2)
}

type Item struct {
	ItemCreateTime int64
	SaleNum        int
}

func (i *Item) String() string {
	str, _ := json.Marshal(i)
	return string(str)
}

func TestHeapMultiPriority(t *testing.T) {
	var items []Item
	for i := 0; i < 10; i++ {
		items = append(items, Item{
			ItemCreateTime: time.Now().Add(time.Duration(rand.IntN(100)) * time.Second).Unix(),
			SaleNum:        rand.IntN(100),
		})
	}
	t.Log(items)

}

func TestOmitempty(t *testing.T) {
	type A struct {
		Bl bool `json:"bl,omitempty"`
	}
	marshal, _ := json.Marshal(A{false})
	t.Log(string(marshal))
}

func lastStoneWeight(stones []int) int {
	hp := New(len(stones), func(t1, t2 int) bool {
		return t1 > t2
	})

	for _, v := range stones {
		hp.Add(v)
	}
	for len(hp.Vals) > 1 {
		s1, _ := hp.Remove()
		s2, _ := hp.Remove()
		if dif := s2 - s1; dif > 0 {
			hp.Add(dif)
		} else if dif < 0 {
			hp.Add(-dif)
		}
	}
	ret, _ := hp.Remove()
	return ret
}

func TestItemSort(t *testing.T) {
	tmp := `{
    "5764698538646900064": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "NOT_FOUND"
        }
    },
    "5764698676438175060": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "PRESENT",
            "Int64Value": 9,
            "OldFeatureId": ""
        }
    },
    "5764698730183986482": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "NOT_FOUND"
        }
    },
    "5764698754003438920": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "PRESENT",
            "Int64Value": 3,
            "OldFeatureId": ""
        }
    },
    "5764699057863984659": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "NOT_FOUND"
        }
    },
    "5764699087706457641": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "PRESENT",
            "Int64Value": 5,
            "OldFeatureId": ""
        }
    },
    "5764699209190280066": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "PRESENT",
            "Int64Value": 9,
            "OldFeatureId": ""
        }
    },
    "5764699224000366123": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "PRESENT",
            "Int64Value": 7,
            "OldFeatureId": ""
        }
    },
    "5764699243705206282": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "PRESENT",
            "Int64Value": 6,
            "OldFeatureId": ""
        }
    },
    "5764699247635269165": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "PRESENT",
            "Int64Value": 7,
            "OldFeatureId": ""
        }
    },
    "5764699328799245852": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "PRESENT",
            "Int64Value": 10,
            "OldFeatureId": ""
        }
    },
    "5764699333811439053": {
        "item:fcTagItem:sales_item_cnt_30": {
            "ValueType": "INT64",
            "ValueStatus": "NOT_FOUND"
        }
    }
}`

	type it struct {
		ItemID     string
		Int64Value int64 `json:"Int64Value"`
	}
	m := make(map[string]map[string]it)
	var its []it

	t.Log(json.Unmarshal(([]byte)(tmp), &m))
	for k, v := range m {
		i := v["item:fcTagItem:sales_item_cnt_30"]
		i.ItemID = k
		its = append(its, i)
	}

	hp := New(5, func(t1, t2 it) bool {
		return t1.Int64Value > t2.Int64Value
	})

	for _, v := range its {
		hp.Add(v)
	}

	t.Log(hp.Vals)

}

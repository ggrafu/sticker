package cache

import (
	"reflect"
	"testing"
	"time"

	"github.com/ggrafu/sticker/utils"
)

func TestAdd(t *testing.T) {
	c := Cache{}
	t1 := time.Now()
	t2 := t1.Add(time.Hour)
	t3 := t2.Add(time.Hour)

	c.Add(t2, 2)
	conds := []struct {
		name     string
		expected bool
		actual   bool
	}{
		{
			"cache head is not nil",
			c.Head != nil,
			true,
		},
		{
			"cache head's next is nil",
			c.Head.next == nil,
			true,
		},
		{
			"cache head's prev is nil",
			c.Head.prev == nil,
			true,
		},
		{
			"cache head's value is 2",
			c.Head.value == 2,
			true,
		},
	}
	for _, cs := range conds {
		if cs.expected != cs.actual {
			t.Errorf("failed to add first element to cache, failed condition: %s", cs.name)
		}
	}

	c.Add(t3, 3)
	conds = []struct {
		name     string
		expected bool
		actual   bool
	}{
		{
			"cache head is not nil",
			c.Head != nil,
			true,
		},
		{
			"cache head's prev is nil",
			c.Head.prev == nil,
			true,
		},
		{
			"cache head's val is 3",
			c.Head.value == 3,
			true,
		},
		{
			"cache head's next val is 2",
			c.Head.next != nil && c.Head.next.value == 2,
			true,
		},
		{
			"cache head's next.prev val is 3",
			c.Head.next.prev != nil && c.Head.next.prev.value == 3,
			true,
		},
		{
			"cache head's next.next val is nil",
			c.Head.next.next == nil,
			true,
		},
	}
	for _, cs := range conds {
		if cs.expected != cs.actual {
			t.Errorf("failed to add second element to cache, failed condition: %s", cs.name)
		}
	}

	c.Add(t1, 1)
	conds = []struct {
		name     string
		expected bool
		actual   bool
	}{
		{
			"cache head's value is 3",
			c.Head.value == 3,
			true,
		},
		{
			"cache head's next value is 2",
			c.Head.next != nil && c.Head.next.value == 2,
			true,
		},
		{
			"cache head.next.next is not nil",
			c.Head.next.next != nil,
			true,
		},
		{
			"cache head.next.next value is 1",
			c.Head.next.next.value == 1,
			true,
		},
		{
			"cache head.next.prev value is 2",
			c.Head.next.next.prev.value == 2,
			true,
		},
		{
			"cache head.next.next.next value is nil",
			c.Head.next.next.next == nil,
			true,
		},
	}
	for _, cs := range conds {
		if cs.expected != cs.actual {
			t.Errorf("failed to add third element to cache, failed condition: %s", cs.name)
		}
	}

}

func TestGetLastElements(t *testing.T) {
	c := Cache{}

	if len(c.GetLastElements(0)) != 0 || len(c.GetLastElements(1)) != 0 {
		t.Error("failed to get element from empty cache")
	}

	c.Add(time.Now(), 1)
	c.Add(time.Now(), 2)
	c.Add(time.Now(), 3)

	conds := []struct {
		name     string
		expected bool
		actual   bool
	}{
		{
			"zero elements from full cache",
			reflect.DeepEqual(c.GetLastElements(0), []float32{}),
			true,
		},
		{
			"one element from full cache",
			reflect.DeepEqual(c.GetLastElements(1), []float32{3}),
			true,
		},
		{
			"three elements from full cache",
			reflect.DeepEqual(c.GetLastElements(3), []float32{3, 2, 1}),
			true,
		},
		{
			"four elements from full cache",
			reflect.DeepEqual(c.GetLastElements(4), []float32{3, 2, 1}),
			true,
		},
	}
	for _, cs := range conds {
		if cs.expected != cs.actual {
			t.Errorf("failed to get elements from cache, failed condition: %s", cs.name)
		}
	}

}

func TestUpdate(t *testing.T) {
	ts := utils.APIData{
		Metadata: map[string]interface{}{
			"test": "test",
		},
		TimeSeries: map[string]utils.Record{
			"2022-12-20": {
				Close: "1",
			},
			"2022-12-22": {
				Close: "2",
			},
			"2022-12-24": {
				Close: "3",
			},
		},
	}
	c := Cache{}
	c.Update(&ts)

	n := c.Head
	for i := 3; i > 0; i-- {
		if n.value != float32(i) {
			t.Errorf("failed to update cache: %f != %f", n.value, float32(i))
		}
		n = n.next
	}
}

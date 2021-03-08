package sqlgen

import (
	"math/rand"
	"reflect"
)

// ColumnTypeGroup stores column types -> Columns
type ColumnTypeGroup = map[ColumnType][]*Column

// GroupColumnsByColumnTypes groups all the Columns in given tables by ColumnType.
// This can be used in ON clause in JOIN.
func GroupColumnsByColumnTypes(tables ...*Table) ColumnTypeGroup {
	group := make(ColumnTypeGroup)
	for _, t := range tables {
		for _, c := range t.Columns {
			if _, ok := group[c.Tp]; ok {
				group[c.Tp] = append(group[c.Tp], c)
			} else {
				group[c.Tp] = []*Column{c}
			}
		}
	}
	return group
}

func RandColumnPairWithSameType(c ColumnTypeGroup) (*Column, *Column) {
	keys := reflect.ValueOf(c).MapKeys()
	randKey := keys[rand.Intn(len(keys))].Interface().(ColumnType)
	list := c[randKey]
	cnt := rand.Perm(len(list))
	cnt = cnt[:2]
	col1 := list[cnt[0]]
	col2 := list[cnt[1]]
	return col1, col2
}

func FilterUniqueColumns(c ColumnTypeGroup) ColumnTypeGroup {
	for k, list := range c {
		if len(list) <= 1 {
			delete(c, k)
		}
	}
	return c
}

// SwapOutParameterizedColumns substitute random Columns with `?` for prepare statements.
// It returns the substituted column in order.
func SwapOutParameterizedColumns(cols []*Column) []*Column {
	if len(cols) == 0 {
		return nil
	}
	var result []*Column
	for {
		chosenIdx := rand.Intn(len(cols))
		if cols[chosenIdx].Name != "?" {
			result = append(result, cols[chosenIdx])
			cols[chosenIdx] = &Column{Name: "?"}
		}
		if RandomBool() {
			break
		}
	}
	return result
}

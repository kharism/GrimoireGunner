package scene

import (
	"fmt"
	"testing"
)

func TestTraversal(t *testing.T) {
	levelLayout := GenerateLayout1()
	tiers := [][]*LevelNode{}
	TraverseBfs(levelLayout.Root, &tiers)
	for _, c := range tiers {
		ids := []string{}
		for _, d := range c {
			ids = append(ids, d.Id)
		}
		fmt.Println(ids)
	}

}

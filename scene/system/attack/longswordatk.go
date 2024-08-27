package attack

import (
	"github.com/kharism/grimoiregunner/scene/component"
	"github.com/yohamta/donburi/ecs"
)

// cost 2 EN to execute. 80 dmg 2 tiles in front
func NewLongSwordAttack(ensetgetter ENSetGetter, ecs *ecs.ECS, playerScrLoc component.ScreenPosComponentData, playerGridLoc component.GridPosComponentData) {
	en := ensetgetter.GetEn()
	if en >= 200 {
		ensetgetter.SetEn(en - 200)
		newLongSwordAttack(ecs, playerScrLoc, playerGridLoc, 80)
	}
}

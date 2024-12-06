package assets

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed images/cannoneer.png
var cannoneer []byte

//go:embed images/gatlingghoul.png
var gatlingghoul []byte

//go:embed images/gatlingghoul_atk.png
var gatlingghoul_atk []byte

//go:embed images/cannoneer_atk.png
var cannoneer_atk []byte

//go:embed images/bloombomber.png
var bloombomber []byte

//go:embed images/bloombomber_atk.png
var bloombomber_atk []byte

//go:embed images/buzzer_1.png
var buzzer_1 []byte

//go:embed images/buzzer_2.png
var buzzer_2 []byte

//go:embed images/reaper1.png
var reaper []byte

//go:embed images/reaper2.png
var reaper_warmup []byte

//go:embed images/reaper3.png
var reaper_cooldown []byte

//go:embed images/demon_1.png
var demon []byte

//go:embed images/demon_2.png
var demon_warmup []byte

//go:embed images/demon_3.png
var demon_cooldown []byte

//go:embed images/slime.png
var slime []byte

//go:embed images/slime2.png
var slime_atk []byte

//go:embed images/yeti_1.png
var yeti []byte

//go:embed images/yeti_2.png
var yeti_warmup []byte

//go:embed images/yeti_3.png
var yeti_cooldown []byte

//go:embed images/yeti_4.png
var yeti_warmup_2 []byte

//go:embed images/yeti_5.png
var yeti_cooldown_2 []byte

//go:embed images/hammerghoul.png
var hammerghoul []byte

//go:embed images/hammerghoul_atk.png
var hammerghoul_warmup []byte

//go:embed images/hammerghoul_atk2.png
var hammerghoul_cooldown []byte

//go:embed images/swordsman.png
var swordswomen []byte

//go:embed images/swordsman2.png
var swordswomen2 []byte

//go:embed images/swordsman3.png
var swordswomen3 []byte

//go:embed images/swordsman_atk1.png
var swordswomen4 []byte

//go:embed images/poacher.png
var poacher []byte

//go:embed images/poacher2.png
var poacher2 []byte

//go:embed images/poacher3.png
var poacher3 []byte

//go:embed images/yanma.png
var yanma1 []byte

//go:embed images/yanma_2.png
var yanma2 []byte

//go:embed images/yanma_option.png
var yanma_option []byte

//go:embed images/pyro-eyes.png
var pyro_eyes []byte

//go:embed images/pyro-eyes_2.png
var pyro_eyes_2 []byte

//go:embed images/lightning_imp_1.png
var lightning_imp_1 []byte

//go:embed images/lightning_imp_2.png
var lightning_imp_2 []byte

var Cannoneer *ebiten.Image
var CannoneerAtk *ebiten.Image
var BloombomberAtk *ebiten.Image
var Bloombomber *ebiten.Image
var Buzzer1 *ebiten.Image
var Buzzer2 *ebiten.Image
var GatlingghoulAtk *ebiten.Image
var Gatlingghoul *ebiten.Image
var Reaper *ebiten.Image
var ReaperWarmup *ebiten.Image
var ReaperCooldown *ebiten.Image
var Hammerghoul *ebiten.Image
var HammerghoulWarmup *ebiten.Image
var HammerghoulCooldown *ebiten.Image
var Slime *ebiten.Image
var Slime2 *ebiten.Image
var Demon *ebiten.Image
var DemonWarmup *ebiten.Image
var DemonCooldown *ebiten.Image
var Swordswomen *ebiten.Image
var SwordswomenShoot *ebiten.Image
var SwordswomenWarmup *ebiten.Image
var SwordswomenCooldown *ebiten.Image
var Poacher *ebiten.Image
var PoacherWarmup *ebiten.Image
var PoacherCooldown *ebiten.Image
var PyroEyes *ebiten.Image
var PyroEyesWarmup *ebiten.Image
var Yeti *ebiten.Image
var YetiWarmup *ebiten.Image
var YetiCooldown *ebiten.Image
var YetiWarmup2 *ebiten.Image
var YetiCooldown2 *ebiten.Image
var Yanma *ebiten.Image
var YanmaAttack *ebiten.Image
var YanmaOption *ebiten.Image
var LightningImp *ebiten.Image
var LightningImpWarmup *ebiten.Image

func init() {
	if Cannoneer == nil {
		imgReader := bytes.NewReader(cannoneer)
		Cannoneer, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if CannoneerAtk == nil {
		imgReader := bytes.NewReader(cannoneer_atk)
		CannoneerAtk, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Bloombomber == nil {
		imgReader := bytes.NewReader(bloombomber)
		Bloombomber, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BloombomberAtk == nil {
		imgReader := bytes.NewReader(bloombomber_atk)
		BloombomberAtk, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Slime == nil {
		imgReader := bytes.NewReader(slime)
		Slime, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Buzzer1 == nil {
		imgReader := bytes.NewReader(buzzer_1)
		Buzzer1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Buzzer2 == nil {
		imgReader := bytes.NewReader(buzzer_2)
		Buzzer2, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Slime2 == nil {
		imgReader := bytes.NewReader(slime_atk)
		Slime2, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Poacher == nil {
		imgReader := bytes.NewReader(poacher)
		Poacher, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if PoacherCooldown == nil {
		imgReader := bytes.NewReader(poacher2)
		PoacherCooldown, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if PoacherWarmup == nil {
		imgReader := bytes.NewReader(poacher3)
		PoacherWarmup, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Demon == nil {
		imgReader := bytes.NewReader(demon)
		Demon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if DemonWarmup == nil {
		imgReader := bytes.NewReader(demon_warmup)
		DemonWarmup, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if DemonCooldown == nil {
		imgReader := bytes.NewReader(demon_cooldown)
		DemonCooldown, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Reaper == nil {
		imgReader := bytes.NewReader(reaper)
		Reaper, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ReaperWarmup == nil {
		imgReader := bytes.NewReader(reaper_warmup)
		ReaperWarmup, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ReaperCooldown == nil {
		imgReader := bytes.NewReader(reaper_cooldown)
		ReaperCooldown, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Hammerghoul == nil {
		imgReader := bytes.NewReader(hammerghoul)
		Hammerghoul, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if HammerghoulWarmup == nil {
		imgReader := bytes.NewReader(hammerghoul_warmup)
		HammerghoulWarmup, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if HammerghoulCooldown == nil {
		imgReader := bytes.NewReader(hammerghoul_cooldown)
		HammerghoulCooldown, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Swordswomen == nil {
		imgReader := bytes.NewReader(swordswomen)
		Swordswomen, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if SwordswomenShoot == nil {
		imgReader := bytes.NewReader(swordswomen4)
		SwordswomenShoot, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if SwordswomenWarmup == nil {
		imgReader := bytes.NewReader(swordswomen2)
		SwordswomenWarmup, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if SwordswomenCooldown == nil {
		imgReader := bytes.NewReader(swordswomen3)
		SwordswomenCooldown, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Yeti == nil {
		imgReader := bytes.NewReader(yeti)
		Yeti, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if YetiWarmup == nil {
		imgReader := bytes.NewReader(yeti_warmup)
		YetiWarmup, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if YetiCooldown == nil {
		imgReader := bytes.NewReader(yeti_cooldown)
		YetiCooldown, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if YetiWarmup2 == nil {
		imgReader := bytes.NewReader(yeti_warmup_2)
		YetiWarmup2, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if YetiCooldown2 == nil {
		imgReader := bytes.NewReader(yeti_cooldown_2)
		YetiCooldown2, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Yanma == nil {
		imgReader := bytes.NewReader(yanma1)
		Yanma, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if YanmaAttack == nil {
		imgReader := bytes.NewReader(yanma2)
		YanmaAttack, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if YanmaOption == nil {
		imgReader := bytes.NewReader(yanma_option)
		YanmaOption, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if PyroEyes == nil {
		imgReader := bytes.NewReader(pyro_eyes)
		PyroEyes, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if PyroEyesWarmup == nil {
		imgReader := bytes.NewReader(pyro_eyes_2)
		PyroEyesWarmup, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if LightningImp == nil {
		imgReader := bytes.NewReader(lightning_imp_1)
		LightningImp, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if LightningImpWarmup == nil {
		imgReader := bytes.NewReader(lightning_imp_2)
		LightningImpWarmup, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
}

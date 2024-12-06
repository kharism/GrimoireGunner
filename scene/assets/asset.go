package assets

import (
	"bytes"
	_ "embed"
	_ "image/jpeg"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kharism/hanashi/core"
)

//go:embed images/tile_blue.png
var blueTilePng []byte

//go:embed images/tile_red.png
var redTilePng []byte

//go:embed images/grid_targetted.png
var grid_targetted []byte

//go:embed images/dmggrid.png
var tileDmgPng []byte

//go:embed images/cardtemplate.png
var cardTemplate []byte

//go:embed images/double_damage.png
var double_damage []byte

//go:embed images/menu_btn_bg.png
var menuBg []byte

//go:embed images/cardpick.png
var cardPick []byte

//go:embed images/cardpick_red.png
var cardPick_red []byte

//go:embed images/icicle.png
var icicle []byte

//go:embed images/basicsprite2.png
var player1Stand []byte

//go:embed images/attacksprite2.png
var player1attack []byte

//go:embed images/magibullet.png
var projectile1 []byte

//go:embed images/fx/chargeshot.png
var chargeshot_projectile []byte

//go:embed images/dagger.png
var projectile2 []byte

//go:embed images/elec_sphere.png
var elec_sphere []byte

//go:embed images/bomb1.png
var bomb1 []byte

//go:embed images/bomb2.png
var bomb2 []byte

//go:embed images/bamboo_lance.png
var bamboo_lance []byte

//go:embed images/boulder.png
var boulder []byte

//go:embed images/wall.png
var wall []byte

//go:embed images/lightning_bolt.png
var lightningbolt []byte

//go:embed fonts/PixelOperator8-bold.ttf
var PixelFontTTF []byte

//go:embed "fonts/Unispace Bd.otf"
var UnispaceBdTTF []byte

//go:embed fonts/monogram.ttf
var MonogramTTF []byte

//go:embed images/forrest_path_bg.png
var bg_forest []byte

//go:embed images/bg_mountain.png
var bg_mountain []byte

//go:embed images/bg_cave.png
var bg_cave []byte

//go:embed images/opening_screen.png
var bg_opening []byte

//go:embed images/restbg.png
var bg_rest []byte

//go:embed images/workbench.png
var bg_workbench []byte

//go:embed images/beartrap.png
var beartrap []byte

//go:embed images/wideslash.png
var wideslash_fx []byte

//go:embed images/fx/longsword.png
var sword_fx []byte

//go:embed images/fx/explosion.png
var explosion_fx []byte

//go:embed images/fx/dust.png
var dust_fx []byte

//go:embed images/fist.png
var fist_fx []byte

//go:embed images/fx/shockwave.png
var shockwave_fx []byte

//go:embed images/fx/buckshot.png
var buckshot_fx []byte

//go:embed images/fx/flametower.png
var flametower_fx []byte

//go:embed images/fx/spore.png
var spore_fx []byte

//go:embed images/fx_charge.png
var chargeshot_fx []byte

//go:embed images/fx/hit.png
var hit_fx []byte

//go:embed images/fx/flamethrower.png
var flamethrower_fx []byte

//go:embed images/fx/heal.png
var heal_fx []byte

//go:embed images/icon_longsword.png
var longsword_icon []byte

//go:embed images/icon_widesword.png
var widesword_icon []byte

//go:embed images/icon_gatling.png
var gatling_icon []byte

//go:embed images/icon_buckshot.png
var buckshot_icon []byte

//go:embed images/icon_lightning.png
var lightning_icon []byte

//go:embed images/icon_na.png
var na_icon []byte

//go:embed images/icon_wall.png
var wall_icon []byte

//go:embed images/icon_fist.png
var fist_icon []byte

//go:embed images/icon_shockwave.png
var shockwave_icon []byte

//go:embed images/icon_flamethrower.png
var flamethrower_icon []byte

//go:embed images/icon_sporebomb.png
var sporebomb_icon []byte

//go:embed images/icon_bomb.png
var bomb_icon []byte

//go:embed images/icon_cannon.png
var cannon_icon []byte

//go:embed images/icon_bamboo_lance.png
var bamboolance_icon []byte

//go:embed images/icon_heal.png
var heal_icon []byte

//go:embed images/icon_atkup.png
var atkup_icon []byte

//go:embed images/icon_workbench.png
var workbench_icon []byte

//go:embed images/icon_hpup.png
var hpup_icon []byte

//go:embed images/icon_enup.png
var enup_icon []byte

//go:embed images/icon_medkit.png
var medkit_icon []byte

//go:embed images/icon_shotgun.png
var shotgun_icon []byte

//go:embed images/icon_pushgun.png
var pushgun_icon []byte

//go:embed images/icon_icespike.png
var icespike_icon []byte

//go:embed images/icon_chargeshot.png
var chargeshot_icon []byte

//go:embed images/icon_battle.png
var battle_icon []byte

//go:embed images/icon_rest.png
var rest_icon []byte

//go:embed images/icon_firewall.png
var firewall_icon []byte

//go:embed images/stageclear.png
var stageclear []byte

//go:embed images/gameover_basic.png
var gameover []byte

var BlueTile *ebiten.Image
var RedTile *ebiten.Image
var DamageGrid *ebiten.Image
var TargetedGrid *ebiten.Image
var Bomb1 *ebiten.Image
var Bomb2 *ebiten.Image
var BearTrap *ebiten.Image
var BambooLance *ebiten.Image
var Flamehtrower *ebiten.Image
var Wall *ebiten.Image
var Bg *ebiten.Image
var BgRest *ebiten.Image
var BgWorkbench *ebiten.Image
var BgForrest *ebiten.Image
var BgMountain *ebiten.Image
var BgCave *ebiten.Image
var BgOpening *ebiten.Image
var GameOver *ebiten.Image
var DoubleDamage *ebiten.Image
var MenuButtonBg *ebiten.Image
var Player1Stand *ebiten.Image
var Player1Attack *ebiten.Image
var Projectile1 *ebiten.Image
var Projectile2 *ebiten.Image
var ElecSphere *ebiten.Image
var Icicle *ebiten.Image
var LightningBolt *ebiten.Image
var Fist *ebiten.Image
var Boulder *ebiten.Image

var LightningIcon *ebiten.Image
var LongSwordIcon *ebiten.Image
var WideSwordIcon *ebiten.Image
var ShockwaveIcon *ebiten.Image
var BuckshotIcon *ebiten.Image
var BombIcon *ebiten.Image
var FirewallIcon *ebiten.Image
var GatlingIcon *ebiten.Image
var CannonIcon *ebiten.Image
var ShotgunIcon *ebiten.Image
var FlamethrowerIcon *ebiten.Image
var IcespikeIcon *ebiten.Image
var PushgunIcon *ebiten.Image
var NAIcon *ebiten.Image
var WallIcon *ebiten.Image
var BattleIcon *ebiten.Image
var HealIcon *ebiten.Image
var BambooLanceIcon *ebiten.Image
var AtkUp *ebiten.Image
var SporebombIcon *ebiten.Image
var HPUpIcon *ebiten.Image
var ENUpIcon *ebiten.Image
var MedkitIcon *ebiten.Image
var ChargeshotIcon *ebiten.Image
var RestIcon *ebiten.Image
var FistIcon *ebiten.Image
var WorkbenchIcon *ebiten.Image

var PixelFont *text.GoTextFaceSource
var MonogramFont *text.GoTextFaceSource
var UnispaceFont *text.GoTextFaceSource

var FontFace *text.GoTextFace
var MonogramFace *text.GoTextFace
var UnispaceFace *text.GoTextFace

var SwordAtkRaw *ebiten.Image
var ExplosionRaw *ebiten.Image
var HitRaw *ebiten.Image
var ShockWaveFxRaw *ebiten.Image
var DustFxRaw *ebiten.Image
var WideslashRaw *ebiten.Image
var BuckShotRaw *ebiten.Image
var FlametowerRaw *ebiten.Image
var SporeRaw *ebiten.Image
var HealFx *ebiten.Image
var ChargeshotRaw *ebiten.Image
var ChargeshotFx *ebiten.Image

var StageClear *ebiten.Image
var CardTemplate *ebiten.Image
var CardPick *ebiten.Image
var CardPickRed *ebiten.Image

// audio stuff

//go:embed bgm/test_melody_orchestra.mp3
var Menumusic []byte

//go:embed bgm/test1.mp3
var BattleMusic1 []byte

//go:embed bgm/area12-131883.mp3
var BattleMusic2 []byte

//go:embed bgm/chiptune-grooving-142242.mp3
var BattleMusic3 []byte

//go:embed bgm/upbeat.mp3
var IntermissionMusic []byte

//go:embed bgm/fanfare3.mp3
var Fanfare []byte

// SFX

//go:embed sfx/menu.mp3
var MenuMove []byte

//go:embed sfx/magibullet.mp3
var MagibulletFx []byte

//go:embed sfx/impact.mp3
var ImpactFx []byte

//go:embed sfx/explosion.mp3
var ExplosionFx []byte

//go:embed sfx/swosh.mp3
var SlashFx []byte

//go:embed sfx/cursor_move.mp3
var CursorFx []byte

//go:embed sfx/shotgun.mp3
var ShotgunFx []byte

//go:embed sfx/hitscan.mp3
var HitscanFx []byte

//go:embed sfx/heal.mp3
var HealsFx []byte

//go:embed sfx/lightning.mp3
var LightningFx []byte

var TileWidth int
var TileHeight int

var TileStartX = float64(165.0)
var TileStartY = float64(360.0)

// return x,y screen coordinate
func GridCoord2Screen(Row, Col int) (float64, float64) {
	return TileStartX + float64(Col)*float64(TileWidth), TileStartY + float64(Row)*float64(TileHeight)
}

// param screen X,Y coords
// return col,row
func Coord2Grid(X, Y float64) (int, int) {
	col := int(X+float64(TileWidth/2)-TileStartX) / TileWidth
	row := int(Y+float64(TileHeight/2)-TileStartY) / TileHeight
	return col, row
}

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(PixelFontTTF))
	if err != nil {
		log.Fatal(err)
	}
	PixelFont = s
	FontFace = &text.GoTextFace{
		Source: PixelFont,
		Size:   15,
	}
	s2, err := text.NewGoTextFaceSource(bytes.NewReader(MonogramTTF))
	if err != nil {
		log.Fatal(err)
	}
	s3, err := text.NewGoTextFaceSource(bytes.NewReader(UnispaceBdTTF))
	if err != nil {
		log.Fatal(err)
	}
	MonogramFont = s2
	PixelFont = s
	UnispaceFont = s3

	if BlueTile == nil {
		imgReader := bytes.NewReader(blueTilePng)
		BlueTile, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if RedTile == nil {
		imgReader := bytes.NewReader(redTilePng)
		RedTile, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if TileWidth == 0 {
		rect := RedTile.Bounds()
		TileWidth = rect.Dx()
		TileHeight = rect.Dy()
	}
	if TargetedGrid == nil {
		imgReader := bytes.NewReader(grid_targetted)
		TargetedGrid, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if GameOver == nil {
		imgReader := bytes.NewReader(gameover)
		GameOver, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if MenuButtonBg == nil {
		imgReader := bytes.NewReader(menuBg)
		MenuButtonBg, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Player1Stand == nil {
		imgReader := bytes.NewReader(player1Stand)
		Player1Stand, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Player1Attack == nil {
		imgReader := bytes.NewReader(player1attack)
		Player1Attack, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BgRest == nil {
		imgReader := bytes.NewReader(bg_rest)
		BgRest, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BgOpening == nil {
		imgReader := bytes.NewReader(bg_opening)
		BgOpening, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BearTrap == nil {
		imgReader := bytes.NewReader(beartrap)
		BearTrap, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Bomb1 == nil {
		imgReader := bytes.NewReader(bomb1)
		Bomb1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BambooLance == nil {
		imgReader := bytes.NewReader(bamboo_lance)
		BambooLance, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Wall == nil {
		imgReader := bytes.NewReader(wall)
		Wall, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Fist == nil {
		imgReader := bytes.NewReader(fist_fx)
		Fist, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Icicle == nil {
		imgReader := bytes.NewReader(icicle)
		Icicle, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Flamehtrower == nil {
		imgReader := bytes.NewReader(flamethrower_fx)
		Flamehtrower, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ElecSphere == nil {
		imgReader := bytes.NewReader(elec_sphere)
		ElecSphere, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Bomb2 == nil {
		imgReader := bytes.NewReader(bomb2)
		Bomb2, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if DoubleDamage == nil {
		imgReader := bytes.NewReader(double_damage)
		DoubleDamage, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Projectile1 == nil {
		imgReader := bytes.NewReader(projectile1)
		Projectile1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Projectile2 == nil {
		imgReader := bytes.NewReader(projectile2)
		Projectile2, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if LightningBolt == nil {
		imgReader := bytes.NewReader(lightningbolt)
		LightningBolt, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Boulder == nil {
		imgReader := bytes.NewReader(boulder)
		Boulder, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BgForrest == nil {
		imgReader := bytes.NewReader(bg_forest)
		BgForrest, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BgMountain == nil {
		imgReader := bytes.NewReader(bg_mountain)
		BgMountain, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BgWorkbench == nil {
		imgReader := bytes.NewReader(bg_workbench)
		BgWorkbench, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BgCave == nil {
		imgReader := bytes.NewReader(bg_cave)
		BgCave, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if DamageGrid == nil {
		imgReader := bytes.NewReader(tileDmgPng)
		DamageGrid, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}

	if LongSwordIcon == nil {
		imgReader := bytes.NewReader(longsword_icon)
		LongSwordIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if PushgunIcon == nil {
		imgReader := bytes.NewReader(pushgun_icon)
		PushgunIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BombIcon == nil {
		imgReader := bytes.NewReader(bomb_icon)
		BombIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BambooLanceIcon == nil {
		imgReader := bytes.NewReader(bamboolance_icon)
		BambooLanceIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if CannonIcon == nil {
		imgReader := bytes.NewReader(cannon_icon)
		CannonIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if WallIcon == nil {
		imgReader := bytes.NewReader(wall_icon)
		WallIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if IcespikeIcon == nil {
		imgReader := bytes.NewReader(icespike_icon)
		IcespikeIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if SporebombIcon == nil {
		imgReader := bytes.NewReader(sporebomb_icon)
		SporebombIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if FlamethrowerIcon == nil {
		imgReader := bytes.NewReader(flamethrower_icon)
		FlamethrowerIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ShotgunIcon == nil {
		imgReader := bytes.NewReader(shotgun_icon)
		ShotgunIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if MedkitIcon == nil {
		imgReader := bytes.NewReader(medkit_icon)
		MedkitIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if WideSwordIcon == nil {
		imgReader := bytes.NewReader(widesword_icon)
		WideSwordIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if LightningIcon == nil {
		imgReader := bytes.NewReader(lightning_icon)
		LightningIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BattleIcon == nil {
		imgReader := bytes.NewReader(battle_icon)
		BattleIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ChargeshotIcon == nil {
		imgReader := bytes.NewReader(chargeshot_icon)
		ChargeshotIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if FistIcon == nil {
		imgReader := bytes.NewReader(fist_icon)
		FistIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}

	if RestIcon == nil {
		imgReader := bytes.NewReader(rest_icon)
		RestIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if WorkbenchIcon == nil {
		imgReader := bytes.NewReader(workbench_icon)
		WorkbenchIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if AtkUp == nil {
		imgReader := bytes.NewReader(atkup_icon)
		AtkUp, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BuckshotIcon == nil {
		imgReader := bytes.NewReader(buckshot_icon)
		BuckshotIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ShockwaveIcon == nil {
		imgReader := bytes.NewReader(shockwave_icon)
		ShockwaveIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if HPUpIcon == nil {
		imgReader := bytes.NewReader(hpup_icon)
		HPUpIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ENUpIcon == nil {
		imgReader := bytes.NewReader(enup_icon)
		ENUpIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if HealIcon == nil {
		imgReader := bytes.NewReader(heal_icon)
		HealIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if FirewallIcon == nil {
		imgReader := bytes.NewReader(firewall_icon)
		FirewallIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if GatlingIcon == nil {
		imgReader := bytes.NewReader(gatling_icon)
		GatlingIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if NAIcon == nil {
		imgReader := bytes.NewReader(na_icon)
		NAIcon, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}

	if ExplosionRaw == nil {
		imgReader := bytes.NewReader(explosion_fx)
		ExplosionRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if FlametowerRaw == nil {
		imgReader := bytes.NewReader(flametower_fx)
		FlametowerRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if HitRaw == nil {
		imgReader := bytes.NewReader(hit_fx)
		HitRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if BuckShotRaw == nil {
		imgReader := bytes.NewReader(buckshot_fx)
		BuckShotRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if SporeRaw == nil {
		imgReader := bytes.NewReader(spore_fx)
		SporeRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ShockWaveFxRaw == nil {
		imgReader := bytes.NewReader(shockwave_fx)
		ShockWaveFxRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ChargeshotRaw == nil {
		imgReader := bytes.NewReader(chargeshot_projectile)
		ChargeshotRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if ChargeshotFx == nil {
		imgReader := bytes.NewReader(chargeshot_fx)
		ChargeshotFx, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if DustFxRaw == nil {
		imgReader := bytes.NewReader(dust_fx)
		DustFxRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if HealFx == nil {
		imgReader := bytes.NewReader(heal_fx)
		HealFx, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if WideslashRaw == nil {
		imgReader := bytes.NewReader(wideslash_fx)
		WideslashRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Gatlingghoul == nil {
		imgReader := bytes.NewReader(gatlingghoul)
		Gatlingghoul, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if GatlingghoulAtk == nil {
		imgReader := bytes.NewReader(gatlingghoul_atk)
		GatlingghoulAtk, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if StageClear == nil {
		imgReader := bytes.NewReader(stageclear)
		StageClear, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if CardTemplate == nil {
		imgReader := bytes.NewReader(cardTemplate)
		CardTemplate, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if CardPick == nil {
		imgReader := bytes.NewReader(cardPick)
		CardPick, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if CardPickRed == nil {
		imgReader := bytes.NewReader(cardPick_red)
		CardPickRed, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if SwordAtkRaw == nil {
		imgReader := bytes.NewReader(sword_fx)
		SwordAtkRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
		// SwordAtkAnim = &core.AnimatedImage{
		// 	MovableImage:   core.NewMovableImage(atkAnim, core.NewMovableImageParams()),
		// 	SubImageWidth:  200,
		// 	SubImageHeight: 50,
		// 	SubImageStartX: 0,
		// 	SubImageStartY: 0,
		// 	Modulo:         6,
		// }

	}
}

type SpriteParam struct {
	ScreenX, ScreenY float64
	Modulo           int
	Done             func()
}

func NewShockwaveAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(ShockWaveFxRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  100,
		SubImageHeight: 100,
		Modulo:         param.Modulo,
		FrameCount:     5,
		Done:           param.Done,
	}
}
func NewSporeAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(SporeRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  100,
		SubImageHeight: 100,
		Modulo:         param.Modulo,
		FrameCount:     5,
		Done:           param.Done,
	}
}
func NewChargeShotAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(ChargeshotRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  100,
		SubImageHeight: 160,
		Modulo:         param.Modulo,
		FrameCount:     5,
		Done:           param.Done,
	}
}
func NewExplosionAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(ExplosionRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  75,
		SubImageHeight: 75,
		Modulo:         param.Modulo,
		FrameCount:     11,
		Done:           param.Done,
	}
}
func NewDustAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(DustFxRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  100,
		SubImageHeight: 50,
		Modulo:         param.Modulo,
		FrameCount:     4,
		Done:           param.Done,
	}
}
func NewHitAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(HitRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  128,
		SubImageHeight: 128,
		Modulo:         param.Modulo,
		FrameCount:     6,
		Done:           param.Done,
	}
}
func NewWideSlashAtkAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(WideslashRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  100,
		SubImageHeight: 150,
		Modulo:         param.Modulo,
		FrameCount:     5,
		Done:           param.Done,
	}
}
func NewBuckshotAtkAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(BuckShotRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  200,
		SubImageHeight: 150,
		Modulo:         param.Modulo,
		FrameCount:     3,
		Done:           param.Done,
	}
}
func NewSwordAtkAnim(param SpriteParam) *core.AnimatedImage {
	return &core.AnimatedImage{
		MovableImage: core.NewMovableImage(SwordAtkRaw,
			core.NewMovableImageParams().
				WithMoveParam(core.MoveParam{Sx: param.ScreenX, Sy: param.ScreenY}),
		),
		SubImageStartX: 0,
		SubImageStartY: 0,
		SubImageWidth:  200,
		SubImageHeight: 50,
		Modulo:         param.Modulo,
		FrameCount:     6,
		Done:           param.Done,
	}
}

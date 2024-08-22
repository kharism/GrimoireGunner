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

//go:embed images/basicsprite2.png
var player1Stand []byte

//go:embed images/attacksprite.png
var player1attack []byte

//go:embed images/magibullet.png
var projectile1 []byte

//go:embed images/bomb1.png
var bomb1 []byte

//go:embed images/boulder.png
var boulder []byte

//go:embed images/lightning_bolt.png
var lightningbolt []byte

//go:embed fonts/PixelOperator8-bold.ttf
var PixelFontTTF []byte

//go:embed images/bg_forest/bg.png
var bg_forest []byte

//go:embed images/pyro-eyes.png
var pyro_eyes []byte

//go:embed images/cannoneer.png
var cannoneer []byte

//go:embed images/cannoneer_atk.png
var cannoneer_atk []byte

//go:embed images/bloombomber.png
var bloombomber []byte

//go:embed images/bloombomber_atk.png
var bloombomber_atk []byte

//go:embed images/fx/longsword.png
var sword_fx []byte

//go:embed images/fx/explosion.png
var explosion_fx []byte

//go:embed images/fx/hit.png
var hit_fx []byte

var BlueTile *ebiten.Image
var RedTile *ebiten.Image
var DamageGrid *ebiten.Image
var TargetedGrid *ebiten.Image
var Bomb1 *ebiten.Image
var Bg *ebiten.Image
var BgForrest *ebiten.Image
var Player1Stand *ebiten.Image
var Player1Attack *ebiten.Image
var Projectile1 *ebiten.Image
var LightningBolt *ebiten.Image
var Boulder *ebiten.Image
var PyroEyes *ebiten.Image
var Cannoneer *ebiten.Image
var CannoneerAtk *ebiten.Image
var BloombomberAtk *ebiten.Image
var Bloombomber *ebiten.Image

var PixelFont *text.GoTextFaceSource
var FontFace *text.GoTextFace

var SwordAtkRaw *ebiten.Image
var ExplosionRaw *ebiten.Image
var HitRaw *ebiten.Image

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
	if Player1Stand == nil {
		imgReader := bytes.NewReader(player1Stand)
		Player1Stand, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Player1Attack == nil {
		imgReader := bytes.NewReader(player1attack)
		Player1Attack, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Bomb1 == nil {
		imgReader := bytes.NewReader(bomb1)
		Bomb1, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if Projectile1 == nil {
		imgReader := bytes.NewReader(projectile1)
		Projectile1, _, _ = ebitenutil.NewImageFromReader(imgReader)
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
	if DamageGrid == nil {
		imgReader := bytes.NewReader(tileDmgPng)
		DamageGrid, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if PyroEyes == nil {
		imgReader := bytes.NewReader(pyro_eyes)
		PyroEyes, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
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
	if ExplosionRaw == nil {
		imgReader := bytes.NewReader(explosion_fx)
		ExplosionRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
	}
	if HitRaw == nil {
		imgReader := bytes.NewReader(hit_fx)
		HitRaw, _, _ = ebitenutil.NewImageFromReader(imgReader)
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

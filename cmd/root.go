package cmd

import (
	"embed"
	"fmt"
	"image"
	"io/fs"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/beta/freetype/truetype"
	"github.com/corona10/goimagehash"
	"github.com/fogleman/gg"
	"github.com/k1LoW/ffff"
	"github.com/mattn/go-lsd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed ipag.ttf
var jpFontFile embed.FS

var rootCmd = &cobra.Command{
	Use:   "homein",
	Short: "Detecting homograph domains using Levenshtein and Hamming distance",
	Long: `Detecting homograph domains using Levenshtein and Hamming distance

	$ homein [strng1] [string2] [options]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := argsValidation(args); err != nil {
			return err
		}
		return run(args[0], args[1])
	},
}

func init() {
	flags := rootCmd.PersistentFlags()
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("HOMEIN")
	flags.BoolP("enable-output-images", "", false, "outputs generated images")
	viper.BindPFlag("enable-output-images", flags.Lookup("enable-output-images"))

}

func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	return rootCmd.Execute()
}

func levenshtein(s1, s2 string) (int, float64) {
	d := lsd.StringDistance(s1, s2)
	p := 1 - float64(d)/float64(max(utf8.RuneCountInString(s1), utf8.RuneCountInString(s2)))
	return d, p
}

func generateImage(s string, width, high int, face font.Face, savePath string) image.Image {
	dc := gg.NewContext(width, high)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	dc.SetFontFace(face)
	dc.DrawRectangle(0, 0, float64(width), float64(high))
	dc.DrawStringAnchored(s, float64(width/2), float64(high/2), 0.5, 0.5)

	if savePath != "" {
		fmt.Printf("save image: %s\n", savePath)
		dc.SavePNG(savePath)
	}

	return dc.Image()
}

func generateImageHash(img image.Image) (*goimagehash.ImageHash, error) {
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func argsValidation(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("invalid args")
	}
	return nil
}

func run(s1, s2 string) error {
	ld, lp := levenshtein(s1, s2)
	fmt.Printf("levenshtein distance: %d, levenshtein percent: %f\n", ld, lp)
	to := &truetype.Options{
		Size:    1.0,
		DPI:     1200,
		Hinting: font.HintingFull,
	}
	oo := &opentype.FaceOptions{
		Size:    1.0,
		DPI:     1200,
		Hinting: font.HintingFull,
	}
	face, err := ffff.FuzzyFindFace("Monaco", to, oo)
	if err != nil {
		return err
	}
	if isJapanese(fmt.Sprintf("%s%s", s1, s2)) {
		data, err := fs.ReadFile(jpFontFile, "ipag.ttf")
		if err != nil {
			return (err)
		}

		tt, err := opentype.Parse(data)
		if err != nil {
			return err
		}

		face, err = opentype.NewFace(tt, oo)
		if err != nil {
			return err
		}
	}

	h := 30

	img1Path := ""
	img2Path := ""
	if viper.GetBool("enable-output-images") {
		img1Path = fmt.Sprintf("%s.png", s1)
		img2Path = fmt.Sprintf("%s.png", s2)
	}
	img1 := generateImage(s1, len(s1)*15, h, face, img1Path)
	img2 := generateImage(s2, len(s2)*15, h, face, img2Path)
	hash1, err := generateImageHash(img1)
	if err != nil {
		return err
	}
	hash2, err := generateImageHash(img2)
	if err != nil {
		return err
	}
	hd, err := hash1.Distance(hash2)
	if err != nil {
		return err
	}
	fmt.Printf("image hash distance: %d\n", hd)
	return nil
}

func isJapanese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) {
			return true
		}
	}
	return false
}

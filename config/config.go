package config

import (
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
)

type Config struct {
	//ForceChrome means app will always use ChromeDP to get pages
	ForceChrome bool `arg:"--force-chrome, env:IMG_DW_FORCE_CHROME"`

	//Debug will force app to save downloaded HTML code into text files
	Debug bool `arg:"-d"`

	//MaxImgHeight will force app to auto-resize imgs to fit its value
	MaxImgHeight int `arg:"--max-height, env:IMG_DW_MAX_HEIGHT"`

	//MaxImgWidth will force app to auto-resize imgs to fit its value
	MaxImgWidth int `arg:"--max-width, env:IMG_DW_MAX_WIDTH"`

	//MinImgHeight will force app to skip imgs that don't fit into the value
	MinImgHeight int `arg:"--min-height, env:IMG_DW_MIN_HEIGHT"`

	//MinImgWidth will force app to skip imgs that don't fit into the value
	MinImgWidth int `arg:"--min-width, env:IMG_DW_MIN_WIDTH"`

	//ConvertAllToJPG will force app to convert all imgs to JPG
	ConvertAllToJPG bool `arg:"--all-to-jpg, env:IMG_DW_ALL_TO_JPG"`

	//ConvertAllToPNG will force app to convert all imgs to PNG
	ConvertAllToPNG bool `arg:"--all-to-png, env:IMG_DW_ALL_TO_PNG"`

	//OutputPath sets root for all downloads
	OutputPath string `arg:"--output-path, env:IMG_DW_OUTPUT_PATH"`

	//DoNotUseSubfolders will make app avoid creating subfolders for each downloaded page
	DoNotUseSubfolders bool `arg:"--no-subfolders, env:IMG_DW_NO_SUBFOLDERS"`
}

// InitApp parses config, sets defaults & creates root download dir
func InitApp() (Config, error) {
	var c Config
	err := arg.Parse(&c)
	if err != nil {
		return c, err
	}
	if c.OutputPath == "" {
		c.OutputPath = "image_downloader"
	}

	if err := os.MkdirAll(c.OutputPath, os.ModePerm); err != nil {
		return c, fmt.Errorf("can not make dir: %w", err)
	}

	return c, nil
}

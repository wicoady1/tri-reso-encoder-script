package operations

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/floostack/transcoder/ffmpeg"
	"github.com/wicoady1/tri-reso-encoder-script/config"
)

type Operations struct {
	conf *config.Config
}

func NewOperations(conf *config.Config) Operations {
	return Operations{
		conf: conf,
	}
}

func (op *Operations) EncodeVideoInBackground() {
	//check if input and output folder are exists
	//if none generate them
	if _, err := os.Stat("input"); os.IsNotExist(err) {
		err = os.Mkdir("input", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err = os.Mkdir("output", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	//check if file exists in input folder
	var files []string
	root := "input"
	err := filepath.Walk(root, op.visit(&files))
	if err != nil {
		panic(err)
	}

	//encode the file, save the filename first
	for _, file := range files {
		log.Println(file)
		op.encodeFile(file, file)

		//each successful video, delete them
		os.Remove("input/" + file)
	}

	//end
	time.Sleep(1 * time.Second)
}

func (op *Operations) visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if filepath.Ext(path) != ".mp4" {
			return nil
		}
		if info == nil {
			return nil
		}

		*files = append(*files, info.Name())
		return nil
	}
}

func (op *Operations) encodeFile(input string, output string) {
	format := "mp4"
	overwrite := true
	videoFilter := "scale=-1:720"
	videoCodec := "libx264"
	crf := uint32(18)
	preset := "veryslow"
	audioCodec := "copy"

	opts := ffmpeg.Options{
		OutputFormat: &format,
		Overwrite:    &overwrite,
		VideoFilter:  &videoFilter,
		VideoCodec:   &videoCodec,
		Crf:          &crf,
		Preset:       &preset,
		AudioCodec:   &audioCodec,
	}

	ffmpegConf := &ffmpeg.Config{
		FfmpegBinPath:   op.conf.FfmpegBinPath,
		FfprobeBinPath:  op.conf.FfprobeBinPath,
		ProgressEnabled: true,
	}

	progress, err := ffmpeg.
		New(ffmpegConf).
		Input("input/" + input).
		Output("output/" + output).
		WithOptions(opts).
		Start(opts)

	if err != nil {
		log.Fatal(err)
	}

	for msg := range progress {
		log.Printf("%+v", msg)
	}
}

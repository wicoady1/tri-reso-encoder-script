package operations

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/floostack/transcoder/ffmpeg"
	"github.com/wicoady1/tri-reso-encoder-script/config"
)

func EncodeVideoInBackground(conf *config.Config) {
	log.Println("masih jalan...")

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
	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	//encode the file, save the filename first
	for _, file := range files {
		log.Println(file)
		encodeFile(file, file)

		//each successful video, delete them
		os.Remove("input/" + file)
	}

	//end

	time.Sleep(1 * time.Second)
}

func visit(files *[]string) filepath.WalkFunc {
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

func encodeFile(input string, output string) {
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
		FfmpegBinPath:   "/usr/local/bin/ffmpeg",
		FfprobeBinPath:  "/usr/local/bin/ffprobe",
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

package operations

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/floostack/transcoder/ffmpeg"
	"github.com/wicoady1/tri-reso-encoder-script/config"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Operations struct {
	conf         *config.Config
	progressInfo *widgets.Paragraph
	progressBar  *widgets.Gauge
}

func NewOperations(conf *config.Config) Operations {

	progressBar := widgets.NewGauge()

	progressBar.Title = "Encoding Video..."
	progressBar.SetRect(0, 10, 50, 13)
	progressBar.BarColor = ui.ColorGreen
	progressBar.LabelStyle = ui.NewStyle(ui.ColorYellow)
	progressBar.BarColor = ui.ColorRed
	progressBar.BorderStyle.Fg = ui.ColorWhite
	progressBar.TitleStyle.Fg = ui.ColorCyan

	progressInfo := widgets.NewParagraph()
	progressInfo.SetRect(0, 13, 50, 50)
	progressInfo.Border = false

	return Operations{
		conf:         conf,
		progressInfo: progressInfo,
		progressBar:  progressBar,
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

	processedFiles := []string{}
	isProcess := false
	//encode the file, save the filename first
	for _, file := range files {
		isProcess = true
		for k := range op.conf.Resolution {
			op.encodeFile(file, file, op.conf.Resolution[k], op.conf.Bitrate[k])
		}

		//each successful video, delete them
		os.Remove("input/" + file)
		processedFiles = append(processedFiles, file)
	}

	if isProcess {
		op.progressInfo.Text = fmt.Sprintf("File\n%+v\nsuccessfully processed", processedFiles)
		ui.Render(op.progressInfo)

		op.progressBar.Label = "Encoding Completed"
		ui.Render(op.progressBar)
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

func (op *Operations) encodeFile(input string, output string, res int, bitrate string) {
	format := op.conf.Format
	overwrite := true
	videoFilter := fmt.Sprintf("scale=-1:%d", res)
	videoCodec := "libx264"
	crf := uint32(18)
	preset := "veryslow"

	opts := ffmpeg.Options{
		OutputFormat: &format,
		Overwrite:    &overwrite,
		VideoFilter:  &videoFilter,
		VideoCodec:   &videoCodec,
		Crf:          &crf,
		Preset:       &preset,
		VideoBitRate: &bitrate,
	}

	ffmpegConf := &ffmpeg.Config{
		FfmpegBinPath:   op.conf.FfmpegBinPath,
		FfprobeBinPath:  op.conf.FfprobeBinPath,
		ProgressEnabled: true,
	}

	progress, err := ffmpeg.
		New(ffmpegConf).
		Input("input/" + input).
		Output(fmt.Sprintf("output/%s_%d.%s", output, res, format)).
		WithOptions(opts).
		Start(opts)

	if err != nil {
		log.Print(err)
	}

	for msg := range progress {
		msg.GetCurrentBitrate()
		msg.GetCurrentTime()
		msg.GetFramesProcessed()
		msg.GetProgress()
		msg.GetSpeed()

		op.progressBar.Percent = int(msg.GetProgress())
		op.progressBar.Label = fmt.Sprintf("%.1f%%", msg.GetProgress())
		ui.Render(op.progressBar)

		op.progressInfo.Text = fmt.Sprintf(`File: %s
		Current Resolution: %d
		Current Bitrate: %s
		Current Time: %s
		Current Frames: %s
		Current Progress: %.2f%%
		Current Speed: %s`, input, res, msg.GetCurrentBitrate(), msg.GetCurrentTime(), msg.GetFramesProcessed(), msg.GetProgress(), msg.GetSpeed())
		ui.Render(op.progressInfo)
	}
}

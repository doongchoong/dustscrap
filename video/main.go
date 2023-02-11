package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const file_fmt = "img%04d.png"

type RenameInfo struct {
	Start string   `json:"start"`
	End   string   `json:"end"`
	Files []string `json:"files"`
}

func RenameImageSequence(path, start, end string) error {
	// set start ~ end range
	if start <= "0000_00_00_0000Z.png" || start > "9999_99_99_9999Z.png" {
		start = "0000_00_00_0000Z.png"
	}
	if end > "9999_99_99_9999Z.png" || end <= "0000_00_00_0000Z.png" {
		end = "9999_99_99_9999Z.png"
	}
	// walk path
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	var strs []string

	for _, file := range files {
		strs = append(strs, file.Name())
		if strings.HasPrefix(file.Name(), "img") {
			return errors.New("Need to Rollback ImageSequence names")
		}
	}
	sort.Strings(strs)

	start_idx := 0
	end_idx := 0

	for i, v := range strs {
		if v < start {
			start_idx = i
		}
		if v < end {
			end_idx = i
		}
	}
	start = strs[start_idx]
	end = strs[end_idx]
	renameInfo := &RenameInfo{start, end, strs[start_idx:(end_idx + 1)]}

	b, err := json.Marshal(renameInfo)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(path, "renameinfo.json"))
	if err != nil {
		return (err)
	}
	_, err = f.Write(b)
	if err != nil {
		return (err)
	}
	f.Close()

	// change names
	for i, v := range renameInfo.Files {
		err = os.Rename(
			filepath.Join(path, v),
			filepath.Join(path, fmt.Sprintf(file_fmt, i)),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func RollbackImageSequence(path string) error {
	// check frames folder
	if _, err := os.Stat(filepath.Join(path, "renameinfo.json")); os.IsNotExist(err) {
		return nil
	}

	//open config file
	f, err := os.Open(filepath.Join(path, "renameinfo.json"))
	if err != nil {
		return err
	}

	byt, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	renameInfo := &RenameInfo{}

	//json
	err = json.Unmarshal(byt, &renameInfo)
	if err != nil {
		return err
	}

	// json file close
	f.Close()

	// change names
	for i, v := range renameInfo.Files {
		src := filepath.Join(path, fmt.Sprintf("img%04d.png", i))
		dst := filepath.Join(path, v)
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}
		err = os.Rename(
			src,
			dst,
		)
		if err != nil {
			return err
		}
	}

	// delete json
	err = os.Remove(filepath.Join(path, "renameinfo.json"))
	if err != nil {
		return err
	}

	return nil
}

func GenVideoByFFMPEG(outputf, binpath, imgpath string, quality, scale int) error {

	var cmd *exec.Cmd

	ext := filepath.Ext(outputf)

	if ext == ".avi" {
		cmd = exec.Command(
			binpath,
			"-f", "image2",
			"-r", "24",
			"-i", filepath.Join(imgpath, file_fmt),
			"-q:v", fmt.Sprintf("%d", quality), //  1 (lossless), 4 (quard
			//"-vf", fmt.Sprintf("scale=-1:%d", scale),
			outputf,
		)
	} else if ext == ".gif" {
		cmd = exec.Command(
			binpath,
			"-f", "image2",
			"-r", "24",
			"-i", filepath.Join(imgpath, file_fmt),
			outputf,
		)
	} else {
		log.Println("!!! output file avi or gif ")
		return nil
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	fmt.Println("start")
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("done")
	return nil
}

func main() {
	imgpath := flag.String("imgpath", "frames", "image sequence directory")
	binpath := flag.String("binpath", "tools/ffmpeg", "ffmpeg path")
	quality := flag.Int("quality", 4, "1:lossless")
	scale := flag.Int("scale", 720, "height scale")
	outputf := flag.String("out", "out.avi", "output file (avi or gif)")
	start := flag.String("start", "", "start image file ")
	end := flag.String("end", "", "end image file")

	flag.Parse()

	var err error
	err = RollbackImageSequence(*imgpath)
	if err != nil {
		panic(err)
	}
	err = RenameImageSequence(*imgpath, *start, *end)
	if err != nil {
		panic(err)
	}
	err = GenVideoByFFMPEG(*outputf, *binpath, *imgpath, *quality, *scale)
	if err != nil {
		panic(err)
	}
	err = RollbackImageSequence(*imgpath)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	// 컨픽 json 파일 로드
	configfnm := flag.String("config", "config.json", "Input config json file")
	flag.Parse()

	conf := &Config{}
	err := conf.Load(*configfnm)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Load config => " + *configfnm)

	// URLs 생성
	urls, err := GenerateUrls(conf)
	if err != nil {
		log.Fatal(err)
	}
	// Frame 폴더 체크
	if _, err := os.Stat(conf.Frames_path); os.IsNotExist(err) {
		log.Fatal("frames directory not exist (location of saved capture images)")
	}

	// 시작
	Scraping(urls, conf)
}

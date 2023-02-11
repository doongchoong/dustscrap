package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// 정규식으로 URL 정합성 판단
func CheckFromToUrl(conf *Config) error {
	pat := `^#\d\d\d\d/(0\d|10|11|12)/([012]\d|30|31)/([01]\d|20|21|22|23)00Z$`
	isMat, err := regexp.MatchString(pat, conf.From_url)
	if err != nil {
		return err
	}
	if isMat == false {
		return errors.New("from_url format error: ex) #2019/02/23/0000Z")
	}
	isMat, err = regexp.MatchString(pat, conf.To_url)
	if err != nil {
		return err
	}
	if isMat == false {
		return errors.New("to_url format error: ex) #2019/02/23/2300Z")
	}
	fmt.Println("from to url match cehk")
	return nil
}

// Url 생성
func GenerateUrls(conf *Config) ([]string, error) {
	var ret []string
	// check from -to pattern
	err := CheckFromToUrl(conf)
	if err != nil {
		return nil, err
	}
	// get real time from from-to strings
	fromTim, err := time.Parse(time.RFC3339,
		fmt.Sprintf("%s-%s-%sT%s:00:00+00:00",
			conf.From_url[1:5],
			conf.From_url[6:8],
			conf.From_url[9:11],
			conf.From_url[12:14]))
	if err != nil {
		return nil, err
	}
	toTim, err := time.Parse(time.RFC3339,
		fmt.Sprintf("%s-%s-%sT%s:00:00+00:00",
			conf.To_url[1:5],
			conf.To_url[6:8],
			conf.To_url[9:11],
			conf.To_url[12:14]))
	if err != nil {
		return nil, err
	}

	// generate url
	tt := fromTim
	url := fmt.Sprintf("https://earth.nullschool.net/%s/%s/%s",
		fmt.Sprintf("#%d/%02d/%02d/%02d00Z", tt.Year(), tt.Month(), tt.Day(), tt.Hour()),
		conf.Mode_url,
		conf.Projection_url,
	)
	ret = append(ret, url)
	for {
		tt = tt.Add(time.Hour)
		if tt.After(toTim) {
			break
		}
		url = fmt.Sprintf("https://earth.nullschool.net/%s/%s/%s",
			fmt.Sprintf("#%d/%02d/%02d/%02d00Z", tt.Year(), tt.Month(), tt.Day(), tt.Hour()),
			conf.Mode_url,
			conf.Projection_url,
		)
		ret = append(ret, url)
	}

	return ret, nil
}

func Scraping(urls []string, conf *Config) {

	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// 테스트용 크롬 인스턴스 강제 수행시간 지정
	//ctx, cancel = context.WithTimeout(ctx, 25*time.Second)
	//defer cancel()

	// URL 순회
	for _, url := range urls {

		log.Print("Try capture :" + url)

		// url 이동
		if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
			log.Fatal(err)
		}
		// 렌더링을 기다린다.
		for cnt := 0; cnt < 500; cnt++ {
			var res string
			if err := chromedp.Run(ctx, chromedp.Text(".field", &res, chromedp.ByQuery)); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Millisecond)
			if strings.Contains(res, "Download") {
				continue
			}
			if strings.Contains(res, "Rendering") {
				continue
			}
			break
		}

		// 스크린샷을 찍는다.
		var buf []byte

		err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 100))
		if err != nil {
			log.Fatal(err)
		}

		// 파일명 생성
		fnm := strings.Replace(strings.Replace(url[30:46], "/", "_", -1), "#", "_", -1)
		savednm := filepath.Join(conf.Frames_path, fmt.Sprintf("%s.png", fnm))

		// 타임스탬프
		if conf.Time_stamp {
			buf, err = AddTimeStamp(buf, 30, 30, 24, image.White, fnm)
			if err != nil {
				log.Fatal(err)
			}
		}

		// 파일저장
		if err := os.WriteFile(savednm, buf, 0o644); err != nil {
			log.Fatal(err)
		}

		log.Println("...Done.")

		// 컨픽 설정한 초 만큼 기다린다.
		time.Sleep(time.Duration(conf.Wait_seconds) * time.Second)
	}

}

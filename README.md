# dustscrap

[어스널스쿨](https://earth.nullschool.net/)의 화면을 연속적으로 캡쳐하여 
미세먼지 확산 영상을 제작하는 프로그램 입니다.
Go 언어로 작성되었습니다.

## Scrap

json파일 형태의 컨픽파일을 설정한 뒤 
프로그램을 실행하여 캡쳐합니다. 

```json
{
	"from_url" : "#2023/01/05/0100Z",
	"to_url" : "#2023/02/09/0000Z",
	"mode_url" : "particulates/surface/level/anim=off/overlay=pm2.5",
	"projection_url" : "orthographic=-234.04,36.14,1377",
	"wait_seconds": 1,
	"frames_path" : "../frm_pm25",
	"time_stamp": true
}
```
pm2.5 미세먼지 모드에서 캡쳐하기 위한 json설정 입니다.



```sh
#실행 예
$ scrap --config=pm25.json
```


## Video

여러장의 이미지를 FFmpeg 를 이용하여 동영상을 만들기 위한 
프로그램 입니다. 

```sh
$ video --imgpath=../frm_pm25 --start=2023_01_01_0000Z.png --end=2023_02_09_0000Z.png
```
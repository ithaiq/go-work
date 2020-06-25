package image_merge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type UploadImgRsp struct {
	RetCode string `json:"RetCode"`
	RetMsg  string `json:"RetMsg"`
}

func MergeGroupLogo(userPicList []string) (error, string) {

	var src []io.Reader

	for _, pic := range userPicList {
		if data, err := DownLoadUserPic(pic); err == nil {
			src = append(src, bytes.NewReader(data))
		}
	}
	fi := bytes.NewBuffer([]byte{})

	if err := Merge(fi, src); err != nil {
		return err, ""
	}
	err, url := UpLoadImgToGroupLogo(fi)
	if err != nil {
		return err, ""
	}

	return nil, url
}
func DownLoadUserPic(downUrl string) ([]byte, error) {
	res, err := http.Get(downUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("DownLoadUserPic has error:%v", err)
		return nil, err
	}
	return bytes, nil
}
func UpLoadImgToGroupLogo(buf *bytes.Buffer) (error, string) {
	// 上传文件POST
	// 下面构造一个文件buf作为POST的BODY
	newBuf := new(bytes.Buffer)
	w := multipart.NewWriter(newBuf)
	fw, _ := w.CreateFormFile("file", "file.jpg") //这里的uploadFile必须和服务器端的FormFile-name一致
	io.Copy(fw, buf)
	w.Close()

	resp, err := http.Post("", w.FormDataContentType(), newBuf)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, ""
	}
	//example
	rsp := &UploadImgRsp{}
	if err := json.Unmarshal(body, rsp); err != nil {
		return err, ""
	} else if rsp.RetCode != "0" {
		return errors.New("UpLoadImgToGroupLogo fail"), ""
	}
	return nil, rsp.RetMsg
}

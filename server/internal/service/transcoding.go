package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/spf13/viper"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/domain/model"
	"interastral-peace.com/alnitak/internal/global"
	"interastral-peace.com/alnitak/utils"
)

type TranscodingTarget struct {
	Resolution  string // 分辨率
	BitrateRate string // 码率
	FPS         string // 帧率
}

// 生成封面
func GenerateCover(inputFile, outputFile string) error {
	command := []string{"-i", inputFile, "-vframes", "1", "-y", outputFile}

	_, err := utils.RunCmd(exec.Command("ffmpeg", command...))
	if err != nil {
		utils.ErrorLog("提取封面失败", "transcoding", err.Error())
		return err
	}

	return nil
}

// 获取视频信息
func ProcessVideoInfo(input string) (*dto.TranscodingInfo, error) {
	var transcodingInfo dto.TranscodingInfo
	videoData, err := getVideoInfo(input)
	if err != nil {
		utils.ErrorLog("读取视频信息失败", "transcoding", err.Error())
		return &transcodingInfo, err
	}

	if videoData.Stream[0].CodecName != "h264" {
		utils.ErrorLog("视频编码不为h264", "transcoding", "")
		return &transcodingInfo, errors.New("not h264")
	}

	//计算最大分辨率
	transcodingInfo.Width = videoData.Stream[0].Width
	transcodingInfo.Height = videoData.Stream[0].Height

	//获取视频时长
	transcodingInfo.Duration, _ = strconv.ParseFloat(videoData.Stream[0].Duration, 64)

	return &transcodingInfo, err
}

func VideoTransCoding(transcodingInfo *dto.TranscodingInfo) {
	var wg sync.WaitGroup
	targets := getTranscodingTarget(transcodingInfo)
	wg.Add(len(targets))
	for _, v := range targets {
		c := v // 处理协程引用循环变量问题
		go func() {
			fileName := c.Resolution + "_" + c.BitrateRate + "_" + c.FPS
			tsFileName := transcodingInfo.OutputDir + fileName + ".ts"
			// 压缩视频
			err := pressingVideo(transcodingInfo.InputFile, tsFileName, c.Resolution)
			if err != nil {
				wg.Done()
				return
			}
			// 切片
			m3u8File, err := generateVideoSlices(tsFileName, transcodingInfo.OutputDir, fileName)
			if err != nil {
				wg.Done()
				return
			}
			// 读取m3u8写入数据库
			err = saveM3u8File(transcodingInfo, fileName, m3u8File)
			if err != nil {
				wg.Done()
				return
			}

			//删除临时文件
			os.Remove(tsFileName)
			os.Remove(m3u8File)

			wg.Done()
		}()
	}

	wg.Wait()

	// 上传oss
	if viper.GetString("storage.oss_type") != "local" {
		fmt.Println("开始上传OSS")
		files, err := os.ReadDir(transcodingInfo.OutputDir)
		if err != nil {
			utils.ErrorLog("读取视频文件夹失败", "oss", err.Error())
			completeTransCoding(transcodingInfo.VideoID, transcodingInfo.ResourceID, global.PROCESSING_FAIL)
			return
		}

		for _, f := range files {
			fmt.Println("上传", f.Name())

			if f.Name() == "upload.mp4" && !viper.GetBool("storage.upload_mp4_file") {
				continue
			}

			objectKey := "video/" + transcodingInfo.DirName + "/" + f.Name()
			filePath := "./upload/" + objectKey
			if err := global.Storage.PutObjectFromFile(objectKey, filePath); err != nil {
				utils.ErrorLog("文件上传OSS失败", "oss", err.Error())
			}
		}
	}

	// 更新状态
	completeTransCoding(transcodingInfo.VideoID, transcodingInfo.ResourceID, global.WAITING_REVIEW)
}

// 获取宽度支持的最大分辨率
func getWidthRes(width int) int {
	//1920*1080
	if width >= 1920 {
		return 1080
	}
	// 1280*720
	if width >= 1280 {
		return 720
	}
	// 720*480
	if width >= 720 {
		return 480
	}
	return 360
}

// 获取高度支持的最大分辨率
func getHeigthRes(height int) int {
	//1920*1080
	if height >= 1080 {
		return 1080
	}
	// 1280*720
	if height >= 720 {
		return 720
	}
	// 720*480
	if height >= 480 {
		return 480
	}
	return 360
}

// 获取转码目标
func getTranscodingTarget(videoInfo *dto.TranscodingInfo) []TranscodingTarget {
    targets := make([]TranscodingTarget, 0)
    maxRresolution := utils.Max(getWidthRes(videoInfo.Width), getHeigthRes(videoInfo.Height))

    if maxRresolution >= 1080 {
        targets = append(targets, TranscodingTarget{Resolution: "1920x1080", FPS: "30"})
    } else if maxRresolution >= 720 {
        targets = append(targets, TranscodingTarget{Resolution: "1080x720", FPS: "30"})
    } else if maxRresolution >= 480 {
        targets = append(targets, TranscodingTarget{Resolution: "854x480", FPS: "30"})
    } else if maxRresolution >= 360 {
        targets = append(targets, TranscodingTarget{Resolution: "640x360", FPS: "30"})
    }

    return targets
}

// 获取视频信息
func getVideoInfo(input string) (info global.VideoInfo, err error) {
	cmd := exec.Command("ffprobe", "-i", input, "-v", "quiet", "-print_format", "json", "-show_format", "-show_streams")
	out, err := utils.RunCmd(cmd)
	if err != nil {
		return info, err
	}

	// 反序列化
	err = json.Unmarshal(out.Bytes(), &info)
	if err != nil {
		return info, err
	}

	return info, nil
}

// 压缩视频

func pressingVideo(inputFile, outputFile, quality string) error {
command := []string{"-i", inputFile, "-crf", "23", "-s", quality,"-c:v", "libx264", "-r", "30", "-preset", "fast","-c:a", "aac", "-f", "mpegts", outputFile,
	}

	_, err := utils.RunCmd(exec.Command("ffmpeg", command...))
	if err != nil {
		utils.ErrorLog(
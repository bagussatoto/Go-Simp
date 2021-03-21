package space

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	database "github.com/JustHumanz/Go-Simp/pkg/database"
	network "github.com/JustHumanz/Go-Simp/pkg/network"

	log "github.com/sirupsen/logrus"
)

func (Space *CheckSctruct) Check(limit string) *CheckSctruct {
	var (
		Videotype string
		PushVideo SpaceVideo
		NewVideo  Vlist
	)
	body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/space/arc/search?mid="+strconv.Itoa(Space.SpaceID)+"&ps="+limit, nil)
	if curlerr != nil {
		log.Error(curlerr)
	}

	err := json.Unmarshal(body, &PushVideo)
	if err != nil {
		log.Error(err)
	}

	for _, video := range PushVideo.Data.List.Vlist {
		if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|翻唱|mv|歌曲)", strings.ToLower(video.Title)); Cover || video.Typeid == 31 {
			Videotype = "Covering"
		} else {
			Videotype = "Streaming"
		}

		Data := database.LiveStream{
			VideoID: video.Bvid,
			Type:    Videotype,
			Title:   video.Title,
			Thumb:   "https:" + video.Pic,
			Desc:    video.Description,
			Schedul: time.Unix(int64(video.Created), 0).In(loc),
			Viewers: strconv.Itoa(video.Play),
			Length:  video.Length,
			Member:  Space.Member,
		}
		new, id := Data.CheckVideo()
		if new {
			Data.InputSpaceVideo()
			video.Pic = "https:" + video.Pic
			video.VideoType = Videotype
			NewVideo = append(NewVideo, video)
		} else {
			Data.UpdateView(id)
		}
	}
	Space.VideoList = NewVideo
	return Space
}

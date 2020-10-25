package database

import (
	log "github.com/sirupsen/logrus"
)

//Get Youtube data from status
func YtGetStatus(Group, Member int64, Status string) []YtDbData {
	var (
		Data  []YtDbData
		list  YtDbData
		limit int
	)
	if Group != 0 && Status != "live" {
		limit = 3
	} else {
		limit = 2525
	}
	rows, err := DB.Query(`call GetYt(?,?,?,?)`, Member, Group, limit, Status)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&list.ID, &list.Group, &list.ChannelID, &list.NameEN, &list.NameJP, &list.YoutubeAvatar, &list.VideoID, &list.Title, &list.Thumb, &list.Desc, &list.Schedul, &list.End, &list.Region, &list.Viewers)
		if err != nil {
			log.Error(err)
		}
		Data = append(Data, list)
	}
	return Data

}

//Input youtube new video
func (Data YtDbData) InputYt(MemberID int64) {
	stmt, err := DB.Prepare(`INSERT INTO Youtube (VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers,Length,VtuberMember_id) values(?,?,?,?,?,?,?,?,?,?,?,?)`)
	BruhMoment(err, "", false)
	defer stmt.Close()

	res, err := stmt.Exec(Data.VideoID, Data.Type, Data.Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, MemberID)
	BruhMoment(err, "", false)

	_, err = res.LastInsertId()
	BruhMoment(err, "", false)

}

//Check new video or not
func CheckVideoID(VideoID string) YtDbData {
	var Data YtDbData
	rows, err := DB.Query(`SELECT id,VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers FROM Vtuber.Youtube Where VideoID=?;`, VideoID)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.VideoID, &Data.Type, &Data.Status, &Data.Title, &Data.Thumb, &Data.Desc, &Data.Published, &Data.Schedul, &Data.End, &Data.Viewers)
		BruhMoment(err, "", false)
	}
	if Data.ID == 0 {
		return YtDbData{}
	} else {
		return Data
	}
}

//Update youtube data
func (Data YtDbData) UpdateYt(Status string) {
	_, err := DB.Exec(`Update Youtube set Type=?,Status=?,Title=?,Thumbnails=?,Description=?,PublishedAt=?,ScheduledStart=?,EndStream=?,Viewers=?,Length=? where id=?`, Data.Type, Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, Data.ID)
	BruhMoment(err, "", false)
}

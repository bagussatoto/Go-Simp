package main

import (
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

//BiliBiliMessage message handler
func BiliBiliMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := configfile.BotPrefix.Bilibili
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, Prefix) {
		loc, _ := time.LoadLocation("Asia/Shanghai") /*Use CST*/
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 1 {
			Payload := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
			if CommandArray[0] == Prefix+Live {
				for _, FindGroupArry := range Payload {
					VTuberGroup, err := FindGropName(FindGroupArry)
					if err != nil {
						Member := FindVtuber(FindGroupArry, 0)
						if Member == (database.Member{}) {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
							return
						} else {
							LiveBili := database.BilGet(0, Member.ID, "Live")
							FixName := engine.FixName(Member.EnName, Member.JpName)
							if len(LiveBili) > 0 {
								Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}
								for _, LiveData := range LiveBili {
									diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
									view, err := strconv.Atoi(LiveData.Viewers)
									if err != nil {
										log.Error(err)
									}
									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetTitle(FixName).
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription(LiveData.Desc).
										SetThumbnail(Member.BiliBiliAvatar).
										SetImage(LiveData.Thumb).
										SetURL("https://live.bilibili.com/"+strconv.Itoa(Member.BiliRoomID)).
										AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
										AddField("Online", engine.NearestThousandFormat(float64(view))).
										SetColor(Color).
										SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("It looks like `"+FixName+"` doesn't have a livestream right now").
									SetImage(config.WorryIMG).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						LiveBili := database.BilGet(VTuberGroup.ID, 0, "Live")
						if len(LiveBili) > 0 {
							Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
							if err != nil {
								log.Error(err)
							}

							for _, LiveData := range LiveBili {
								LiveData.AddMember(FindVtuber("", LiveData.Member.ID))
								FixName := engine.FixName(LiveData.Member.EnName, LiveData.Member.JpName)
								diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
								view, err := strconv.Atoi(LiveData.Viewers)
								if err != nil {
									log.Error(err)
								}
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription(LiveData.Desc).
									SetThumbnail(LiveData.Member.BiliBiliAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://live.bilibili.com/"+strconv.Itoa(LiveData.Member.BiliRoomID)).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Online", engine.NearestThousandFormat(float64(view))).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("It looks like `"+VTuberGroup.GroupName+"` doesn't have a livestream right now").
								SetImage(config.WorryIMG).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}

			} else if CommandArray[0] == Prefix+Past || CommandArray[0] == Prefix+"last" {
				for _, FindGroupArry := range Payload {
					VTuberGroup, err := FindGropName(FindGroupArry)
					if err != nil {
						Member := FindVtuber(FindGroupArry, 0)
						if Member == (database.Member{}) {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
							return
						} else {
							LiveBili := database.BilGet(0, Member.ID, "Past")
							FixName := engine.FixName(Member.JpName, Member.EnName)
							if len(LiveBili) > 0 {
								Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}

								for _, LiveData := range LiveBili {
									diff := LiveData.Schedul.In(loc).Sub(time.Now().In(loc))
									view, err := strconv.Atoi(LiveData.Viewers)
									if err != nil {
										log.Error(err)
									}
									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetTitle(FixName).
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription(LiveData.Desc).
										SetThumbnail(Member.BiliBiliAvatar).
										SetImage(LiveData.Thumb).
										SetURL("https://live.bilibili.com/"+strconv.Itoa(Member.BiliRoomID)).
										AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
										AddField("Online", engine.NearestThousandFormat(float64(view))).
										SetColor(Color).
										SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("It looks like `"+FixName+"` doesn't have a Past livestream right now").
									SetImage(config.WorryIMG).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						LiveBili := database.BilGet(VTuberGroup.ID, 0, "Past")
						if len(LiveBili) > 0 {
							Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
							if err != nil {
								log.Error(err)
							}

							for _, LiveData := range LiveBili {
								LiveData.AddMember(FindVtuber("", LiveData.Member.ID))
								FixName := engine.FixName(LiveData.Member.EnName, LiveData.Member.JpName)
								view, err := strconv.Atoi(LiveData.Viewers)
								if err != nil {
									log.Error(err)
								}
								diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription(LiveData.Desc).
									SetThumbnail(LiveData.Member.BiliBiliAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://live.bilibili.com/"+strconv.Itoa(LiveData.Member.BiliRoomID)).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Online", engine.NearestThousandFormat(float64(view))).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("It looks like `"+VTuberGroup.GroupName+"` doesn't have a Past livestream right now").
								SetImage(config.WorryIMG).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Prefix+"` command")
			if err != nil {
				log.Error(err)
			}
		}
	}
}

//BiliBiliSpace message hadler
func BiliBiliSpace(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := "sp_" + configfile.BotPrefix.Bilibili
	if strings.HasPrefix(m.Content, Prefix) {
		loc, _ := time.LoadLocation("Asia/Shanghai") /*Use CST*/
		Payload := m.Content[len(Prefix):]
		if Payload != "" {
			for _, FindGroupArry := range strings.Split(strings.TrimSpace(Payload), ",") {
				VTuberGroup, err := FindGropName(FindGroupArry)
				if err != nil {
					Member := FindVtuber(FindGroupArry, 0)
					if Member == (database.Member{}) {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
						return
					} else {
						SpaceBili := database.SpaceGet(0, Member.ID)
						FixName := engine.FixName(Member.EnName, Member.JpName)
						if len(SpaceBili) > 0 {
							Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
							if err != nil {
								log.Error(err)
							}

							for _, SpaceData := range SpaceBili {
								diff := time.Now().In(loc).Sub(SpaceData.Schedul.In(loc))
								view, err := strconv.Atoi(SpaceData.Viewers)
								if err != nil {
									log.Error(err)
								}
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetTitle(FixName).
									SetDescription(SpaceData.Title).
									SetImage(SpaceData.Thumb).
									SetThumbnail(Member.BiliBiliAvatar).
									SetURL("https://www.bilibili.com/video/"+SpaceData.VideoID).
									AddField("Type", SpaceData.Type).
									AddField("Video uploaded", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Duration", SpaceData.Length).
									AddField("Viewers now", engine.NearestThousandFormat(float64(view))).
									SetFooter(SpaceData.Schedul.In(loc).Format(time.RFC822), config.BiliBiliIMG).
									InlineAllFields().
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("It looks like `"+FixName+"` doesn't have a video in space.bilibili").
								SetImage(config.WorryIMG).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}
					}
					break
				} else {
					SpaceBili := database.SpaceGet(VTuberGroup.ID, 0)
					if len(SpaceBili) > 0 {
						Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
						if err != nil {
							log.Error(err)
						}

						for _, SpaceData := range SpaceBili {
							SpaceData.AddMember(FindVtuber("", SpaceData.Member.ID))

							FixName := engine.FixName(SpaceData.Member.EnName, SpaceData.Member.JpName)
							diff := time.Now().In(loc).Sub(SpaceData.Schedul.In(loc))
							view, err := strconv.Atoi(SpaceData.Viewers)
							if err != nil {
								log.Error(err)
							}
							_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetTitle(FixName).
								SetDescription(SpaceData.Title).
								SetImage(SpaceData.Thumb).
								SetThumbnail(SpaceData.Member.BiliBiliAvatar).
								SetURL("https://www.bilibili.com/video/"+SpaceData.VideoID).
								AddField("Type", SpaceData.Type).
								AddField("Video uploaded", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
								AddField("Duration", SpaceData.Length).
								AddField("Viewers now", engine.NearestThousandFormat(float64(view))).
								SetFooter(SpaceData.Schedul.In(loc).Format(time.RFC822), config.BiliBiliIMG).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}

						}
					} else {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("It looks like `"+VTuberGroup.GroupName+"` doesn't have a video in space.bilibili").
							SetImage(config.WorryIMG).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						return
					}
				}
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Prefix+"` command")
			if err != nil {
				log.Error(err)
			}
		}
	}
}

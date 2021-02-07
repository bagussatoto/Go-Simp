package space

import (
	"strconv"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	log "github.com/sirupsen/logrus"
)

var (
	loc *time.Location
	Bot *discordgo.Session
)

//Start start twitter module
func Start(BotInit *discordgo.Session, cronInit *cron.Cron) {
	loc, _ = time.LoadLocation("Asia/Shanghai") /*Use CST*/
	Bot = BotInit
	cronInit.AddFunc(config.BiliBiliSpace, CheckSpaceVideo)
	log.Info("Enable space bilibili module")
}

func CheckSpaceVideo() {
	for _, GroupData := range engine.GroupData {
		if GroupData.GroupName != "Hololive" {
			wg := new(sync.WaitGroup)
			for i, MemberData := range database.GetMembers(GroupData.ID) {
				wg.Add(1)
				go func(Group database.Group, Member database.Member, wg *sync.WaitGroup) {
					defer wg.Done()
					if Member.BiliBiliID != 0 {
						log.WithFields(log.Fields{
							"Group":      Group.GroupName,
							"Vtuber":     Member.EnName,
							"BiliBiliID": Member.BiliBiliID,
						}).Info("Checking Space BiliBili")

						if Group.GroupName == "Independen" {
							Group.IconURL = ""
						}

						Data := &CheckSctruct{
							Member: Member,
							Group:  Group,
						}
						Data.Check(strconv.Itoa(config.BotConf.LimitConf.SpaceBiliBili)).SendNude()

					}
				}(GroupData, MemberData, wg)
				if i%5 == 0 {
					wg.Wait()
				}
			}
			wg.Wait()
		}
	}
}
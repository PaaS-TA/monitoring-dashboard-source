package service

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"kr/paasta/monitoring-batch/dao"
	mod "kr/paasta/monitoring-batch/model"
	cb "kr/paasta/monitoring-batch/model/base"
	"kr/paasta/monitoring-batch/util"
	"sync"
)

func UpdateSnsAlarmTarget(f *BackendServices) {

	alarmSns, err := dao.GetCommonDao(f.Influxclient).GetAlarmSns("", f.MonitoringDbClient)
	if err != nil {
		fmt.Println("Failed to get sns_id(ChatRoomId)! :", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(alarmSns))
	for _, v := range alarmSns {
		go func(wg *sync.WaitGroup, v mod.AlarmSns) {
			defer wg.Done()

			if v.SnsType == cb.SNS_TYPE_TELEGRAM {
				bot, err := tgbotapi.NewBotAPI(v.Token)
				if err != nil {
					fmt.Println(err)
				} else {
					bot.Debug = true
					var updateConfig tgbotapi.UpdateConfig
					updateConfig.Offset = 0
					updateConfig.Timeout = 30
					updates, err := bot.GetUpdates(updateConfig)
					if err != nil {
						fmt.Println(err)
					} else {
						var chatIdList []int64
						for _, update := range updates {
							if update.Message == nil {
								continue
							}
							chatIdList = append(chatIdList, update.Message.Chat.ID)
						}
						chatIdList = util.RemoveDuplicates(chatIdList)
						for _, chatId := range chatIdList {
							var alarmSnsTarget mod.AlarmSnsTarget
							alarmSnsTarget.ChannelId = v.ChannelId
							alarmSnsTarget.TargetId = chatId
							err := dao.GetCommonDao(f.Influxclient).UpdateSnsAlarmTargets(alarmSnsTarget, f.MonitoringDbClient)
							if err != nil {
								return
							}
						}
					}
				}
			}
		}(&wg, v)
	}
	wg.Wait()
}

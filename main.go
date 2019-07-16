// bot_telega project main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron"
)

var db *sql.DB

func main() {
	var GlobalUpdate = 0
	os.Setenv("HTTP_PROXY", "socks5://streamfest:Stream.Fest1@vpnwork.space:10800")
	var err error
	db, err = sql.Open("mysql", "root:ilovegg@/psn")

	if err != nil {
		log.Println(err.Error())
	}

	bot, err := tgbotapi.NewBotAPI("833358732:AAGq9fXr5KjapQqinHBKLx1Oqr2kfsfWjY0")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	c := cron.New()
	//c.AddFunc("0 */5 * * * *", func() {
	//	var onlieId string
	//		rows, err := db.Query("SELECT onlineId FROM `users`")
	//	for rows.Next() {
	//		err = rows.Scan(&onlieId)
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		log.Println(onlieId)
	//		updateBase, err := exec.Command("/root/psn_bot/PSNbot", "-username", onlieId, "-fast", "yes").Output()
	//		log.Println(string(updateBase))
	//		if err != nil {
	//			log.Println(string(updateBase), fmt.Sprint(err))
	//		}
	//	}
	//	log.Println("updateDone")
	//})
	c.AddFunc("0 */38 * * * *", func() {
		GlobalUpdate = 1
		var onlieId string
		rows, err := db.Query("SELECT onlineId FROM `users`")
		for rows.Next() {
			err = rows.Scan(&onlieId)
			if err != nil {
				log.Println(err)
			}
			log.Println(onlieId)
			updateBase, err := exec.Command("/root/psn_bot/PSNbot", "-username", onlieId, "-fast", "no").Output()
			log.Println(string(updateBase))
			if err != nil {
				log.Println(string(updateBase), fmt.Sprint(err))
			}
		}
		log.Println("updateDone")
		GlobalUpdate = 0
	})
	c.AddFunc("0 */5 * * * *", func() {
		fmt.Println("Global: ", GlobalUpdate)
		if GlobalUpdate == 0 {

			var onlieId string
			rows, err := db.Query("SELECT onlineId FROM `users`")
			for rows.Next() {
				err = rows.Scan(&onlieId)
				if err != nil {
					log.Println(err)
				}
				log.Println(onlieId)
				updateBase, err := exec.Command("/root/psn_bot/PSNbot", "-username", onlieId, "-fast", "yes").Output()
				log.Println(string(updateBase))
				if err != nil {
					log.Println(string(updateBase), fmt.Sprint(err))
				}
			}
			log.Println("updateDone")
			var gameId string
			var trophyName string
			var trophyType string
			var trophyDetail string
			var trophyIcon string
			var onlineId string
			var gameName string
			var earnrate string
			rows, err = db.Query("SELECT gameId,trophyName,trophyType,trophyDetail,trophyIcon,onlineId,earnrate FROM `trophyqueue` WHERE send = 0 ORDER BY trophyType = 'platinum'")
			for rows.Next() {
				err = rows.Scan(&gameId, &trophyName, &trophyType, &trophyDetail, &trophyIcon, &onlineId, &earnrate)
				if err != nil {
					log.Println(err)
				}
				rows2, _ := db.Query("SELECT name FROM `games` WHERE NpCommunicationId = '" + gameId + "'")
				for rows2.Next() {
					err = rows2.Scan(&gameName)
				}
				log.Println(gameId, trophyName, trophyType, trophyDetail, trophyIcon, onlineId)
				chat_row, _ := db.Query("SELECT * FROM `registredchats`")
				for chat_row.Next() {
					var chat_id int64
					err = chat_row.Scan(&chat_id)
					if err != nil {
						log.Println(err)
					}
					//msg2 := tgbotapi.NewMessage(chat_id, "  /top"+gameId)
					trophyUrl := "http://storage1.goodgame.ru/psn/troph.php?usrname=" + onlineId + "&game=" + gameName + "&detail=" + url.QueryEscape(escape2(trophyDetail)) + "&img=" + trophyIcon + "&type=" + trophyType + "&title=" + url.QueryEscape(escape2(trophyName)) + "&earnrate=" + url.QueryEscape(earnrate)
					imgMsg := tgbotapi.NewMessage(chat_id, ` @`+onlineId+` — <a href="`+trophyUrl+`"> `+gameName+`</a>`)
					imgMsg.ParseMode = "HTML"
					bot.Send(imgMsg)
					if trophyType == "platinum" {
						bot.Send(tgbotapi.NewStickerShare(chat_id, "BQADAgADTwsAAkKvaQABElnJclGri9EC"))
					}
					//bot.Send(msg2)
					log.Println(chat_id)
				}
				db.Exec("UPDATE `trophyqueue` SET `send`='1' WHERE (`trophyName`='" + MysqlRealEscapeString(trophyName) + "') AND (`onlineId`='" + onlineId + "')")
			}
			var message string
			var attach string
			rows, err = db.Query("SELECT onlineId,message,attach FROM `user_msg` WHERE send = 0")
			log.Println(err)
			for rows.Next() {
				err = rows.Scan(&onlieId, &message, &attach)
				if err != nil {
					log.Println(err)
				}
				log.Println("Image ", attach)
				chat_row, _ := db.Query("SELECT * FROM `registredchats`")
				for chat_row.Next() {
					var chat_id int64
					err = chat_row.Scan(&chat_id)
					if err != nil {
						log.Println(err)
					}
					log.Println(chat_id)
					msg := tgbotapi.NewMessage(chat_id, ``)
					if len(message) > 1 {
						msg = tgbotapi.NewMessage(chat_id, `<a href="http://ftp.goodgame.ru/`+attach+`"> `+onlieId+`  присылает картинку </a> с cообщением: `+message)
					} else {
						msg = tgbotapi.NewMessage(chat_id, `<a href="http://ftp.goodgame.ru/`+attach+`"> `+onlieId+`  присылает картинку </a>`)
					}

					msg.ParseMode = "HTML"
					bot.Send(msg)
				}
				db.Exec("UPDATE `user_msg` SET `send`='1' WHERE (`attach`='" + attach + "') AND (`onlineId`='" + onlieId + "')")
			}
		}
	})
	//a := cron.New()
	//a.AddFunc("0 */5 * * * *", func() {
	//	var onlieId string
	//	var message string
	//	var attach string
	//	rows, err := db.Query("SELECT onlineId,message,attach FROM `user_msg` WHERE send = 0")
	//	log.Println(err)
	//	for rows.Next() {
	//		err = rows.Scan(&onlieId, &message, &attach)
	//		if err != nil {
	//			log.Println(err)
	//		}
	//		log.Println("Image ", attach)
	//		chat_row, _ := db.Query("SELECT * FROM `registredchats`")
	//		for chat_row.Next() {
	//			var chat_id int64
	//			err = chat_row.Scan(&chat_id)
	//			if err != nil {
	//				log.Println(err)
	//			}
	//			log.Println(chat_id)
	//			msg := tgbotapi.NewMessage(chat_id, ``)
	//			if len(message) > 1 {
	//				msg = tgbotapi.NewMessage(chat_id, `<a href="http://ftp.goodgame.ru/`+attach+`"> `+onlieId+`  присылает картинку </a> с cообщением: `+message)
	//			} else {
	//				msg = tgbotapi.NewMessage(chat_id, `<a href="http://ftp.goodgame.ru/`+attach+`"> `+onlieId+`  присылает картинку </a>`)
	//			}
	//
	//			msg.ParseMode = "HTML"
	//			bot.Send(msg)
	//		}
	//		db.Exec("UPDATE `user_msg` SET `send`='1' WHERE (`attach`='" + attach + "') AND (`onlineId`='" + onlieId + "')")
	//	}
	//})
	//	a.Start()
	c.Start()
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		log.Println(update.Message.Chat.ID)
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		bez_sobaki := strings.Split(update.Message.Text, "@")
		text := strings.Fields(bez_sobaki[0])
		var about = regexp.MustCompile(`/aboutme`)
		var about_other = regexp.MustCompile(`/aboutother`)
		var top = regexp.MustCompile(`/top*`)
		var find = regexp.MustCompile(`/find`)
		var registerChat = regexp.MustCompile(`/registerchat`)
		var last = regexp.MustCompile(`/last`)
		var frnd = regexp.MustCompile(`/friendme`)
		var nk = regexp.MustCompile(`/nearkomplete`)
		var ug = regexp.MustCompile(`/uniqgames`)
		var up = regexp.MustCompile(`/uniqplats`)
		if len(text) > 0 {
			switch {
			case up.MatchString(text[0]):
				var psn string
				rows, err := db.Query("SELECT psn FROM `telegram-psn` WHERE telegram_username = '" + update.Message.From.UserName + "'")
				for rows.Next() {
					err = rows.Scan(&psn)
					if err != nil {
						log.Println(err)
					}
					log.Println(psn)
				}
				go uniqplats(psn, update, bot)
			case ug.MatchString(text[0]):
				var psn string
				rows, err := db.Query("SELECT psn FROM `telegram-psn` WHERE telegram_username = '" + update.Message.From.UserName + "'")
				for rows.Next() {
					err = rows.Scan(&psn)
					if err != nil {
						log.Println(err)
					}
					log.Println(psn)
				}
				go uniqgames(psn, update, bot)
			case about.MatchString(text[0]):
				var psn string
				rows, err := db.Query("SELECT psn FROM `telegram-psn` WHERE telegram_username = '" + update.Message.From.UserName + "'")
				for rows.Next() {
					err = rows.Scan(&psn)
					if err != nil {
						log.Println(err)
					}
					log.Println(psn)
				}
				go aboutother(psn, update, bot)
			case about_other.MatchString(text[0]):
				if len(text) > 1 {
					var psn string = text[1]
					go aboutother(psn, update, bot)
				}
			case top.MatchString(text[0]):
				go gaameTop(text[0], update, bot)
			case find.MatchString(text[0]):
				if len(text) > 1 {
					var NpCommunicationId string
					var name string
					var search_results string
					findstr := strings.Replace(update.Message.Text, "/find ", "", -1)
					log.Println(findstr, len(findstr))
					rows, err := db.Query("SELECT NpCommunicationId, name FROM `games` WHERE `name` LIKE '%" + MysqlRealEscapeString(findstr) + "%'")
					if err != nil {
						fmt.Println(err)
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Запрос неверный")
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					} else {
						for rows.Next() {
							err = rows.Scan(&NpCommunicationId, &name)
							if err != nil {
								log.Println(err)
							}
							search_results = search_results + "\n" + escape2(name) + "       " + "/top" + NpCommunicationId
						}
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, search_results)
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
					}

				}
			case registerChat.MatchString(text[0]):
				if update.Message.From.UserName == "JIEgOKOJI" {
					log.Println(update.Message.Chat.ID)
					res, err := db.Exec("INSERT INTO `registredchats` (`chatId`) VALUES ('" + strconv.FormatInt(update.Message.Chat.ID, 10) + "')")
					log.Println(res, err)

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "chat registred "+strconv.FormatInt(update.Message.Chat.ID, 10))
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)

				}
			case last.MatchString(text[0]):
				var lastGame string
				rows, err := db.Query("SELECT gameId FROM `trophyqueue` ORDER BY earnedDate DESC LiMIT 1")
				log.Println(err)
				for rows.Next() {
					err = rows.Scan(&lastGame)
					if err != nil {
						log.Println(err)
					}
				}
				go gaameTop(lastGame, update, bot)
			case frnd.MatchString(text[0]):
				var result string
				var psn string
				rowse, err := db.Query("SELECT psn FROM `telegram-psn` WHERE telegram_username = '" + update.Message.From.UserName + "'")
				for rowse.Next() {
					err = rowse.Scan(&psn)
					if err != nil {
						log.Println(err)
					}
					log.Println(psn)
				}
				rows, err := db.Query("SELECT count(psn) FROM `friends` WHERE psn ='" + psn + "'")
				log.Println(err)
				for rows.Next() {
					err = rows.Scan(&result)
					if err != nil {
						log.Println(err)
					}
				}
				if result == "0" {
					updateBase, err := exec.Command("/root/psn_bot/PSNbot", "-username", psn, "-fast", "friend").Output()
					log.Println(string(updateBase))
					if err != nil {
						log.Println(string(updateBase), fmt.Sprint(err))
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Запрос отправлен от пользователя ledokol322")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вероятно уже друзья или запрос был отправлен, проверьте уведомления на предмет запроса от пользователя ledokol322")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)

				}
			case nk.MatchString(text[0]):
				var psn string
				rows, err := db.Query("SELECT psn FROM `telegram-psn` WHERE telegram_username = '" + update.Message.From.UserName + "'")
				for rows.Next() {
					err = rows.Scan(&psn)
					if err != nil {
						log.Println(err)
					}
					log.Println(psn)
				}
				go nearKomplete(psn, update, bot)
			}
		}
	}
}
func gaameTop(text string, update tgbot.Update, bot *tgbot.BotAPI) {
	var onlieId string
	var progress string
	var platinum string
	var gameName string
	var image string
	var platform string
	arrr := []string{}
	progress_name := make(map[string]string)
	platinum_name := make(map[string]int)
	m := make(map[string]int)
	s := strings.Split(text, "N")
	var whole_top_new string = "<pre>"
	if len(s) > 1 {
		log.Println(s[1])
		rows, err := db.Query("SELECT name,Image,platform FROM `games` WHERE NpCommunicationId = 'N" + s[1] + "'")
		for rows.Next() {
			err = rows.Scan(&gameName, &image, &platform)
		}
		rows, err = db.Query("SELECT onlieId, progress, Platinum FROM `user_earnings` WHERE gameId = 'N" + s[1] + "' ORDER BY Platinum = 1 DESC, CAST(progress as unsigned) DESC")
		for rows.Next() {
			err = rows.Scan(&onlieId, &progress, &platinum)
			if err != nil {
				log.Println(err)
			}
			arrr = append(arrr, onlieId)
			progress_name[onlieId] = progress
			m[onlieId] = len(onlieId)
			if platinum == "1" {
				platinum_name[onlieId] = 1
			} else {
				platinum_name[onlieId] = 0

			}

		}
		for _, v := range arrr {
			var nm string
			if len(v) > 13 {
				m[v] = 13
				nm = v[:13]
			} else {
				nm = v
			}
			if platinum_name[v] > 0 {
				whole_top_new += "\n" + nm + center("", 15-m[v]) + ":" + progress_name[v] + "\xF0\x9F\x8F\x86"
			} else {
				whole_top_new += "\n" + nm + center("", 15-m[v]) + ":" + progress_name[v]
			}

		}
		whole_top_new += "</pre>"
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*"+MysqlRealEscapeString(gameName)+"* ("+platform+") \n ")
		imgMsg := tgbotapi.NewMessage(update.Message.Chat.ID, `<a href="`+image+`">`+MysqlRealEscapeString(gameName)+`  (`+platform+`)</a>`)
		msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, whole_top_new)
		imgMsg.ReplyToMessageID = update.Message.MessageID
		//msg.ParseMode = "Markdown"
		msg2.ParseMode = "HTML"
		imgMsg.ParseMode = "HTML"
		//bot.Send(msg)
		bot.Send(imgMsg)
		bot.Send(msg2)
	}
}
func uniqgames(psn string, update tgbot.Update, bot *tgbot.BotAPI) {
	var oneGame string
	var platform string
	var uniqGames string = "<pre>"
	rows, err := db.Query("SELECT gameId FROM (SELECT DISTINCT gameId,onlieId FROM user_earnings GROUP BY gameId HAVING COUNT(*) <=1) q1 WHERE onlieId = '" + psn + "'")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&oneGame)
		rowsgame, _ := db.Query("SELECT name, platform FROM `games` WHERE NpCommunicationId = '" + oneGame + "'")
		for rowsgame.Next() {
			err = rowsgame.Scan(&oneGame, &platform)
			uniqGames = uniqGames + "\n " + oneGame + " " + platform
		}
	}
	uniqGames += "</pre>"
	msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if len(uniqGames) > 14 {
		msg2 = tgbotapi.NewMessage(update.Message.Chat.ID, "уникальные игры: "+uniqGames)
	} else {
		msg2 = tgbotapi.NewMessage(update.Message.Chat.ID, "уникальных игр нет")
	}
	msg2.ReplyToMessageID = update.Message.MessageID
	msg2.ParseMode = "HTML"
	bot.Send(msg2)
}
func uniqplats(psn string, update tgbot.Update, bot *tgbot.BotAPI) {
	var oneGame string
	var platform string
	var uniqPlat string = "<pre>"
	rows, _ := db.Query("SELECT gameId FROM (SELECT gameId,onlieId,Platinum FROM user_earnings WHERE Platinum = 1 GROUP BY gameId HAVING COUNT(Platinum) <=1) q1 WHERE onlieId = '" + psn + "'")
	for rows.Next() {
		err := rows.Scan(&oneGame)
		if err != nil {
			fmt.Println(err)
		}
		rowsgame, _ := db.Query("SELECT name, platform FROM `games` WHERE NpCommunicationId = '" + oneGame + "'")
		for rowsgame.Next() {
			err = rowsgame.Scan(&oneGame, &platform)
			//uniqPlat += 1
			uniqPlat = uniqPlat + "\n " + oneGame + " " + platform
		}
	}
	uniqPlat += "</pre>"
	msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if len(uniqPlat) > 14 {
		msg2 = tgbotapi.NewMessage(update.Message.Chat.ID, "уникальные платины в играх: "+uniqPlat)
	} else {
		msg2 = tgbotapi.NewMessage(update.Message.Chat.ID, "уникальных игр нет")
	}
	msg2.ReplyToMessageID = update.Message.MessageID
	msg2.ParseMode = "HTML"
	bot.Send(msg2)
}
func aboutother(psn string, update tgbot.Update, bot *tgbot.BotAPI) {
	var totalgames string
	var platinum string
	var plus string
	var lastPlayed string
	var lastPlayed_name string
	var oneGame string
	var uniqGames int = 0
	var platform string
	var uniqPlat int = 0
	rows, err := db.Query("SELECT totalgames,platinum,plus,lastPlayed  FROM `users` WHERE onlineId = '" + psn + "'")
	log.Println("psn: ", psn, err)
	for rows.Next() {
		err = rows.Scan(&totalgames, &platinum, &plus, &lastPlayed)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println(totalgames, platinum, plus, lastPlayed)
	rows, err = db.Query("SELECT name FROM `games` WHERE NpCommunicationId = '" + lastPlayed + "'")
	for rows.Next() {
		err = rows.Scan(&lastPlayed_name)
	}
	rows, err = db.Query("SELECT gameId FROM (SELECT DISTINCT gameId,onlieId FROM user_earnings GROUP BY gameId HAVING COUNT(*) <=1) q1 WHERE onlieId = '" + psn + "'")
	for rows.Next() {
		err = rows.Scan(&oneGame)
		rowsgame, _ := db.Query("SELECT name, platform FROM `games` WHERE NpCommunicationId = '" + oneGame + "'")
		for rowsgame.Next() {
			err = rowsgame.Scan(&oneGame, &platform)
			uniqGames += 1
			//uniqGames = uniqGames + "\n " + oneGame + " " + platform
		}
	}
	rows, err = db.Query("SELECT gameId FROM (SELECT gameId,onlieId,Platinum FROM user_earnings WHERE Platinum = 1 GROUP BY gameId HAVING COUNT(Platinum) <=1) q1 WHERE onlieId = '" + psn + "'")
	for rows.Next() {
		err = rows.Scan(&oneGame)
		rowsgame, _ := db.Query("SELECT name, platform FROM `games` WHERE NpCommunicationId = '" + oneGame + "'")
		for rowsgame.Next() {
			err = rowsgame.Scan(&oneGame, &platform)
			uniqPlat += 1
			//uniqPlat = uniqPlat + "\n " + oneGame + " " + platform
		}
	}
	//log.Println(uniqPlat)
	if plus == "1" {
		plus = "есть"
	} else {
		plus = "нет"
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "   Информация о профиле: *"+MysqlRealEscapeString(psn)+"*\n Количество игр: *"+totalgames+"*\n Количество платин: *"+platinum+"*\n Подписка PSplus: *"+plus+"*\n В последний раз играл в :*"+lastPlayed_name+"*\n")
	msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if uniqGames > 0 {
		msg2 = tgbotapi.NewMessage(update.Message.Chat.ID, "количество уникальных игр: "+strconv.Itoa(uniqGames))
	} else {
		msg2 = tgbotapi.NewMessage(update.Message.Chat.ID, "уникальных игр нет")
	}

	msg3 := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	if uniqPlat > 0 {
		msg3 = tgbotapi.NewMessage(update.Message.Chat.ID, "\n количество уникальных платин: "+strconv.Itoa(uniqPlat))
	} else {
		msg3 = tgbotapi.NewMessage(update.Message.Chat.ID, "\n Уникальных платин не имеет")
	}

	msg.ReplyToMessageID = update.Message.MessageID
	msg.ParseMode = "Markdown"

	bot.Send(msg)
	bot.Send(msg2)
	bot.Send(msg3)

}
func nearKomplete(psn string, update tgbot.Update, bot *tgbot.BotAPI) {
	var gameid string
	var progress string
	var result string = "<pre>"
	rows, err := db.Query("SELECT gameId,progress FROM `user_earnings` WHERE onlieId = '" + psn + "' AND Platinum = '0' AND progress != '100' ORDER BY CAST(progress as unsigned) DESC LIMIT 6")
	log.Println(err)
	for rows.Next() {
		err = rows.Scan(&gameid, &progress)
		rowsgame, _ := db.Query("SELECT name, platform FROM `games` WHERE NpCommunicationId = '" + gameid + "'")
		for rowsgame.Next() {
			var name string
			var platform string
			err = rowsgame.Scan(&name, &platform)
			result += "\n" + escape2(name) + center("", 43-len(name)) + ":" + progress + " (" + platform + ")"
		}
	}
	result += "</pre>"
	msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, "Ближайшие игры на завершение: \n"+result)
	msg2.ReplyToMessageID = update.Message.MessageID
	msg2.ParseMode = "HTML"
	bot.Send(msg2)

}
func MysqlRealEscapeString(value string) string {
	replace := map[string]string{"'": `\'`, "\\0": "\\\\0", "\n": "\\n", "\r": "\\r", `"`: `\"`, "\x1a": "\\Z", "_": `\_`, `\`: ``}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}
func escape2(value string) string {
	replace := map[string]string{"'": "", `\`: ``}

	for b, a := range replace {
		value = strings.Replace(value, b, a, -1)
	}

	return value
}
func center(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

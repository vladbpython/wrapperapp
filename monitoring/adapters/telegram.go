package adapters

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/vladbpython/wrapperapp/validators"
)

const DefaultHost = "https://api.telegram.org/bot"

//Структура модели состояния ответа бота
type BotModelOk struct {
	Ok bool `json:"ok"` // Статус успешного ответа от бота
}

//Структура модели информации бота
type BotInfoModel struct {
	ID                      int    `json:"id"`                          //Идентификатор
	CanJoinGroups           bool   `json:"is_bot"`                      //Статус присоединения к группам
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"` //Статус читатить все сообщения в группе
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`     //Статутс поддержка Inline запросов
	Username                string `json:"username"`                    // Имя пользователя
}

//Структура модели бота
type BotModel struct {
	BotModelOk              // модель состояния ответа бота
	Result     BotInfoModel `json:"result"` //модели информации бота
}

type BotModelSendMessage struct {
	ChatID    int    `json:"chat_id"`
	Message   string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type Telegram struct {
	AppName    string
	Host       string
	token      string
	Merchants  []int
	HttpClient *HttpClient
}

func (a *Telegram) setUrl(urlSuffix string) string {
	return a.Host + a.token + urlSuffix

}

func (a *Telegram) onError(err error) error {
	return fmt.Errorf("%s %s", a.AppName, err)
}

func (a *Telegram) auth() error {
	url := a.setUrl("/getMe")
	model := BotModel{}
	err := a.HttpClient.MakeRequestWithModel("GET", url, nil, nil, nil, &model, true, true)
	if err != nil {
		return a.onError(err)
	}
	if !model.Ok && model.Result.ID == 0 {
		return a.onError(errors.New("invalid token"))
	}
	return nil
}

func (a *Telegram) Initializate() error {

	return a.auth()
}

func (a *Telegram) addMerchants(merchants ...interface{}) error {
	var i uint

	for _, merchant := range merchants {
		switch merchant.(type) {
		case int:
			a.Merchants[i] = merchant.(int)
		case int64:
			a.Merchants[i] = int(merchant.(int64))
		case int32:
			a.Merchants[i] = int(merchant.(int32))
		case int16:
			a.Merchants[i] = int(merchant.(int16))
		case int8:
			a.Merchants[i] = int(merchant.(int8))
		case uint:
			a.Merchants[i] = int(merchant.(uint))
		case uint64:
			a.Merchants[i] = int(merchant.(uint64))
		case uint32:
			a.Merchants[i] = int(merchant.(uint32))
		case uint16:
			a.Merchants[i] = int(merchant.(uint16))
		case uint8:
			a.Merchants[i] = int(merchant.(uint8))
		default:
			return a.onError(fmt.Errorf("Invalid merchant: %v, value of marchant must be integer", merchant))
		}
		i++
	}
	return nil

}

func (a *Telegram) sendToMerchant(chatID int, message string) error {
	url := a.setUrl("/sendMessage")
	newMessage := BotModelSendMessage{
		ChatID:    chatID,
		Message:   message,
		ParseMode: "HTML",
	}
	data, _ := json.Marshal(&newMessage)
	responseData, statusCode, err := a.HttpClient.MakeRequest("POST", url, nil, nil, data, true, true)
	if err != nil {
		return a.onError(err)
	} else if statusCode >= 200 && statusCode <= 300 {
		if !validators.ValidateResponseStatusCode(statusCode) {
			return a.onError(errors.New(string(responseData)))
		}
	}
	return nil

}

func (a *Telegram) SendData(message string) error {
	for _, merchantID := range a.Merchants {
		err := a.sendToMerchant(merchantID, message)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewTelegramAdapter(appName string, cfg *ConfigAdapter) (*Telegram, error) {

	host := DefaultHost
	if cfg.Host != "" {
		host = strings.TrimRight(cfg.Host, "/")
	}

	adapter := &Telegram{
		AppName:    appName,
		Host:       host,
		token:      cfg.Token,
		Merchants:  make([]int, len(cfg.Merchants)),
		HttpClient: NewHttpClient("rest", 10, nil),
	}
	err := adapter.addMerchants(cfg.Merchants...)

	return adapter, err

}

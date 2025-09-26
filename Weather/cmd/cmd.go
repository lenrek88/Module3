package cmd

import (
	"fmt"
	"os"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"

	"lenrek88/weather"
	"strings"
)

type Cmd struct {
	weather *weather.Informer
}

func NewCmd(w *weather.Informer) *Cmd {
	return &Cmd{
		weather: w,
	}
}

func (c *Cmd) Run() {
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
		prompt.OptionMaxSuggestion(10),
	)

	p.Run()

}

func (c *Cmd) executor(input string) {
	if len(input) < 2 {
		fmt.Println("command cannot be empty  \n")
		return
	}
	parts, err := shlex.Split(input)
	if err != nil {
		fmt.Println("error executing command \n")
		return
	}
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "today":
		fmt.Println("today command called")

		if len(parts) < 4 {
			fmt.Println("Формат: today \"Moscow\" \"Metric\" \"ru\" \n")
			return
		}

		city := parts[1]
		unit := parts[2]
		lang := parts[3]

		resp, err := c.weather.TodayHandler(city, unit, lang)

		if err != nil {
			fmt.Println("failed Informer TodayWeather: " + err.Error() + "\n")
		} else {
			fmt.Println("Дата: ", resp.Forecasts[0].DateTime, "\n Температура: ", resp.Forecasts[0].Temperature, "Ощущается как: ", resp.Forecasts[0].FeelsLike, "Влажность: ", resp.Forecasts[0].Humidity, "Скорость ветра: ", resp.Forecasts[0].WindSpeed)
		}

	case "weekly":
		fmt.Println("weekly command called")

		if len(parts) < 4 {
			fmt.Println("Формат: today \"Moscow\" \"Metric\" \"ru\" \n")
			return
		}

		city := parts[1]
		unit := parts[2]
		lang := parts[3]

		resp, err := c.weather.WeeklyHandler(city, unit, lang)
		if err != nil {
			fmt.Println("failed Informer WeeklyWeather: " + err.Error() + "\n")
		} else {
			for _, e := range resp.Forecasts {
				fmt.Println("Дата: ", e.DateTime, "\n Температура: ", e.Temperature, "Ощущается как: ", e.FeelsLike, "Влажность: ", e.Humidity, "Скорость ветра: ", e.WindSpeed)
			}
		}

	case "help":
		fmt.Println("help command called")
		fmt.Println("	today - Погода в данный момент \n" +
			"	weekly - Погода на неделю \n" +
			"	help - Показать справку \n" +
			"	exit - Выход из программы")
	case "exit":

		fmt.Println("exit command called")
		logger.Info("application is exiting")
		os.Exit(0)

	default:
		fmt.Println("Неизвестная команда:")
		fmt.Println("Введите 'help' для списка команд")
	}

}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	suggestions := []prompt.Suggest{
		{Text: "today", Description: "Погода в данный момент"},
		{Text: "weekly", Description: "Погода на неделю"},
		{Text: "help", Description: "Показать справку"},
		{Text: "exit", Description: "Выйти из программы"},
	}
	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

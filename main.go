package main

import (
	"bufio"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"letters/game"
	"letters/solver"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

var wrongCommandError = errors.New("wrong command")

func loadWords(url string, lines chan<- string, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик горутин
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Download dictionary failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Download dictionary failed. Status Code:", resp.StatusCode)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lines <- line
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Read dictionary error:", err)
	}
}

func filterByLength(lines chan string, length int, filtered chan<- string) {
	for l := range lines {
		if utf8.RuneCountInString(l) == length {
			filtered <- l
		}
	}
	close(filtered) // Закрываем канал после обработки
}

func selectUniqueLines(strings <-chan string) map[string]struct{} {
	uniqueStrings := make(map[string]struct{})
	for str := range strings {
		uniqueStrings[str] = struct{}{}
	}
	return uniqueStrings
}

func writeToFile(lines []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func readFromFile(fileName string) []string {
	file, err := os.Open(fileName)
	if err != nil {
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var result []string
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return result
}

func main() {
	url := "https://raw.githubusercontent.com/LussRus/Rus_words/master/cp1251/txt/nouns/summary.txt"
	url2 := "https://gist.githubusercontent.com/kissarat/bd30c324439cee668f0ac76732d6c825/raw/147eecc9a86ec7f97f6dd442c2eda0641ddd78dc/russian-mnemonic-words.txt"
	letters5TmpPath := "/var/tmp/letters5"
	lettersChan := make(chan string)
	letters5chan := make(chan string)
	var letters5 []string

	if _, err := os.Stat(letters5TmpPath); err == nil {
		letters5 = readFromFile(letters5TmpPath)
	} else {
		var wg sync.WaitGroup

		wg.Add(2) // Добавляем количество горутин для загрузки словарей
		go loadWords(url, lettersChan, &wg)
		go loadWords(url2, lettersChan, &wg)

		go func() {
			wg.Wait()          // Ждем завершения всех загрузок
			close(lettersChan) // Закрываем канал после завершения всех горутин
		}()

		go filterByLength(lettersChan, 5, letters5chan)

		uniqueLines := selectUniqueLines(letters5chan)

		letters5 = make([]string, len(uniqueLines), len(uniqueLines))
		for line := range uniqueLines {
			letters5 = append(letters5, line)
		}
		fmt.Printf("Unique 5-letter lines: %d\n", len(uniqueLines))

		err = writeToFile(letters5, letters5TmpPath)
		if err != nil {
			fmt.Println("Failed to save temp file")
		}
	}
	g := game.NewGame()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN environment variable not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		player := update.Message.From.ID
		if update.Message.Text == "/start" {
			g.AddPlayer(player, solver.NewSolver(letters5, 5))
			sendResponse("New game started!", update, bot)
			continue
		}
		notContainPrefix := "-"
		if strings.HasPrefix(update.Message.Text, notContainPrefix) {
			notContains := strings.TrimLeft(update.Message.Text, notContainPrefix)
			count, err := g.AddNotContains(player, []rune(notContains))
			if err != nil {
				sendResponse(fmt.Sprintf("Error: %v", err), update, bot)
			}
			sendResponse(fmt.Sprintf("Sutable words count: %d", count), update, bot)
			continue
		}
		containPrefix := "+"
		if strings.HasPrefix(update.Message.Text, containPrefix) {
			contains := strings.TrimLeft(update.Message.Text, containPrefix)
			count, err := g.AddContains(player, []rune(contains))
			if err != nil {
				sendResponse(fmt.Sprintf("Error: %v", err), update, bot)
			}
			sendResponse(fmt.Sprintf("Sutable words count: %d", count), update, bot)
			continue
		}
		if correctPosition(update, g, player, bot) {
			continue
		}
		if incorrectPosition(update, g, player, bot) {
			continue
		}
		if update.Message.Text == "/result" {
			result, err := g.GetResult(player)
			message := ""
			if err != nil {
				message = "Error!"
			}
			message = strings.Join(result, ",\n")
			sendResponse(message, update, bot)

			if err != nil {
				log.Println(err)
			}
			continue
		}
		sendResponse("Wrong command", update, bot)
	}
}

var replyKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/start"),
		tgbotapi.NewKeyboardButton("/result"),
	),
)

func sendResponse(message string, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = replyKeyboard
	_, err := bot.Send(msg)

	if err != nil {
		log.Println(err)
	}
}

func correctPosition(update tgbotapi.Update, g *game.Game, player int, bot *tgbotapi.BotAPI) bool {
	processed := false
	re := regexp.MustCompile(`^(?P<pos>\d)\+(?P<char>[а-я])$`)
	match := re.FindStringSubmatch(update.Message.Text)

	if len(match) > 0 {
		charIndex := re.SubexpIndex("char")
		posIndex := re.SubexpIndex("pos")
		if charIndex != -1 && posIndex != -1 {
			processed = true
			c := runeAt(match[charIndex], 0)
			p := match[posIndex]
			pos, err := strconv.Atoi(p)
			if err != nil {
				sendResponse(fmt.Sprintf("Error: %v", err), update, bot)
			}
			count, err := g.AddCorrectPosition(player, c, pos)
			if err != nil {
				sendResponse(fmt.Sprintf("Error: %v", err), update, bot)
			}
			sendResponse(fmt.Sprintf("Sutable words count: %d", count), update, bot)
		}
	}
	return processed
}

func runeAt(s string, i int) rune {
	runeSlice := []rune(s)
	if i >= len(runeSlice) || i < 0 {
		return -1
	}
	return runeSlice[i]
}

func incorrectPosition(update tgbotapi.Update, g *game.Game, player int, bot *tgbotapi.BotAPI) bool {
	processed := false
	c, pos, err := getIncorrectCharPositionFromRequest(update.Message.Text)

	if err != nil {
		if !errors.Is(err, wrongCommandError) {
			sendResponse(fmt.Sprintf("Error: %v", err), update, bot)
		}
	} else {
		count, err := g.AddIncorrectPosition(player, c, pos)
		processed = true
		if err != nil {
			sendResponse(fmt.Sprintf("Error: %v", err), update, bot)
		}
		sendResponse(fmt.Sprintf("Sutable words count: %d", count), update, bot)
	}
	return processed
}

func getIncorrectCharPositionFromRequest(message string) (rune, int, error) {
	re := regexp.MustCompile(`^(?P<pos>\d)-(?P<char>[а-я])$`)
	match := re.FindStringSubmatch(message)

	if len(match) > 0 {
		charIndex := re.SubexpIndex("char")
		posIndex := re.SubexpIndex("pos")
		if charIndex != -1 && posIndex != -1 {
			c := []rune(match[charIndex])[0]
			p := match[posIndex]
			pos, err := strconv.Atoi(p)

			return c, pos, err
		}
	}
	return 0, 0, wrongCommandError
}

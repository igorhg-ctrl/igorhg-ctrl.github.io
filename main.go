package main

import (
  "bufio"
  "fmt"
  "log"
  "net/http"
  "os"
  "strings"

  "time"

  "github.com/faiface/beep"
  "github.com/faiface/beep/mp3"
  "github.com/faiface/beep/speaker"
)

func main() {
  // Список доступных песен
  songs := map[string]string{
    "1": "https://github.com/RoboLask/kardana_music_repo/raw/main/Rammstein/07.%20Spieluhr.mp3",
    "2": "https://github.com/RoboLask/kardana_music_repo/raw/main/Rammstein/06.%20Mutter.mp3",
    "3": "https://github.com/RoboLask/kardana_music_repo/raw/main/Rammstein/06.%20Du%20Riechst%20So%20Gut.mp3",
    // Добавьте больше песен по необходимости
  }

  fmt.Println("Выберите песню для воспроизведения:")
  for key, value := range songs {
    fmt.Printf("%s: %s\n", key, getSongName(value))
  }

  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Введите номер песни: ")
  choice, _ := reader.ReadString('\n')
  choice = strings.TrimSpace(choice)

  url, exists := songs[choice]
  if !exists {
    fmt.Println("Неверный выбор.")
    return
  }

  // Загрузите MP3 файл
  resp, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer resp.Body.Close()

  // Декодируйте MP3 данные
  streamer, format, err := mp3.Decode(resp.Body)
  if err != nil {
    log.Fatal(err)
  }
  defer streamer.Close()

  // Инициализируйте динамик
  err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
  if err != nil {
    log.Fatal(err)
  }

  // Воспроизведите аудио
  done := make(chan bool)
  speaker.Play(beep.Seq(streamer, beep.Callback(func() {
    done <- true
  })))

  // Блокируйте выполнение до завершения воспроизведения
  <-done

  fmt.Println("Воспроизведение завершено")
}

// getSongName извлекает имя песни из URL
func getSongName(url string) string {
  parts := strings.Split(url, "/")
  return parts[len(parts)-1]
}
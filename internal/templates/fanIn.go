package templates

import (
	"sync"
)

// FanIn объединяет все каналы в один канал.
func FanIn(doneCh chan struct{}, sliceURLs []string) chan string {
	// Создаем выходной канал для объединения всех результатов.
	outCh := make(chan string, len(sliceURLs))

	// Используем WaitGroup для ожидания завершения всех горутин.
	var wg sync.WaitGroup

	// Запускаем горутину для каждого URL.
	for _, url := range sliceURLs {
		ch := url
		wg.Add(1)

		go func() {
			defer wg.Done()

			select {
			case <-doneCh:
				return
			case outCh <- ch:
			}
		}()
	}

	// Закрываем выходной канал после завершения работы всех горутин.
	go func() {
		wg.Wait()
		close(outCh)
	}()

	return outCh
}

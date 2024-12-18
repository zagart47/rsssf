package router

import (
	"github.com/stretchr/testify/assert"
	"rsssf/entity"
	"testing"
	"time"
)

func TestFetchRSSWithHabrGoHub(t *testing.T) {
	// URL RSS-ленты с Habr
	habrRSSURL := "https://habr.com/ru/rss/hub/go/all/?fl=ru"

	// Канал для получения результатов
	ch := make(chan []entity.Post, 1)

	// Вызываем функцию fetchRSS с реальным URL
	go fetchRSS(habrRSSURL, ch)

	// Ждем результата
	select {
	case posts := <-ch:
		// Проверяем, что результат не пустой
		assert.Greater(t, len(posts), 0, "Должен быть хотя бы один пост")

		// Проверяем структуру первого поста
		firstPost := posts[0]
		assert.NotEmpty(t, firstPost.Title, "Заголовок поста не должен быть пустым")
		assert.NotEmpty(t, firstPost.Link, "Ссылка поста не должна быть пустой")
		assert.NotEmpty(t, firstPost.Content, "Содержимое поста не должно быть пустым")
		assert.NotZero(t, firstPost.PubTime, "Время публикации поста не должно быть нулевым")

		// Пример проверки даты публикации
		assert.True(t, firstPost.PubTime > 0, "Время публикации должно быть больше 0")

	case <-time.After(10 * time.Second):
		// Если результат не получен за 10 секунд, тест провален
		t.Fatal("fetchRSS не вернул результат за 10 секунд")
	}
}

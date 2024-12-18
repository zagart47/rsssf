package router

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log"
	"net/http"
	"rsssf/config"
	"rsssf/entity"
	"rsssf/service"
	"strconv"
	"sync"
	"time"
)

type Router struct {
	mux.Router
	Services service.Services
}

func NewRouter(services service.Services) Router {
	return Router{
		Services: services,
	}
}

func (r *Router) InitHandlers() {
	r.HandleFunc("/news/{n}", r.GetNews).Methods(http.MethodGet, http.MethodOptions)
}

// GetNews показывает последние новости, по умолчанию 10
func (r *Router) GetNews(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		log.Printf("Invalid parameter 'n': %v", err)
		http.Error(w, "Invalid parameter 'n'", http.StatusBadRequest)
		return
	}
	if n <= 0 {
		log.Printf("Parameter 'n' must be greater than 0, got: %d", n)
		http.Error(w, "Parameter 'n' must be greater than 0", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	posts, err := r.Services.Poster.GetPosts(ctx, n)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	postsForPublication := convertPosts(posts)

	// Проверяем, не истек ли контекст
	if ctx.Err() != nil {
		log.Printf("Request timed out: %v", ctx.Err())
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
		return
	}

	// Устанавливаем заголовки
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	// Рендерим шаблон
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	err = tmpl.Execute(w, postsForPublication)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

// convertPosts конвертирует новости из одной структуры в другую, пригодную для публикации
func convertPosts(posts []entity.Post) []entity.PostForPublic {
	postsForPublication := make([]entity.PostForPublic, len(posts))
	for i := range posts {
		postsForPublication[i] = entity.PostForPublic{
			Link:        template.URL(posts[i].Link),
			Title:       posts[i].Title,
			ContentHTML: template.HTML(posts[i].Content),
			PubTime:     time.Unix(posts[i].PubTime, 0).Format("02.01.2006 15:04:05"),
		}
	}
	return postsForPublication
}

// UpdateNews в цикле каждые 5 минут получает и обновляет в БД новости из RSS
func (r *Router) UpdateNews() {
	for {
		rss := config.Configs.RSS
		ch := make(chan []entity.Post, len(rss))
		var wg sync.WaitGroup // Используем WaitGroup для ожидания завершения горутин

		for _, url := range rss {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				fetchRSS(url, ch)
			}(url) // Передаем url как аргумент
		}

		// Запускаем горутину для закрытия канала после завершения всех горутин
		go func() {
			wg.Wait()
			close(ch)
		}()

		var allNews []entity.Post
		for items := range ch {
			allNews = append(allNews, items...)
		}

		// Создаем контекст с таймаутом
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		err := r.Services.Poster.AddPosts(ctx, allNews)
		if err != nil {
			log.Println("Error adding posts:", err)
		}

		cancel()
		// ждем до следующего обновления
		time.Sleep(config.Configs.Timeout)
	}
}

func fetchRSS(url string, ch chan<- []entity.Post) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body from %s: %v\n", url, err)
		return
	}

	var rss entity.RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		fmt.Printf("Error unmarshalling XML from %s: %v\n", url, err)
		return
	}
	layouts := []string{
		"Mon, 02 Jan 2006 15:04:05 GMT",
		"Mon, 2 Jan 2006 15:04:05 -0700",
	}
	items := rss.Channel.Items
	var posts []entity.Post

	for _, item := range items {
		var post entity.Post
		post.Title = item.Title
		post.Content = item.Description
		post.Link = item.Link

		var postTime time.Time
		var err error
		for _, layout := range layouts {
			postTime, err = time.Parse(layout, item.PubDate)
			if err == nil {
				break // Дата успешно распарсена, выходим из цикла
			}
		}
		if err != nil {
			fmt.Printf("Error parsing date from %v: %v\n", item.PubDate, err)
			continue // Пропускаем элемент, если дата не распарсена
		}

		post.PubTime = postTime.Unix()
		posts = append(posts, post)
	}
	ch <- posts
}

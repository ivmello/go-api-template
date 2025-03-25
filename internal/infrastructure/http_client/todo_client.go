package http_client

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	TodoAPIBaseURL = "https://jsonplaceholder.typicode.com"
)

// TodoClient is a client for the JSONPlaceholder Todo API
type TodoClient struct {
	client *Client
}

// Todo represents a todo item from the JSONPlaceholder API
type Todo struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Post represents a post item from the JSONPlaceholder API
type Post struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// NewTodoClient creates a new TodoClient
func NewTodoClient(timeout time.Duration) *TodoClient {
	return &TodoClient{
		client: NewClient(timeout),
	}
}

// GetTodo fetches a single todo by ID
func (c *TodoClient) GetTodo(ctx context.Context, id int) (*Todo, error) {
	url := fmt.Sprintf("%s/todos/%d", TodoAPIBaseURL, id)
	
	resp, err := c.client.Get(ctx, url, &RequestConfig{
		Headers: map[string]string{
			"Accept": "application/json",
		},
		Timeout: 3 * time.Second, // Custom timeout for this specific request
	})
	
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	var todo Todo
	if err := UnmarshalResponse(resp, &todo); err != nil {
		return nil, err
	}
	
	return &todo, nil
}

// GetPost fetches a single post by ID
func (c *TodoClient) GetPost(ctx context.Context, id int) (*Post, error) {
	url := fmt.Sprintf("%s/posts/%d", TodoAPIBaseURL, id)
	
	resp, err := c.client.Get(ctx, url, &RequestConfig{
		Headers: map[string]string{
			"Accept": "application/json",
		},
	})
	
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	var post Post
	if err := UnmarshalResponse(resp, &post); err != nil {
		return nil, err
	}
	
	return &post, nil
}

// FetchMultiple demonstrates fetching multiple resources concurrently
func (c *TodoClient) FetchMultiple(ctx context.Context, todoIDs []int, postIDs []int) ([]Todo, []Post, error) {
	var wg sync.WaitGroup
	
	// Create a buffered channel for todos
	todosChan := make(chan Todo, len(todoIDs))
	// Create a buffered channel for posts
	postsChan := make(chan Post, len(postIDs))
	// Error channel
	errChan := make(chan error, len(todoIDs)+len(postIDs))
	
	// Fetch todos concurrently
	for _, id := range todoIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			todo, err := c.GetTodo(ctx, id)
			if err != nil {
				errChan <- fmt.Errorf("error fetching todo %d: %w", id, err)
				return
			}
			
			todosChan <- *todo
		}(id)
	}
	
	// Fetch posts concurrently
	for _, id := range postIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			post, err := c.GetPost(ctx, id)
			if err != nil {
				errChan <- fmt.Errorf("error fetching post %d: %w", id, err)
				return
			}
			
			postsChan <- *post
		}(id)
	}
	
	// Wait for all goroutines to finish
	wg.Wait()
	close(todosChan)
	close(postsChan)
	close(errChan)
	
	// Check for errors
	if len(errChan) > 0 {
		// Return the first error
		err := <-errChan
		return nil, nil, err
	}
	
	// Collect results
	todos := make([]Todo, 0, len(todoIDs))
	for todo := range todosChan {
		todos = append(todos, todo)
	}
	
	posts := make([]Post, 0, len(postIDs))
	for post := range postsChan {
		posts = append(posts, post)
	}
	
	return todos, posts, nil
}
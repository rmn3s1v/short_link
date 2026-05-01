package service_test

import (
	"regexp"
	"short-link/cmd/internal/repository"
	"short-link/cmd/internal/service"
	"sync"
	"testing"
)

func TestShortenIdempotent(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.New(repo)

	url := "https://google.com"

	s1, err := svc.Shorten(url)
	if err != nil {
		t.Fatal(err)
	}

	s2, err := svc.Shorten(url)
	if err != nil {
		t.Fatal(err)
	}

	if s1 != s2 {
		t.Fatal("not idempotent")
	}
}

func TestShortenFormat(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.New(repo)

	short, err := svc.Shorten("https://google.com")
	if err != nil {
		t.Fatal(err)
	}

	if len(short) != 10 {
		t.Fatalf("expected short URL length 10, got %d", len(short))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(short) {
		t.Fatalf("unexpected short URL format: %q", short)
	}
}

func TestResolve(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.New(repo)

	url := "https://example.com"

	short, err := svc.Shorten(url)
	if err != nil {
		t.Fatal(err)
	}

	resolved, err := svc.Resolve(short)
	if err != nil {
		t.Fatal(err)
	}

	if resolved != url {
		t.Fatalf("expected %s, got %s", url, resolved)
	}
}

func TestShortenInvalidURL(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.New(repo)

	_, err := svc.Shorten("not-a-url")

	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestConcurrentShorten(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.New(repo)

	url := "https://google.com"

	var wg sync.WaitGroup
	results := make(chan string, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s, err := svc.Shorten(url)
			if err != nil {
				t.Errorf("shorten failed: %v", err)
				return
			}
			results <- s
		}()
	}

	wg.Wait()
	close(results)

	var first string
	for r := range results {
		if first == "" {
			first = r
		}
		if r != first {
			t.Fatal("not consistent under concurrency")
		}
	}
}

func TestDifferentURLs(t *testing.T) {
	repo := repository.NewMemoryRepo()
	svc := service.New(repo)

	s1, err := svc.Shorten("https://a.com")
	if err != nil {
		t.Fatal(err)
	}

	s2, err := svc.Shorten("https://b.com")
	if err != nil {
		t.Fatal(err)
	}

	if s1 == s2 {
		t.Fatal("different URLs should have different shorts")
	}
}

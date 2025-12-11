package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080"

var sessionCookie string

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateEventRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Date        string   `json:"date"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Fee         float32  `json:"fee"`
	Tags        []string `json:"tags,omitempty"`
}

var dummyUsers = []RegisterRequest{
	{Name: "Alice Johnson", Email: "alice@example.com", Password: "Password123@"},
	{Name: "Bob Smith", Email: "bob@example.com", Password: "Password123@"},
	{Name: "Charlie Brown", Email: "charlie@example.com", Password: "Password123@"},
	{Name: "Diana Prince", Email: "diana@example.com", Password: "Password123@"},
	{Name: "Edward Stark", Email: "edward@example.com", Password: "Password123@"},
	{Name: "Fiona Green", Email: "fiona@example.com", Password: "Password123@"},
	{Name: "George Miller", Email: "george@example.com", Password: "Password123@"},
	{Name: "Hannah White", Email: "hannah@example.com", Password: "Password123@"},
	{Name: "Ivan Petrov", Email: "ivan@example.com", Password: "Password123@"},
	{Name: "Julia Roberts", Email: "julia@example.com", Password: "Password123@"},
}

var dummyEvents = []CreateEventRequest{
	{
		Name:        "Tech Conference 2025",
		Description: "Annual technology conference featuring the latest innovations in AI and cloud computing.",
		Latitude:    52.2297,
		Longitude:   21.0122,
		Fee:         49.99,
		Tags:        []string{"Tech", "Conference", "Networking"},
	},
	{
		Name:        "Summer Music Festival",
		Description: "A three-day outdoor music festival featuring local and international artists.",
		Latitude:    52.4064,
		Longitude:   16.9252,
		Fee:         75.00,
		Tags:        []string{"Music", "Outdoor", "Party"},
	},
	{
		Name:        "Charity Run for Health",
		Description: "5K charity run to raise awareness for mental health. All proceeds go to local organizations.",
		Latitude:    50.0647,
		Longitude:   19.9450,
		Fee:         15.00,
		Tags:        []string{"Sports", "Charity", "Health & Wellness"},
	},
	{
		Name:        "Startup Networking Night",
		Description: "Connect with fellow entrepreneurs, investors, and tech enthusiasts over drinks and appetizers.",
		Latitude:    52.2297,
		Longitude:   21.0122,
		Fee:         0,
		Tags:        []string{"Networking", "Tech"},
	},
	{
		Name:        "Art Exhibition: Modern Visions",
		Description: "Explore contemporary art from emerging Polish artists in this exclusive gallery showing.",
		Latitude:    51.7592,
		Longitude:   19.4560,
		Fee:         10.00,
		Tags:        []string{"Art", "Education"},
	},
	{
		Name:        "Gaming Tournament",
		Description: "Compete in our annual esports tournament featuring popular titles. Prizes for winners!",
		Latitude:    54.3520,
		Longitude:   18.6466,
		Fee:         25.00,
		Tags:        []string{"Gaming", "Tech"},
	},
	{
		Name:        "Cooking Workshop: Italian Cuisine",
		Description: "Learn to prepare authentic Italian dishes with Chef Marco. All ingredients provided.",
		Latitude:    52.4064,
		Longitude:   16.9252,
		Fee:         85.00,
		Tags:        []string{"Workshop", "Food & Drink", "Education"},
	},
	{
		Name:        "Yoga in the Park",
		Description: "Free outdoor yoga session for all levels. Bring your own mat and water.",
		Latitude:    52.2297,
		Longitude:   21.0122,
		Fee:         0,
		Tags:        []string{"Health & Wellness", "Outdoor"},
	},
	{
		Name:        "Developer Meetup",
		Description: "Monthly meetup for developers to share knowledge, discuss new technologies, and network.",
		Latitude:    51.1079,
		Longitude:   17.0385,
		Fee:         0,
		Tags:        []string{"Tech", "Meetup", "Networking"},
	},
	{
		Name:        "Food Truck Festival",
		Description: "Sample delicious food from 30+ food trucks representing cuisines from around the world.",
		Latitude:    50.0647,
		Longitude:   19.9450,
		Fee:         5.00,
		Tags:        []string{"Food & Drink", "Outdoor"},
	},
	{
		Name:        "Photography Workshop",
		Description: "Learn professional photography techniques from award-winning photographers.",
		Latitude:    52.2297,
		Longitude:   21.0122,
		Fee:         120.00,
		Tags:        []string{"Workshop", "Art", "Education"},
	},
	{
		Name:        "Basketball Tournament",
		Description: "3v3 street basketball tournament. Form your team and compete for prizes!",
		Latitude:    54.3520,
		Longitude:   18.6466,
		Fee:         30.00,
		Tags:        []string{"Sports", "Outdoor"},
	},
}

func main() {
	log.Println("Starting database seeding...")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	log.Println("Creating users...")
	for _, user := range dummyUsers {
		err := registerUser(client, user)
		if err != nil {
			log.Printf("Failed to register user %s: %v", user.Email, err)
		} else {
			log.Printf("Successfully registered user: %s", user.Email)
		}
	}

	log.Println("\nCreating events...")
	for i, user := range dummyUsers {
		err := loginUser(client, LoginRequest{Email: user.Email, Password: user.Password})
		if err != nil {
			log.Printf("Failed to login as %s: %v", user.Email, err)
			continue
		}

		err = becomeHost(client)
		if err != nil {
			log.Printf("Failed to become host as %s: %v", user.Email, err)
			logoutUser(client)
			continue
		}
		numEvents := 1 + rand.Intn(2)
		for j := 0; j < numEvents && (i*2+j) < len(dummyEvents); j++ {
			eventIdx := (i*2 + j) % len(dummyEvents)
			event := dummyEvents[eventIdx]
			daysFromNow := rand.Intn(60) + 1
			eventDate := time.Now().AddDate(0, 0, daysFromNow)
			event.Date = eventDate.Format(time.RFC3339)

			err := createEvent(client, event)
			if err != nil {
				log.Printf("Failed to create event '%s': %v", event.Name, err)
			} else {
				log.Printf("Successfully created event: %s (by %s)", event.Name, user.Email)
			}
		}

		logoutUser(client)
	}

	log.Println("\nSeeding completed!")
}

func registerUser(client *http.Client, user RegisterRequest) error {
	body, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	resp, err := client.Post(baseURL+"/api/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func loginUser(client *http.Client, login LoginRequest) error {
	body, err := json.Marshal(login)
	if err != nil {
		return fmt.Errorf("failed to marshal login: %w", err)
	}

	resp, err := client.Post(baseURL+"/api/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session-id" {
			sessionCookie = cookie.Value
			break
		}
	}

	return nil
}

func becomeHost(client *http.Client) error {
	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/become-host", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionCookie})

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func logoutUser(client *http.Client) error {
	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/logout", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionCookie})

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	sessionCookie = ""
	return nil
}

func createEvent(client *http.Client, event CreateEventRequest) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/events/add", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionCookie})

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

package gg_eleven 

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GenerateSoundEffect(text string) ([]byte, error) {
	fmt.Println("Sending request to Elevenlabs")
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ELEVENLABS_API_KEY environment variable not set.")
		os.Exit(1)
	}

	body, err := json.Marshal(map[string]string{
		"text": text,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		os.Exit(1)
	}

	req, _ := http.NewRequest("POST", "https://api.elevenlabs.io/v1/sound-generation", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error from ElevenLabs API (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}

	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR ERROR ERROR")
		os.Exit(1)
	}
	return audioBytes, nil
}

func TextToSpeech(voiceID string, text string)  []byte {
	fmt.Println("Sending request to Elevenlabs")
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ELEVENLABS_API_KEY environment variable not set.")
		os.Exit(1)
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)

	payload := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_flash_v2_5", // You can change this model
		// "voice_settings": map[string]float64{ // Optional: uncomment and adjust
		// 	"stability":        0.75,
		// 	"similarity_boost": 0.75,
		//  "style":            0.0, // Set style to 0.0 if using eleven_multilingual_v2 with non-English text
		//  "use_speaker_boost": true,
		// },
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	log.Printf("ELEVENLABS: Got response, status code: %d", resp.StatusCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request to ElevenLabs: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error from ElevenLabs API (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}
	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR ERROR ERROR")
		os.Exit(1)
	}
	return audioBytes
}


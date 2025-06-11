package gg_eleven 

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GenerateSoundEffect(text string) ([]byte, error) {
	fmt.Printf("Sending request to Elevenlabs")
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ELEVENLABS_API_KEY environment variable not set.")
		return nil, errors.New("SET YOUR API KEY, MORON")
	}

	body, err := json.Marshal(map[string]string{
		"text": text,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling JSON: %v\n", err)
		return nil, err
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
		return nil, err
	}

	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR ERROR ERROR")
		return nil, err
	}
	return audioBytes, nil
}

func TextToSpeech(voiceID string, text string)  ([]byte, error) {
	fmt.Println("Sending request to Elevenlabs")
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ELEVENLABS_API_KEY environment variable not set.")
		return nil, errors.New("SET YOUR API KEY, MORON")
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
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
		return nil, err
	}

	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	log.Printf("ELEVENLABS: Got response, status code: %d", resp.StatusCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making request to ElevenLabs: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "Error from ElevenLabs API (Status %d): %s\n", resp.StatusCode, string(bodyBytes))
		return nil, err
	}
	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR ERROR ERROR")
		return nil, err
	}
	return audioBytes, nil
}


func GetVoiceIDs() ([]string, error) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ELEVENLABS_API_KEY environment variable not set.")
		return nil, errors.New("SET YOUR API KEY, MORON")
	}
	url := "https://api.elevenlabs.io/v1/voices"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("xi-api-key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	var ids []string
	voices, ok := result["voices"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse voices")
	}

	for _, v := range voices {
		voiceMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := voiceMap["voice_id"].(string); ok {
			ids = append(ids, id)
		}
	}

	return ids, nil
}


func AddSharedVoice(publicUserID, voiceID, newName string) error {
	fmt.Printf("Attempting to add voice %s from userID %s and set name to %s", voiceID, publicUserID, newName)
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ELEVENLABS_API_KEY environment variable not set.")
		return errors.New("SET YOUR API KEY, MORON")
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/voices/add/%s/%s", publicUserID, voiceID)
	
	payload := map[string]string{
		"new_name": newName,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Xi-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp)
		body, _:= io.ReadAll(resp.Body)
		fmt.Println(body)
		return fmt.Errorf("failed to add voice, status code: %d", resp.StatusCode)
	}

	return nil
}

package livekit

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	lksdk "github.com/livekit/server-sdk-go/v2"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

type WebhookEvent struct {
	Event       string `json:"event"`
	Room        string `json:"room"`
	Participant struct {
		Identity string `json:"identity"`
		Type     string `json:"type"`
	} `json:"participant"`
}

// ConnectToRoom connects to LiveKit as sender and receiver and returns the rooms and room name
func ConnectToRoom(host, apikey, apiSecret string) (senderRoom *lksdk.Room, receiverRoom *lksdk.Room, roomName string, err error) {
	roomName = "airagent" + uuid.NewString()
	senderRoom, err = lksdk.ConnectToRoom(host, lksdk.ConnectInfo{
		APIKey:              apikey,
		APISecret:           apiSecret,
		RoomName:            roomName,
		ParticipantIdentity: "airagent",
	}, nil)
	if err != nil {
		return nil, nil, roomName, err
	}
	defer func() {
		if err != nil {
			senderRoom.Disconnect()
		}
	}()
	log.Println("Connected to room as sender:", senderRoom.Name())

	receiverRoom, err = lksdk.ConnectToRoom(host, lksdk.ConnectInfo{
		APIKey:              apikey,
		APISecret:           apiSecret,
		RoomName:            roomName,
		ParticipantIdentity: "receiver",
	}, nil)
	if err != nil {
		return nil, nil, roomName, err
	}
	log.Println("Connected to room as receiver:", receiverRoom.Name())

	return senderRoom, receiverRoom, roomName, nil
}

// ConnectSingleRoom connects to a LiveKit room as a single participant with the given identity
func ConnectSingleRoom(host, apiKey, apiSecret, roomName, identity string) (*lksdk.Room, error) {
	room, err := lksdk.ConnectToRoom(host, lksdk.ConnectInfo{
		APIKey:              apiKey,
		APISecret:           apiSecret,
		RoomName:            roomName,
		ParticipantIdentity: identity,
	}, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to room as", identity, "in room:", roomName)
	return room, nil
}

// SubscribeToRemoteAudio subscribes to remote audio tracks and calls the handler for each audio packet
func SubscribeToRemoteAudio(room *lksdk.Room, onAudio func([]byte)) {
	room.LocalParticipant.Callback.OnTrackSubscribed = func(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, participant *lksdk.RemoteParticipant) {
		if track.Kind() == webrtc.RTPCodecTypeAudio {
			go func() {
				for {
					buf := make([]byte, 1500)
					n, _, err := track.Read(buf)
					if err != nil {
						log.Println("Error reading audio sample:", err)
						return
					}
					onAudio(buf[:n])
				}
			}()
		}
	}
}

// PublishProcessedAudio publishes audio data as a new track to the room
func PublishProcessedAudio(room *lksdk.Room, audioChan <-chan []byte) error {
	track, err := lksdk.NewLocalSampleTrack(webrtc.RTPCodecCapability{
		MimeType: "audio/opus",
	})
	if err != nil {
		return err
	}
	_, err = room.LocalParticipant.PublishTrack(track, nil)
	if err != nil {
		return err
	}
	go func() {
		for pkt := range audioChan {
			err := track.WriteSample(media.Sample{Data: pkt, Duration: 20 * time.Millisecond}, nil)
			if err != nil {
				log.Println("Error writing audio sample:", err)
				return
			}
		}
	}()
	return nil
}

// BridgeAudioWithWebSocket bridges audio between LiveKit and a WebSocket client
func BridgeAudioWithWebSocket(room *lksdk.Room, wsURL string, stopChan <-chan struct{}) error {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	audioChan := make(chan []byte, 100)
	// Forward LiveKit audio to WS client
	go SubscribeToRemoteAudio(room, func(audio []byte) {
		conn.WriteMessage(websocket.BinaryMessage, audio)
	})

	// Forward WS client audio to LiveKit
	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				close(audioChan)
				return
			}
			audioChan <- data
		}
	}()

	// Publish audio to LiveKit
	err = PublishProcessedAudio(room, audioChan)
	if err != nil {
		return err
	}

	// Wait for stop signal (e.g., SIP participant left)
	<-stopChan
	return nil
}

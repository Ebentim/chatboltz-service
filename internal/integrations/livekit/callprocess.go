package livekit

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

type CallProcessConfig struct {
	APIURL      string
	APIKey      string
	APISecret   string
	RoomPrefix  string
	TrunkName   string
	PhoneNumber string
	Username    string
	Password    string
	RuleName    string
	RoomName    string
	RequirePin  bool
	Pin         string
}

// SetupTrunkAndRule checks for existing trunk and rule, creates them if missing, and attaches the rule to the trunk
func SetupTrunkAndRule(client *lksdk.SIPClient, numbers []string, cfg CallProcessConfig) (*livekit.SIPInboundTrunkInfo, *livekit.SIPDispatchRuleInfo, error) {

	trunkIds := []string{"air-agent" + uuid.NewString()}
	sipTrunkInfo, err := GetSIPTrunkByName(client, trunkIds)
	if err != nil {
		return nil, nil, err
	}
	if sipTrunkInfo == nil {
		trunk, err := CreateSIPTrunk(client, cfg.TrunkName, "air-agent2.0", "airagenttest", trunkIds[0], numbers)
		if err != nil {
			return nil, nil, err
		}
		log.Println("Created new trunk: \n", trunk)
	} else {
		log.Println("Trunk already exists: \n", sipTrunkInfo)
	}

	dispatchIds := []string{"air-agent" + uuid.NewString()}
	sipDispatchRule, err := GetSIPDispatchRuleByName(client, dispatchIds)
	if err != nil {
		return nil, nil, err
	}

	if sipDispatchRule == nil {
		t := time.Now().UTC().Format(time.RFC3339)
		rule, err := CreateSIPDispatchRule(client, cfg.RuleName, cfg.RoomPrefix, t, "airagent2.0", "airagent2.0 trial")
		if err != nil {
			return nil, nil, err
		}
		log.Println("CREATED NEW DISPATCH RULE", rule)
	}
	return nil, nil, nil
}

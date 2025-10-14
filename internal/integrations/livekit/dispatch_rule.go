package livekit

import (
	"context"
	"strings"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

func CreateSIPDispatchRule(Client *lksdk.SIPClient, ruleName, roomPrefix, roomName, agentName, agentMetadata string) (*livekit.SIPDispatchRule, error) {
	request := &livekit.CreateSIPDispatchRuleRequest{
		Name: ruleName,
		Rule: &livekit.SIPDispatchRule{
			Rule: &livekit.SIPDispatchRule_DispatchRuleIndividual{
				DispatchRuleIndividual: &livekit.SIPDispatchRuleIndividual{
					RoomPrefix: roomPrefix,
				},
			},
		},

		RoomConfig: &livekit.RoomConfiguration{
			Agents: []*livekit.RoomAgentDispatch{
				{
					AgentName: agentName,
					Metadata:  agentMetadata,
				},
			},
			Name:            roomName,
			MaxParticipants: 3,
		},
	}
	dispatchRuleInfo, err := Client.CreateSIPDispatchRule(context.Background(), request)
	if err != nil {
		return nil, err
	}
	if dispatchRuleInfo == nil || dispatchRuleInfo.Rule == nil {
		return nil, nil
	}
	return dispatchRuleInfo.Rule, nil
}

func GetSIPDispatchRuleByName(Client *lksdk.SIPClient, trunkIds []string) ([]*livekit.SIPDispatchRuleInfo, error) {
	response, err := Client.ListSIPDispatchRule(context.Background(), &livekit.ListSIPDispatchRuleRequest{
		TrunkIds: trunkIds,
	})
	if err != nil {
		return nil, err
	}
	var matchedRules []*livekit.SIPDispatchRuleInfo

	for _, resp := range response.Items {
		for _, id := range trunkIds {
			if strings.Contains(strings.ToLower(resp.SipDispatchRuleId), strings.ToLower(id)) {
				matchedRules = append(matchedRules, resp)
				break
			}
		}
	}

	return matchedRules, nil
}

func DeleteSIPDispatchRule(Client *lksdk.SIPClient, ruleId string) error {
	request := &livekit.DeleteSIPDispatchRuleRequest{
		SipDispatchRuleId: ruleId,
	}

	_, err := Client.DeleteSIPDispatchRule(context.Background(), request)

	if err != nil {
		return err
	}
	return nil
}

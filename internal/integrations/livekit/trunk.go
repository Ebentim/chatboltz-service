package livekit

import (
	"context"
	"strings"

	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go/v2"
)

func CreateSIPTrunk(Client *lksdk.SIPClient, name, username, password, trunkId string, numbers []string) (*livekit.SIPInboundTrunkInfo, error) {
	trunkInfo := &livekit.SIPInboundTrunkInfo{
		AuthPassword: password,
		AuthUsername: username,
		Name:         name,
		Numbers:      numbers,
		SipTrunkId:   trunkId,
	}

	request := &livekit.CreateSIPInboundTrunkRequest{
		Trunk: trunkInfo,
	}

	trunk, err := Client.CreateSIPInboundTrunk(context.Background(), request)

	if err != nil {
		return nil, err
	}
	return trunk, nil
}

func GetSIPTrunkByName(sipClient *lksdk.SIPClient, trunkIds []string) ([]*livekit.SIPInboundTrunkInfo, error) {
	response, err := sipClient.ListSIPInboundTrunk(context.Background(), &livekit.ListSIPInboundTrunkRequest{
		TrunkIds: trunkIds,
	})
	if err != nil {
		return nil, err
	}

	var matchedTrunks []*livekit.SIPInboundTrunkInfo

	for _, resp := range response.Items {
		for _, id := range trunkIds {
			if strings.Contains(strings.ToLower(resp.SipTrunkId), strings.ToLower(id)) {
				matchedTrunks = append(matchedTrunks, resp)
				break
			}
		}
	}

	return matchedTrunks, nil
}

func DeleteSIPTrunk(Client *lksdk.SIPClient, trunkId string) error {
	request := &livekit.DeleteSIPTrunkRequest{
		SipTrunkId: trunkId,
	}
	_, err := Client.DeleteSIPTrunk(context.Background(), request)

	if err != nil {
		return err
	}
	return nil
}

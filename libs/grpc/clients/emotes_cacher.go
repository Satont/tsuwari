package clients

import (
	"fmt"
	"log"

	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/twirapp/twir/libs/grpc/emotes_cacher"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewEmotesCacher(env string) emotes_cacher.EmotesCacherClient {
	serverAddress := createClientAddr(env, "emotes-cacher", constants.EMOTES_CACHER_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(
				`{"loadBalancingConfig": [{"%s":{}}]}`,
				roundrobin.Name,
			),
		),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := emotes_cacher.NewEmotesCacherClient(conn)
	return c
}

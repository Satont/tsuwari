package clients

import (
	"fmt"
	"log"

	"github.com/satont/twir/libs/grpc/generated/scheduler"
	"github.com/satont/twir/libs/grpc/servers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewScheduler(env string) scheduler.SchedulerClient {
	serverAddress := createClientAddr(env, "scheduler", servers.SCHEDULER_SERVER_PORT)

	conn, err := grpc.Dial(
		serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":{}}]}`, roundrobin.Name)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := scheduler.NewSchedulerClient(conn)
	return c
}

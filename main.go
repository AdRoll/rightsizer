package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SemanticSugar/rightsizer/clients"
	"github.com/SemanticSugar/rightsizer/models"
	"github.com/SemanticSugar/rightsizer/services"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func main() {
	defaultDuration, _ := time.ParseDuration("336h")

	app := &cli.Command{
		Name:                   "rigthsizer",
		Usage:                  "Right size your AWS ECS services.",
		Version:                "3.0.0",
		HideHelpCommand:        true,
		ArgsUsage:              "<cluster> <service>",
		UseShortOptionHandling: true,

		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "time-frame",
				Aliases: []string{"t"},
				Value:   defaultDuration,
				Usage:   "Time `DURATION` to draw stats from",
			},
			&cli.StringFlag{
				Name:    "region",
				Aliases: []string{"r"},
				Usage:   "AWS `REGION` to use",
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() < 2 {
				return errors.New("invocation requires <cluster> and <service> parameters")
			}
			clusterName := cmd.Args().Get(0)
			serviceName := cmd.Args().Get(1)

			// Parse that time frame
			timeFrame := cmd.Duration("time-frame")
			if timeFrame < time.Hour {
				return errors.New("cannot see into the future just yet")
			}

			region := cmd.String("region")
			if region == "" {
				region = os.Getenv("AWS_REGION")
			}

			if region == "" {
				return errors.New("cannot determine AWS region, checked the --region flag and the AWS_REGION environment variable")
			}

			cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
			if err != nil {
				return fmt.Errorf("unable to load SDK config, %v", err)
			}

			awsCloudwatchClient := cloudwatch.NewFromConfig(cfg)
			cloutwatchClient := clients.NewCloudWatchClient(awsCloudwatchClient)
			usageService := services.NewUsageService(cloutwatchClient)

			awsEcsClient := ecs.NewFromConfig(cfg)
			ecsClient := clients.NewECSClient(awsEcsClient)
			allocationService := services.NewAllocationService(ecsClient)

			usage, err := usageService.GetUsage(ctx, &services.GetUsageInput{
				ClusterName: clusterName,
				ServiceName: serviceName,
				TimeFrame:   timeFrame,
			})

			if err != nil {
				return fmt.Errorf("failed to get usage: %w", err)
			}

			allocation, err := allocationService.GetAllocation(ctx, &services.GetAllocationInput{
				ClusterName: clusterName,
				ServiceName: serviceName,
			})

			if err != nil {
				return fmt.Errorf("failed to get allocation: %w", err)
			}

			newAllocation := allocation.Fix(usage, &models.Usage{CPU: 90, Memory: 90})

			bytes, err := yaml.Marshal(newAllocation)
			if err != nil {
				return fmt.Errorf("failed to marshal allocation: %w", err)
			}
			fmt.Printf("%s", string(bytes))

			return nil
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

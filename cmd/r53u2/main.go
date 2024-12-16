/*

          ______ _____        ___
   _____ / ____/|__  / __  __|__ \
  / ___//___ \   /_ < / / / /__/ /
 / /   ____/ / ___/ // /_/ // __/
/_/   /_____/ /____/ \__,_//____/

Route53Updater2
by michaelpeterswa (nw.codes)
2022

	"software development is clearly still a black art"
		- President William J. Clinton

*/

package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/alpineworks/ootel"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/robfig/cron/v3"
	"nw.codes/r53u2/internal/config"
	"nw.codes/r53u2/internal/ip"
	"nw.codes/r53u2/internal/logging"
	"nw.codes/r53u2/internal/util"
	"nw.codes/r53u2/internal/zones"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}

	slogLevel, err := logging.LogLevelToSlogLevel(logLevel)
	if err != nil {
		log.Fatalf("could not convert log level: %s", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	})))

	r53u2Config, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to parse config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx := context.Background()

	ootelClient := ootel.NewOotelClient(
		ootel.WithMetricConfig(
			ootel.NewMetricConfig(
				r53u2Config.MetricsEnabled,
				r53u2Config.MetricsPort,
			),
		),
		ootel.WithTraceConfig(
			ootel.NewTraceConfig(
				r53u2Config.TracingEnabled,
				r53u2Config.TracingSampleRate,
				r53u2Config.TracingService,
				r53u2Config.TracingVersion,
			),
		),
	)

	shutdown, err := ootelClient.Init(ctx)
	if err != nil {
		slog.Error("could not initialize ootel client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer func() {
		_ = shutdown(ctx)
	}()

	// ensure that dns records are updated on first check
	var previouslyStoredIP net.IP

	slog.Debug("starting aws session")
	awsSession, err := session.NewSession()
	if err != nil {
		slog.Error("failed to start aws session", slog.String("error", err.Error()))
	}

	slog.Debug("creating route53 client")
	r53 := route53.New(awsSession)

	ipClient := ip.NewIPClient(r53u2Config.CheckIPProvider, r53u2Config.CheckIPTimeout)

	c := cron.New()

	slog.Debug("adding cron function")
	_, err = c.AddFunc(r53u2Config.CronSchedule, func() {
		slog.Debug("pre ip-client get")
		currentIP, err := ipClient.Get()
		if err != nil {
			slog.Error("failed to acquire current ip address", slog.String("error", err.Error()))
			return
		}
		slog.Debug("acquired current ip address", slog.String("ip", currentIP.String()))
		if !currentIP.Equal(previouslyStoredIP) {
			hostedZones, err := r53.ListHostedZones(&route53.ListHostedZonesInput{
				MaxItems: aws.String("100"),
			})
			if err != nil {
				slog.Error("failed to list hosted zones", slog.String("error", err.Error()))
				return
			}

			// skipping pagination because it doesn't apply to me at this moment
			// (with MaxItems set in the request, pagination will not occur when zones <= 100)
			if *hostedZones.IsTruncated {
				slog.Warn("list of hosted zones is truncated", slog.Bool("isTruncated", *hostedZones.IsTruncated))
			}

			// match domains in the settings to hosted zones on Route53 and only update zones common to both listss
			for _, zone := range hostedZones.HostedZones {
				for _, domain := range r53u2Config.Domains {
					if util.GetURLFromZoneName(*zone.Name) == domain {
						err := zones.UpdateHostedZone(r53, zone, currentIP.String())
						if err != nil {
							slog.Error("failed to update hosted zone", slog.String("domain", domain))
						}
					}
				}
			}
			slog.Info("updated ip for route53 zones", slog.Int("zones", len(hostedZones.HostedZones)), slog.String("previous-ip", previouslyStoredIP.String()), slog.String("new-ip", currentIP.String()))
			previouslyStoredIP = currentIP
		}
	})
	if err != nil {
		slog.Error("failed to add cron function", slog.String("error", err.Error()))
		os.Exit(1)
	}

	c.Start()
	select {}
}

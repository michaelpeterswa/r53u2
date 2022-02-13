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
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	"nw.codes/r53u2/internal/ip"
	"nw.codes/r53u2/internal/logging"
	"nw.codes/r53u2/internal/settings"
	"nw.codes/r53u2/internal/util"
	"nw.codes/r53u2/internal/zones"
)

func main() {
	logger, err := logging.InitZap()
	if err != nil {
		log.Panicf("could not acquire zap logger: %s", err.Error())
	}
	logger.Info("r53u2 init...")

	r53u2Settings := settings.InitSettings(logger, "config/settings.yaml")
	err = settings.SetAWSEnvironment(r53u2Settings.AWS)
	if err != nil {
		logger.Error("failed to set AWS environment variables", zap.Error(err))
	}

	// ensure that dns records are updated on first check
	previouslyStoredIP := ""

	awsSession, err := session.NewSession()
	if err != nil {
		logger.Error("failed to start AWS session", zap.Error(err))
	}

	r53 := route53.New(awsSession)

	c := cron.New()
	err = c.AddFunc(r53u2Settings.CheckInterval, func() {
		currentIP, err := ip.Get(r53u2Settings.IPProvider)
		if err != nil {
			logger.Error("failed to acquire current ip address", zap.Error(err))
		}
		if currentIP != previouslyStoredIP {
			hostedZones, err := r53.ListHostedZones(&route53.ListHostedZonesInput{
				MaxItems: aws.String("100"),
			})
			if err != nil {
				logger.Error("failed to list hosted zones", zap.Error(err))
			}

			// skipping pagination because it doesn't apply to me at this moment
			// (with MaxItems set in the request, pagination will not occur when zones <= 100)
			if *hostedZones.IsTruncated {
				logger.Warn("list of hosted zones is truncated", zap.Bool("isTruncated", *hostedZones.IsTruncated))
			}

			// match domains in the settings to hosted zones on Route53 and only update zones common to both listss
			for _, zone := range hostedZones.HostedZones {
				for _, domain := range r53u2Settings.Domains {
					if util.GetURLFromZoneName(*zone.Name) == domain {
						err := zones.UpdateHostedZone(r53, zone, currentIP)
						if err != nil {
							logger.Error("failed to update hosted zone", zap.String("domain", domain))
						}
					}
				}
			}
			logger.Info("updated ip for route53 zones", zap.Int("zones", len(hostedZones.HostedZones)), zap.String("previous-ip", previouslyStoredIP), zap.String("new-ip", currentIP))
			previouslyStoredIP = currentIP
		}
	})
	if err != nil {
		logger.Fatal("failed to add cron function", zap.Error(err))
	}

	c.Start()
	select {}
}

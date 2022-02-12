package zones

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"nw.codes/r53u2/internal/util"
)

func UpdateHostedZone(r53 *route53.Route53, hz *route53.HostedZone, ip string) error {
	url := util.GetURLFromZoneName(*hz.Name)

	update := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: hz.Id,
		ChangeBatch: &route53.ChangeBatch{
			Comment: aws.String("updating ip with r53u2..."),
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(url),
						Type: aws.String("A"),
						TTL:  aws.Int64(300),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
					},
				},
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(fmt.Sprintf("*.%s", url)),
						Type: aws.String("A"),
						TTL:  aws.Int64(300),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
					},
				},
			},
		},
	}
	_, err := r53.ChangeResourceRecordSets(update)
	if err != nil {
		return err
	}

	return nil
}

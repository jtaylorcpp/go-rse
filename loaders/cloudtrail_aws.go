package loaders

import (
	"fmt"
	"rse"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
)

type CLoudTrailLoader struct{}

func (c CLoudTrailLoader) GetParamters() (int, int, float64) {
	return DefaultParameters()
}

func (c CLoudTrailLoader) Load() (*rse.Matrix, error) {

	servicesInUse := map[string]bool{}
	serviceActionsInUse := map[string]bool{}
	users := map[string]bool{}
	events := []*cloudtrail.Event{}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := cloudtrail.New(sess)

	input := &cloudtrail.LookupEventsInput{EndTime: aws.Time(time.Now())}

	err := svc.LookupEventsPages(input, func(output *cloudtrail.LookupEventsOutput, last bool) bool {
		for _, event := range output.Events {
			servicesInUse[*event.EventSource] = true
			serviceActionsInUse[fmt.Sprintf("%s:%s", *event.EventSource, *event.EventName)] = true
			users[*event.Username] = true
			events = append(events, event)
		}
		return !last
	})

	if err != nil {
		return nil, err
	}

	/*
		build NxM matrix w/ N being users M being services, service:Action
	*/

	m := rse.NewEmpty(len(users), len(servicesInUse)+len(serviceActionsInUse))

	// create labels and indices
	userLabels := make([]rse.LabelSet, len(users), 0)
	serviceLabels := make([]rse.LabelSet, len(servicesInUse)+len(serviceActionsInUse), 0)

	userIndex := map[string]int{}
	serviceIndex := map[string]int{}

	for k, _ := range users {
		userLabels = append(userLabels, rse.LabelSet{k, int64(len(userLabels))})
		userIndex[k] = len(userLabels) - 1
	}

	for k, _ := range servicesInUse {
		serviceLabels = append(serviceLabels, rse.LabelSet{k, int64(len(serviceLabels))})
		serviceIndex[k] = len(serviceLabels) - 1
	}

	for k, _ := range serviceActionsInUse {
		serviceLabels = append(serviceLabels, rse.LabelSet{k, int64(len(serviceLabels))})
		serviceIndex[k] = len(serviceLabels) - 1
	}

	m.WithRowLabels(userLabels...).WithColumnLabels(serviceLabels...)

	for _, event := range events {

	}

	return nil, nil
}

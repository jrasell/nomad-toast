package notifier

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/nomad/structs"
	"github.com/nlopes/slack"
)

func (n *Notifier) formatDeploymentMessage(d *api.Deployment) {

	f := make([]slack.AttachmentField, 1+len(d.TaskGroups))
	f[0] = slack.AttachmentField{Title: "Job", Value: d.JobID, Short: true}
	f[1] = slack.AttachmentField{Title: "Status", Value: d.StatusDescription, Short: true}

	for name, tg := range d.TaskGroups {

		dt := fmt.Sprintf("%s %s\n", deadlineText, n.formatTimeToSecond(tg.RequireProgressBy))
		da := fmt.Sprintf("%s %v\n", desiredAllocsText, tg.DesiredTotal)
		pa := fmt.Sprintf("%s %v\n", placedAllocsText, tg.PlacedAllocs)
		ua := fmt.Sprintf("%s %v\n", unhealthyAllocsText, tg.UnhealthyAllocs)

		f = append(f, slack.AttachmentField{
			Title: fmt.Sprintf("Task Group: %s", name),
			Value: dt + da + pa + ua,
			Short: false,
		})
	}

	m := &slack.Attachment{Fallback: "", Fields: f}
	m.Fields = f
	m.Title = fmt.Sprintf("Nomad Deployment Notification: %s", strings.ToUpper(n.nomadRegion))

	if n.config.ui.HashiUIEnabled {
		m.TitleLink = fmt.Sprintf("%s/nomad/%s/deployments/%s/info", n.config.ui.HashiUIHost, n.nomadRegion, d.ID)
	}

	switch d.Status {
	case structs.DeploymentStatusRunning, structs.DeploymentStatusPaused:
		m.Color = warningColour
	case structs.DeploymentStatusSuccessful:
		m.Color = goodColour
	default:
		m.Color = dangerColour
	}

	n.sendNotification(d.ID, *m)
}

func (n *Notifier) formatAllocationMessage(d *api.AllocationListStub) {

	f := make([]slack.AttachmentField, 4+len(d.TaskStates))
	f[0] = slack.AttachmentField{Title: "Job", Value: d.JobID, Short: true}
	f[1] = slack.AttachmentField{Title: "AllocID", Value: d.ID[0:7], Short: true}
	f[2] = slack.AttachmentField{Title: "Status", Value: d.ClientStatus, Short: true}
	f[3] = slack.AttachmentField{Title: "NodeID", Value: d.NodeID[0:7], Short: true}

	for n, ts := range d.TaskStates {

		cs := fmt.Sprintf("%s %s\n", taskClientStatus, d.ClientStatus)
		t := fmt.Sprintf("%s %s\n", taskStateText, ts.State)
		r := fmt.Sprintf("%s %v\n", taskRestartText, ts.Restarts)
		ff := fmt.Sprintf("%s %v\n", taskFailedText, ts.Failed)
		te := fmt.Sprintf("%s %s\n", taskEventText, ts.Events[len(ts.Events)-1].DisplayMessage)

		f = append(f, slack.AttachmentField{
			Title: fmt.Sprintf("Task Group: %s", n),
			Value: cs + t + r + ff + te,
			Short: false,
		})
	}

	m := &slack.Attachment{Fallback: "", Fields: f}
	m.Fields = f
	m.Title = fmt.Sprintf("Nomad Allocation Notification: %s", strings.ToUpper(n.nomadRegion))

	if n.config.ui.HashiUIEnabled {
		m.TitleLink = fmt.Sprintf("%s/nomad/%s/allocations/%s/info", n.config.ui.HashiUIHost, n.nomadRegion, d.ID)
	}

	switch d.ClientStatus {
	case structs.AllocClientStatusPending:
		m.Color = warningColour
	case structs.AllocClientStatusRunning, structs.AllocClientStatusComplete:
		m.Color = goodColour
	default:
		m.Color = dangerColour
	}

	n.sendNotification(d.ID, *m)
}

func (n *Notifier) formatTimeToSecond(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

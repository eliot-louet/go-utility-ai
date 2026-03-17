package utility_providers

import (
	"github.com/eliot-louet/go-utility-ai/ai"
	"github.com/eliot-louet/go-utility-ai/ai/mailbox"
)

type SentMailProvider struct {
	MailType mailbox.MailType
	CachedID ai.TargetProviderID
	buf      []ai.Target
}

type HaveMailbox interface {
	GetMailbox() *mailbox.Mailbox
}

func NewSentMailProvider(mailType mailbox.MailType) *SentMailProvider {
	provider := &SentMailProvider{
		MailType: mailType,
		CachedID: ai.TargetProviderID("sent_mail:" + string(mailType)),
	}

	return provider
}

func (p *SentMailProvider) Targets(ctx *ai.Context) []ai.Target {
	ag := ctx.Self.(HaveMailbox)

	p.buf = p.buf[:0]
	ag.GetMailbox().ForEachByType(p.MailType, func(mail mailbox.Mail) bool {
		p.buf = append(p.buf, ai.Target{
			ID:   mail.FromID,
			Type: ai.TargetTypeActor,
		})

		return true
	})

	return p.buf
}

func (p *SentMailProvider) ID() ai.TargetProviderID {
	return p.CachedID
}

func (p *SentMailProvider) ShouldCache() bool {
	return false
}

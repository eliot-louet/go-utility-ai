package mailbox

type MailType string

type Mail struct {
	FromID int
	Time   int64
	Expiry int64
	Type   MailType
}

type MailKey struct {
	FromID int
	Type   MailType
}

type Mailbox struct {
	messages map[MailKey]Mail
	Time     int64
}

func NewMailbox() *Mailbox {
	return &Mailbox{
		messages: make(map[MailKey]Mail),
		Time:     0,
	}
}

func (mb *Mailbox) CleanExpired() {
	for key, mail := range mb.messages {
		if mail.Time+mail.Expiry < mb.Time {
			delete(mb.messages, key)
		}
	}
}

func (mb *Mailbox) Get(senderID int, mtype MailType) (Mail, bool) {
	key := MailKey{FromID: senderID, Type: mtype}
	mail, exists := mb.messages[key]

	if !exists {
		return Mail{}, false
	}

	if mail.Time+mail.Expiry >= mb.Time {
		return mail, true
	} else {
		delete(mb.messages, key)
		return Mail{}, false
	}
}

// Add inserts or replaces a message with the same sender and type
func (mb *Mailbox) Add(mail Mail) {
	mail.Time = mb.Time

	key := MailKey{FromID: mail.FromID, Type: mail.Type}
	mb.messages[key] = mail
}

func (mb *Mailbox) ForEach(yield func(Mail) bool) {
	for key, m := range mb.messages {
		if m.Time+m.Expiry >= mb.Time {
			if !yield(m) {
				return
			}
		} else {
			delete(mb.messages, key)
		}
	}
}

func (mb *Mailbox) ForEachByType(mtype MailType, yield func(Mail) bool) {
	for key, m := range mb.messages {
		if m.Type != mtype {
			continue
		}

		if m.Time+m.Expiry >= mb.Time {
			if !yield(m) {
				return
			}
		} else {
			delete(mb.messages, key)
		}
	}
}

func (mb *Mailbox) ForEachBySender(senderID int, yield func(Mail) bool) {
	for key, m := range mb.messages {
		if m.FromID == senderID {
			if m.Time+m.Expiry >= mb.Time {
				if !yield(m) {
					return
				}
			} else {
				delete(mb.messages, key)
			}
		}
	}
}

func (mb *Mailbox) ForEachBySenderAndType(senderID int, mtype MailType, yield func(Mail) bool) {
	for key, m := range mb.messages {
		if m.FromID == senderID && m.Type == mtype {
			if m.Time+m.Expiry >= mb.Time {
				if !yield(m) {
					return
				}
			} else {
				delete(mb.messages, key)
			}
		}
	}
}

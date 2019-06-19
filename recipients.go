package apub

func RecipientMap(o *Object) map[string]bool {
	rec := make(map[string]bool)

	subObj := o.Object("object")
	addRecipients(rec, o)
	addRecipients(rec, o.Object("target"))
	addRecipients(rec, subObj)
	addRecipients(rec, subObj.Object("inReplyTo"))
	addRecipients(rec, o.Object("inReplyTo"))
	return rec
}

func Recipients(o *Object) []string {
	rec := RecipientMap(o)
	ids := make([]string, 0, len(rec))
	for id := range rec {
		ids = append(ids, id)
	}
	return ids
}

func addRecipients(recipients map[string]bool, o *Object) {
	for _, id := range o.BTo() {
		recipients[id] = true
	}
	for _, id := range o.BCC() {
		recipients[id] = true
	}
	for _, id := range o.To() {
		recipients[id] = true
	}
	for _, id := range o.CC() {
		recipients[id] = true
	}
	for _, id := range o.Audience() {
		recipients[id] = true
	}
	for _, id := range o.AttributedTo() {
		recipients[id] = true
	}
	if id := o.Str("actor"); len(id) > 0 {
		recipients[id] = true
	}

	for _, tag := range o.Tags() {
		if id := tag.DefaultValue(); len(id) > 0 && recipientTagTypes[tag.Type()] {
			recipients[id] = true
		}
	}
}

var recipientTagTypes = map[string]bool{
	"Application":       true,
	"Collection":        true,
	"Group":             true,
	"Mention":           true,
	"OrderedCollection": true,
	"Organization":      true,
	"Person":            true,
	"Service":           true,
}

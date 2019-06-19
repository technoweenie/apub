package apub

func CreateActivity(o *Object) *Object {
	act := map[string]interface{}{
		"@context": "https://www.w3.org/ns/activitystreams",
		"type":     "Create",
	}
	obj := make(map[string]interface{})
	for k, v := range o.data {
		if createActivityAttrs[k] {
			act[k] = v
		}
		if createActivityIgnored[k] {
			continue
		}
		obj[k] = v
	}

	act["object"] = obj
	return New(act)
}

var createActivityAttrs = map[string]bool{
	"audience":  true,
	"bcc":       true,
	"bto":       true,
	"cc":        true,
	"to":        true,
	"published": true,
}

var createActivityIgnored = map[string]bool{
	"@context": true,
}

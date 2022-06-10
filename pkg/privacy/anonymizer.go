package privacy

import (
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/envelope"
	"github.com/silverton-io/honeypot/pkg/util"
)

// Anonymize
func anonymize(val string) string {
	anonymizedVal := util.Md5(val)
	return anonymizedVal
}

func AnonymizeEnvelopes(envelopes []envelope.Envelope, c config.Privacy) []envelope.Envelope {
	var envs []envelope.Envelope
	for _, e := range envelopes {
		if c.Anonymize.Device.Ip {
			anonymizedIp := anonymize(e.Device.Ip)
			e.Device.Ip = anonymizedIp
		}
		if c.Anonymize.Device.Useragent {
			anonymizedUa := anonymize(e.Device.Useragent)
			e.Device.Useragent = anonymizedUa
		}
		if c.Anonymize.User.Id && e.User.Id != nil {
			anonymizedUserId := anonymize(*e.User.Id)
			e.User.Id = &anonymizedUserId
			e.User.AnonymousId = &anonymizedUserId
		}
		envs = append(envs, e)
	}
	return envs
}

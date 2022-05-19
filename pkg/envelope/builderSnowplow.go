package envelope

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/tidwall/gjson"
)

func buildSnowplowEnvelope(c *gin.Context, e snowplow.SnowplowEvent, m *meta.CollectorMeta) Envelope {
	n := buildCommonEnvelope(c, m)
	// Event Meta
	n.EventMeta.Protocol = protocol.SNOWPLOW
	// Pipeline
	n.Pipeline.Source.GeneratedTstamp = e.DvceCreatedTstamp
	n.Pipeline.Source.SentTstamp = e.DvceSentTstamp
	// Device
	n.Device.Id = e.DomainUserid
	n.Device.Os = Os{Timezone: e.OsTimezone}
	n.Device.Browser = Browser{
		Lang:           e.BrLang,
		Cookies:        e.BrCookies,
		ColorDepth:     e.BrColordepth,
		Charset:        e.DocCharset,
		ViewportSize:   e.ViewportSize,
		ViewportHeight: e.BrViewHeight,
		ViewportWidth:  e.BrViewWidth,
		DocumentSize:   e.DocSize,
		DocumentHeight: e.DocHeight,
		DocumentWidth:  e.DocWidth,
	}
	n.Device.Screen = Screen{
		Resolution: e.DvceScreenResolution,
		Height:     e.DvceScreenHeight,
		Width:      e.DvceScreenWidth,
	}
	// User
	n.User = User{
		Id:          e.Userid,
		Fingerprint: e.UserFingerprint,
	}
	// Session
	n.Session = Session{
		Id:  e.DomainSessionId,
		Idx: e.DomainSessionIdx,
	}
	// Page
	n.Page.Page.Url = *e.PageUrl
	n.Page.Page.Title = e.PageTitle
	n.Page.Page.Scheme = *e.PageUrlScheme
	n.Page.Page.Host = *e.PageUrlHost
	n.Page.Page.Port = *e.PageUrlPort
	n.Page.Page.Path = *e.PageUrlPath
	n.Page.Page.Query = e.PageUrlQuery
	n.Page.Page.Fragment = e.PageUrlFragment
	n.Page.Page.Medium = e.MktMedium
	n.Page.Page.Source = e.MktMedium
	n.Page.Page.Term = e.MktTerm
	n.Page.Page.Content = e.MktContent
	n.Page.Page.Campaign = e.MktCampaign
	// Page Referrer
	n.Page.Referrer.Url = *e.PageReferrer
	n.Page.Referrer.Scheme = *e.RefrUrlScheme
	n.Page.Referrer.Host = *e.RefrUrlHost
	n.Page.Referrer.Port = *e.RefrUrlPort
	n.Page.Referrer.Path = *e.RefrUrlPath
	n.Page.Referrer.Query = e.RefrUrlQuery
	n.Page.Referrer.Fragment = e.RefrUrlFragment
	n.Page.Referrer.Medium = e.RefrMedium
	n.Page.Referrer.Source = e.RefrSource
	n.Page.Referrer.Term = e.RefrTerm
	n.Page.Referrer.Content = e.RefrContent
	n.Page.Referrer.Campaign = e.RefrCampaign
	// Contexts
	n.Contexts = *e.Contexts
	n.Payload = e.SelfDescribingEvent
	return n
}

func BuildSnowplowEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	if c.Request.Method == "POST" {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not read request body")
			return envelopes
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := snowplow.BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), *conf)
			e := buildSnowplowEnvelope(c, spEvent, m)
			envelopes = append(envelopes, e)
		}
	} else {
		params := util.MapUrlParams(c)
		spEvent := snowplow.BuildEventFromMappedParams(c, params, *conf)
		e := buildSnowplowEnvelope(c, spEvent, m)
		envelopes = append(envelopes, e)
	}
	return envelopes
}

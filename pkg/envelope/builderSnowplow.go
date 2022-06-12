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
	n.EventMeta.Schema = *e.SelfDescribingEvent.SchemaName()
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
	n.Web.Page.Url = *e.PageUrl
	n.Web.Page.Title = e.PageTitle
	n.Web.Page.Scheme = *e.PageUrlScheme
	n.Web.Page.Host = *e.PageUrlHost
	n.Web.Page.Port = *e.PageUrlPort
	n.Web.Page.Path = *e.PageUrlPath
	n.Web.Page.Query = e.PageUrlQuery
	n.Web.Page.Fragment = e.PageUrlFragment
	n.Web.Page.Medium = e.MktMedium
	n.Web.Page.Source = e.MktMedium
	n.Web.Page.Term = e.MktTerm
	n.Web.Page.Content = e.MktContent
	n.Web.Page.Campaign = e.MktCampaign
	// Page Referrer
	n.Web.Referrer.Url = *e.PageReferrer
	n.Web.Referrer.Scheme = *e.RefrUrlScheme
	n.Web.Referrer.Host = *e.RefrUrlHost
	n.Web.Referrer.Port = *e.RefrUrlPort
	n.Web.Referrer.Path = *e.RefrUrlPath
	n.Web.Referrer.Query = e.RefrUrlQuery
	n.Web.Referrer.Fragment = e.RefrUrlFragment
	n.Web.Referrer.Medium = e.RefrMedium
	n.Web.Referrer.Source = e.RefrSource
	n.Web.Referrer.Term = e.RefrTerm
	n.Web.Referrer.Content = e.RefrContent
	n.Web.Referrer.Campaign = e.RefrCampaign
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

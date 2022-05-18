package envelope

import (
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/config"
	"github.com/silverton-io/honeypot/pkg/event"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/meta"
	"github.com/silverton-io/honeypot/pkg/pixel"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"

	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/silverton-io/honeypot/pkg/webhook"
	"github.com/tidwall/gjson"
)

func buildCommonEnvelope(c *gin.Context, m *meta.CollectorMeta) Envelope {
	nid := c.GetString("identity")
	envelope := Envelope{
		EventMeta: EventMeta{
			Uuid: uuid.New(),
		},
		Pipeline: Pipeline{
			Source: Source{
				Ip: c.ClientIP(),
			},
			Collector: Collector{
				Tstamp:  time.Now().UTC(),
				Name:    &m.Name,
				Version: &m.Version,
			},
			Relay: Relay{
				Relayed: false,
			},
		},
		Device: Device{
			Nid: &nid,
		},
		User:       User{},
		Session:    Session{},
		Page:       Page{},
		Validation: Validation{},
		Contexts:   event.Contexts{},
	}
	return envelope
}

func buildSnowplowEnvelope(c *gin.Context, e snowplow.SnowplowEvent, m *meta.CollectorMeta) Envelope {
	n := buildCommonEnvelope(c, m)
	// Event Meta
	n.EventMeta.Protocol = protocol.SNOWPLOW
	// Pipeline
	n.Pipeline.Source.GeneratedTstamp = e.DvceCreatedTstamp
	n.Pipeline.Source.SentTstamp = e.DvceSentTstamp
	// Device
	n.Device.Ip = e.UserIpAddress
	n.Device.Useragent = e.Useragent
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

func BuildGenericEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	advancingCookieVal, _ := c.Cookie(conf.Cookie.Name)
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		genEvent := generic.BuildEvent(e, conf.Generic)
		envelope := Envelope{
			EventMeta: EventMeta{
				Uuid:     uuid.New(),
				Protocol: protocol.GENERIC,
			},
			Pipeline: Pipeline{
				Source: Source{
					Ip: c.ClientIP(),
				},
				Collector: Collector{
					Tstamp: time.Now().UTC(),
				},
				Relay: Relay{
					Relayed: false,
				},
			},
			Device: Device{
				Nid: &advancingCookieVal,
			},
			Validation: Validation{},
			Payload:    genEvent,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func BuildCloudeventEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	advancingCookieVal, _ := c.Cookie(conf.Cookie.Name)
	for _, ce := range gjson.ParseBytes(reqBody).Array() {
		cEvent, _ := cloudevents.BuildEvent(ce)
		envelope := Envelope{
			EventMeta: EventMeta{
				Uuid:     uuid.New(),
				Protocol: protocol.GENERIC,
			},
			Pipeline: Pipeline{
				Source: Source{
					Ip: c.ClientIP(),
				},
				Collector: Collector{
					Tstamp: time.Now().UTC(),
				},
				Relay: Relay{
					Relayed: false,
				},
			},
			Device: Device{
				Nid: &advancingCookieVal,
			},
			Validation: Validation{},
			Payload:    cEvent,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func BuildWebhookEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	advancingCookieVal, _ := c.Cookie(conf.Cookie.Name)
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		whEvent, err := webhook.BuildEvent(c, e)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not build WebhookEvent")
		}
		envelope := Envelope{
			EventMeta: EventMeta{
				Uuid:     uuid.New(),
				Protocol: protocol.GENERIC,
			},
			Pipeline: Pipeline{
				Source: Source{
					Ip: c.ClientIP(),
				},
				Collector: Collector{
					Tstamp: time.Now().UTC(),
				},
				Relay: Relay{
					Relayed: false,
				},
			},
			Device: Device{
				Nid: &advancingCookieVal,
			},
			Validation: Validation{},
			Payload:    whEvent,
		}
		envelopes = append(envelopes, envelope)
	}
	return envelopes
}

func BuildPixelEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []Envelope {
	var envelopes []Envelope
	params := util.MapUrlParams(c)
	var urlNames = make(map[string]string)
	for _, i := range conf.Pixel.Paths {
		urlNames[i.Path] = i.Name
	}
	name := urlNames[c.Request.URL.Path]
	pEvent, err := pixel.BuildEvent(c, name, params)
	if err != nil {
		log.Error().Err(err).Msg("could not build PixelEvent")
	}
	advancingCookieVal, _ := c.Cookie(conf.Cookie.Name)
	envelope := Envelope{
		EventMeta: EventMeta{
			Uuid:     uuid.New(),
			Protocol: protocol.GENERIC,
		},
		Pipeline: Pipeline{
			Source: Source{
				Ip: c.ClientIP(),
			},
			Collector: Collector{
				Tstamp: time.Now().UTC(),
			},
			Relay: Relay{
				Relayed: false,
			},
		},
		Device: Device{
			Nid: &advancingCookieVal,
		},
		Validation: Validation{},
		Payload:    pEvent,
	}
	envelopes = append(envelopes, envelope)
	return envelopes
}

// func BuildRelayEnvelopesFromRequest(c *gin.Context) []Envelope {
// 	var envelopes []Envelope
// 	reqBody, err := ioutil.ReadAll(c.Request.Body)
// 	if err != nil {
// 		log.Error().Stack().Err(err).Msg("could not read request body")
// 		envelope := Envelope{}
// 		envelopes = append(envelopes, envelope)
// 		return envelopes
// 	}
// 	relayedEvents := gjson.ParseBytes(reqBody)
// 	for _, relayedEvent := range relayedEvents.Array() {
// 		uid := uuid.New()
// 		now := time.Now().UTC()
// 		envelope.Pipeline.Relay = Relay{
// 			Relayed: true,
// 			Id:      &uid,
// 			Tstamp:  &now,
// 		}
// 		envelopes = append(envelopes, envelope)
// 	}
// 	return envelopes
// }

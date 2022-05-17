package envelope

import (
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/honeypot/pkg/cloudevents"
	"github.com/silverton-io/honeypot/pkg/generic"
	"github.com/silverton-io/honeypot/pkg/handler"
	"github.com/silverton-io/honeypot/pkg/pixel"
	"github.com/silverton-io/honeypot/pkg/protocol"
	"github.com/silverton-io/honeypot/pkg/snowplow"
	"github.com/silverton-io/honeypot/pkg/util"
	"github.com/silverton-io/honeypot/pkg/webhook"
	"github.com/tidwall/gjson"
)

func buildSnowplowEnvelope(c *gin.Context, e snowplow.SnowplowEvent, h handler.EventHandlerParams) Envelope {

	envelope := Envelope{
		EventMeta: EventMeta{
			Protocol: protocol.SNOWPLOW,
			Uuid:     uuid.New(),
		},
		Pipeline: Pipeline{
			Source: Source{
				Ip:              *e.UserIpAddress,
				GeneratedTstamp: e.DvceCreatedTstamp,
				SentTstamp:      e.DvceSentTstamp,
				Name:            &e.NameTracker,
				Version:         e.TrackerVersion,
			},
			Collector: Collector{
				Tstamp: time.Now().UTC(),
				// Name:   "honeypot", // FIXME!
				// Version:
			},
			Relay: Relay{
				Relayed: false,
			},
		},
		Device: Device{
			Ip:        e.UserIpAddress,
			Useragent: e.Useragent,
			Id:        e.DomainUserid,
			Nid:       e.NetworkUserid,
			Os:        Os{}, // FIXME
			Browser: Browser{
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
			},
			Screen: Screen{
				Resolution: e.DvceScreenResolution,
				Height:     e.DvceScreenHeight,
				Width:      e.DvceScreenWidth,
			},
		},
		User: User{
			Id:          e.Userid,
			Fingerprint: e.UserFingerprint,
		},
		Session: Session{
			Id:  e.DomainSessionId,
			Idx: e.DomainSessionIdx,
		},
		Page: Page{
			Page: PageAttrs{
				Url:    *e.PageUrl,
				Title:  e.PageTitle,
				Scheme: *e.PageUrlScheme,
				Host:   *e.PageUrlHost,
				Port:   *e.PageUrlPort,
				Path:   *e.PageUrlPath,
				// Query:  e.PageUrlQuery,
				Fragment: e.PageUrlFragment,
				// Medium:  ,
				// Source: ,
				// Term: ,
				// Content: ,
				// Campaign: ,
			},
			Referrer: PageAttrs{
				Url:    *e.PageReferrer,
				Scheme: *e.RefrUrlScheme,
				Host:   *e.RefrUrlHost,
				Port:   *e.RefrUrlPort,
				Path:   *e.RefrUrlPath,
				// Query:  e.RefrUrlQuery,
				Fragment: e.RefrUrlFragment,
				// Medium:  ,
				// Source: ,
				// Term: ,
				// Content: ,
				// Campaign: ,
			},
		},
		Validation: Validation{},
		Contexts:   *e.Contexts,
		Payload:    e.SelfDescribingEvent,
	}
	return envelope
}

func BuildSnowplowEnvelopesFromRequest(c *gin.Context, h handler.EventHandlerParams) []Envelope {
	var envelopes []Envelope
	if c.Request.Method == "POST" {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Stack().Err(err).Msg("could not read request body")
			return envelopes
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := snowplow.BuildEventFromMappedParams(c, event.Value().(map[string]interface{}), h)
			e := buildSnowplowEnvelope(c, spEvent, h)
			envelopes = append(envelopes, e)
		}
	} else {
		params := util.MapUrlParams(c)
		spEvent := snowplow.BuildEventFromMappedParams(c, params, h)
		e := buildSnowplowEnvelope(c, spEvent, h)
		envelopes = append(envelopes, e)
	}
	return envelopes
}

func BuildGenericEnvelopesFromRequest(c *gin.Context, h handler.EventHandlerParams) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	advancingCookieVal, _ := c.Cookie(h.Config.Cookie.Name)
	for _, e := range gjson.ParseBytes(reqBody).Array() {
		genEvent := generic.BuildEvent(e, h.Config.Generic)
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

func BuildCloudeventEnvelopesFromRequest(c *gin.Context, h handler.EventHandlerParams) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	advancingCookieVal, _ := c.Cookie(h.Config.Cookie.Name)
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

func BuildWebhookEnvelopesFromRequest(c *gin.Context, h handler.EventHandlerParams) []Envelope {
	var envelopes []Envelope
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Stack().Err(err).Msg("could not read request body")
		return envelopes
	}
	advancingCookieVal, _ := c.Cookie(h.Config.Cookie.Name)
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

func BuildPixelEnvelopesFromRequest(c *gin.Context, h handler.EventHandlerParams) []Envelope {
	var envelopes []Envelope
	params := util.MapUrlParams(c)
	var urlNames = make(map[string]string)
	for _, i := range h.Config.Pixel.Paths {
		urlNames[i.Path] = i.Name
	}
	name := urlNames[c.Request.URL.Path]
	pEvent, err := pixel.BuildEvent(c, name, params)
	if err != nil {
		log.Error().Err(err).Msg("could not build PixelEvent")
	}
	advancingCookieVal, _ := c.Cookie(h.Config.Cookie.Name)
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

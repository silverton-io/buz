// Copyright (c) 2022 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package snowplow

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/silverton-io/buz/pkg/config"
	"github.com/silverton-io/buz/pkg/envelope"
	"github.com/silverton-io/buz/pkg/meta"
	"github.com/silverton-io/buz/pkg/protocol"
	"github.com/silverton-io/buz/pkg/util"
	"github.com/tidwall/gjson"
)

func buildSnowplowEnvelope(c *gin.Context, conf *config.Config, e SnowplowEvent, m *meta.CollectorMeta) envelope.Envelope {
	n := envelope.BuildCommonEnvelope(c, conf.Middleware, m)
	// Event Meta
	n.EventMeta.Protocol = protocol.SNOWPLOW
	n.EventMeta.Schema = *e.SelfDescribingEvent.SchemaName()
	// Pipeline
	n.Pipeline.Source.GeneratedTstamp = e.DvceCreatedTstamp
	n.Pipeline.Source.SentTstamp = e.DvceSentTstamp
	// Device
	if e.DomainUserid != nil {
		n.Device.Id = *e.DomainUserid
	}
	n.Device.Os = &envelope.Os{Timezone: e.OsTimezone}
	n.Device.Browser = &envelope.Browser{
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
	n.Device.Screen = &envelope.Screen{
		Resolution: e.DvceScreenResolution,
		Height:     e.DvceScreenHeight,
		Width:      e.DvceScreenWidth,
	}
	// User
	n.User = &envelope.User{
		Id:          e.Userid,
		Fingerprint: e.UserFingerprint,
	}
	// Session
	n.Session = &envelope.Session{
		Id:  e.DomainSessionId,
		Idx: e.DomainSessionIdx,
	}
	// Page
	n.Web = &envelope.Web{
		Page:     envelope.PageAttrs{},
		Referrer: envelope.PageAttrs{},
	}
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
	n.Contexts = e.Contexts
	n.Payload = e.SelfDescribingEvent.Data
	return n
}

func buildEnvelopesFromRequest(c *gin.Context, conf *config.Config, m *meta.CollectorMeta) []envelope.Envelope {
	var envelopes []envelope.Envelope
	if c.Request.Method == "POST" {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Error().Err(err).Msg("🔴 could not read request body")
			return envelopes
		}
		payloadData := gjson.GetBytes(body, "data")
		for _, event := range payloadData.Array() {
			spEvent := buildEventFromMappedParams(c, event.Value().(map[string]interface{}), *conf)
			e := buildSnowplowEnvelope(c, conf, spEvent, m)
			envelopes = append(envelopes, e)
		}
	} else {
		params := util.MapUrlParams(c)
		spEvent := buildEventFromMappedParams(c, params, *conf)
		e := buildSnowplowEnvelope(c, conf, spEvent, m)
		envelopes = append(envelopes, e)
	}
	return envelopes
}

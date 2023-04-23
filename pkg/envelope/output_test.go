// Copyright (c) 2023 Silverton Data, Inc.
// You may use, distribute, and modify this code under the terms of the Apache-2.0 license, a copy of
// which may be found at https://github.com/silverton-io/buz/blob/main/LICENSE

package envelope

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOutputLocationFromEnvelope(t *testing.T) {
	d, _ := time.Parse("2006-01-02", "2023-03-10")
	fakeEnvelope := Envelope{
		IsValid:      true,
		Vendor:       "com.vendor",
		Namespace:    "customer.activity",
		Version:      "1.1",
		BuzTimestamp: d,
	}
	expectedPath := "isValid=true/vendor=com.vendor/namespace=customer.activity/version=1.1/year=2023/month=3/day=10/"
	expectedDatabaseFqn := "com_vendor.customer_activity_1"
	expectedNamespace := "com.vendor.customer.activity.1.1"

	outputLocation := NewOutputLocationFromEnvelope(&fakeEnvelope)
	assert.Equal(t, expectedPath, outputLocation.Path)
	assert.Equal(t, expectedDatabaseFqn, outputLocation.DatabaseFqn)
	assert.Equal(t, expectedNamespace, outputLocation.Namespace)
}

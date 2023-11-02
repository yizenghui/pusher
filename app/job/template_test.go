// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"log"
	"testing"
	"time"
)

func Test_MessageTemplateSend(t *testing.T) {

	msg := Template{
		OpenID:      `opendid`,
		ReqID:       `xx`,
		TemplateID:  `xx`,
		Key1:        `xx`,
		Value1:      `xx`,
		Key2:        `xx`,
		Value2:      `xx`,
		Key3:        `xx`,
		Value3:      `xx`,
		Key4:        `xx`,
		Value4:      `xx`,
		Key5:        `xx`,
		Value5:      `xx`,
		AccessToken: ``,
	}
	// `62__O0B3S26sBLWoxc1feUXDR20eKrZL1xMp1j7ratjnX8QBWiOSs_bjC3w1PS_0i_3rR5C9vFsjrqoiiYWkXJetobz2j-kS58o1LCUjRnfNWFQvVZJSB8bpAfToclmHepIt_qx24ixGD18LdznBECfACAXAY`
	JobQueue <- Job{
		Player: &msg,
	}

	log.Println(`JobQueue len `, len(JobQueue))
	if len(JobQueue) > 0 {
		close(JobQueue)
	}

	time.Sleep(time.Second)

	t.Fatal(msg)

}

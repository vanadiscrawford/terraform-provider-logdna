package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var alertDefaults = cloneDefaults(rsDefaults["alert"])

func TestAlert_ErrorProviderUrl(t *testing.T) {
	pcArgs := []string{serviceKey, "https://api.logdna.co"}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("alert", "new", pcArgs, alertDefaults, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"https://api.logdna.co/v1/config/presetalert\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestAlert_ErrorResourceName(t *testing.T) {
	args := cloneDefaults(chnlDefaults["alert"])
	args["name"] = ""

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("alert", "new", nilLst, args, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_ErrorsChannel(t *testing.T) {
	imArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	imArgs["email"]["immediate"] = `"not a bool"`
	immdte := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, imArgs, nilLst)

	opArgs := map[string]map[string]string{"pagerduty": cloneDefaults(chnlDefaults["pagerduty"])}
	opArgs["pagerduty"]["operator"] = `1000`
	opratr := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, opArgs, nilLst)

	trArgs := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	trArgs["webhook"]["terminal"] = `"invalid"`
	trmnal := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, trArgs, nilLst)

	tiArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	tiArgs["email"]["triggerinterval"] = `18`
	tintvl := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, tiArgs, nilLst)

	tlArgs := map[string]map[string]string{"slack": cloneDefaults(chnlDefaults["slack"])}
	tlArgs["slack"]["triggerlimit"] = `0`
	tlimit := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, tlArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      immdte,
				ExpectError: regexp.MustCompile(`"\\"channels\[0\].immediate\\" must be a boolean"`),
			},
			{
				Config:      opratr,
				ExpectError: regexp.MustCompile(`"\\"channels\[0\].operator\\" must be one of \[presence, absence\]"`),
			},
			{
				Config:      trmnal,
				ExpectError: regexp.MustCompile(`"\\"channels\[0\].terminal\\" must be a boolean"`),
			},
			{
				Config:      tintvl,
				ExpectError: regexp.MustCompile(`"\\"channels\[0\].triggerinterval\\" must be one of \[15m, 30m, 1h, 6h, 12h, 24h\]"`),
			},
			{
				Config:      tlimit,
				ExpectError: regexp.MustCompile(`Error: ".*channel.0.triggerlimit" must be between 1 and 100,000 inclusive`),
			},
		},
	})
}

func TestAlert_ErrorsEmailChannel(t *testing.T) {
	msArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	msArgs["email"]["emails"] = ""
	misngE := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, msArgs, nilLst)

	inArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	inArgs["email"]["emails"] = `"not an array of strings"`
	invldE := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, inArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      misngE,
				ExpectError: regexp.MustCompile("The argument \"emails\" is required, but no definition was found."),
			},
			{
				Config:      invldE,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"emails\": list of string required"),
			},
		},
	})
}

func TestAlert_ErrorsPagerDutyChannel(t *testing.T) {
	chArgs := map[string]map[string]string{"pagerduty": cloneDefaults(chnlDefaults["pagerduty"])}
	chArgs["pagerduty"]["key"] = ""

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("alert", "new", nilLst, alertDefaults, chArgs, nilLst),
				ExpectError: regexp.MustCompile("The argument \"key\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_ErrorsSlackChannel(t *testing.T) {
	ulInvd := map[string]map[string]string{"slack": cloneDefaults(chnlDefaults["slack"])}
	ulInvd["slack"]["url"] = `"this is not a valid url"`
	ulCfgE := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, ulInvd, nilLst)

	ulMsng := map[string]map[string]string{"slack": cloneDefaults(chnlDefaults["slack"])}
	ulMsng["slack"]["url"] = ""
	ulCfgM := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, ulMsng, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      ulCfgE,
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\]\.url\\" must be a valid uri"`),
			},
			{
				Config:      ulCfgM,
				ExpectError: regexp.MustCompile("The argument \"url\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_ErrorsWebhookChannel(t *testing.T) {
	btArgs := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	btArgs["webhook"]["bodytemplate"] = `"{\"test\": }"`
	btCfgE := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, btArgs, nilLst)

	mdArgs := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	mdArgs["webhook"]["method"] = `"false"`
	mdCfgE := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, mdArgs, nilLst)

	ulInvd := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	ulInvd["webhook"]["url"] = `"this is not a valid url"`
	ulCfgE := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, ulInvd, nilLst)

	ulMsng := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	ulMsng["webhook"]["url"] = ""
	ulCfgM := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, ulMsng, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      btCfgE,
				ExpectError: regexp.MustCompile("Error: bodytemplate is not a valid JSON string"),
			},
			{
				Config:      mdCfgE,
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\].method\\" must be one of \[post, put, patch, get, delete\]"`),
			},
			{
				Config:      ulCfgE,
				ExpectError: regexp.MustCompile(`"message":"\\"channels\[0\]\.url\\" must be a valid uri"`),
			},
			{
				Config:      ulCfgM,
				ExpectError: regexp.MustCompile("The argument \"url\" is required, but no definition was found."),
			},
		},
	})
}

func TestAlert_Basic(t *testing.T) {
	chArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	iniCfg := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, chArgs, nilLst)

	rsArgs := cloneDefaults(rsDefaults["alert"])
	rsArgs["name"] = `"test2"`
	updCfg := fmtTestConfigResource("alert", "new", nilLst, rsArgs, chArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: iniCfg,
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: updCfg,
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test2"),
				),
			},
			{
				ResourceName:      "logdna_alert.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAlert_BulkChannels(t *testing.T) {
	emArgs := map[string]map[string]string{
		"email":  cloneDefaults(chnlDefaults["email"]),
		"email1": cloneDefaults(chnlDefaults["email"]),
	}
	emsCfg := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, emArgs, nilLst)

	pdArgs := map[string]map[string]string{
		"pagerduty":  cloneDefaults(chnlDefaults["pagerduty"]),
		"pagerduty1": cloneDefaults(chnlDefaults["pagerduty"]),
	}
	pdsCfg := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, pdArgs, nilLst)

	slArgs := map[string]map[string]string{
		"slack":  cloneDefaults(chnlDefaults["slack"]),
		"slack1": cloneDefaults(chnlDefaults["slack"]),
	}
	slsCfg := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, slArgs, nilLst)

	wbArgs := map[string]map[string]string{
		"webhook":  cloneDefaults(chnlDefaults["webhook"]),
		"webhook1": cloneDefaults(chnlDefaults["webhook"]),
	}
	wbsCfg := fmtTestConfigResource("alert", "new", nilLst, alertDefaults, wbArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: emsCfg,
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.1.%", "7"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: pdsCfg,
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.1.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: slsCfg,
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.1.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: wbsCfg,
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.%", "9"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.1.%", "9"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "0"),
				),
			},
		},
	})
}

func TestAlert_MultipleChannels(t *testing.T) {
	chArgs := map[string]map[string]string{
		"email":     cloneDefaults(chnlDefaults["email"]),
		"pagerduty": cloneDefaults(chnlDefaults["pagerduty"]),
		"slack":     cloneDefaults(chnlDefaults["slack"]),
		"webhook":   cloneDefaults(chnlDefaults["webhook"]),
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmtTestConfigResource("alert", "new", nilLst, alertDefaults, chArgs, nilLst),
				Check: resource.ComposeTestCheckFunc(
					testAlertExists("logdna_alert.new"),
					resource.TestCheckResourceAttr("logdna_alert.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.key", "Your PagerDuty API key goes here"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "pagerduty_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.triggerinterval", "30m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "slack_channel.0.url", "https://hooks.slack.com/services/identifier/secret"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.%", "9"),
					// The JSON will have newlines per our API which uses JSON.stringify(obj, null, 2) as the value
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.bodytemplate", "{\n  \"fields\": {\n    \"description\": \"{{ matches }} matches found for {{ name }}\",\n    \"issuetype\": {\n      \"name\": \"Bug\"\n    },\n    \"project\": {\n      \"key\": \"test\"\n    },\n    \"summary\": \"Alert from {{ name }}\"\n  }\n}"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.headers.%", "2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.headers.hello", "test3"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.headers.test", "test2"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.method", "post"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_alert.new", "webhook_channel.0.url", "https://yourwebhook/endpoint"),
				),
			},
			{
				ResourceName:      "logdna_alert.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAlertExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set")
		}
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		return nil
	}
}

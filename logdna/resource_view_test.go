package logdna

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const ctgies = `["DEMOCATEGORY1", "DemoCategory2"]`

var viewDefaults = cloneDefaults(rsDefaults["view"])

func TestView_ErrorProviderUrl(t *testing.T) {
	pcArgs := []string{serviceKey, "https://api.logdna.co"}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("view", "new", pcArgs, viewDefaults, nilOpt, nilLst),
				ExpectError: regexp.MustCompile("Error: error during HTTP request: Post \"https://api.logdna.co/v1/config/view\": dial tcp: lookup api.logdna.co"),
			},
		},
	})
}

func TestView_ErrorsResourceFields(t *testing.T) {
	nme := cloneDefaults(rsDefaults["view"])
	nme["name"] = ""
	nmeCfg := fmtTestConfigResource("view", "new", nilLst, nme, nilOpt, nilLst)

	app := cloneDefaults(rsDefaults["view"])
	app["apps"] = `"invalid apps value"`
	appCfg := fmtTestConfigResource("view", "new", nilLst, app, nilOpt, nilLst)

	ctg := cloneDefaults(rsDefaults["view"])
	ctg["categories"] = `"invalid categories value"`
	ctgCfg := fmtTestConfigResource("view", "new", nilLst, ctg, nilOpt, nilLst)

	hst := cloneDefaults(rsDefaults["view"])
	hst["hosts"] = `"invalid hosts value"`
	hstCfg := fmtTestConfigResource("view", "new", nilLst, hst, nilOpt, nilLst)

	lvl := cloneDefaults(rsDefaults["view"])
	lvl["levels"] = `"invalid levels value"`
	lvlCfg := fmtTestConfigResource("view", "new", nilLst, lvl, nilOpt, nilLst)

	tgs := cloneDefaults(rsDefaults["view"])
	tgs["tags"] = `"invalid tags value"`
	tgsCfg := fmtTestConfigResource("view", "new", nilLst, tgs, nilOpt, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      nmeCfg,
				ExpectError: regexp.MustCompile("The argument \"name\" is required, but no definition was found."),
			},
			{
				Config:      appCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"apps\": list of string required."),
			},
			{
				Config:      ctgCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"categories\": list of string required."),
			},
			{
				Config:      hstCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"hosts\": list of string required."),
			},
			{
				Config:      lvlCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"levels\": list of string required."),
			},
			{
				Config:      tgsCfg,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"tags\": list of string required."),
			},
		},
	})
}

func TestView_ErrorsChannel(t *testing.T) {
	imArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	imArgs["email"]["immediate"] = `"not a bool"`
	immdte := fmtTestConfigResource("view", "new", nilLst, viewDefaults, imArgs, nilLst)

	opArgs := map[string]map[string]string{"pagerduty": cloneDefaults(chnlDefaults["pagerduty"])}
	opArgs["pagerduty"]["operator"] = `1000`
	opratr := fmtTestConfigResource("view", "new", nilLst, viewDefaults, opArgs, nilLst)

	trArgs := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	trArgs["webhook"]["terminal"] = `"invalid"`
	trmnal := fmtTestConfigResource("view", "new", nilLst, viewDefaults, trArgs, nilLst)

	tiArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	tiArgs["email"]["triggerinterval"] = `18`
	tintvl := fmtTestConfigResource("view", "new", nilLst, viewDefaults, tiArgs, nilLst)

	tlArgs := map[string]map[string]string{"slack": cloneDefaults(chnlDefaults["slack"])}
	tlArgs["slack"]["triggerlimit"] = `0`
	tlimit := fmtTestConfigResource("view", "new", nilLst, viewDefaults, tlArgs, nilLst)

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

func TestView_ErrorsEmailChannel(t *testing.T) {
	msArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	msArgs["email"]["emails"] = ""
	misngE := fmtTestConfigResource("view", "new", nilLst, viewDefaults, msArgs, nilLst)

	inArgs := map[string]map[string]string{"email": cloneDefaults(chnlDefaults["email"])}
	inArgs["email"]["emails"] = `"not an array of strings"`
	invldE := fmtTestConfigResource("view", "new", nilLst, viewDefaults, inArgs, nilLst)

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

func TestView_ErrorsPagerDutyChannel(t *testing.T) {
	chArgs := map[string]map[string]string{"pagerduty": cloneDefaults(chnlDefaults["pagerduty"])}
	chArgs["pagerduty"]["key"] = ""

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      fmtTestConfigResource("view", "new", nilLst, viewDefaults, chArgs, nilLst),
				ExpectError: regexp.MustCompile("The argument \"key\" is required, but no definition was found."),
			},
		},
	})
}

func TestView_ErrorsSlackChannel(t *testing.T) {
	ulInvd := map[string]map[string]string{"slack": cloneDefaults(chnlDefaults["slack"])}
	ulInvd["slack"]["url"] = `"this is not a valid url"`
	ulCfgE := fmtTestConfigResource("view", "new", nilLst, viewDefaults, ulInvd, nilLst)

	ulMsng := map[string]map[string]string{"slack": cloneDefaults(chnlDefaults["slack"])}
	ulMsng["slack"]["url"] = ""
	ulCfgM := fmtTestConfigResource("view", "new", nilLst, viewDefaults, ulMsng, nilLst)

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

func TestView_ErrorsWebhookChannel(t *testing.T) {
	btArgs := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	btArgs["webhook"]["bodytemplate"] = `"{\"test\": }"`
	btCfgE := fmtTestConfigResource("view", "new", nilLst, viewDefaults, btArgs, nilLst)

	hdArgs := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	hdArgs["webhook"]["headers"] = `["headers", "invalid", "array"]`
	hdCfgE := fmtTestConfigResource("view", "new", nilLst, viewDefaults, hdArgs, nilLst)

	mdArgs := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	mdArgs["webhook"]["method"] = `"false"`
	mdCfgE := fmtTestConfigResource("view", "new", nilLst, viewDefaults, mdArgs, nilLst)

	ulInvd := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	ulInvd["webhook"]["url"] = `"this is not a valid url"`
	ulCfgE := fmtTestConfigResource("view", "new", nilLst, viewDefaults, ulInvd, nilLst)

	ulMsng := map[string]map[string]string{"webhook": cloneDefaults(chnlDefaults["webhook"])}
	ulMsng["webhook"]["url"] = ""
	ulCfgM := fmtTestConfigResource("view", "new", nilLst, viewDefaults, ulMsng, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      btCfgE,
				ExpectError: regexp.MustCompile("Error: bodytemplate is not a valid JSON string"),
			},
			{
				Config:      hdCfgE,
				ExpectError: regexp.MustCompile("Inappropriate value for attribute \"headers\": map of string required"),
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

func TestView_Basic(t *testing.T) {
	iniCfg := fmtTestConfigResource("view", "new", nilLst, viewDefaults, nilOpt, nilLst)

	rsArgs := cloneDefaults(rsDefaults["view"])
	rsArgs["name"] = `"test2"`
	rsArgs["query"] = `"test2"`
	updCfg := fmtTestConfigResource("view", "new", nilLst, rsArgs, nilOpt, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: iniCfg,
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
				),
			},
			{
				Config: updCfg,
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test2"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test2"),
				),
			},
			{
				ResourceName:      "logdna_view.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestView_BulkChannels(t *testing.T) {
	emArgs := map[string]map[string]string{
		"email":  cloneDefaults(chnlDefaults["email"]),
		"email1": cloneDefaults(chnlDefaults["email"]),
	}
	emsCfg := fmtTestConfigResource("view", "new", nilLst, viewDefaults, emArgs, nilLst)

	pdArgs := map[string]map[string]string{
		"pagerduty":  cloneDefaults(chnlDefaults["pagerduty"]),
		"pagerduty1": cloneDefaults(chnlDefaults["pagerduty"]),
	}
	pdsCfg := fmtTestConfigResource("view", "new", nilLst, viewDefaults, pdArgs, nilLst)

	slArgs := map[string]map[string]string{
		"slack":  cloneDefaults(chnlDefaults["slack"]),
		"slack1": cloneDefaults(chnlDefaults["slack"]),
	}
	slsCfg := fmtTestConfigResource("view", "new", nilLst, viewDefaults, slArgs, nilLst)

	wbArgs := map[string]map[string]string{
		"webhook":  cloneDefaults(chnlDefaults["webhook"]),
		"webhook1": cloneDefaults(chnlDefaults["webhook"]),
	}
	wbsCfg := fmtTestConfigResource("view", "new", nilLst, viewDefaults, wbArgs, nilLst)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: emsCfg,
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.%", "7"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.1.%", "7"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: pdsCfg,
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.1.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: slsCfg,
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.1.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "0"),
				),
			},
			{
				Config: wbsCfg,
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.%", "9"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.1.%", "9"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "0"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "0"),
				),
			},
		},
	})
}

func TestView_MultipleChannels(t *testing.T) {
	chArgs := map[string]map[string]string{
		"email":     cloneDefaults(chnlDefaults["email"]),
		"pagerduty": cloneDefaults(chnlDefaults["pagerduty"]),
		"slack":     cloneDefaults(chnlDefaults["slack"]),
		"webhook":   cloneDefaults(chnlDefaults["webhook"]),
	}

	dependencies := []string{"logdna_category.cat_1", "logdna_category.cat_2"}

	cat1Args := map[string]string{
		"name": `"DemoCategory1"`,
		"type": `"views"`,
	}
	cat2Args := map[string]string{
		"name": `"DemoCategory2"`,
		"type": `"views"`,
	}

	rsArgs := cloneDefaults(rsDefaults["view"])
	rsArgs["apps"] = `["app1", "app2"]`
	rsArgs["categories"] = ctgies
	rsArgs["hosts"] = `["host1", "host2"]`
	rsArgs["levels"] = `["fatal", "critical"]`
	rsArgs["tags"] = `["tags1", "tags2"]`
	iniCfg := fmt.Sprintf(
		"%s\n%s\n%s",
		fmtTestConfigResource("view", "new", nilLst, rsArgs, chArgs, dependencies),
		fmtResourceBlock("category", "cat_1", cat1Args, nilOpt, nilLst),
		fmtResourceBlock("category", "cat_2", cat2Args, nilOpt, nilLst),
	)

	rsUptd := cloneDefaults(rsDefaults["view"])
	rsUptd["apps"] = `["app3", "app4"]`
	rsUptd["categories"] = ctgies
	rsUptd["hosts"] = `["host3", "host4"]`
	rsUptd["levels"] = `["error", "warning"]`
	rsUptd["tags"] = `["tags3", "tags4"]`
	rsUptd["name"] = `"test2"`
	rsUptd["query"] = `"query2"`
	updCfg := fmt.Sprintf(
		"%s\n%s\n%s",
		fmtTestConfigResource("view", "new", nilLst, rsUptd, chArgs, dependencies),
		fmtResourceBlock("category", "cat_1", cat1Args, nilOpt, nilLst),
		fmtResourceBlock("category", "cat_2", cat2Args, nilOpt, nilLst),
	)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: iniCfg,
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "test"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.0", "app1"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.1", "app2"),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "categories.0", "DemoCategory1"), // This value on the server is mixed case
					resource.TestCheckResourceAttr("logdna_view.new", "categories.1", "DemoCategory2"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", "host1"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.1", "host2"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.0", "fatal"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", "critical"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.#", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.0", "tags1"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.1", "tags2"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.emails.0", "test@logdna.com"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.timezone", "Pacific/Samoa"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "email_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.key", "Your PagerDuty API key goes here"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "pagerduty_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.%", "6"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.operator", "absence"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.triggerinterval", "30m"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "slack_channel.0.url", "https://hooks.slack.com/services/identifier/secret"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.#", "1"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.%", "9"),
					// The JSON will have newlines per our API which uses JSON.stringify(obj, null, 2) as the value
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.bodytemplate", "{\n  \"fields\": {\n    \"description\": \"{{ matches }} matches found for {{ name }}\",\n    \"issuetype\": {\n      \"name\": \"Bug\"\n    },\n    \"project\": {\n      \"key\": \"test\"\n    },\n    \"summary\": \"Alert from {{ name }}\"\n  }\n}"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.headers.%", "2"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.headers.hello", "test3"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.headers.test", "test2"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.immediate", "false"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.method", "post"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.operator", "presence"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.terminal", "true"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.triggerinterval", "15m"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.triggerlimit", "15"),
					resource.TestCheckResourceAttr("logdna_view.new", "webhook_channel.0.url", "https://yourwebhook/endpoint"),
				),
			},
			{
				Config: updCfg,
				Check: resource.ComposeTestCheckFunc(
					testViewExists("logdna_view.new"),
					resource.TestCheckResourceAttr("logdna_view.new", "name", "test2"),
					resource.TestCheckResourceAttr("logdna_view.new", "query", "query2"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.0", "app3"),
					resource.TestCheckResourceAttr("logdna_view.new", "apps.1", "app4"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.0", "host3"),
					resource.TestCheckResourceAttr("logdna_view.new", "hosts.1", "host4"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.0", "error"),
					resource.TestCheckResourceAttr("logdna_view.new", "levels.1", "warning"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.0", "tags3"),
					resource.TestCheckResourceAttr("logdna_view.new", "tags.1", "tags4"),
				),
			},
			{
				ResourceName:      "logdna_view.new",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testViewExists(n string) resource.TestCheckFunc {
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

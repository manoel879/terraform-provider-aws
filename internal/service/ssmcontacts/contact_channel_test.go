package ssmcontacts_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssmcontacts"
	"github.com/aws/aws-sdk-go-v2/service/ssmcontacts/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfssmcontacts "github.com/hashicorp/terraform-provider-aws/internal/service/ssmcontacts"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testContactChannel_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	ctx := context.Background()

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	contactResourceName := "aws_ssmcontacts_contact.test"
	channelResourceName := "aws_ssmcontacts_contact_channel.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccContactPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMContactsEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContactChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactChannelConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(contactResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttr(channelResourceName, "activation_status", "NOT_ACTIVATED"),
					resource.TestCheckResourceAttr(channelResourceName, "delivery_address.0.simple_address", "default@example.com"),
					resource.TestCheckResourceAttr(channelResourceName, "name", rName),
					resource.TestCheckResourceAttr(channelResourceName, "type", "EMAIL"),
					resource.TestCheckResourceAttrPair(channelResourceName, "contact_id", contactResourceName, "arn"),
					acctest.MatchResourceAttrRegionalARN(channelResourceName, "arn", "ssm-contacts", regexp.MustCompile("contact-channel/test-contact-for-"+rName+"/.")),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// We need to explicitly test destroying this resource instead of just using CheckDestroy,
				// because CheckDestroy will run after the replication set has been destroyed and destroying
				// the replication set will destroy all other resources.
				Config: testAccContactChannelConfig_none(),
				Check:  testAccCheckContactChannelDestroy,
			},
		},
	})
}

func testContactChannel_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	ctx := context.Background()

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	channelResourceName := "aws_ssmcontacts_contact_channel.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccContactPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMContactsEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContactChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactChannelConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactChannelExists(channelResourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfssmcontacts.ResourceContactChannel(), channelResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testContactChannel_contactId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	ctx := context.Background()

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	testContactOneResourceName := "aws_ssmcontacts_contact.test_contact_one"
	testContactTwoResourceName := "aws_ssmcontacts_contact.test_contact_two"
	channelResourceName := "aws_ssmcontacts_contact_channel.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccContactPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMContactsEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContactChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactChannelConfig_withTwoContacts(rName, testContactOneResourceName+".arn"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(testContactOneResourceName),
					testAccCheckContactExists(testContactTwoResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttrPair(channelResourceName, "contact_id", testContactOneResourceName, "arn"),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactChannelConfig_withTwoContacts(rName, testContactTwoResourceName+".arn"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(testContactOneResourceName),
					testAccCheckContactExists(testContactTwoResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttrPair(channelResourceName, "contact_id", testContactTwoResourceName, "arn"),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testContactChannel_deliveryAddress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	ctx := context.Background()

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	contactResourceName := "aws_ssmcontacts_contact.test"
	channelResourceName := "aws_ssmcontacts_contact_channel.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccContactPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMContactsEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContactChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactChannelConfig(rName, rName, "EMAIL", "first@example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(contactResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttr(channelResourceName, "activation_status", "NOT_ACTIVATED"),
					resource.TestCheckResourceAttr(channelResourceName, "delivery_address.0.simple_address", "first@example.com"),
					resource.TestCheckResourceAttr(channelResourceName, "type", "EMAIL"),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactChannelConfig(rName, rName, "EMAIL", "second@example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(contactResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttr(channelResourceName, "activation_status", "NOT_ACTIVATED"),
					resource.TestCheckResourceAttr(channelResourceName, "delivery_address.0.simple_address", "second@example.com"),
					resource.TestCheckResourceAttr(channelResourceName, "type", "EMAIL"),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testContactChannel_name(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	ctx := context.Background()

	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix + "1")
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix + "2")
	contactResourceName := "aws_ssmcontacts_contact.test"
	channelResourceName := "aws_ssmcontacts_contact_channel.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccContactPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMContactsEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContactChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactChannelConfig(rName1, "update-name-test", "EMAIL", "test@example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(contactResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttr(channelResourceName, "name", rName1),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactChannelConfig(rName2, "update-name-test", "EMAIL", "test@example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(contactResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttr(channelResourceName, "name", rName2),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testContactChannel_type(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	ctx := context.Background()

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	contactResourceName := "aws_ssmcontacts_contact.test"
	channelResourceName := "aws_ssmcontacts_contact_channel.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			testAccContactPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SSMContactsEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckContactChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactChannelConfig_defaultEmail(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(contactResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttr(channelResourceName, "activation_status", "NOT_ACTIVATED"),
					resource.TestCheckResourceAttr(channelResourceName, "delivery_address.0.simple_address", "default@example.com"),
					resource.TestCheckResourceAttr(channelResourceName, "type", "EMAIL"),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactChannelConfig_defaultSms(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(contactResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttr(channelResourceName, "activation_status", "NOT_ACTIVATED"),
					resource.TestCheckResourceAttr(channelResourceName, "delivery_address.0.simple_address", "+12065550100"),
					resource.TestCheckResourceAttr(channelResourceName, "type", "SMS"),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactChannelConfig_defaultVoice(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContactExists(contactResourceName),
					testAccCheckContactChannelExists(channelResourceName),
					resource.TestCheckResourceAttr(channelResourceName, "activation_status", "NOT_ACTIVATED"),
					resource.TestCheckResourceAttr(channelResourceName, "delivery_address.0.simple_address", "+12065550199"),
					resource.TestCheckResourceAttr(channelResourceName, "type", "VOICE"),
				),
			},
			{
				ResourceName:      channelResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckContactChannelDestroy(s *terraform.State) error {
	ctx := context.Background()
	conn := acctest.Provider.Meta().(*conns.AWSClient).SSMContactsClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_ssmcontacts_contact_channel" {
			continue
		}

		input := &ssmcontacts.GetContactChannelInput{
			ContactChannelId: aws.String(rs.Primary.ID),
		}
		_, err := conn.GetContactChannel(ctx, input)

		if err != nil {
			// Getting resources may return validation exception when the replication set has been destroyed
			var ve *types.ValidationException
			if errors.As(err, &ve) {
				continue
			}

			var nfe *types.ResourceNotFoundException
			if errors.As(err, &nfe) {
				continue
			}

			return err
		}

		return create.Error(names.SSMContacts, create.ErrActionCheckingDestroyed, tfssmcontacts.ResNameContactChannel, rs.Primary.ID, errors.New("not destroyed"))
	}

	return nil
}

func testAccCheckContactChannelExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.SSMContacts, create.ErrActionCheckingExistence, tfssmcontacts.ResNameContactChannel, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.SSMContacts, create.ErrActionCheckingExistence, tfssmcontacts.ResNameContactChannel, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SSMContactsClient()
		_, err := conn.GetContactChannel(ctx, &ssmcontacts.GetContactChannelInput{
			ContactChannelId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return create.Error(names.SSMContacts, create.ErrActionCheckingExistence, tfssmcontacts.ResNameContactChannel, rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccContactChannelConfig_basic(rName string) string {
	return testAccContactChannelConfig_defaultEmail(rName)
}

func testAccContactChannelConfig_none() string {
	return testAccContactChannelConfig_base()
}

func testAccContactChannelConfig_defaultEmail(rName string) string {
	return testAccContactChannelConfig(rName, rName, "EMAIL", "default@example.com")
}

func testAccContactChannelConfig_defaultSms(rName string) string {
	return testAccContactChannelConfig(rName, rName, "SMS", "+12065550100")
}

func testAccContactChannelConfig_defaultVoice(rName string) string {
	return testAccContactChannelConfig(rName, rName, "VOICE", "+12065550199")
}

func testAccContactChannelConfig(rName string, contactAliasDisambiguator string, channelType string, address string) string {
	return acctest.ConfigCompose(
		testAccContactChannelConfig_base(),
		fmt.Sprintf(`
resource "aws_ssmcontacts_contact" "test" {
  alias = "test-contact-for-%[1]s"
  type  = "PERSONAL"

  depends_on = [aws_ssmincidents_replication_set.test]
}

resource "aws_ssmcontacts_contact_channel" "test" {
  contact_id = aws_ssmcontacts_contact.test.arn

  delivery_address {
    simple_address = %[4]q
  }

  name = %[1]q
  type = %[3]q
}
`, rName, contactAliasDisambiguator, channelType, address))
}

func testAccContactChannelConfig_withTwoContacts(rName, contactArn string) string {
	return acctest.ConfigCompose(
		testAccContactChannelConfig_base(),
		fmt.Sprintf(`
resource "aws_ssmcontacts_contact" "test_contact_one" {
  alias = "test-contact-one-for-%[1]s"
  type  = "PERSONAL"

  depends_on = [aws_ssmincidents_replication_set.test]
}

resource "aws_ssmcontacts_contact" "test_contact_two" {
  alias = "test-contact-two-for-%[1]s"
  type = "PERSONAL"

  depends_on = [aws_ssmincidents_replication_set.test]
}

resource "aws_ssmcontacts_contact_channel" "test" {
  contact_id = %[2]s

  delivery_address {
    simple_address = "test@example.com"
  }

  name = %[1]q
  type = "EMAIL"
}
`, rName, contactArn))
}

func testAccContactChannelConfig_base() string {
	return fmt.Sprintf(`
resource "aws_ssmincidents_replication_set" "test" {
  region {
    name = %[1]q
  }
}
`, acctest.Region())
}

package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediapackage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAWSMediaPackageEndpoint_basic(t *testing.T) {
	resourceName := "aws_media_package_endpoint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSMediaPackage(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsMediaPackageEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMediaPackageEndpointConfig(acctest.RandString(5)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsMediaPackageEndpointExists(resourceName),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "mediapackage", regexp.MustCompile(`channels/.+`)),
					resource.TestMatchResourceAttr(resourceName, "hls_ingest.0.ingest_endpoints.0.password", regexp.MustCompile("^[0-9a-f]*$")),
					resource.TestMatchResourceAttr(resourceName, "hls_ingest.0.ingest_endpoints.0.url", regexp.MustCompile("^https://")),
					resource.TestMatchResourceAttr(resourceName, "hls_ingest.0.ingest_endpoints.0.username", regexp.MustCompile("^[0-9a-f]*$")),
					resource.TestMatchResourceAttr(resourceName, "hls_ingest.0.ingest_endpoints.1.password", regexp.MustCompile("^[0-9a-f]*$")),
					resource.TestMatchResourceAttr(resourceName, "hls_ingest.0.ingest_endpoints.1.url", regexp.MustCompile("^https://")),
					resource.TestMatchResourceAttr(resourceName, "hls_ingest.0.ingest_endpoints.1.username", regexp.MustCompile("^[0-9a-f]*$")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSMediaPackageEndpoint_description(t *testing.T) {
	resourceName := "aws_media_package_endpoint.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSMediaPackage(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsMediaPackageEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMediaPackageEndpointConfigDescription(rName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsMediaPackageEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMediaPackageEndpointConfigDescription(rName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsMediaPackageEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
		},
	})
}

func TestAccAWSMediaPackageEndpoint_tags(t *testing.T) {
	resourceName := "aws_media_package_endpoint.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSMediaPackage(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsMediaPackageEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMediaPackageEndpointConfigWithTags(rName, "Environment", "test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsMediaPackageEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMediaPackageEndpointConfigWithTags(rName, "Environment", "test1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsMediaPackageEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test1"),
				),
			},
			{
				Config: testAccMediaPackageEndpointConfigWithTags(rName, "Update", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsMediaPackageEndpointExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Update", "true"),
				),
			},
		},
	})
}

func testAccCheckAwsMediaPackageEndpointDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).mediapackageconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_media_package_endpoint" {
			continue
		}

		input := &mediapackage.DescribeChannelInput{
			Id: aws.String(rs.Primary.ID),
		}

		_, err := conn.DescribeChannel(input)
		if err == nil {
			return fmt.Errorf("MediaPackage Channel (%s) not deleted", rs.Primary.ID)
		}

		if !isAWSErr(err, mediapackage.ErrCodeNotFoundException, "") {
			return err
		}
	}

	return nil
}

func testAccCheckAwsMediaPackageEndpointExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*AWSClient).mediapackageconn

		input := &mediapackage.DescribeChannelInput{
			Id: aws.String(rs.Primary.ID),
		}

		_, err := conn.DescribeChannel(input)

		return err
	}
}

func testAccPreCheckAWSMediaPackage(t *testing.T) {
	conn := testAccProvider.Meta().(*AWSClient).mediapackageconn

	input := &mediapackage.ListChannelsInput{}

	_, err := conn.ListChannels(input)

	if testAccPreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccMediaPackageEndpointConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_media_package_endpoint" "test" {
  channel_id = "tf_mediachannel_%s"
}
`, rName)
}

func testAccMediaPackageEndpointConfigDescription(rName, description string) string {
	return fmt.Sprintf(`
resource "aws_media_package_endpoint" "test" {
  channel_id  = %q
  description = %q
}
`, rName, description)
}

func testAccMediaPackageEndpointConfigWithTags(rName, key, value string) string {
	return fmt.Sprintf(`
resource "aws_media_package_endpoint" "test" {
  channel_id = "%[1]s"

  tags = {
	  Name = "%[1]s"
	  %[2]s = "%[3]s"
  }
}
`, rName, key, value)
}

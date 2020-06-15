package aws

import (
	// "fmt"
	// "regexp"
	// "time"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/service/medialive"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAwsMediaLiveChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsMediaLiveChannelCreate,
		Read:   resourceAwsMediaLiveChannelRead,
		Update: resourceAwsMediaLiveChannelUpdate,
		Delete: resourceAwsMediaLiveChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAwsMediaLiveChannelCreate(d *schema.ResourceData, meta interface{}) error {
	// conn := meta.(*AWSClient).medialiveconn

	// input := &medialive.CreateChannelInput{
	// 	Id:          aws.String(d.Get("channel_id").(string)),
	// 	Description: aws.String(d.Get("description").(string)),
	// }

	// if attr, ok := d.GetOk("tags"); ok {
	// 	input.Tags = tagsFromMapGeneric(attr.(map[string]interface{}))
	// }

	// _, err := conn.CreateChannel(input)
	// if err != nil {
	// 	return fmt.Errorf("error creating MediaLive Channel: %s", err)
	// }

	// d.SetId(d.Get("channel_id").(string))
	return resourceAwsMediaLiveChannelRead(d, meta)
}

func resourceAwsMediaLiveChannelRead(d *schema.ResourceData, meta interface{}) error {
	// conn := meta.(*AWSClient).medialiveconn

	// input := &medialive.DescribeChannelInput{
	// 	Id: aws.String(d.Id()),
	// }
	// resp, err := conn.DescribeChannel(input)
	// if err != nil {
	// 	return fmt.Errorf("error describing MediaLive Channel: %s", err)
	// }
	// d.Set("arn", resp.Arn)
	// d.Set("channel_id", resp.Id)
	// d.Set("description", resp.Description)

	// if err := d.Set("hls_ingest", flattenMediaLiveHLSIngest(resp.HlsIngest)); err != nil {
	// 	return fmt.Errorf("error setting hls_ingest: %s", err)
	// }

	// if err := d.Set("tags", tagsToMapGeneric(resp.Tags)); err != nil {
	// 	return fmt.Errorf("error setting tags: %s", err)
	// }

	return nil
}

func resourceAwsMediaLiveChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	// conn := meta.(*AWSClient).medialiveconn

	// input := &medialive.UpdateChannelInput{
	// 	Id:          aws.String(d.Id()),
	// 	Description: aws.String(d.Get("description").(string)),
	// }

	// _, err := conn.UpdateChannel(input)
	// if err != nil {
	// 	return fmt.Errorf("error updating MediaLive Channel: %s", err)
	// }

	// if err := setTagsMediaLive(conn, d, d.Get("arn").(string)); err != nil {
	// 	return fmt.Errorf("error updating MediaLive Channel (%s) tags: %s", d.Id(), err)
	// }

	return resourceAwsMediaLiveChannelRead(d, meta)
}

func resourceAwsMediaLiveChannelDelete(d *schema.ResourceData, meta interface{}) error {
	// conn := meta.(*AWSClient).medialiveconn

	// input := &medialive.DeleteChannelInput{
	// 	Id: aws.String(d.Id()),
	// }
	// _, err := conn.DeleteChannel(input)
	// if err != nil {
	// 	if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
	// 		return nil
	// 	}
	// 	return fmt.Errorf("error deleting MediaLive Channel: %s", err)
	// }

	// dcinput := &medialive.DescribeChannelInput{
	// 	Id: aws.String(d.Id()),
	// }
	// err = resource.Retry(5*time.Minute, func() *resource.RetryError {
	// 	_, err := conn.DescribeChannel(dcinput)
	// 	if err != nil {
	// 		if isAWSErr(err, medialive.ErrCodeNotFoundException, "") {
	// 			return nil
	// 		}
	// 		return resource.NonRetryableError(err)
	// 	}
	// 	return resource.RetryableError(fmt.Errorf("MediaLive Channel (%s) still exists", d.Id()))
	// })
	// if isResourceTimeoutError(err) {
	// 	_, err = conn.DescribeChannel(dcinput)
	// }
	// if err != nil {
	// 	return fmt.Errorf("error waiting for MediaLive Channel (%s) deletion: %s", d.Id(), err)
	// }

	return nil
}

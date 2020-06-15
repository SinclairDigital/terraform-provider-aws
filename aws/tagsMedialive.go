package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/medialive"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func setTagsMediaLive(conn *medialive.MediaLive, d *schema.ResourceData, arn string) error {
	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTagsGeneric(o, n)

		// Set tags
		if len(remove) > 0 {
			log.Printf("[DEBUG] Removing tags: %#v", remove)
			keys := make([]*string, 0, len(remove))
			for k := range remove {
				keys = append(keys, aws.String(k))
			}

			_, err := conn.DeleteTags(&medialive.DeleteTagsInput{
				ResourceArn: aws.String(arn),
				TagKeys:     keys,
			})
			if err != nil {
				return err
			}
		}
		if len(create) > 0 {
			log.Printf("[DEBUG] Creating tags: %#v", create)
			_, err := conn.CreateTags(&medialive.CreateTagsInput{
				ResourceArn: aws.String(arn),
				Tags:        create,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

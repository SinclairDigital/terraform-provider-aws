package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediapackage"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAwsMediaPackageEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsMediaPackageEndpointCreate,
		Read:   resourceAwsMediaPackageEndpointRead,
		Update: resourceAwsMediaPackageEndpointUpdate,
		Delete: resourceAwsMediaPackageEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// type CreateOriginEndpointOutput struct {
			//   Arn *string `locationName:"arn" type:"string"`
			//   ChannelId *string `locationName:"channelId" type:"string"`
			//   CmafPackage *CmafPackage `locationName:"cmafPackage" type:"structure"`
			//   DashPackage *DashPackage `locationName:"dashPackage" type:"structure"`
			//   Description *string `locationName:"description" type:"string"`
			//   HlsPackage *HlsPackage `locationName:"hlsPackage" type:"structure"`
			//   Id *string `locationName:"id" type:"string"`
			//   ManifestName *string `locationName:"manifestName" type:"string"`
			//   MssPackage *MssPackage `locationName:"mssPackage" type:"structure"`
			//   Origination *string `locationName:"origination" type:"string" enum:"Origination"`
			//   StartoverWindowSeconds *int64 `locationName:"startoverWindowSeconds" type:"integer"`
			//   Tags map[string]*string `locationName:"tags" type:"map"`
			//   TimeDelaySeconds *int64 `locationName:"timeDelaySeconds" type:"integer"`
			//   Url *string `locationName:"url" type:"string"`
			//   Whitelist []*string `locationName:"whitelist" type:"list"`
			// }
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"channel_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"endpoint_id": { // Id
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"manifest_name": { // ManifestName
				Type:     schema.TypeString,
				Optional: true,
				Default:  "index",
			},
			"origination": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ALLOW",
				ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY"}, false),
			},
			"startover_window_seconds": { // StartoverWindowSeconds
				Type:     schema.TypeInt,
				Optional: true,
			},
			"tags": tagsSchema(),
			"time_delay_seconds": { // TimeDelaySeconds
				Type:     schema.TypeInt,
				Optional: true,
			},
			"url": { // Url
				Type:     schema.TypeString,
				Computed: true,
			},
			"whitelist": { // Whitelist
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"type": {
				// Type of delivery package (other parameters are only used by some of the types available)
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "HLS",
				ValidateFunc: validation.StringInSlice([]string{"CMAF", "DASH", "HLS", "MSS"}, false),
			},
			"encryption": { // ALL
				// Whether the stream should be encrypted
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// A Common Media Application Format (CMAF) packaging configuration.
			// type CmafPackage struct { // type DashEncryption struct {
			// 	Encryption *CmafEncryption `locationName:"encryption" type:"structure"`
			// 	HlsManifests []*HlsManifest `locationName:"hlsManifests" type:"list"`
			// 	SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
			// 	SegmentPrefix *string `locationName:"segmentPrefix" type:"string"`
			// 	StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
			// }
			"segment_duration_seconds": { // ALL
				// Duration (in seconds) of each segment. Actual segments will berounded to
				// the nearest multiple of the source segment duration.
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			"segment_prefix": { // CMAF
				// An optional custom string that is prepended to the name of each segment.
				// If not specified, it defaults to the ChannelId.
				Type:     schema.TypeString,
				Optional: true,
			},

			// type (CMAF|DASH|HLS|MSS)Encryption struct {
			//   ConstantInitializationVector *string `locationName:"constantInitializationVector" type:"string"`
			//   EncryptionMethod *string `locationName:"encryptionMethod" type:"string" enum:"EncryptionMethod"`
			//   KeyRotationIntervalSeconds *int64 `locationName:"keyRotationIntervalSeconds" type:"integer"`
			//   RepeatExtXKey *bool `locationName:"repeatExtXKey" type:"boolean"`
			//   SpekeKeyProvider *SpekeKeyProvider `locationName:"spekeKeyProvider" type:"structure" required:"true"`
			// }
			"constant_initialization_vector": { // HLS
				// A constant initialization vector for encryption. When not specified
				// the initialization vector will be periodically rotated.
				Type:     schema.TypeString,
				Optional: true,
			},
			"encryption_method": { // HLS
				// The encryption method to use.
				Type:     schema.TypeString,
				Optional: true,
			},
			"key_rotation_interval_seconds": { // CMAF, DASH, HLS
				// Interval (in seconds) between each encryption key rotation.
				Type:     schema.TypeInt,
				Optional: true,
			},
			"repeat_ext_x_key": { // HLS
				// When enabled, the EXT-X-KEY tag will be repeated in output manifests.
				Type:     schema.TypeBool,
				Optional: true,
			},

			"speke_key_provider": { // ALL
				// A configuration for accessing an external Secure Packager and Encoder Key
				// Exchange (SPEKE) service that will provide encryption keys. Required if encryption
				// is enabled
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					// type SpekeKeyProvider struct {
					//   CertificateArn *string `locationName:"certificateArn" type:"string"`
					//   ResourceId *string `locationName:"resourceId" type:"string" required:"true"`
					//   RoleArn *string `locationName:"roleArn" type:"string" required:"true"`
					//   SystemIds []*string `locationName:"systemIds" type:"list" required:"true"`
					//   Url *string `locationName:"url" type:"string" required:"true"`
					// }
					Schema: map[string]*schema.Schema{
						"certificate_arn": {
							// An Amazon Resource Name (ARN) of a Certificate Manager certificatethat MediaPackage
							// will use for enforcing secure end-to-end datatransfer with the key provider
							// service.
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_id": {
							// The resource ID to include in key requests.
							Type:     schema.TypeString,
							Required: true,
						},
						"role_arn": {
							// An Amazon Resource Name (ARN) of an IAM role that AWS ElementalMediaPackage
							// will assume when accessing the key provider service.
							Type:     schema.TypeString,
							Required: true,
						},
						"system_ids": {
							// The system IDs to include in key requests.
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 2,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"url": {
							// The URL of the external key provider service.
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			// type HlsManifest struct {
			//   AdMarkers *string `locationName:"adMarkers" type:"string" enum:"AdMarkers"`
			//   Id *string `locationName:"id" type:"string" required:"true"`
			//   IncludeIframeOnlyStream *bool `locationName:"includeIframeOnlyStream" type:"boolean"`
			//   ManifestName *string `locationName:"manifestName" type:"string"`
			//   PlaylistType *string `locationName:"playlistType" type:"string" enum:"PlaylistType"`
			//   PlaylistWindowSeconds *int64 `locationName:"playlistWindowSeconds" type:"integer"`
			//   ProgramDateTimeIntervalSeconds *int64 `locationName:"programDateTimeIntervalSeconds" type:"integer"`
			//   Url *string `locationName:"url" type:"string"`
			// }
			"hls_manifests": { // CMAF
				// A list of HLS manifest configurations, required for CMAF
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_markers": {
							// This setting controls how ad markers are included in the packaged OriginEndpoint.
							// "NONE" will omit all SCTE-35 ad markers from the output.
							// "PASSTHROUGH" causes the manifest to contain a copy of the SCTE-35 admarkers
							// (comments) taken directly from the input HTTP Live Streaming (HLS) manifest.
							// "SCTE35_ENHANCED" generates ad markers and blackout tags based on SCTE-35messages
							// in the input source.
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "NONE",
							ValidateFunc: validation.StringInSlice([]string{"NONE", "PASSTHROUGH", "SCTE35_ENHANCED"}, false),
						},
						"id": {
							// The ID of the manifest. The ID must be unique within the OriginEndpoint and
							// it cannot be changed after it is created.
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"include_iframe_only_stream": {
							// When enabled, an I-Frame only stream will be included in the output.
							Type:     schema.TypeBool,
							Optional: true,
						},
						"manifest_name": {
							// An optional short string appended to the end of the OriginEndpoint URL. If
							// not specified, defaults to the manifestName for the OriginEndpoint.
							Type:     schema.TypeString,
							Optional: true,
						},
						"playlist_type": {
							// The HTTP Live Streaming (HLS) playlist type. When either "EVENT" or "VOD"
							// is specified, a corresponding EXT-X-PLAYLIST-TYPEentry will be included in
							// the media playlist.
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"EVENT", "VOD"}, false),
						},
						"playlist_window_seconds": {
							// Time window (in seconds) contained in each parent manifest.
							Type:     schema.TypeInt,
							Optional: true,
							Default:  60,
						},
						"program_date_time_interval_seconds": {
							// The interval (in seconds) between each EXT-X-PROGRAM-DATE-TIME taginserted
							// into manifests. Additionally, when an interval is specifiedID3Timed Metadata
							// messages will be generated every 5 seconds using theingest time of the content.If
							// the interval is not specified, or set to 0, thenno EXT-X-PROGRAM-DATE-TIME
							// tags will be inserted into manifests and noID3Timed Metadata messages will
							// be generated. Note that irrespectiveof this parameter, if any ID3 Timed Metadata
							// is found in HTTP Live Streaming (HLS) input,it will be passed through to
							// HLS output.
							Type:     schema.TypeInt,
							Optional: true,
						},
						"url": {
							// The URL of the packaged OriginEndpoint for consumption.
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			// type StreamSelection struct {
			//   MaxVideoBitsPerSecond *int64 `locationName:"maxVideoBitsPerSecond" type:"integer"`
			//   MinVideoBitsPerSecond *int64 `locationName:"minVideoBitsPerSecond" type:"integer"`
			//   StreamOrder *string `locationName:"streamOrder" type:"string" enum:"StreamOrder"`
			// }
			"max_video_bits_per_second": { // ALL
				// The maximum video bitrate (bps) to include in output.
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_video_bits_per_second": { // ALL
				// The minimum video bitrate (bps) to include in output.
				Type:     schema.TypeInt,
				Optional: true,
			},
			"stream_order": { // ALL
				// A directive that determines the order of streams in the output.
				Type:     schema.TypeString,
				Optional: true,
			},

			// A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
			// type DashPackage struct {
			//   AdTriggers []*string `locationName:"adTriggers" type:"list"`
			//   AdsOnDeliveryRestrictions *string `locationName:"adsOnDeliveryRestrictions" type:"string" enum:"AdsOnDeliveryRestrictions"`
			//   Encryption *DashEncryption `locationName:"encryption" type:"structure"`
			//   ManifestLayout *string `locationName:"manifestLayout" type:"string" enum:"ManifestLayout"`
			//   ManifestWindowSeconds *int64 `locationName:"manifestWindowSeconds" type:"integer"`
			//   MinBufferTimeSeconds *int64 `locationName:"minBufferTimeSeconds" type:"integer"`
			//   MinUpdatePeriodSeconds *int64 `locationName:"minUpdatePeriodSeconds" type:"integer"`
			//   PeriodTriggers []*string `locationName:"periodTriggers" type:"list"`
			//   Profile *string `locationName:"profile" type:"string" enum:"Profile"`
			//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
			//   SegmentTemplateFormat *string `locationName:"segmentTemplateFormat" type:"string" enum:"SegmentTemplateFormat"`
			//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
			//   SuggestedPresentationDelaySeconds *int64 `locationName:"suggestedPresentationDelaySeconds" type:"integer"`
			// }
			"ad_triggers": { // DASH, HLS
				// A list of SCTE-35 message types that are treated as ad markers in the output.
				// If empty, no ad markers are output. Specify multiple items to create ad markers
				// for all of the includedmessage types.
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ads_on_delivery_restrictions": { // DASH, HLS
				// This setting allows the delivery restriction flags on SCTE-35 segmentation
				// descriptors todetermine whether a message signals an ad.
				// NONE: no SCTE-35 messages becomeads.
				// RESTRICTED: SCTE-35 messages of the types specified in AdTriggers thatcontain
				// delivery restrictions will be treated as ads.
				// UNRESTRICTED: SCTE-35 messages of the types specified in AdTriggers that do not
				// contain delivery restrictions willbe treated as ads.
				// BOTH: all SCTE-35 messages of the types specified inAdTriggers will be treated
				// as ads. Note that Splice Insert messages do not have these flagsand are always
				// treated as ads if specified in AdTriggers.
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "RESTRICTED", "UNRESTRICTED", "BOTH"}, false),
			},
			"manifest_layout": { // DASH
				// Determines the position of some tags in the Media Presentation Description
				// (MPD). When set to FULL, elements like SegmentTemplate and ContentProtection
				// are included in each Representation. When set to COMPACT, duplicate elements
				// are combined and presented at the AdaptationSet level.
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPACT", "FULL"}, false),
			},
			"manifest_window_seconds": { // DASH, MSS
				// Time window (in seconds) contained in each manifest.
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_buffer_time_seconds": { // DASH
				// Minimum duration (in seconds) that a player will buffer media before starting
				// the presentation.
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_update_period_seconds": { // DASH
				// Minimum duration (in seconds) between potential changes to the Dynamic Adaptive
				// Streaming over HTTP (DASH) Media Presentation Description (MPD).
				Type:     schema.TypeInt,
				Optional: true,
			},
			"period_triggers": { // DASH
				// A list of triggers that controls when the outgoing Dynamic Adaptive Streaming
				// over HTTP (DASH)Media Presentation Description (MPD) will be partitioned
				// into multiple periods. If empty, the content will notbe partitioned into
				// more than one period. If the list contains "ADS", new periods will be created
				// wherethe Channel source contains SCTE-35 ad markers.
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"profile": { // DASH
				// The Dynamic Adaptive Streaming over HTTP (DASH) profile type. When set to
				// "HBBTV_1_5", HbbTV 1.5 compliant output is enabled.
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "HBBTV_1_5"}, false),
			},
			"segment_template_format": { // DASH
				// Determines the type of SegmentTemplate included in the Media Presentation
				// Description (MPD).
				// NUMBER_WITH_TIMELINE: a full timeline is presented in each SegmentTemplate,
				// with $Number$ media URLs.
				// TIME_WITH_TIMELINE: a full timeline is presented in each SegmentTemplate,
				// with $Time$ media URLs.
				// NUMBER_WITH_DURATION: only a duration is included in each SegmentTemplate,
				// with $Number$ media URLs.
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NUMBER_WITH_TIMELINE", "TIME_WITH_TIMELINE", "NUMBER_WITH_DURATION"}, false),
			},
			"suggested_presentation_delay_seconds": { // DASH
				// Duration (in seconds) to delay live content before presentation.
				Type:     schema.TypeInt,
				Optional: true,
			},

			// An HTTP Live Streaming (HLS) packaging configuration.
			// type HlsPackage struct {
			//   AdMarkers *string `locationName:"adMarkers" type:"string" enum:"AdMarkers"`
			//   AdTriggers []*string `locationName:"adTriggers" type:"list"`
			//   AdsOnDeliveryRestrictions *string `locationName:"adsOnDeliveryRestrictions" type:"string" enum:"AdsOnDeliveryRestrictions"`
			//   Encryption *HlsEncryption `locationName:"encryption" type:"structure"`
			//   IncludeIframeOnlyStream *bool `locationName:"includeIframeOnlyStream" type:"boolean"`
			//   PlaylistType *string `locationName:"playlistType" type:"string" enum:"PlaylistType"`
			//   PlaylistWindowSeconds *int64 `locationName:"playlistWindowSeconds" type:"integer"`
			//   ProgramDateTimeIntervalSeconds *int64 `locationName:"programDateTimeIntervalSeconds" type:"integer"`
			//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
			//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
			//   UseAudioRenditionGroup *bool `locationName:"useAudioRenditionGroup" type:"boolean"`
			// }

			"ad_markers": { // HLS
				// This setting controls how ad markers are included in the packaged OriginEndpoint.
				// "NONE" will omit all SCTE-35 ad markers from the output.
				// "PASSTHROUGH" causes the manifest to contain a copy of the SCTE-35 admarkers
				// (comments) taken directly from the input HTTP Live Streaming (HLS) manifest.
				// "SCTE35_ENHANCED" generates ad markers and blackout tags based on SCTE-35messages
				// in the input source.
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "PASSTHROUGH", "SCTE35_ENHANCED"}, false),
			},
			"include_iframe_only_stream": { // HLS
				// When enabled, an I-Frame only stream will be included in the output.
				Type:     schema.TypeBool,
				Optional: true,
			},
			"playlist_type": { // HLS
				// The HTTP Live Streaming (HLS) playlist type.When either "EVENT" or "VOD"
				// is specified, a corresponding EXT-X-PLAYLIST-TYPEentry will be included in
				// the media playlist.
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"EVENT", "VOD"}, false),
			},
			"playlist_window_seconds": { // HLS
				// Time window (in seconds) contained in each parent manifest.
				Type:     schema.TypeInt,
				Optional: true,
			},
			"program_date_time_interval_seconds": { // HLS
				// The interval (in seconds) between each EXT-X-PROGRAM-DATE-TIME taginserted
				// into manifests. Additionally, when an interval is specifiedID3Timed Metadata
				// messages will be generated every 5 seconds using theingest time of the content.If
				// the interval is not specified, or set to 0, thenno EXT-X-PROGRAM-DATE-TIME
				// tags will be inserted into manifests and noID3Timed Metadata messages will
				// be generated. Note that irrespectiveof this parameter, if any ID3 Timed Metadata
				// is found in HTTP Live Streaming (HLS) input,it will be passed through to
				// HLS output.
				Type:     schema.TypeInt,
				Optional: true,
			},
			"use_audio_rendition_group": { // HLS
				// When enabled, audio streams will be placed in rendition groups in the output.
				Type:     schema.TypeBool,
				Optional: true,
			},

			// A Microsoft Smooth Streaming (MSS) packaging configuration.
			// type MssPackage struct {
			//   Encryption *MssEncryption `locationName:"encryption" type:"structure"`
			//   ManifestWindowSeconds *int64 `locationName:"manifestWindowSeconds" type:"integer"`
			//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
			//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
			// }

		},
	}
}

// type CreateOriginEndpointInput struct {
//   // ChannelId is a required field
//   ChannelId *string `locationName:"channelId" type:"string" required:"true"`
//   // A Common Media Application Format (CMAF) packaging configuration.
//   CmafPackage *CmafPackageCreateOrUpdateParameters `locationName:"cmafPackage" type:"structure"`
//   // A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
//   DashPackage *DashPackage `locationName:"dashPackage" type:"structure"`
//   Description *string `locationName:"description" type:"string"`
//   // An HTTP Live Streaming (HLS) packaging configuration.
//   HlsPackage *HlsPackage `locationName:"hlsPackage" type:"structure"`
//   // Id is a required field
//   Id *string `locationName:"id" type:"string" required:"true"`
//   ManifestName *string `locationName:"manifestName" type:"string"`
//   // A Microsoft Smooth Streaming (MSS) packaging configuration.
//   MssPackage *MssPackage `locationName:"mssPackage" type:"structure"`
//   Origination *string `locationName:"origination" type:"string" enum:"Origination"`
//   StartoverWindowSeconds *int64 `locationName:"startoverWindowSeconds" type:"integer"`
//   // A collection of tags associated with a resource
//   Tags map[string]*string `locationName:"tags" type:"map"`
//   TimeDelaySeconds *int64 `locationName:"timeDelaySeconds" type:"integer"`
//   Whitelist []*string `locationName:"whitelist" type:"list"`
//   // contains filtered or unexported fields
// }

// type CreateOriginEndpointOutput struct {
//   Arn *string `locationName:"arn" type:"string"`
//   ChannelId *string `locationName:"channelId" type:"string"`
//   // A Common Media Application Format (CMAF) packaging configuration.
//   CmafPackage *CmafPackage `locationName:"cmafPackage" type:"structure"`
//   // A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
//   DashPackage *DashPackage `locationName:"dashPackage" type:"structure"`
//   Description *string `locationName:"description" type:"string"`
//   // An HTTP Live Streaming (HLS) packaging configuration.
//   HlsPackage *HlsPackage `locationName:"hlsPackage" type:"structure"`
//   Id *string `locationName:"id" type:"string"`
//   ManifestName *string `locationName:"manifestName" type:"string"`
//   // A Microsoft Smooth Streaming (MSS) packaging configuration.
//   MssPackage *MssPackage `locationName:"mssPackage" type:"structure"`
//   Origination *string `locationName:"origination" type:"string" enum:"Origination"`
//   StartoverWindowSeconds *int64 `locationName:"startoverWindowSeconds" type:"integer"`
//   // A collection of tags associated with a resource
//   Tags map[string]*string `locationName:"tags" type:"map"`
//   TimeDelaySeconds *int64 `locationName:"timeDelaySeconds" type:"integer"`
//   Url *string `locationName:"url" type:"string"`
//   Whitelist []*string `locationName:"whitelist" type:"list"`
//   // contains filtered or unexported fields
// }
func resourceAwsMediaPackageEndpointCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).mediapackageconn

	input := &mediapackage.CreateOriginEndpointInput{
		ChannelId:              aws.String(d.Get("channel_id").(string)),
		CmafPackage:            buildCmafPackage(d),
		DashPackage:            buildDashPackage(d),
		Description:            aws.String(d.Get("description").(string)),
		HlsPackage:             buildHlsPackage(d),
		Id:                     aws.String(d.Id()),
		ManifestName:           aws.String(d.Get("manifest_name").(string)),
		MssPackage:             buildMssPackage(d),
		Origination:            aws.String(d.Get("origination").(string)),
		StartoverWindowSeconds: aws.Int64(int64(d.Get("startover_window_seconds").(int))),
		TimeDelaySeconds:       aws.Int64(int64(d.Get("time_delay_seconds").(int))),
		Whitelist:              aws.StringSlice(d.Get("whitelist").([]string)),
	}

	if attr, ok := d.GetOk("tags"); ok {
		input.Tags = tagsFromMapGeneric(attr.(map[string]interface{}))
	}

	_, err := conn.CreateOriginEndpoint(input)
	if err != nil {
		return fmt.Errorf("error creating MediaPackage Endpoint: %s", err)
	}

	d.SetId(d.Get("id").(string))
	return resourceAwsMediaPackageEndpointRead(d, meta)
}

// type DescribeOriginEndpointInput struct {
//   // Id is a required field
//   Id *string `location:"uri" locationName:"id" type:"string" required:"true"`
//   // contains filtered or unexported fields
// }

// type DescribeOriginEndpointOutput struct {
//   Arn *string `locationName:"arn" type:"string"`
//   ChannelId *string `locationName:"channelId" type:"string"`
//   // A Common Media Application Format (CMAF) packaging configuration.
//   CmafPackage *CmafPackage `locationName:"cmafPackage" type:"structure"`
//   // A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
//   DashPackage *DashPackage `locationName:"dashPackage" type:"structure"`
//   Description *string `locationName:"description" type:"string"`
//   // An HTTP Live Streaming (HLS) packaging configuration.
//   HlsPackage *HlsPackage `locationName:"hlsPackage" type:"structure"`
//   Id *string `locationName:"id" type:"string"`
//   ManifestName *string `locationName:"manifestName" type:"string"`
//   // A Microsoft Smooth Streaming (MSS) packaging configuration.
//   MssPackage *MssPackage `locationName:"mssPackage" type:"structure"`
//   Origination *string `locationName:"origination" type:"string" enum:"Origination"`
//   StartoverWindowSeconds *int64 `locationName:"startoverWindowSeconds" type:"integer"`
//   // A collection of tags associated with a resource
//   Tags map[string]*string `locationName:"tags" type:"map"`
//   TimeDelaySeconds *int64 `locationName:"timeDelaySeconds" type:"integer"`
//   Url *string `locationName:"url" type:"string"`
//   Whitelist []*string `locationName:"whitelist" type:"list"`
//   // contains filtered or unexported fields
// }
func resourceAwsMediaPackageEndpointRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).mediapackageconn

	input := &mediapackage.DescribeOriginEndpointInput{
		Id: aws.String(d.Id()),
	}
	resp, err := conn.DescribeOriginEndpoint(input)
	if err != nil {
		return fmt.Errorf("error describing MediaPackage Endpoint: %s", err)
	}

	extractMediaPackageEndpointValues(d, resp)

	return nil
}

// type UpdateOriginEndpointInput struct {
//   // A Common Media Application Format (CMAF) packaging configuration.
//   CmafPackage *CmafPackageCreateOrUpdateParameters `locationName:"cmafPackage" type:"structure"`
//   // A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
//   DashPackage *DashPackage `locationName:"dashPackage" type:"structure"`
//   Description *string `locationName:"description" type:"string"`
//   // An HTTP Live Streaming (HLS) packaging configuration.
//   HlsPackage *HlsPackage `locationName:"hlsPackage" type:"structure"`
//   // Id is a required field
//   Id *string `location:"uri" locationName:"id" type:"string" required:"true"`
//   ManifestName *string `locationName:"manifestName" type:"string"`
//   // A Microsoft Smooth Streaming (MSS) packaging configuration.
//   MssPackage *MssPackage `locationName:"mssPackage" type:"structure"`
//   Origination *string `locationName:"origination" type:"string" enum:"Origination"`
//   StartoverWindowSeconds *int64 `locationName:"startoverWindowSeconds" type:"integer"`
//   TimeDelaySeconds *int64 `locationName:"timeDelaySeconds" type:"integer"`
//   Whitelist []*string `locationName:"whitelist" type:"list"`
//   // contains filtered or unexported fields
// }

// type UpdateOriginEndpointOutput struct {
//   Arn *string `locationName:"arn" type:"string"`
//   ChannelId *string `locationName:"channelId" type:"string"`
//   // A Common Media Application Format (CMAF) packaging configuration.
//   CmafPackage *CmafPackage `locationName:"cmafPackage" type:"structure"`
//   // A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
//   DashPackage *DashPackage `locationName:"dashPackage" type:"structure"`
//   Description *string `locationName:"description" type:"string"`
//   // An HTTP Live Streaming (HLS) packaging configuration.
//   HlsPackage *HlsPackage `locationName:"hlsPackage" type:"structure"`
//   Id *string `locationName:"id" type:"string"`
//   ManifestName *string `locationName:"manifestName" type:"string"`
//   // A Microsoft Smooth Streaming (MSS) packaging configuration.
//   MssPackage *MssPackage `locationName:"mssPackage" type:"structure"`
//   Origination *string `locationName:"origination" type:"string" enum:"Origination"`
//   StartoverWindowSeconds *int64 `locationName:"startoverWindowSeconds" type:"integer"`
//   // A collection of tags associated with a resource
//   Tags map[string]*string `locationName:"tags" type:"map"`
//   TimeDelaySeconds *int64 `locationName:"timeDelaySeconds" type:"integer"`
//   Url *string `locationName:"url" type:"string"`
//   Whitelist []*string `locationName:"whitelist" type:"list"`
//   // contains filtered or unexported fields
// }
func resourceAwsMediaPackageEndpointUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).mediapackageconn

	input := &mediapackage.UpdateOriginEndpointInput{
		CmafPackage:            buildCmafPackage(d),
		DashPackage:            buildDashPackage(d),
		Description:            aws.String(d.Get("description").(string)),
		HlsPackage:             buildHlsPackage(d),
		Id:                     aws.String(d.Id()),
		ManifestName:           aws.String(d.Get("manifest_name").(string)),
		MssPackage:             buildMssPackage(d),
		Origination:            aws.String(d.Get("origination").(string)),
		StartoverWindowSeconds: aws.Int64(int64(d.Get("startover_window_seconds").(int))),
		TimeDelaySeconds:       aws.Int64(int64(d.Get("time_delay_seconds").(int))),
		Whitelist:              aws.StringSlice(d.Get("whitelist").([]string)),
	}

	_, err := conn.UpdateOriginEndpoint(input)
	if err != nil {
		return fmt.Errorf("error updating MediaPackage Endpoint: %s", err)
	}

	if err := setTagsMediaPackage(conn, d, d.Get("arn").(string)); err != nil {
		return fmt.Errorf("error updating MediaPackage Endpoint (%s) tags: %s", d.Id(), err)
	}

	return resourceAwsMediaPackageEndpointRead(d, meta)
}

// type DeleteOriginEndpointInput struct {
//   // Id is a required field
//   Id *string `location:"uri" locationName:"id" type:"string" required:"true"`
//   // contains filtered or unexported fields
// }

// type DeleteOriginEndpointOutput struct {
//   // contains filtered or unexported fields
// }
func resourceAwsMediaPackageEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).mediapackageconn

	input := &mediapackage.DeleteOriginEndpointInput{
		Id: aws.String(d.Id()),
	}
	_, err := conn.DeleteOriginEndpoint(input)
	if err != nil {
		if isAWSErr(err, mediapackage.ErrCodeNotFoundException, "") {
			return nil
		}
		return fmt.Errorf("error deleting MediaPackage Endpoint: %s", err)
	}

	dcinput := &mediapackage.DescribeOriginEndpointInput{
		Id: aws.String(d.Id()),
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.DescribeOriginEndpoint(dcinput)
		if err != nil {
			if isAWSErr(err, mediapackage.ErrCodeNotFoundException, "") {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("MediaPackage Endpoint (%s) still exists", d.Id()))
	})
	if isResourceTimeoutError(err) {
		_, err = conn.DescribeOriginEndpoint(dcinput)
	}
	if err != nil {
		return fmt.Errorf("error waiting for MediaPackage Endpoint (%s) deletion: %s", d.Id(), err)
	}

	return nil
}

// A Common Media Application Format (CMAF) packaging configuration.
// type CmafPackage struct { // type DashEncryption struct {
// 	Encryption *CmafEncryption `locationName:"encryption" type:"structure"`
// 	HlsManifests []*HlsManifestCreateOrUpdateParameters `locationName:"hlsManifests" type:"list"`
// 	SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
// 	SegmentPrefix *string `locationName:"segmentPrefix" type:"string"`
// 	StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
// }
func buildCmafPackage(d *schema.ResourceData) *mediapackage.CmafPackageCreateOrUpdateParameters {
	var pack = &mediapackage.CmafPackageCreateOrUpdateParameters{
		SegmentDurationSeconds: aws.Int64(int64(d.Get("segment_duration_seconds").(int))),
		SegmentPrefix:          aws.String(d.Get("segment_prefix").(string)),

		Encryption:      buildCmafEncryption(d),
		HlsManifests:    buildHlsManifestList(d),
		StreamSelection: buildStreamSelection(d),
	}
	return pack
}

// A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
// type DashPackage struct {
//   AdTriggers []*string `locationName:"adTriggers" type:"list"`
//   AdsOnDeliveryRestrictions *string `locationName:"adsOnDeliveryRestrictions" type:"string" enum:"AdsOnDeliveryRestrictions"`
//   Encryption *DashEncryption `locationName:"encryption" type:"structure"`
//   ManifestLayout *string `locationName:"manifestLayout" type:"string" enum:"ManifestLayout"`
//   ManifestWindowSeconds *int64 `locationName:"manifestWindowSeconds" type:"integer"`
//   MinBufferTimeSeconds *int64 `locationName:"minBufferTimeSeconds" type:"integer"`
//   MinUpdatePeriodSeconds *int64 `locationName:"minUpdatePeriodSeconds" type:"integer"`
//   PeriodTriggers []*string `locationName:"periodTriggers" type:"list"`
//   Profile *string `locationName:"profile" type:"string" enum:"Profile"`
//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
//   SegmentTemplateFormat *string `locationName:"segmentTemplateFormat" type:"string" enum:"SegmentTemplateFormat"`
//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
//   SuggestedPresentationDelaySeconds *int64 `locationName:"suggestedPresentationDelaySeconds" type:"integer"`
// }
func buildDashPackage(d *schema.ResourceData) *mediapackage.DashPackage {
	var pack = &mediapackage.DashPackage{
		AdTriggers:                        aws.StringSlice(d.Get("ad_triggers").([]string)),
		AdsOnDeliveryRestrictions:         aws.String(d.Get("ads_on_delivery_restrictions").(string)),
		ManifestLayout:                    aws.String(d.Get("manifest_layout").(string)),
		ManifestWindowSeconds:             aws.Int64(int64(d.Get("manifest_window_seconds").(int))),
		MinBufferTimeSeconds:              aws.Int64(int64(d.Get("min_buffer_time_seconds").(int))),
		MinUpdatePeriodSeconds:            aws.Int64(int64(d.Get("min_update_period_seconds").(int))),
		PeriodTriggers:                    aws.StringSlice(d.Get("period_triggers").([]string)),
		Profile:                           aws.String(d.Get("profile").(string)),
		SegmentDurationSeconds:            aws.Int64(int64(d.Get("segment_duration_seconds").(int))),
		SegmentTemplateFormat:             aws.String(d.Get("segment_template_format").(string)),
		SuggestedPresentationDelaySeconds: aws.Int64(int64(d.Get("suggested_presentation_delay_seconds").(int))),

		Encryption:      buildDashEncryption(d),
		StreamSelection: buildStreamSelection(d),
	}
	return pack
}

// An HTTP Live Streaming (HLS) packaging configuration.
// type HlsPackage struct {
//   AdMarkers *string `locationName:"adMarkers" type:"string" enum:"AdMarkers"`
//   AdTriggers []*string `locationName:"adTriggers" type:"list"`
//   AdsOnDeliveryRestrictions *string `locationName:"adsOnDeliveryRestrictions" type:"string" enum:"AdsOnDeliveryRestrictions"`
//   Encryption *HlsEncryption `locationName:"encryption" type:"structure"`
//   IncludeIframeOnlyStream *bool `locationName:"includeIframeOnlyStream" type:"boolean"`
//   PlaylistType *string `locationName:"playlistType" type:"string" enum:"PlaylistType"`
//   PlaylistWindowSeconds *int64 `locationName:"playlistWindowSeconds" type:"integer"`
//   ProgramDateTimeIntervalSeconds *int64 `locationName:"programDateTimeIntervalSeconds" type:"integer"`
//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
//   UseAudioRenditionGroup *bool `locationName:"useAudioRenditionGroup" type:"boolean"`
// }
func buildHlsPackage(d *schema.ResourceData) *mediapackage.HlsPackage {
	var pack = &mediapackage.HlsPackage{
		AdMarkers:                      aws.String(d.Get("ad_markers").(string)),
		AdTriggers:                     aws.StringSlice(d.Get("ad_triggers").([]string)),
		AdsOnDeliveryRestrictions:      aws.String(d.Get("ads_on_delivery_restrictions").(string)),
		IncludeIframeOnlyStream:        aws.Bool(d.Get("include_iframe_only_stream").(bool)),
		PlaylistType:                   aws.String(d.Get("playlist_type").(string)),
		PlaylistWindowSeconds:          aws.Int64(int64(d.Get("playlist_window_seconds").(int))),
		ProgramDateTimeIntervalSeconds: aws.Int64(int64(d.Get("program_date_time_interval_seconds").(int))),
		SegmentDurationSeconds:         aws.Int64(int64(d.Get("segment_duration_seconds").(int))),
		UseAudioRenditionGroup:         aws.Bool(d.Get("use_audio_rendition_group").(bool)),

		Encryption:      buildHlsEncryption(d),
		StreamSelection: buildStreamSelection(d),
	}
	return pack
}

// A Microsoft Smooth Streaming (MSS) packaging configuration.
// type MssPackage struct {
//   Encryption *MssEncryption `locationName:"encryption" type:"structure"`
//   ManifestWindowSeconds *int64 `locationName:"manifestWindowSeconds" type:"integer"`
//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
// }
func buildMssPackage(d *schema.ResourceData) *mediapackage.MssPackage {
	var pack = &mediapackage.MssPackage{
		ManifestWindowSeconds:  aws.Int64(int64(d.Get("manifest_window_seconds").(int))),
		SegmentDurationSeconds: aws.Int64(int64(d.Get("segment_duration_seconds").(int))),

		Encryption:      buildMssEncryption(d),
		StreamSelection: buildStreamSelection(d),
	}
	return pack
}

// type (CMAF|DASH|HLS|MSS)Encryption struct {
//   ConstantInitializationVector *string `locationName:"constantInitializationVector" type:"string"`
//   EncryptionMethod *string `locationName:"encryptionMethod" type:"string" enum:"EncryptionMethod"`
//   KeyRotationIntervalSeconds *int64 `locationName:"keyRotationIntervalSeconds" type:"integer"`
//   RepeatExtXKey *bool `locationName:"repeatExtXKey" type:"boolean"`
//   SpekeKeyProvider *SpekeKeyProvider `locationName:"spekeKeyProvider" type:"structure" required:"true"`
// }
func buildCmafEncryption(d *schema.ResourceData) *mediapackage.CmafEncryption {
	var encryption = &mediapackage.CmafEncryption{
		KeyRotationIntervalSeconds: aws.Int64(int64(d.Get("key_rotation_interval_seconds").(int))),
		SpekeKeyProvider:           buildSpekeKeyProvider(d),
	}
	return encryption
}

func buildDashEncryption(d *schema.ResourceData) *mediapackage.DashEncryption {
	var encryption = &mediapackage.DashEncryption{
		KeyRotationIntervalSeconds: aws.Int64(int64(d.Get("key_rotation_interval_seconds").(int))),
		SpekeKeyProvider:           buildSpekeKeyProvider(d),
	}
	return encryption
}

func buildHlsEncryption(d *schema.ResourceData) *mediapackage.HlsEncryption {
	var encryption = &mediapackage.HlsEncryption{
		ConstantInitializationVector: aws.String(d.Get("constant_initialization_vector").(string)),
		EncryptionMethod:             aws.String(d.Get("encryption_method").(string)),
		KeyRotationIntervalSeconds:   aws.Int64(int64(d.Get("key_rotation_interval_seconds").(int))),
		RepeatExtXKey:                aws.Bool(d.Get("repeat_ext_x_key").(bool)),
		SpekeKeyProvider:             buildSpekeKeyProvider(d),
	}
	return encryption
}

func buildMssEncryption(d *schema.ResourceData) *mediapackage.MssEncryption {
	var encryption = &mediapackage.MssEncryption{
		SpekeKeyProvider: buildSpekeKeyProvider(d),
	}
	return encryption
}

// type SpekeKeyProvider struct {
//   CertificateArn *string `locationName:"certificateArn" type:"string"`
//   ResourceId *string `locationName:"resourceId" type:"string" required:"true"`
//   RoleArn *string `locationName:"roleArn" type:"string" required:"true"`
//   SystemIds []*string `locationName:"systemIds" type:"list" required:"true"`
//   Url *string `locationName:"url" type:"string" required:"true"`
// }
func buildSpekeKeyProvider(d *schema.ResourceData) *mediapackage.SpekeKeyProvider {
	var providerMap = (d.Get("speke_key_provider").([]interface{})[0].(map[string]interface{}))
	var provider = &mediapackage.SpekeKeyProvider{
		CertificateArn: aws.String(providerMap["certificate_arn"].(string)),
		ResourceId:     aws.String(providerMap["resource_id"].(string)),
		RoleArn:        aws.String(providerMap["role_arn"].(string)),
		SystemIds:      aws.StringSlice(providerMap["system_ids"].([]string)),
		Url:            aws.String(providerMap["url"].(string)),
	}
	return provider
}

type HlsManifestCreateOrUpdateParameters struct {
	AdMarkers                      *string   `locationName:"adMarkers" type:"string" enum:"AdMarkers"`
	AdTriggers                     []*string `locationName:"adTriggers" type:"list"`
	AdsOnDeliveryRestrictions      *string   `locationName:"adsOnDeliveryRestrictions" type:"string" enum:"AdsOnDeliveryRestrictions"`
	Id                             *string   `locationName:"id" type:"string" required:"true"`
	IncludeIframeOnlyStream        *bool     `locationName:"includeIframeOnlyStream" type:"boolean"`
	ManifestName                   *string   `locationName:"manifestName" type:"string"`
	PlaylistType                   *string   `locationName:"playlistType" type:"string" enum:"PlaylistType"`
	PlaylistWindowSeconds          *int64    `locationName:"playlistWindowSeconds" type:"integer"`
	ProgramDateTimeIntervalSeconds *int64    `locationName:"programDateTimeIntervalSeconds" type:"integer"`
}

func buildHlsManifestList(d *schema.ResourceData) []*mediapackage.HlsManifestCreateOrUpdateParameters {
	var manifestList = d.Get("hls_manifests").([]interface{})

	var hlsManifests []*mediapackage.HlsManifestCreateOrUpdateParameters
	for _, m := range manifestList {
		manifest := buildHlsManifest(m.(map[string]interface{}))

		hlsManifests = append(hlsManifests, manifest)
	}

	return hlsManifests
}

// type HlsManifest struct {
//   AdMarkers *string `locationName:"adMarkers" type:"string" enum:"AdMarkers"`
//   AdTriggers []*string `locationName:"adTriggers" type:"list"`
//   AdsOnDeliveryRestrictions *string `locationName:"adsOnDeliveryRestrictions" type:"string" enum:"AdsOnDeliveryRestrictions"`
//   Id *string `locationName:"id" type:"string" required:"true"`
//   IncludeIframeOnlyStream *bool `locationName:"includeIframeOnlyStream" type:"boolean"`
//   ManifestName *string `locationName:"manifestName" type:"string"`
//   PlaylistType *string `locationName:"playlistType" type:"string" enum:"PlaylistType"`
//   PlaylistWindowSeconds *int64 `locationName:"playlistWindowSeconds" type:"integer"`
//   ProgramDateTimeIntervalSeconds *int64 `locationName:"programDateTimeIntervalSeconds" type:"integer"`
// }
func buildHlsManifest(m map[string]interface{}) *mediapackage.HlsManifestCreateOrUpdateParameters {
	var manifest = &mediapackage.HlsManifestCreateOrUpdateParameters{
		AdMarkers:                      aws.String(m["ad_markers"].(string)),
		AdTriggers:                     aws.StringSlice(m["ad_triggers"].([]string)),
		AdsOnDeliveryRestrictions:      aws.String(m["ads_on_delivery_restrictions"].(string)),
		Id:                             aws.String(m["id"].(string)),
		IncludeIframeOnlyStream:        aws.Bool(m["include_iframe_only_stream"].(bool)),
		ManifestName:                   aws.String(m["manifest_name"].(string)),
		PlaylistType:                   aws.String(m["playlist_type"].(string)),
		PlaylistWindowSeconds:          aws.Int64(int64(m["playlist_window_seconds"].(int))),
		ProgramDateTimeIntervalSeconds: aws.Int64(int64(m["program_date_time_interval_seconds"].(int))),
	}
	return manifest
}

// type StreamSelection struct {
//   MaxVideoBitsPerSecond *int64 `locationName:"maxVideoBitsPerSecond" type:"integer"`
//   MinVideoBitsPerSecond *int64 `locationName:"minVideoBitsPerSecond" type:"integer"`
//   StreamOrder *string `locationName:"streamOrder" type:"string" enum:"StreamOrder"`
// }
func buildStreamSelection(d *schema.ResourceData) *mediapackage.StreamSelection {
	var selection = &mediapackage.StreamSelection{
		MaxVideoBitsPerSecond: aws.Int64(int64(d.Get("max_video_bits_per_second").(int))),
		MinVideoBitsPerSecond: aws.Int64(int64(d.Get("min_video_bits_per_second").(int))),
		StreamOrder:           aws.String(d.Get("stream_order").(string)),
	}
	return selection
}

// type OriginEndpointDetail interface {
// 	func SetArn(v string) *interface{}
// 	func SetChannelId(v string) *interface{}
// 	func SetCmafPackage(v *CmafPackage) *interface{}
// 	func SetDashPackage(v *DashPackage) *interface{}
// 	func SetDescription(v string) *interface{}
// 	func SetHlsPackage(v *HlsPackage) *interface{}
// 	func SetId(v string) *interface{}
// 	func SetManifestName(v string) *interface{}
// 	func SetMssPackage(v *MssPackage) *interface{}
// 	func SetOrigination(v string) *interface{}
// 	func SetStartoverWindowSeconds(v int64) *interface{}
// 	func SetTags(v map[string]*string) *interface{}
// 	func SetTimeDelaySeconds(v int64) *interface{}
// 	func SetUrl(v string) *interface{}
// 	func SetWhitelist(v []*string) *interface{}
// }

// type CreateOriginEndpointOutput struct {
//   Arn *string `locationName:"arn" type:"string"`
//   ChannelId *string `locationName:"channelId" type:"string"`
//   // A Common Media Application Format (CMAF) packaging configuration.
//   CmafPackage *CmafPackage `locationName:"cmafPackage" type:"structure"`
//   // A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
//   DashPackage *DashPackage `locationName:"dashPackage" type:"structure"`
//   Description *string `locationName:"description" type:"string"`
//   // An HTTP Live Streaming (HLS) packaging configuration.
//   HlsPackage *HlsPackage `locationName:"hlsPackage" type:"structure"`
//   Id *string `locationName:"id" type:"string"`
//   ManifestName *string `locationName:"manifestName" type:"string"`
//   // A Microsoft Smooth Streaming (MSS) packaging configuration.
//   MssPackage *MssPackage `locationName:"mssPackage" type:"structure"`
//   Origination *string `locationName:"origination" type:"string" enum:"Origination"`
//   StartoverWindowSeconds *int64 `locationName:"startoverWindowSeconds" type:"integer"`
//   // A collection of tags associated with a resource
//   Tags map[string]*string `locationName:"tags" type:"map"`
//   TimeDelaySeconds *int64 `locationName:"timeDelaySeconds" type:"integer"`
//   Url *string `locationName:"url" type:"string"`
//   Whitelist []*string `locationName:"whitelist" type:"list"`
//   // contains filtered or unexported fields
// }
// func addOriginEndpointDetails(d *schema.ResourceData, obj *OriginEndpointDetail) error {
// 	obj.SetArn(aws.String(d.Get("arn").(string)))
// 	obj.SetChannelId(aws.String(d.Get("channel_id").(string)))
// 	obj.SetCmafPackage(buildCmafPackage(d))
// 	obj.SetDashPackage(buildDashPackage(d))
// 	obj.SetDescription(aws.String(d.Get("description").(string)))
// 	obj.SetHlsPackage(buildHlsPackage(d))
// 	obj.SetId(aws.String(d.Get("id").(string)))
// 	obj.SetManifestName(aws.String(d.Get("manifest_name").(string)))
// 	obj.SetMssPackage(buildMssPackage(d))
// 	obj.SetOrigination(aws.String(d.Get("origination").(string)))
// 	obj.SetStartoverWindowSeconds(aws.Int64(d.Get("startover_window_seconds").(int))))
// 	obj.SetTimeDelaySeconds(aws.Int64(d.Get("time_delay_seconds").(int))))
// 	obj.SetUrl(aws.String(d.Get("url").(string)))
// 	obj.SetWhitelist(aws.StringSlice(d.Get("whitelist").([]string)))

// 	if err := setTagsMediaPackage(conn, d, d.Get("arn").(string)); err != nil {
// 		return fmt.Errorf("error updating MediaPackage Endpoint (%s) tags: %s", d.Id(), err)
// 	}
// }

func extractMediaPackageEndpointValues(d *schema.ResourceData, resp *mediapackage.DescribeOriginEndpointOutput) error {
	d.Set("arn", resp.Arn)
	d.Set("channel_id", resp.ChannelId)
	d.Set("description", resp.Description)
	d.Set("endpoint_id", resp.Id)
	d.Set("manifest_name", resp.ManifestName)
	d.Set("origination", resp.Origination)
	d.Set("startover_window_seconds", resp.StartoverWindowSeconds)
	d.Set("time_delay_seconds", resp.TimeDelaySeconds)
	d.Set("url", resp.Url)
	d.Set("whitelist", aws.StringValueSlice(resp.Whitelist))

	var typeVal string
	if resp.CmafPackage != nil {
		typeVal = "CMAF"
		if err := extractCmafPackageValues(d, resp.CmafPackage); err != nil {
			return fmt.Errorf("error parsing CmafPackage: %s", err)
		}
	} else if resp.DashPackage != nil {
		typeVal = "DASH"
		if err := extractDashPackageValues(d, resp.DashPackage); err != nil {
			return fmt.Errorf("error parsing DashPackage: %s", err)
		}
	} else if resp.HlsPackage != nil {
		typeVal = "HLS"
		if err := extractHlsPackageValues(d, resp.HlsPackage); err != nil {
			return fmt.Errorf("error parsing HlsPackage: %s", err)
		}
	} else if resp.MssPackage != nil {
		typeVal = "MSS"
		if err := extractMssPackageValues(d, resp.MssPackage); err != nil {
			return fmt.Errorf("error parsing MssPackage: %s", err)
		}
	}

	if typeVal == "" {
		return fmt.Errorf("MediaPackage Endpoint (%s) does not have a packaging configuration or packaging type is unknown", d.Id())
	}

	if err := d.Set("tags", tagsToMapGeneric(resp.Tags)); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

// type SpekeKeyProvider struct {
//   CertificateArn *string `locationName:"certificateArn" type:"string"`
//   ResourceId *string `locationName:"resourceId" type:"string" required:"true"`
//   RoleArn *string `locationName:"roleArn" type:"string" required:"true"`
//   SystemIds []*string `locationName:"systemIds" type:"list" required:"true"`
//   Url *string `locationName:"url" type:"string" required:"true"`
// }
func dereferenceSpekeKeyProviderValues(provider *mediapackage.SpekeKeyProvider) map[string]interface{} {
	return map[string]interface{}{
		"certificate_arn": aws.StringValue(provider.CertificateArn),
		"resource_id":     aws.StringValue(provider.ResourceId),
		"role_arn":        aws.StringValue(provider.RoleArn),
		"system_ids":      aws.StringValueSlice(provider.SystemIds),
		"url":             aws.StringValue(provider.Url),
	}
}

func extractSpekeKeyProviderValues(d *schema.ResourceData, skProvider *mediapackage.SpekeKeyProvider) error {
	if skProvider == nil {
		return nil
	}

	var provider []map[string]interface{}
	values := dereferenceSpekeKeyProviderValues(skProvider)
	provider = append(provider, values)

	d.Set("speke_key_provider", provider)

	return nil
}

// type CMAFEncryption struct {
//   KeyRotationIntervalSeconds *int64 `locationName:"keyRotationIntervalSeconds" type:"integer"`
//   SpekeKeyProvider *SpekeKeyProvider `locationName:"spekeKeyProvider" type:"structure" required:"true"`
// }
func extractCmafEncryptionValues(d *schema.ResourceData, crypt *mediapackage.CmafEncryption) error {
	if crypt == nil {
		return nil
	}

	d.Set("key_rotation_interval_seconds", crypt.KeyRotationIntervalSeconds)
	extractSpekeKeyProviderValues(d, crypt.SpekeKeyProvider)

	return nil
}

// type HlsManifest struct {
//   AdMarkers *string `locationName:"adMarkers" type:"string" enum:"AdMarkers"`
//   Id *string `locationName:"id" type:"string" required:"true"`
//   IncludeIframeOnlyStream *bool `locationName:"includeIframeOnlyStream" type:"boolean"`
//   ManifestName *string `locationName:"manifestName" type:"string"`
//   PlaylistType *string `locationName:"playlistType" type:"string" enum:"PlaylistType"`
//   PlaylistWindowSeconds *int64 `locationName:"playlistWindowSeconds" type:"integer"`
//   ProgramDateTimeIntervalSeconds *int64 `locationName:"programDateTimeIntervalSeconds" type:"integer"`
//   Url *string `locationName:"url" type:"string"`
// }
func dereferenceManifestValues(manifest *mediapackage.HlsManifest) map[string]interface{} {
	return map[string]interface{}{
		"ad_markers":                         aws.StringValue(manifest.AdMarkers),
		"id":                                 aws.StringValue(manifest.Id),
		"include_iframe_only_stream":         aws.BoolValue(manifest.IncludeIframeOnlyStream),
		"manifest_name":                      aws.StringValue(manifest.ManifestName),
		"playlist_type":                      aws.StringValue(manifest.PlaylistType),
		"playlist_window_seconds":            aws.Int64Value(manifest.PlaylistWindowSeconds),
		"program_date_time_interval_seconds": aws.Int64Value(manifest.ProgramDateTimeIntervalSeconds),
		"url":                                aws.StringValue(manifest.Url),
	}
}

func extractHlsManifestList(d *schema.ResourceData, manifestList []*mediapackage.HlsManifest) error {
	var hlsManifests []map[string]interface{}
	for _, m := range manifestList {
		manifest := dereferenceManifestValues(m)

		hlsManifests = append(hlsManifests, manifest)
	}

	d.Set("hls_manifests", hlsManifests)

	return nil
}

// type StreamSelection struct {
//   MaxVideoBitsPerSecond *int64 `locationName:"maxVideoBitsPerSecond" type:"integer"`
//   MinVideoBitsPerSecond *int64 `locationName:"minVideoBitsPerSecond" type:"integer"`
//   StreamOrder *string `locationName:"streamOrder" type:"string" enum:"StreamOrder"`
// }
func extractStreamSelectionValues(d *schema.ResourceData, stream *mediapackage.StreamSelection) error {
	d.Set("max_video_bits_per_second", stream.MaxVideoBitsPerSecond)
	d.Set("min_video_bits_per_second", stream.MinVideoBitsPerSecond)
	d.Set("stream_order", stream.StreamOrder)

	return nil
}

// A Common Media Application Format (CMAF) packaging configuration.
// type CmafPackage struct { // type DashEncryption struct {
// 	Encryption *CmafEncryption `locationName:"encryption" type:"structure"`
// 	HlsManifests []*HlsManifest `locationName:"hlsManifests" type:"list"`
// 	SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
// 	SegmentPrefix *string `locationName:"segmentPrefix" type:"string"`
// 	StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
// }
func extractCmafPackageValues(d *schema.ResourceData, pack *mediapackage.CmafPackage) error {
	if pack == nil {
		return nil
	}

	extractCmafEncryptionValues(d, pack.Encryption)
	extractHlsManifestList(d, pack.HlsManifests)
	d.Set("segment_duration_seconds", pack.SegmentDurationSeconds)
	d.Set("segment_prefix", pack.SegmentPrefix)
	extractStreamSelectionValues(d, pack.StreamSelection)

	return nil
}

// type DASHEncryption struct {
//   KeyRotationIntervalSeconds *int64 `locationName:"keyRotationIntervalSeconds" type:"integer"`
//   SpekeKeyProvider *SpekeKeyProvider `locationName:"spekeKeyProvider" type:"structure" required:"true"`
// }
func extractDashEncryptionValues(d *schema.ResourceData, crypt *mediapackage.DashEncryption) error {
	if crypt == nil {
		return nil
	}

	d.Set("key_rotation_interval_seconds", crypt.KeyRotationIntervalSeconds)
	extractSpekeKeyProviderValues(d, crypt.SpekeKeyProvider)

	return nil
}

// A Dynamic Adaptive Streaming over HTTP (DASH) packaging configuration.
// type DashPackage struct {
//   AdTriggers []*string `locationName:"adTriggers" type:"list"`
//   AdsOnDeliveryRestrictions *string `locationName:"adsOnDeliveryRestrictions" type:"string" enum:"AdsOnDeliveryRestrictions"`
//   Encryption *DashEncryption `locationName:"encryption" type:"structure"`
//   ManifestLayout *string `locationName:"manifestLayout" type:"string" enum:"ManifestLayout"`
//   ManifestWindowSeconds *int64 `locationName:"manifestWindowSeconds" type:"integer"`
//   MinBufferTimeSeconds *int64 `locationName:"minBufferTimeSeconds" type:"integer"`
//   MinUpdatePeriodSeconds *int64 `locationName:"minUpdatePeriodSeconds" type:"integer"`
//   PeriodTriggers []*string `locationName:"periodTriggers" type:"list"`
//   Profile *string `locationName:"profile" type:"string" enum:"Profile"`
//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
//   SegmentTemplateFormat *string `locationName:"segmentTemplateFormat" type:"string" enum:"SegmentTemplateFormat"`
//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
//   SuggestedPresentationDelaySeconds *int64 `locationName:"suggestedPresentationDelaySeconds" type:"integer"`
// }
func extractDashPackageValues(d *schema.ResourceData, pack *mediapackage.DashPackage) error {
	if pack == nil {
		return nil
	}

	d.Set("ad_triggers", aws.StringValueSlice(pack.AdTriggers))
	d.Set("ads_on_delivery_restrictions", pack.AdsOnDeliveryRestrictions)
	extractDashEncryptionValues(d, pack.Encryption)
	d.Set("manifest_layout", pack.ManifestLayout)
	d.Set("manifest_window_seconds", pack.ManifestWindowSeconds)
	d.Set("min_buffer_time_seconds", pack.MinBufferTimeSeconds)
	d.Set("min_update_period_seconds", pack.MinUpdatePeriodSeconds)
	d.Set("period_triggers", aws.StringValueSlice(pack.PeriodTriggers))
	d.Set("profile", pack.Profile)
	d.Set("segment_duration_seconds", pack.SegmentDurationSeconds)
	d.Set("segment_template_format", pack.SegmentTemplateFormat)
	extractStreamSelectionValues(d, pack.StreamSelection)
	d.Set("suggested_presentation_delay_seconds", pack.SuggestedPresentationDelaySeconds)

	return nil
}

// type HLSEncryption struct {
//   ConstantInitializationVector *string `locationName:"constantInitializationVector" type:"string"`
//   EncryptionMethod *string `locationName:"encryptionMethod" type:"string" enum:"EncryptionMethod"`
//   KeyRotationIntervalSeconds *int64 `locationName:"keyRotationIntervalSeconds" type:"integer"`
//   RepeatExtXKey *bool `locationName:"repeatExtXKey" type:"boolean"`
//   SpekeKeyProvider *SpekeKeyProvider `locationName:"spekeKeyProvider" type:"structure" required:"true"`
// }
func extractHlsEncryptionValues(d *schema.ResourceData, crypt *mediapackage.HlsEncryption) error {
	if crypt == nil {
		return nil
	}

	d.Set("constant_initialization_vector", crypt.ConstantInitializationVector)
	d.Set("encryption_method", crypt.EncryptionMethod)
	d.Set("key_rotation_interval_seconds", crypt.KeyRotationIntervalSeconds)
	d.Set("repeat_ext_x_key", crypt.RepeatExtXKey)
	extractSpekeKeyProviderValues(d, crypt.SpekeKeyProvider)

	return nil
}

// An HTTP Live Streaming (HLS) packaging configuration.
// type HlsPackage struct {
//   AdMarkers *string `locationName:"adMarkers" type:"string" enum:"AdMarkers"`
//   AdTriggers []*string `locationName:"adTriggers" type:"list"`
//   AdsOnDeliveryRestrictions *string `locationName:"adsOnDeliveryRestrictions" type:"string" enum:"AdsOnDeliveryRestrictions"`
//   Encryption *HlsEncryption `locationName:"encryption" type:"structure"`
//   IncludeIframeOnlyStream *bool `locationName:"includeIframeOnlyStream" type:"boolean"`
//   PlaylistType *string `locationName:"playlistType" type:"string" enum:"PlaylistType"`
//   PlaylistWindowSeconds *int64 `locationName:"playlistWindowSeconds" type:"integer"`
//   ProgramDateTimeIntervalSeconds *int64 `locationName:"programDateTimeIntervalSeconds" type:"integer"`
//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
//   UseAudioRenditionGroup *bool `locationName:"useAudioRenditionGroup" type:"boolean"`
// }
func extractHlsPackageValues(d *schema.ResourceData, pack *mediapackage.HlsPackage) error {
	if pack == nil {
		return nil
	}

	d.Set("ad_markers", pack.AdMarkers)
	d.Set("ad_triggers", aws.StringValueSlice(pack.AdTriggers))
	d.Set("ads_on_delivery_restrictions", pack.AdsOnDeliveryRestrictions)
	extractHlsEncryptionValues(d, pack.Encryption)
	d.Set("include_iframe_only_stream", pack.IncludeIframeOnlyStream)
	d.Set("playlist_type", pack.PlaylistType)
	d.Set("playlist_window_seconds", pack.PlaylistWindowSeconds)
	d.Set("program_date_time_interval_seconds", pack.ProgramDateTimeIntervalSeconds)
	d.Set("segment_duration_seconds", pack.SegmentDurationSeconds)
	extractStreamSelectionValues(d, pack.StreamSelection)
	d.Set("use_audio_rendition_group", pack.UseAudioRenditionGroup)

	return nil
}

// type MSSEncryption struct {
//   SpekeKeyProvider *SpekeKeyProvider `locationName:"spekeKeyProvider" type:"structure" required:"true"`
// }
func extractMssEncryptionValues(d *schema.ResourceData, crypt *mediapackage.MssEncryption) error {
	if crypt == nil {
		return nil
	}

	extractSpekeKeyProviderValues(d, crypt.SpekeKeyProvider)

	return nil
}

// A Microsoft Smooth Streaming (MSS) packaging configuration.
// type MssPackage struct {
//   Encryption *MssEncryption `locationName:"encryption" type:"structure"`
//   ManifestWindowSeconds *int64 `locationName:"manifestWindowSeconds" type:"integer"`
//   SegmentDurationSeconds *int64 `locationName:"segmentDurationSeconds" type:"integer"`
//   StreamSelection *StreamSelection `locationName:"streamSelection" type:"structure"`
// }
func extractMssPackageValues(d *schema.ResourceData, pack *mediapackage.MssPackage) error {
	if pack == nil {
		return nil
	}

	extractMssEncryptionValues(d, pack.Encryption)
	d.Set("manifest_window_seconds", pack.ManifestWindowSeconds)
	d.Set("segment_duration_seconds", pack.SegmentDurationSeconds)
	extractStreamSelectionValues(d, pack.StreamSelection)

	return nil
}

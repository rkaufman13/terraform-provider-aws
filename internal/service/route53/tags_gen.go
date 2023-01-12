// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package route53

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// ListTags lists route53 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn route53iface.Route53API, identifier string, resourceType string) (tftags.KeyValueTags, error) {
	return ListTagsWithContext(context.Background(), conn, identifier, resourceType)
}

func ListTagsWithContext(ctx context.Context, conn route53iface.Route53API, identifier string, resourceType string) (tftags.KeyValueTags, error) {
	input := &route53.ListTagsForResourceInput{
		ResourceId:   aws.String(identifier),
		ResourceType: aws.String(resourceType),
	}

	output, err := conn.ListTagsForResourceWithContext(ctx, input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.ResourceTagSet.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns route53 service tags.
func Tags(tags tftags.KeyValueTags) []*route53.Tag {
	result := make([]*route53.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &route53.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from route53 service tags.
func KeyValueTags(tags []*route53.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(m)
}

// UpdateTags updates route53 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn route53iface.Route53API, identifier string, resourceType string, oldTags interface{}, newTags interface{}) error {
	return UpdateTagsWithContext(context.Background(), conn, identifier, resourceType, oldTags, newTags)
}
func UpdateTagsWithContext(ctx context.Context, conn route53iface.Route53API, identifier string, resourceType string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)
	removedTags := oldTags.Removed(newTags)
	updatedTags := oldTags.Updated(newTags)

	// Ensure we do not send empty requests
	if len(removedTags) == 0 && len(updatedTags) == 0 {
		return nil
	}

	input := &route53.ChangeTagsForResourceInput{
		ResourceId:   aws.String(identifier),
		ResourceType: aws.String(resourceType),
	}

	if len(updatedTags) > 0 {
		input.AddTags = Tags(updatedTags.IgnoreAWS())
	}

	if len(removedTags) > 0 {
		input.RemoveTagKeys = aws.StringSlice(removedTags.Keys())
	}

	_, err := conn.ChangeTagsForResourceWithContext(ctx, input)

	if err != nil {
		return fmt.Errorf("tagging resource (%s): %w", identifier, err)
	}

	return nil
}
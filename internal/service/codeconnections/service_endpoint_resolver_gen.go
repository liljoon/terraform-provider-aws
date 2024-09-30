// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package codeconnections

import (
	"context"
	"fmt"
	"net"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	codeconnections_sdkv2 "github.com/aws/aws-sdk-go-v2/service/codeconnections"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
)

var _ codeconnections_sdkv2.EndpointResolverV2 = resolverSDKv2{}

type resolverSDKv2 struct {
	defaultResolver codeconnections_sdkv2.EndpointResolverV2
}

func newEndpointResolverSDKv2() resolverSDKv2 {
	return resolverSDKv2{
		defaultResolver: codeconnections_sdkv2.NewDefaultEndpointResolverV2(),
	}
}

func (r resolverSDKv2) ResolveEndpoint(ctx context.Context, params codeconnections_sdkv2.EndpointParameters) (endpoint smithyendpoints.Endpoint, err error) {
	params = params.WithDefaults()
	useFIPS := aws_sdkv2.ToBool(params.UseFIPS)

	if eps := params.Endpoint; aws_sdkv2.ToString(eps) != "" {
		tflog.Debug(ctx, "setting endpoint", map[string]any{
			"tf_aws.endpoint": endpoint,
		})

		if useFIPS {
			tflog.Debug(ctx, "endpoint set, ignoring UseFIPSEndpoint setting")
			params.UseFIPS = aws_sdkv2.Bool(false)
		}

		return r.defaultResolver.ResolveEndpoint(ctx, params)
	} else if useFIPS {
		ctx = tflog.SetField(ctx, "tf_aws.use_fips", useFIPS)

		endpoint, err = r.defaultResolver.ResolveEndpoint(ctx, params)
		if err != nil {
			return endpoint, err
		}

		tflog.Debug(ctx, "endpoint resolved", map[string]any{
			"tf_aws.endpoint": endpoint.URI.String(),
		})

		hostname := endpoint.URI.Hostname()
		_, err = net.LookupHost(hostname)
		if err != nil {
			if dnsErr, ok := errs.As[*net.DNSError](err); ok && dnsErr.IsNotFound {
				tflog.Debug(ctx, "default endpoint host not found, disabling FIPS", map[string]any{
					"tf_aws.hostname": hostname,
				})
				params.UseFIPS = aws_sdkv2.Bool(false)
			} else {
				err = fmt.Errorf("looking up codeconnections endpoint %q: %s", hostname, err)
				return
			}
		} else {
			return endpoint, err
		}
	}

	return r.defaultResolver.ResolveEndpoint(ctx, params)
}

func withBaseEndpoint(endpoint string) func(*codeconnections_sdkv2.Options) {
	return func(o *codeconnections_sdkv2.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}
}

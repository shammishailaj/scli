// Code generated by smithy-go-codegen DO NOT EDIT.

package ec2

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	internalauth "github.com/aws/aws-sdk-go-v2/internal/auth"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Allocates an Elastic IP address to your Amazon Web Services account. After you
// allocate the Elastic IP address you can associate it with an instance or network
// interface. After you release an Elastic IP address, it is released to the IP
// address pool and can be allocated to a different Amazon Web Services account.
// You can allocate an Elastic IP address from an address pool owned by Amazon Web
// Services or from an address pool created from a public IPv4 address range that
// you have brought to Amazon Web Services for use with your Amazon Web Services
// resources using bring your own IP addresses (BYOIP). For more information, see
// Bring Your Own IP Addresses (BYOIP) (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-byoip.html)
// in the Amazon Elastic Compute Cloud User Guide. If you release an Elastic IP
// address, you might be able to recover it. You cannot recover an Elastic IP
// address that you released after it is allocated to another Amazon Web Services
// account. To attempt to recover an Elastic IP address that you released, specify
// it in this operation. For more information, see Elastic IP Addresses (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/elastic-ip-addresses-eip.html)
// in the Amazon Elastic Compute Cloud User Guide. You can allocate a carrier IP
// address which is a public IP address from a telecommunication carrier, to a
// network interface which resides in a subnet in a Wavelength Zone (for example an
// EC2 instance).
func (c *Client) AllocateAddress(ctx context.Context, params *AllocateAddressInput, optFns ...func(*Options)) (*AllocateAddressOutput, error) {
	if params == nil {
		params = &AllocateAddressInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "AllocateAddress", params, optFns, c.addOperationAllocateAddressMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*AllocateAddressOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type AllocateAddressInput struct {

	// The Elastic IP address to recover or an IPv4 address from an address pool.
	Address *string

	// The ID of a customer-owned address pool. Use this parameter to let Amazon EC2
	// select an address from the address pool. Alternatively, specify a specific
	// address from the address pool.
	CustomerOwnedIpv4Pool *string

	// The network ( vpc ).
	Domain types.DomainType

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have the
	// required permissions, the error response is DryRunOperation . Otherwise, it is
	// UnauthorizedOperation .
	DryRun *bool

	// A unique set of Availability Zones, Local Zones, or Wavelength Zones from which
	// Amazon Web Services advertises IP addresses. Use this parameter to limit the IP
	// address to this location. IP addresses cannot move between network border
	// groups. Use DescribeAvailabilityZones (https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeAvailabilityZones.html)
	// to view the network border groups. You cannot use a network border group with
	// EC2 Classic. If you attempt this operation on EC2 Classic, you receive an
	// InvalidParameterCombination error.
	NetworkBorderGroup *string

	// The ID of an address pool that you own. Use this parameter to let Amazon EC2
	// select an address from the address pool. To specify a specific address from the
	// address pool, use the Address parameter instead.
	PublicIpv4Pool *string

	// The tags to assign to the Elastic IP address.
	TagSpecifications []types.TagSpecification

	noSmithyDocumentSerde
}

type AllocateAddressOutput struct {

	// The ID that represents the allocation of the Elastic IP address.
	AllocationId *string

	// The carrier IP address. This option is only available for network interfaces
	// that reside in a subnet in a Wavelength Zone.
	CarrierIp *string

	// The customer-owned IP address.
	CustomerOwnedIp *string

	// The ID of the customer-owned address pool.
	CustomerOwnedIpv4Pool *string

	// The network ( vpc ).
	Domain types.DomainType

	// The set of Availability Zones, Local Zones, or Wavelength Zones from which
	// Amazon Web Services advertises IP addresses.
	NetworkBorderGroup *string

	// The Elastic IP address.
	PublicIp *string

	// The ID of an address pool.
	PublicIpv4Pool *string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationAllocateAddressMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsEc2query_serializeOpAllocateAddress{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsEc2query_deserializeOpAllocateAddress{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addAllocateAddressResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opAllocateAddress(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addendpointDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opAllocateAddress(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ec2",
		OperationName: "AllocateAddress",
	}
}

type opAllocateAddressResolveEndpointMiddleware struct {
	EndpointResolver EndpointResolverV2
	BuiltInResolver  builtInParameterResolver
}

func (*opAllocateAddressResolveEndpointMiddleware) ID() string {
	return "ResolveEndpointV2"
}

func (m *opAllocateAddressResolveEndpointMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	if awsmiddleware.GetRequiresLegacyEndpoints(ctx) {
		return next.HandleSerialize(ctx, in)
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
	}

	if m.EndpointResolver == nil {
		return out, metadata, fmt.Errorf("expected endpoint resolver to not be nil")
	}

	params := EndpointParameters{}

	m.BuiltInResolver.ResolveBuiltIns(&params)

	var resolvedEndpoint smithyendpoints.Endpoint
	resolvedEndpoint, err = m.EndpointResolver.ResolveEndpoint(ctx, params)
	if err != nil {
		return out, metadata, fmt.Errorf("failed to resolve service endpoint, %w", err)
	}

	req.URL = &resolvedEndpoint.URI

	for k := range resolvedEndpoint.Headers {
		req.Header.Set(
			k,
			resolvedEndpoint.Headers.Get(k),
		)
	}

	authSchemes, err := internalauth.GetAuthenticationSchemes(&resolvedEndpoint.Properties)
	if err != nil {
		var nfe *internalauth.NoAuthenticationSchemesFoundError
		if errors.As(err, &nfe) {
			// if no auth scheme is found, default to sigv4
			signingName := "ec2"
			signingRegion := m.BuiltInResolver.(*builtInResolver).Region
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)

		}
		var ue *internalauth.UnSupportedAuthenticationSchemeSpecifiedError
		if errors.As(err, &ue) {
			return out, metadata, fmt.Errorf(
				"This operation requests signer version(s) %v but the client only supports %v",
				ue.UnsupportedSchemes,
				internalauth.SupportedSchemes,
			)
		}
	}

	for _, authScheme := range authSchemes {
		switch authScheme.(type) {
		case *internalauth.AuthenticationSchemeV4:
			v4Scheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4)
			var signingName, signingRegion string
			if v4Scheme.SigningName == nil {
				signingName = "ec2"
			} else {
				signingName = *v4Scheme.SigningName
			}
			if v4Scheme.SigningRegion == nil {
				signingRegion = m.BuiltInResolver.(*builtInResolver).Region
			} else {
				signingRegion = *v4Scheme.SigningRegion
			}
			if v4Scheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4Scheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)
			break
		case *internalauth.AuthenticationSchemeV4A:
			v4aScheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4A)
			if v4aScheme.SigningName == nil {
				v4aScheme.SigningName = aws.String("ec2")
			}
			if v4aScheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4aScheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, *v4aScheme.SigningName)
			ctx = awsmiddleware.SetSigningRegion(ctx, v4aScheme.SigningRegionSet[0])
			break
		case *internalauth.AuthenticationSchemeNone:
			break
		}
	}

	return next.HandleSerialize(ctx, in)
}

func addAllocateAddressResolveEndpointMiddleware(stack *middleware.Stack, options Options) error {
	return stack.Serialize.Insert(&opAllocateAddressResolveEndpointMiddleware{
		EndpointResolver: options.EndpointResolverV2,
		BuiltInResolver: &builtInResolver{
			Region:       options.Region,
			UseDualStack: options.EndpointOptions.UseDualStackEndpoint,
			UseFIPS:      options.EndpointOptions.UseFIPSEndpoint,
			Endpoint:     options.BaseEndpoint,
		},
	}, "ResolveEndpoint", middleware.After)
}
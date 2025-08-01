package settings

import (
	"context"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/query"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	settings "github.com/zitadel/zitadel/pkg/grpc/settings/v2beta"
)

func (s *Server) GetLoginSettings(ctx context.Context, req *connect.Request[settings.GetLoginSettingsRequest]) (*connect.Response[settings.GetLoginSettingsResponse], error) {
	current, err := s.query.LoginPolicyByID(ctx, true, object.ResourceOwnerFromReq(ctx, req.Msg.GetCtx()), false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.GetLoginSettingsResponse{
		Settings: loginSettingsToPb(current),
		Details: &object_pb.Details{
			Sequence:      current.Sequence,
			CreationDate:  timestamppb.New(current.CreationDate),
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.OrgID,
		},
	}), nil
}

func (s *Server) GetPasswordComplexitySettings(ctx context.Context, req *connect.Request[settings.GetPasswordComplexitySettingsRequest]) (*connect.Response[settings.GetPasswordComplexitySettingsResponse], error) {
	current, err := s.query.PasswordComplexityPolicyByOrg(ctx, true, object.ResourceOwnerFromReq(ctx, req.Msg.GetCtx()), false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.GetPasswordComplexitySettingsResponse{
		Settings: passwordComplexitySettingsToPb(current),
		Details: &object_pb.Details{
			Sequence:      current.Sequence,
			CreationDate:  timestamppb.New(current.CreationDate),
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.ResourceOwner,
		},
	}), nil
}

func (s *Server) GetPasswordExpirySettings(ctx context.Context, req *connect.Request[settings.GetPasswordExpirySettingsRequest]) (*connect.Response[settings.GetPasswordExpirySettingsResponse], error) {
	current, err := s.query.PasswordAgePolicyByOrg(ctx, true, object.ResourceOwnerFromReq(ctx, req.Msg.GetCtx()), false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.GetPasswordExpirySettingsResponse{
		Settings: passwordExpirySettingsToPb(current),
		Details: &object_pb.Details{
			Sequence:      current.Sequence,
			CreationDate:  timestamppb.New(current.CreationDate),
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.ResourceOwner,
		},
	}), nil
}

func (s *Server) GetBrandingSettings(ctx context.Context, req *connect.Request[settings.GetBrandingSettingsRequest]) (*connect.Response[settings.GetBrandingSettingsResponse], error) {
	current, err := s.query.ActiveLabelPolicyByOrg(ctx, object.ResourceOwnerFromReq(ctx, req.Msg.GetCtx()), false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.GetBrandingSettingsResponse{
		Settings: brandingSettingsToPb(current, s.assetsAPIDomain(ctx)),
		Details: &object_pb.Details{
			Sequence:      current.Sequence,
			CreationDate:  timestamppb.New(current.CreationDate),
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.ResourceOwner,
		},
	}), nil
}

func (s *Server) GetDomainSettings(ctx context.Context, req *connect.Request[settings.GetDomainSettingsRequest]) (*connect.Response[settings.GetDomainSettingsResponse], error) {
	current, err := s.query.DomainPolicyByOrg(ctx, true, object.ResourceOwnerFromReq(ctx, req.Msg.GetCtx()), false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.GetDomainSettingsResponse{
		Settings: domainSettingsToPb(current),
		Details: &object_pb.Details{
			Sequence:      current.Sequence,
			CreationDate:  timestamppb.New(current.CreationDate),
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.ResourceOwner,
		},
	}), nil
}

func (s *Server) GetLegalAndSupportSettings(ctx context.Context, req *connect.Request[settings.GetLegalAndSupportSettingsRequest]) (*connect.Response[settings.GetLegalAndSupportSettingsResponse], error) {
	current, err := s.query.PrivacyPolicyByOrg(ctx, true, object.ResourceOwnerFromReq(ctx, req.Msg.GetCtx()), false)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.GetLegalAndSupportSettingsResponse{
		Settings: legalAndSupportSettingsToPb(current),
		Details: &object_pb.Details{
			Sequence:      current.Sequence,
			CreationDate:  timestamppb.New(current.CreationDate),
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.ResourceOwner,
		},
	}), nil
}

func (s *Server) GetLockoutSettings(ctx context.Context, req *connect.Request[settings.GetLockoutSettingsRequest]) (*connect.Response[settings.GetLockoutSettingsResponse], error) {
	current, err := s.query.LockoutPolicyByOrg(ctx, true, object.ResourceOwnerFromReq(ctx, req.Msg.GetCtx()))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.GetLockoutSettingsResponse{
		Settings: lockoutSettingsToPb(current),
		Details: &object_pb.Details{
			Sequence:      current.Sequence,
			CreationDate:  timestamppb.New(current.CreationDate),
			ChangeDate:    timestamppb.New(current.ChangeDate),
			ResourceOwner: current.ResourceOwner,
		},
	}), nil
}

func (s *Server) GetActiveIdentityProviders(ctx context.Context, req *connect.Request[settings.GetActiveIdentityProvidersRequest]) (*connect.Response[settings.GetActiveIdentityProvidersResponse], error) {
	links, err := s.query.IDPLoginPolicyLinks(ctx, object.ResourceOwnerFromReq(ctx, req.Msg.GetCtx()), &query.IDPLoginPolicyLinksSearchQuery{}, false)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&settings.GetActiveIdentityProvidersResponse{
		Details:           object.ToListDetails(links.SearchResponse),
		IdentityProviders: identityProvidersToPb(links.Links),
	}), nil
}

func (s *Server) GetGeneralSettings(ctx context.Context, _ *connect.Request[settings.GetGeneralSettingsRequest]) (*connect.Response[settings.GetGeneralSettingsResponse], error) {
	instance := authz.GetInstance(ctx)
	return connect.NewResponse(&settings.GetGeneralSettingsResponse{
		SupportedLanguages: domain.LanguagesToStrings(i18n.SupportedLanguages()),
		DefaultOrgId:       instance.DefaultOrganisationID(),
		DefaultLanguage:    instance.DefaultLanguage().String(),
	}), nil
}

func (s *Server) GetSecuritySettings(ctx context.Context, req *connect.Request[settings.GetSecuritySettingsRequest]) (*connect.Response[settings.GetSecuritySettingsResponse], error) {
	policy, err := s.query.SecurityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.GetSecuritySettingsResponse{
		Settings: securityPolicyToSettingsPb(policy),
	}), nil
}

func (s *Server) SetSecuritySettings(ctx context.Context, req *connect.Request[settings.SetSecuritySettingsRequest]) (*connect.Response[settings.SetSecuritySettingsResponse], error) {
	details, err := s.command.SetSecurityPolicy(ctx, securitySettingsToCommand(req.Msg))
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&settings.SetSecuritySettingsResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
}

func (s *Server) SetOrganizationSettings(ctx context.Context, req *connect.Request[settings.SetOrganizationSettingsRequest]) (*connect.Response[settings.SetOrganizationSettingsResponse], error) {
	details, err := s.command.SetOrganizationSettings(ctx, organizationSettingsToCommand(req.Msg))
	if err != nil {
		return nil, err
	}
	var setDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		setDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&settings.SetOrganizationSettingsResponse{
		SetDate: setDate,
	}), nil
}

func (s *Server) DeleteOrganizationSettings(ctx context.Context, req *connect.Request[settings.DeleteOrganizationSettingsRequest]) (*connect.Response[settings.DeleteOrganizationSettingsResponse], error) {
	details, err := s.command.DeleteOrganizationSettings(ctx, req.Msg.GetOrganizationId())
	if err != nil {
		return nil, err
	}
	var deletionDate *timestamppb.Timestamp
	if !details.EventDate.IsZero() {
		deletionDate = timestamppb.New(details.EventDate)
	}
	return connect.NewResponse(&settings.DeleteOrganizationSettingsResponse{
		DeletionDate: deletionDate,
	}), nil
}

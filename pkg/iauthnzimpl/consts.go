/*
 * Copyright (c) 2022-present unTill Pro, Ltd.
 */

package iauthnzimpl

import "github.com/voedger/voedger/pkg/appdef"

var (
	qNameViewDeviceProfileWSIDIdx                   = appdef.NewQName(airPackage, "DeviceProfileWSIDIdx")
	qNameCDocWorkspaceKindRestaurant                = appdef.NewQName(airPackage, "Restaurant")
	qNameCDocWorkspaceKindAppWorkspace              = appdef.NewQName(appdef.SysPackage, "AppWorkspace")
	qNameCDocSubscriptionProfile                    = appdef.NewQName(airPackage, "SubscriptionProfile")
	qNameCDocUnTillOrders                           = appdef.NewQName(untillPackage, "orders")
	qNameCDocUnTillPBill                            = appdef.NewQName(untillPackage, "pbill")
	qNameTestDeniedCmd                              = appdef.NewQName(appdef.SysPackage, "TestDeniedCmd")
	qNameTestDeniedQry                              = appdef.NewQName(appdef.SysPackage, "TestDeniedQry")
	qNameTestDeniedCDoc                             = appdef.NewQName(appdef.SysPackage, "TestDeniedCDoc")
	qNameCDocLogin                                  = appdef.NewQName(appdef.SysPackage, "Login")
	qNameCDocChildWorkspace                         = appdef.NewQName(appdef.SysPackage, "ChildWorkspace")
	qNameCDocWorkspaceKindUser                      = appdef.NewQName(appdef.SysPackage, "UserProfile")
	qNameCDocWorkspaceKindDevice                    = appdef.NewQName(appdef.SysPackage, "DeviceProfile")
	qNameCDocWorkspaceDescriptor                    = appdef.NewQName(appdef.SysPackage, "WorkspaceDescriptor")
	qNameCmdUpdateSubscription                      = appdef.NewQName(airPackage, "UpdateSubscription")
	qNameCmdStoreSubscriptionProfile                = appdef.NewQName(airPackage, "StoreSubscriptionProfile")
	qNameCmdLinkDeviceToRestaurant                  = appdef.NewQName(airPackage, "LinkDeviceToRestaurant")
	qNameQryIssuePrincipalToken                     = appdef.NewQName(appdef.SysPackage, "IssuePrincipalToken")
	qNameCmdCreateLogin                             = appdef.NewQName(appdef.SysPackage, "CreateLogin")
	qNameQryEcho                                    = appdef.NewQName(appdef.SysPackage, "Echo")
	qNameQryGRCount                                 = appdef.NewQName(appdef.SysPackage, "GRCount")
	qNameCmdSendEmailVerificationCode               = appdef.NewQName(appdef.SysPackage, "SendEmailVerificationCode")
	qNameCmdResetPasswordByEmail                    = appdef.NewQName(appdef.SysPackage, "ResetPasswordByEmail")
	qNameQryInitiateResetPasswordByEmail            = appdef.NewQName(appdef.SysPackage, "InitiateResetPasswordByEmail")
	qNameQryIssueVerifiedValueTokenForResetPassword = appdef.NewQName(appdef.SysPackage, "IssueVerifiedValueTokenForResetPassword")
	qNameQryDescribePackageNames                    = appdef.NewQName(appdef.SysPackage, "DescribePackageNames")
	qNameQryDescribePackage                         = appdef.NewQName(appdef.SysPackage, "DescribePackage")
	qNameCmdInitiateJoinWorkspace                   = appdef.NewQName(appdef.SysPackage, "InitiateJoinWorkspace")
	qNameCmdInitiateLeaveWorkspace                  = appdef.NewQName(appdef.SysPackage, "InitiateLeaveWorkspace")
	qNameCmdChangePassword                          = appdef.NewQName(appdef.SysPackage, "ChangePassword")
	qNameCmdInitiateInvitationByEmail               = appdef.NewQName(appdef.SysPackage, "InitiateInvitationByEMail")
	qNameQryCollection                              = appdef.NewQName(appdef.SysPackage, "Collection")
	qNameCmdInitiateUpdateInviteRoles               = appdef.NewQName(appdef.SysPackage, "InitiateUpdateInviteRoles")
	qNameCmdInitiateCancelAcceptedInvite            = appdef.NewQName(appdef.SysPackage, "InitiateCancelAcceptedInvite")
	qNameCmdCancelSendInvite                        = appdef.NewQName(appdef.SysPackage, "CancelSendInvite")
	qNameCmdInitChildWorkspace                      = appdef.NewQName(appdef.SysPackage, "InitChildWorkspace")
	qNameCmdEnrichPrincipalToken                    = appdef.NewQName(appdef.SysPackage, "EnrichPrincipalToken")
	qNameCDocUPProfile                              = appdef.NewQName(airPackage, "UPProfile")
	qNameCDocResellerSubscriptionsProfile           = appdef.NewQName(airPackage, "ResellerSubscriptionsProfile")
	qNameCmdCreateUPProfile                         = appdef.NewQName(airPackage, "CreateUPProfile")
	qNameQryGetUPOnboardingPage                     = appdef.NewQName(airPackage, "GetUPOnboardingPage")
	qNameQryGetUPVerificationStatus                 = appdef.NewQName(airPackage, "GetUPVerificationStatus")
	qNameQryGetUPAccountStatus                      = appdef.NewQName(airPackage, "GetUPAccountStatus")
	qNameQryGetUPEventHistory                       = appdef.NewQName(airPackage, "GetUPEventHistory")
	qNameCmdStoreResellerSubscriptionsProfile       = appdef.NewQName(airPackage, "StoreResellerSubscriptionsProfile")
	qNameQryGetHostedAirSubscriptions               = appdef.NewQName(airPackage, "GetHostedAirSubscriptions")
	qNameQryGetUPStatus                             = appdef.NewQName(airPackage, "GetUPStatus")
	qNameQryQueryResellerInfo                       = appdef.NewQName(airPackage, "QueryResellerInfo")
	qNameCmdCreateUntillPayment                     = appdef.NewQName(airPackage, "CreateUntillPayment")
	qNameCmdRegenerateUPProfileApiToken             = appdef.NewQName(airPackage, "RegenerateUPProfileApiToken")
	qNameCmdEnsureUPPredefinedPaymentModesExist     = appdef.NewQName(airPackage, "EnsureUPPredefinedPaymentModesExist")
	qNameQryGetUPTerminals                          = appdef.NewQName(airPackage, "GetUPTerminals")
	qNameQryActivateUPTerminal                      = appdef.NewQName(airPackage, "ActivateUPTerminal")
	qNameQryGetUPPaymentMethods                     = appdef.NewQName(airPackage, "GetUPPaymentMethods")
	qNameQryToggleUPPaymentMethod                   = appdef.NewQName(airPackage, "ToggleUPPaymentMethod")
	qNameQryRequestUPPaymentMethod                  = appdef.NewQName(airPackage, "RequestUPPaymentMethod")
	qNameQryUPTerminalWebhook                       = appdef.NewQName(airPackage, "UPTerminalWebhook")

	// Air roles
	qNameRoleResellersAdmin         = appdef.NewQName(airPackage, "ResellersAdmin")
	qNameRoleUntillPaymentsReseller = appdef.NewQName(airPackage, "UntillPaymentsReseller")
	qNameRoleUntillPaymentsUser     = appdef.NewQName(airPackage, "UntillPaymentsUser")
	qNameRoleAirReseller            = appdef.NewQName(airPackage, "AirReseller")
	qNameRoleUntillPaymentsTerminal = appdef.NewQName(airPackage, "UntillPaymentsTerminal")
)

const (
	field_DeviceProfileWSID     = "DeviceProfileWSID"
	field_ComputersID           = "ComputersID"
	field_RestaurantComputersID = "RestaurantComputersID"
	field_dummy                 = "dummy"
	field_OwnerWSID             = "OwnerWSID"
	airPackage                  = "air"
	untillPackage               = "untill"
	untillChargebeeAgentLogin   = "untillchargebeeagent"
)

const (
	ACPolicy_Deny ACPolicyType = iota
	ACPolicy_Allow
)

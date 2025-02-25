package dto

type CreateAccountResponseDocs = ResponseSuccess[CreateAccountResponseWithOTP]
type LoginResponseDocs = ResponseSuccess[LoginResponseWithOTP]
type LogoutResponseDocs = ResponseSuccess[LogoutResponse]
type RefreshTokenResponseDocs = ResponseSuccess[RefreshTokenResponse]
type ShareVideoResponseDocs = ResponseSuccess[ShareVideoResponse]
type ListVideosResponseDocs = ResponseSuccessPagingation[[]VideoResponse]
type CheckTokenResponseDocs = ResponseSuccess[CheckTokenResponse]

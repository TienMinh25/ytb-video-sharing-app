package dto

type CreateAccountResponseDocs = ResponseSuccess[CreateAccountResponse]
type LoginResponseDocs = ResponseSuccess[LoginResponse]
type LogoutResponseDocs = ResponseSuccess[LogoutResponse]
type RefreshTokenResponseDocs = ResponseSuccess[RefreshTokenResponse]
type ShareVideoResponseDocs = ResponseSuccess[ShareVideoResponse]
type ListVideosResponseDocs = ResponseSuccessPagingation[[]VideoResponse]
type CheckTokenResponseDocs = ResponseSuccess[CheckTokenResponse]

package response

import (
	errConstant "cleaning/internal/constants/error"
	"errors"
)

func messageJa(err error) string {
	switch {
	case errors.Is(err, errConstant.ErrInvalidEmailOrPassword):
		return "メールアドレスまたはパスワードが正しくありません"
	case errors.Is(err, errConstant.ErrInvalidRegistrationToken):
		return "登録用トークンが無効、または有効期限が切れています"
	case errors.Is(err, errConstant.ErrInvalidPermission):
		return "指定された権限に無効なものが含まれています"
	case errors.Is(err, errConstant.ErrAccountIsDeactivated):
		return "このアカウントは無効化されています。管理者にお問い合わせください"
	case errors.Is(err, errConstant.ErrAlreadyDeactivated):
		return "既に無効化されています"
	case errors.Is(err, errConstant.ErrAlreadyActivated):
		return "このアカウントは既に有効化されています"
	case errors.Is(err, errConstant.ErrEmployeeAlreadyAssigned):
		return "この従業員は既にユーザーに登録されています"
	case errors.Is(err, errConstant.ErrCannotUpdateSelf):
		return "自分自身のユーザー情報は変更できません"
	case errors.Is(err, errConstant.ErrPasswordIncorrect):
		return "現在のパスワードが正しくありません"
	case errors.Is(err, errConstant.ErrPasswordIsSame):
		return "新しいパスワードは現在のパスワードと異なるものを設定してください"
	case errors.Is(err, errConstant.ErrEmployeeInActive):
		return "従業員は無効化されているため、操作を実行できません"
	case errors.Is(err, errConstant.ErrUnauthorized):
		return "認証が必要です"
	case errors.Is(err, errConstant.ErrForbidden):
		return "この操作を実行する権限がありません"
	case errors.Is(err, errConstant.ErrNotFound):
		return "該当するデータがありません"
	case errors.Is(err, errConstant.ErrBadRequest):
		return "入力内容に誤りがあります"
	case errors.Is(err, errConstant.ErrSystemRoleImmutable):
		return "システムロールの権限は変更できません"
	case errors.Is(err, errConstant.ErrRoleInUse):
		return "このロールはユーザーに割り当てられているため、無効化できません"
	default:
		return "サーバー内部でエラーが発生しました"
	}
}

package response

import (
	errConstant "cleaning/internal/constants/error"
	"errors"
)

func messageJa(err error) string {
	switch {
	case errors.Is(err, errConstant.ErrInvalidEmailOrPassword):
		return "メールアドレスまたはパスワードが正しくありません"
	case errors.Is(err, errConstant.ErrAccountIsDeactivated):
		return "このアカウントは無効化されています。管理者にお問い合わせください"
	case errors.Is(err, errConstant.ErrAlreadyDeactivated):
		return "このアカウントは既に無効化されています"
	case errors.Is(err, errConstant.ErrAlreadyActivated):
		return "このアカウントは既に有効化されています"
	case errors.Is(err, errConstant.ErrEmployeeAlreadyAssigned):
		return "この従業員は既にユーザーに登録されています"
	case errors.Is(err, errConstant.ErrCannotUpdateSelf):
		return "自分自身のユーザー情報は変更できません"
	case errors.Is(err, errConstant.ErrUnauthorized):
		return "認証が必要です"
	case errors.Is(err, errConstant.ErrForbidden):
		return "この操作を実行する権限がありません"
	case errors.Is(err, errConstant.ErrNotFound):
		return "該当するデータがありません"
	case errors.Is(err, errConstant.ErrBadRequest):
		return "入力内容に誤りがあります"
	default:
		return "サーバー内部でエラーが発生しました"
	}
}

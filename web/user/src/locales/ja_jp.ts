import { I18n } from '../types/locales'

const lang: I18n = {
  layout: {
    header: {
      becomeShopOwner: '出品する',
      cartEmptyMessage: 'カートは空です',
      cartNotEmptyMessage: 'レジに進む',
      signUp: '登録する',
      signIn: 'ログイン',
      changeLocaleText: '英語に切り替える',
    },
  },
  auth: {
    signIn: {
      email: 'メールアドレス',
      password: 'パスワード',
      forgetPasswordLink: 'パスワードを忘れた場合',
      signIn: 'ログイン',
      dontHaveAccount: 'アカウントをお持ちではないですか？',
      signUpLink: '登録する',
    },
    signUp: {
      tel: '電話番号',
      email: 'メールアドレス',
      password: 'パスワード',
      passwordConfirm: 'パスワード（確認用）',
      signUp: 'サインアップ',
      alreadyHas: 'すでにアカウントをお持ちですか？',
    },
    verify: {
      message: '認証コードを入力してください',
      btnText: '認証',
    },
  },
}

export default lang

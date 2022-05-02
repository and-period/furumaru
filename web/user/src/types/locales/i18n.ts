export interface I18n {
  default: {
    header: {
      becomeShopOwner: string
      cartEmptyMessage: string
      cartNotEmptyMessage: string
      signIn: string
      localeText: string
    }
  }
  auth: {
    signIn: {
      email: string
      password: string
      forgetPasswordLink: string
      signIn: string
      dontHaveAccount: string
      signUpLink: string
    }
    signUp: {
      tel: string
      email: string
      password: string
      passwordConfirm: string
      signUp: string
      alreadyHas: string
    }
    verify: {
      message: string
      btnText: string
    }
  }
}

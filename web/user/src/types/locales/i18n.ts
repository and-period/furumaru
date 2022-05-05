export interface I18n {
  layout: {
    header: {
      becomeShopOwner: string
      cartEmptyMessage: string
      cartNotEmptyMessage: string
      signUp: string
      signIn: string
      changeLocaleText: string
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

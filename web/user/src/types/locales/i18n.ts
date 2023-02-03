export interface I18n {
  layout: {
    header: {
      becomeShopOwner: string;
      cartEmptyMessage: string;
      cartNotEmptyMessage: string;
      signUp: string;
      signIn: string;
      changeLocaleText: string;

      topLinkText: string
      searchItemLinkText: string
      allItemLinkText: string
      aboutLinkText: string
    };

    footer: {
      qaLinkText: string
      privacyPolicyLinkText: string
      lawLinkText: string
      inquiryLinkText: string
    }
  };

  auth: {
    signIn: {
      pageName: string;
      email: string;
      password: string;
      forgetPasswordLink: string;
      signIn: string;
      googleButtonText: string;
      facebookButtonText: string;
      lineButtonText: string;
      dontHaveAccount: string;
      signUpLink: string;
    };

    signUp: {
      pageName: string;
      tel: string;
      email: string;
      password: string;
      passwordConfirm: string;
      signUp: string;
      alreadyHas: string;
    };

    verify: {
      pageName: string;
      message: string;
      btnText: string;
    };
  };
}

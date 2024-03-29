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

      cartMenuMessage: string

      notificationTitle: string
      noNotificationItemText: string

      myPageLinkText: string
      viewMyCartText: string
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
      authErrorMessage: string
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

    register: {
      pageName: string
      name: string
      username: string
      btnText: string
      cautionText: string
      privacyPolicyLinkText: string
      termsOfServiceLink: string
    }
  };

  base: {
    about: {
      leadSentence: string
      description: string
      firstPointTitle: string
      firstPointDescription: string
      firstPointLinkText: string
      secondPointTitle: string
      secondPointDescription: string
      thirdPointTitle: string
      thirdPointDescription: string
      forthPointTitle: string
      forthPointDescription: string
      forthPointLinkText: string
    }
  };

  purchase: {
    complete: {
      thanksMessage: string,
      completeMessage: string
    }
  };
}

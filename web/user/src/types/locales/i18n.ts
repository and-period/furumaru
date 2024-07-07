export interface I18n {
  layout: {
    header: {
      becomeShopOwner: string
      cartEmptyMessage: string
      cartNotEmptyMessage: string
      signUp: string
      signIn: string
      changeLocaleText: string

      topLinkText: string
      searchItemLinkText: string
      allItemLinkText: string
      allMarcheLinkText: string
      aboutLinkText: string

      cartMenuMessage: string
      cartTotalPriceText: string
      cartTotalPriceTaxIncludedText: string

      notificationTitle: string
      noNotificationItemText: string

      myPageLinkText: string
      viewMyCartText: string

      numberOfCartsText: string
      shippingFeeAnnotation: string
      shippingFeeAnnotationLinkText: string
      shippingFeeAnnotationCheckText: string
    }

    footer: {
      qaLinkText: string
      privacyPolicyLinkText: string
      lawLinkText: string
      inquiryLinkText: string
    }
  }

  auth: {
    signIn: {
      pageName: string
      authErrorMessage: string
      email: string
      password: string
      forgetPasswordLink: string
      signIn: string
      googleButtonText: string
      facebookButtonText: string
      lineButtonText: string
      dontHaveAccount: string
      signUpLink: string
    }

    signUp: {
      pageName: string
      tel: string
      email: string
      password: string
      passwordConfirm: string
      signUp: string
      alreadyHas: string
    }

    verify: {
      pageName: string
      message: string
      btnText: string
    }

    register: {
      pageName: string
      name: string
      username: string
      btnText: string
      cautionText: string
      privacyPolicyLinkText: string
      termsOfServiceLink: string
    }
  }

  base: {
    top: {
      lineAddFriendImageUrl: string
      lineAddFriendImageAlt: string
      lineCouponText: string
      marcheListSubTitle: string
      liveStreamingText: string
      liveUpcomingText: string
      noMarcheItemFirstText: string
      noMarcheItemSecondText: string
      pastMarcheLinkText: string
      viewMoreText: string
      productsLinkText: string
      archiveListSubTitle: string
      archivedStreamText: string
      archivesLinkText: string
    }

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
  }

  purchase: {
    complete: {
      thanksMessage: string
      completeMessage: string
    }
  }

  items: {
    list: {
      allItemsTitle: string
      forSaleText: string
      soldOutText: string
      outOfSalesText: string
      presalesText: string
      unknownItemText: string
      itemThumbnailAlt: string
      itemPriceTaxIncludedText: string
      quantityLabel: string
      addToCartText: string
      coordinatorLabel: string
      coordinatorThumbnailAlt: string
      addCartSnackbarMessage: string
    },
    details: {
      itemThumbnailAlt: string
      producerLabel: string
      highlightsLabel: string
      itemPriceTaxIncludedText: string
      quantityLabel: string
      addToCartText: string
      addCartSnackbarMessage: string
      expirationDateLabel:string
      expirationDateText: string
      weightLabel:string
      deliveryTypeLabel:string
      deliveryTypeStandard: string
      deliveryTypeRefrigerated: string
      deliveryTypeFrozen: string
      storageTypeLabel: string
      storageTypeUnknown: string
      storageTypeRoomTemperature: string
      storageTypeCoolAndDark: string
      storageTypeRefrigerated: string
      storageTypeFrozen: string
      producerInformationTitle: string
    }
  }
}

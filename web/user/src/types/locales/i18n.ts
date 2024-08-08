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
      blogLinkText: string
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
    cart: {
      cartTitle: string
      cartCountLabel: string
      firstNotice: string
      secondNotice: string
      shipFromLabel: string
      coordinatorLabel: string
      totalPriceLabel: string
      shippingFeeNotice: string
      checkoutButtonText: string
      productNameLabel: string
      productPriceLabel: string
      quantityLabel: string
      subtotalLabel: string
      deleteButtonText: string
      boxTypeLabel: string
      boxSizeLabel: string
      utilizationRateLabel: string
    }
    auth: {
      loginRequiredMessage: string
      loginNewAccountMessage: string
      withAccountTitle: string
      loginAndCheckoutButtonText: string
      usernameLabel: string
      passwordLabel: string
      usernamePlaceholder: string
      passwordPlaceholder: string
      noAccountButtonText: string
      forgetPasswordLink: string
      notSignUpTitle: string
      checkoutWithoutAccountDescription: string
      checkoutWithoutAccountButtonText: string
    }
    guest: {
      checkoutTitle: string
      customerInformationTitle: string
      nameErrorMessage: string
      nameKanaErrorMessage: string
      phoneErrorMessage: string
      postalCodeErrorMessage: string
      cityErrorMessage: string
      addressErrorMessage: string
      emailErrorMessage: string
      emailInvalidErrorMessage: string
      unknownErrorMessage: string
      firstNamePlaceholder: string
      lastNamePlaceholder: string
      firstNameKanaPlaceholder: string
      lastNameKanaPlaceholder: string
      phoneNumberLabel: string
      postalCodeLabel: string
      searchButtonText: string
      prefectureLabel: string
      cityPlaceholder: string
      streetPlaceholder: string
      apartmentPlaceholder: string
      orderDetailsTitle: string
      shipFromLabel: string
      coordinatorLabel: string
      boxCountLabel: string
      quantityLabel: string
      couponPlaceholder: string
      applyButtonText: string
      couponAppliedMessage: string
      couponInvalidMessage: string
      itemTotalPriceLabel: string
      applyCouponLabel: string
      shippingFeeLabel: string
      calculateNextPageMessage: string
      totalPriceLabel: string
      backToCartButtonText: string
      paymentMethodButtonText: string
    }
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
    }
    details: {
      itemThumbnailAlt: string
      producerLabel: string
      highlightsLabel: string
      itemPriceTaxIncludedText: string
      quantityLabel: string
      addToCartText: string
      addCartSnackbarMessage: string
      expirationDateLabel: string
      expirationDateText: string
      weightLabel: string
      deliveryTypeLabel: string
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

  lives: {
    list: {
      allMarcheTitle: string
    }
    details: {
      archivedStreamText: string
      coordinatorThumbnailAlt: string
      coordinatorLabel: string
      hideMarcheDetailsText: string
      showMarcheDetailsText: string
      itemsTabLabel: string
      commentsTabLabel: string
      itemPriceTaxIncludedText: string
      addToCartText: string
      addCartSnackbarMessage: string
      commentPlaceholder: string
      submitButtonText: string
      guestCommentNote: string
      noCommentsText: string
      guestNameLabel: string
    }
  }
}

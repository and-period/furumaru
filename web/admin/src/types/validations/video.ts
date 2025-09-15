import type { ValidationArgs } from '@vuelidate/core'
import { maxLength } from '~/lib/validations'

export const CreateVideoValidationRules: ValidationArgs = {
  title: { required: true, maxLength: maxLength(64) },
  description: { required: true, maxLength: maxLength(2000) },
  coordinatorId: {},
  productIds: {},
  experienceIds: {},
  thumbnailUrl: { required: true },
  videoUrl: { required: true },
  _public: {},
  limited: {},
  publishedAt: {},
  displayProduct: {},
  displayExperience: {},
}

export const UpdateVideoValidationRules: ValidationArgs = {
  title: { required: true, maxLength: maxLength(64) },
  description: { required: true, maxLength: maxLength(2000) },
  coordinatorId: {},
  productIds: {},
  experienceIds: {},
  thumbnailUrl: { required: true },
  videoUrl: { required: true },
  _public: {},
  limited: {},
  publishedAt: {},
  displayProduct: {},
  displayExperience: {},
}

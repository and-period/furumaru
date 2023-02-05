<template>
  <div>
    <v-card-title>コーディネーター編集</v-card-title>

    <v-tabs v-model="tab" grow color="dark">
      <v-tabs-slider color="accent"></v-tabs-slider>
      <v-tab
        v-for="tabItem in tabItems"
        :key="tabItem.value"
        :href="`#${tabItem.value}`"
      >
        {{ tabItem.name }}
      </v-tab>
    </v-tabs>

    <v-tabs-items v-model="tab">
      <v-tab-item value="coordinators">
        <v-skeleton-loader v-if="fetchState.pending" type="article" />

        <the-coordinator-edit-form
          v-else
          :form-data="formData"
          :thumbnail-upload-status="thumbnailUploadStatus"
          :header-upload-status="headerUploadStatus"
          :search-loading="searchLoading"
          :search-error-message="searchErrorMessage"
          @update:thumbnailFile="handleUpdateThumbnail"
          @update:headerFile="handleUpdateHeader"
          @submit="handleSubmit"
          @click:search="searchAddress"
        />
      </v-tab-item>

      <v-tab-item value="relationProducers">
        <v-dialog v-model="dialog" width="500">
          <template #activator="{ on, attrs }">
            <div class="d-flex pt-3 pr-3">
              <v-spacer />
              <v-btn outlined color="primary" v-bind="attrs" v-on="on">
                <v-icon left>mdi-plus</v-icon>
                生産者紐付け
              </v-btn>
            </div>
          </template>

          <v-card>
            <v-card-title class="primaryLight"> 生産者紐付け </v-card-title>

            <v-autocomplete
              v-model="producers"
              chips
              label="関連生産者"
              multiple
              filled
              :items="producerItems"
              item-text="firstname"
              item-value="id"
            >
              <template #selection="data">
                <v-chip close @click:close="remove(data.item.id)">
                  <v-avatar left>
                    <v-img :src="data.item.thumbnailUrl"></v-img>
                  </v-avatar>
                  {{ data.item.firstname }}
                </v-chip>
              </template>
              <template #item="data">
                <v-list-item-avatar>
                  <img :src="data.item.thumbnailUrl" />
                </v-list-item-avatar>
                <v-list-item-content>
                  <v-list-item-title>{{
                    data.item.firstname
                  }}</v-list-item-title>
                  <v-list-item-subtitle>{{
                    data.item.storeName
                  }}</v-list-item-subtitle>
                </v-list-item-content>
              </template>
            </v-autocomplete>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="error" text @click="cancel"> キャンセル </v-btn>
              <v-btn color="primary" outlined @click="relateProducers">
                登録
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>

        <the-related-producer-list
          :table-footer-props="producersOptions"
          @update:items-per-page="handleUpdateProducersItemsPerPage"
          @update:page="handleUpdateProducersPage"
        ></the-related-producer-list>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  reactive,
  ref,
  useFetch,
  useRoute,
  useRouter,
  watch,
} from '@nuxtjs/composition-api'
import { useVuelidate } from '@vuelidate/core'

import { usePagination, useSearchAddress } from '~/lib/hooks'
import {
  kana,
  getErrorMessage,
  required,
  tel,
  maxLength,
} from '~/lib/validations'
import { useCoordinatorStore } from '~/store/coordinator'
import { useProducerStore } from '~/store/producer'
import {
  ProducersResponseProducersInner,
  RelateProducersRequest,
  UpdateCoordinatorRequest,
} from '~/types/api'
import { ImageUploadStatus } from '~/types/props'
import { Coordinator } from '~/types/props/coordinator'

export default defineComponent({
  setup() {
    const tab = ref<string>('coordinators')
    const tabItems: Coordinator[] = [
      { name: '基本情報', value: 'coordinators' },
      { name: '関連生産者', value: 'relationProducers' },
    ]
    const coordinatorStore = useCoordinatorStore()

    const producers = ref<string[]>([])
    const dialog = ref<boolean>(false)

    const producerStore = useProducerStore()

    const producerItems = computed(() => {
      return producerStore.producers
    })

    const relateProducersItems = reactive<{
      offset: number
      relateProducers: ProducersResponseProducersInner[]
    }>({ offset: 0, relateProducers: [] })

    const route = useRoute()
    const id = route.value.params.id
    const router = useRouter()

    const { uploadCoordinatorThumbnail, uploadCoordinatorHeader } =
      useCoordinatorStore()

    const { getCoordinator } = useCoordinatorStore()

    const {
      itemsPerPage: producersItemsPerPage,
      offset: producersOffset,
      options: producersOptions,
      handleUpdateItemsPerPage: handleUpdateProducersItemsPerPage,
      updateCurrentPage: _handleUpdateProducersPage,
    } = usePagination()

    watch(producersItemsPerPage, () => {
      coordinatorStore.fetchRelatedProducers(id, producersItemsPerPage.value, 0)
    })

    const handleUpdateProducersPage = async (page: number) => {
      _handleUpdateProducersPage(page)
      await coordinatorStore.fetchRelatedProducers(
        id,
        producersItemsPerPage.value,
        producersOffset.value
      )
    }

    const formData = reactive<UpdateCoordinatorRequest>({
      storeName: '',
      firstname: '',
      lastname: '',
      firstnameKana: '',
      lastnameKana: '',
      companyName: '',
      thumbnailUrl: '',
      headerUrl: '',
      twitterAccount: '',
      instagramAccount: '',
      facebookAccount: '',
      phoneNumber: '',
      postalCode: '',
      prefecture: '',
      city: '',
      addressLine1: '',
      addressLine2: '',
    })

    const producerData = reactive<RelateProducersRequest>({
      producerIds: [],
    })

    const { fetchState } = useFetch(async () => {
      try {
        const coordinator = await getCoordinator(id)
        formData.storeName = coordinator.storeName
        formData.firstname = coordinator.firstname
        formData.lastname = coordinator.lastname
        formData.firstnameKana = coordinator.firstnameKana
        formData.lastnameKana = coordinator.lastnameKana
        formData.companyName = coordinator.companyName
        formData.thumbnailUrl = coordinator.thumbnailUrl
        formData.headerUrl = coordinator.headerUrl
        formData.twitterAccount = coordinator.twitterAccount
        formData.instagramAccount = coordinator.instagramAccount
        formData.facebookAccount = coordinator.facebookAccount
        formData.phoneNumber = coordinator.phoneNumber.replace('+81', '0')
        formData.postalCode = coordinator.postalCode
        formData.prefecture = coordinator.prefecture
        formData.city = coordinator.city
        formData.addressLine1 = coordinator.addressLine1
        formData.addressLine2 = coordinator.addressLine2

        await Promise.all([
          coordinatorStore.fetchRelatedProducers(
            id,
            producersItemsPerPage.value
          ),
        ])
        relateProducersItems.relateProducers = coordinatorStore.producers
      } catch (err) {
        console.log(err)
      }
    })

    const rules = computed(() => ({
      storeName: { required, maxLength: maxLength(64) },
      companyName: { required, maxLength: maxLength(64) },
      firstname: { required, maxLength: maxLength(16) },
      lastname: { required, maxLength: maxLength(16) },
      firstnameKana: { required, kana },
      lastnameKana: { required, kana },
      phoneNumber: { required, tel },
    }))

    const v$ = useVuelidate(rules, formData)

    const {
      loading: searchLoading,
      errorMessage: searchErrorMessage,
      searchAddressByPostalCode,
    } = useSearchAddress()

    const searchAddress = async () => {
      searchLoading.value = true
      searchErrorMessage.value = ''
      const res = await searchAddressByPostalCode(Number(formData.postalCode))
      if (res) {
        formData.prefecture = res.prefecture
        formData.city = res.city
        formData.addressLine1 = res.addressLine1
      }
    }

    const thumbnailUploadStatus = reactive<ImageUploadStatus>({
      error: false,
      message: '',
    })

    const headerUploadStatus = reactive<ImageUploadStatus>({
      error: false,
      message: '',
    })

    const handleUpdateThumbnail = (files: FileList) => {
      if (files.length > 0) {
        uploadCoordinatorThumbnail(files[0])
          .then((res) => {
            formData.thumbnailUrl = res.url
          })
          .catch(() => {
            thumbnailUploadStatus.error = true
            thumbnailUploadStatus.message = 'アップロードに失敗しました。'
          })
      }
    }

    const handleUpdateHeader = async (files: FileList) => {
      if (files.length > 0) {
        await uploadCoordinatorHeader(files[0])
          .then((res) => {
            formData.headerUrl = res.url
          })
          .catch(() => {
            headerUploadStatus.error = true
            headerUploadStatus.message = 'アップロードに失敗しました。'
          })
      }
    }

    const handleSubmit = async (): Promise<void> => {
      try {
        const result = await v$.value.$validate()
        if (!result) {
          return
        }
        await coordinatorStore.updateCoordinator(
          {
            ...formData,
            phoneNumber: formData.phoneNumber.replace('0', '+81'),
          },
          id
        )
        router.push('/coordinators')
      } catch (error) {
        console.log(error)
      }
    }

    const relateProducers = async (): Promise<void> => {
      producerData.producerIds = producers.value
      try {
        await coordinatorStore.relateProducers(id, producerData)
        dialog.value = false
      } catch (error) {
        console.log(error)
      }
    }

    const remove = (item: string) => {
      producers.value = producers.value.filter((id) => id !== item)
    }

    const cancel = (): void => {
      dialog.value = false
    }

    useFetch(async () => {
      try {
        await producerStore.fetchProducers(20, 0, 'unrelated')
      } catch (err) {
        console.log(err)
      }
    })

    return {
      id,
      fetchState,
      formData,
      producers,
      v$,
      searchLoading,
      searchErrorMessage,
      thumbnailUploadStatus,
      headerUploadStatus,
      tabItems,
      tab,
      dialog,
      producersOptions,
      producerItems,
      getErrorMessage,
      searchAddress,
      handleSubmit,
      handleUpdateThumbnail,
      handleUpdateHeader,
      remove,
      relateProducers,
      cancel,
      handleUpdateProducersItemsPerPage,
      handleUpdateProducersPage,
    }
  },
})
</script>

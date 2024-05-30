<script lang="ts" setup>
import { mdiPaperclip, mdiContentCopy } from '@mdi/js'
import Hls from 'hls.js'
import {
  type Broadcast,
  BroadcastStatus,
  type AuthYoutubeBroadcastRequest,
} from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  pauseDialog: {
    type: Boolean,
    default: false,
  },
  liveMp4Dialog: {
    type: Boolean,
    default: false,
  },
  archiveMp4Dialog: {
    type: Boolean,
    default: false,
  },
  mp4FormData: {
    type: Object as PropType<File[] | undefined>,
    default: (): File[] | undefined => undefined,
  },
  authYoutubeFormData: {
    type: Object as PropType<AuthYoutubeBroadcastRequest>,
    default: (): AuthYoutubeBroadcastRequest => ({
      youtubeHandle: '',
    }),
  },
  selectedTabItem: {
    type: String,
    default: 'schedule',
  },
  broadcast: {
    type: Object as PropType<Broadcast>,
    default: (): Broadcast => ({
      id: '',
      scheduleId: '',
      status: BroadcastStatus.UNKNOWN,
      inputUrl: '',
      outputUrl: '',
      archiveUrl: '',
      youtubeAccount: '',
      youtubeViewerUrl: '',
      youtubeAdminUrl: '',
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  authYoutubeUrl: {
    type: String,
    default: '',
  },
})

const emit = defineEmits<{
  (e: 'update:broadcast', broadcast: Broadcast): void
  (e: 'update:mp4-form-data', file: File[] | undefined): void
  (
    e: 'update:auth-youtube-form-data',
    formData: AuthYoutubeBroadcastRequest,
  ): void
  (e: 'update:pause-dialog', toggle: boolean): void
  (e: 'update:live-mp4-dialog', toggle: boolean): void
  (e: 'update:archive-mp4-dialog', toggle: boolean): void
  (e: 'click:link-youtube'): void
  (e: 'click:activate-static-image'): void
  (e: 'click:deactivate-static-image'): void
  (e: 'submit:pause'): void
  (e: 'submit:unpause'): void
  (e: 'submit:change-input-mp4'): void
  (e: 'submit:change-input-rtmp'): void
  (e: 'submit:upload-archive-mp4'): void
}>()

const statuses = [
  { title: 'リソース未作成', value: BroadcastStatus.DISABLED },
  { title: 'リソース作成中', value: BroadcastStatus.WAITING },
  { title: '配信停止中', value: BroadcastStatus.IDLE },
  { title: '配信中', value: BroadcastStatus.ACTIVE },
  { title: '不明', value: BroadcastStatus.UNKNOWN },
]

const videoRef = ref<HTMLVideoElement>()

const broadcastValue = computed({
  get: (): Broadcast => props.broadcast,
  set: (broadcast: Broadcast): void => emit('update:broadcast', broadcast),
})
const pauseDialogValue = computed({
  get: (): boolean => props.pauseDialog,
  set: (v: boolean): void => emit('update:pause-dialog', v),
})
const archiveMp4DialogValue = computed({
  get: (): boolean => props.archiveMp4Dialog,
  set: (v: boolean): void => emit('update:archive-mp4-dialog', v),
})
const liveMp4DialogValue = computed({
  get: (): boolean => props.liveMp4Dialog,
  set: (v: boolean): void => emit('update:live-mp4-dialog', v),
})
const mp4FormDataValue = computed({
  get: (): File[] | undefined => props.mp4FormData,
  set: (formData: File[] | undefined): void =>
    emit('update:mp4-form-data', formData),
})
const authYoutubeFormDataValue = computed({
  get: (): AuthYoutubeBroadcastRequest => props.authYoutubeFormData,
  set: (formData: AuthYoutubeBroadcastRequest): void =>
    emit('update:auth-youtube-form-data', formData),
})
const authYoutubeUrlValue = computed({
  get: (): string => props.authYoutubeUrl,
  set: (url: string): void => console.log(url),
})

watch(
  (): string => props.selectedTabItem,
  (): void => {
    if (props.selectedTabItem !== 'streaming') {
      return
    }
    onClickVideo()
  },
)

const isLive = (): boolean => {
  return props.broadcast?.status === BroadcastStatus.ACTIVE
}

const isVOD = (): boolean => {
  return props.broadcast?.status === BroadcastStatus.DISABLED
}

const onClickVideo = (): void => {
  if (!videoRef.value || !props.broadcast) {
    return
  }
  if (!isLive()) {
    return
  }

  const video = videoRef.value
  const src = props.broadcast.outputUrl

  if (Hls.isSupported()) {
    const hls = new Hls({ enableWorker: false })
    hls.loadSource(src)
    hls.attachMedia(video)
    video.play()
  }
  else if (video.canPlayType('application/vnd.apple.mpegurl')) {
    video.src = src
    video.play()
  }
}

const onClickLinkYouTube = (): void => {
  emit('click:link-youtube')
}

const onClickPause = (): void => {
  emit('update:pause-dialog', true)
}

const onClickActivateStaticImage = (): void => {
  emit('click:activate-static-image')
}

const onClickDeactivateStaticImage = (): void => {
  emit('click:deactivate-static-image')
}

const onClickChangeMp4Input = (): void => {
  mp4FormDataValue.value = undefined // 初期化
  emit('update:live-mp4-dialog', true)
}

const onClickUploadArchiveMp4 = (): void => {
  mp4FormDataValue.value = undefined // 初期化
  emit('update:archive-mp4-dialog', true)
}

const onClosePauseDialog = (): void => {
  emit('update:pause-dialog', false)
}

const onCloseArchiveMp4Dialog = (): void => {
  emit('update:archive-mp4-dialog', false)
}

const onCloseLiveMp4Dialog = (): void => {
  emit('update:live-mp4-dialog', false)
}

const onSubmitPause = (): void => {
  emit('submit:pause')
}

const onSubmitUnpause = (): void => {
  emit('submit:unpause')
}

const onSubmitChangeMp4Input = (): void => {
  emit('submit:change-input-mp4')
}

const onSubmitChangeRtmpInput = (): void => {
  emit('submit:change-input-rtmp')
}

const onSubmitUploadArchiveMp4 = (): void => {
  emit('submit:upload-archive-mp4')
}

const onClickCopyUrl = (url: string) => {
  navigator.clipboard.writeText(url)
}
</script>

<template>
  <v-dialog
    v-model="pauseDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h7">
        本当に停止しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="info"
          variant="text"
          @click="onClosePauseDialog"
        >
          閉じる
        </v-btn>
        <v-btn
          :loading="loading"
          color="error"
          variant="outlined"
          @click="onSubmitPause"
        >
          停止
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="archiveMp4DialogValue">
    <v-card>
      <v-card-title class="primaryLight">
        アーカイブ映像アップロード
      </v-card-title>
      <v-card-text>
        <v-file-input
          v-model="mp4FormDataValue"
          counter
          label="アーカイブ動画"
          :prepend-icon="mdiPaperclip"
          outlined
          accept="video/mp4"
          :show-size="1000"
        >
          <template #selection="{ fileNames }">
            <template
              v-for="fileName in fileNames"
              :key="fileName"
            >
              <v-chip
                size="small"
                label
                color="primary"
                class="me-2"
              >
                {{ fileName }}
              </v-chip>
            </template>
          </template>
        </v-file-input>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onCloseArchiveMp4Dialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onSubmitUploadArchiveMp4"
        >
          送信
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="liveMp4DialogValue">
    <v-card>
      <v-card-title class="primaryLight">
        ライブ映像切り替え
      </v-card-title>
      <v-card-text>
        <v-file-input
          v-model="mp4FormDataValue"
          counter
          label="ライブ動画"
          :prepend-icon="mdiPaperclip"
          outlined
          accept="video/mp4"
          :show-size="1000"
        >
          <template #selection="{ fileNames }">
            <template
              v-for="fileName in fileNames"
              :key="fileName"
            >
              <v-chip
                size="small"
                label
                color="primary"
                class="me-2"
              >
                {{ fileName }}
              </v-chip>
            </template>
          </template>
        </v-file-input>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onCloseLiveMp4Dialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onSubmitChangeMp4Input"
        >
          送信
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-row>
    <v-col
      sm="12"
      md="12"
      lg="8"
    >
      <v-card>
        <v-card-text>
          <v-container>
            <video
              v-show="isLive()"
              id="video"
              ref="videoRef"
              controls
            />
            <video
              v-show="isVOD()"
              id="video"
              controls
            >
              <source
                :src="broadcast.archiveUrl"
                type="video/mp4"
              >
              <a
                :href="broadcast.archiveUrl"
                type="video/mp4"
              >mp4</a>
            </video>
            <v-btn
              v-show="isLive()"
              @click="onClickVideo"
            >
              映像の更新
            </v-btn>
          </v-container>
        </v-card-text>
      </v-card>

      <v-card class="mt-4">
        <v-card-text>
          <div class="my-4">
            <div
              v-if="broadcastValue.youtubeAccount === ''"
              class="px-0"
            >
              <template v-if="authYoutubeUrlValue === ''">
                <v-text-field
                  v-model="authYoutubeFormDataValue.youtubeHandle"
                  variant="outlined"
                  label="YouTube 連携先ハンドル名"
                  placeholder="@から始まるハンドル名（@含めて入力してください）"
                />
                <v-btn
                  block
                  variant="outlined"
                  color="primary"
                  @click="onClickLinkYouTube"
                >
                  連携する
                </v-btn>
              </template>
              <template v-else>
                <v-text-field
                  v-model="authYoutubeUrlValue"
                  variant="outlined"
                  label="YouTube 連携用URL（配信者へ以下のURLを連携してください）"
                  readonly
                  :append-icon="mdiContentCopy"
                  @click:append="onClickCopyUrl(authYoutubeUrlValue)"
                />
              </template>
            </div>
            <div v-else>
              <v-text-field
                v-model="broadcastValue.youtubeAccount"
                variant="outlined"
                label="YouTube 連携先ハンドル名"
                readonly
              />
              <v-text-field
                v-model="broadcastValue.youtubeViewerUrl"
                variant="outlined"
                label="YouTube 配信視聴画面URL"
                readonly
                :append-icon="mdiContentCopy"
                @click:append="onClickCopyUrl(broadcast.youtubeViewerUrl)"
              />
              <v-text-field
                v-model="broadcastValue.youtubeAdminUrl"
                variant="outlined"
                label="YouTube 配信管理画面URL"
                readonly
                :append-icon="mdiContentCopy"
                @click:append="onClickCopyUrl(broadcast.youtubeAdminUrl)"
              />
            </div>
          </div>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col
      sm="12"
      md="12"
      lg="4"
    >
      <v-card>
        <v-card-text class="px-4">
          <v-select
            v-model="broadcastValue.status"
            label="配信状況"
            :items="statuses"
            item-title="title"
            item-value="value"
            variant="outlined"
            readonly
          />
          <v-text-field
            v-show="isLive()"
            v-model="broadcastValue.inputUrl"
            label="配信エンドポイント：入力側"
            readonly
            variant="outlined"
          />
          <v-text-field
            v-show="isLive()"
            v-model="broadcastValue.outputUrl"
            label="配信エンドポイント：出力側"
            readonly
            variant="outlined"
          />
          <v-text-field
            v-show="isVOD()"
            v-model="broadcastValue.archiveUrl"
            label="オンデマンド配信URL"
            readonly
            variant="outlined"
          />
        </v-card-text>
      </v-card>

      <v-card
        v-show="isLive()"
        class="mt-4"
      >
        <v-card-text>
          <v-list>
            <v-list-item class="px-0">
              <v-list-item-subtitle> ライブ配信の操作 </v-list-item-subtitle>
              <v-btn
                block
                variant="outlined"
                color="primary"
                class="mt-2"
                @click="onClickPause"
              >
                停止する
              </v-btn>
              <v-btn
                block
                variant="outlined"
                color="secondary"
                class="mt-2"
                @click="onSubmitUnpause"
              >
                停止を解除する
              </v-btn>
            </v-list-item>

            <v-list-item class="px-0 mt-4">
              <v-list-item-subtitle>
                入力チャンネルの設定
              </v-list-item-subtitle>
              <v-btn
                block
                variant="outlined"
                color="primary"
                class="mt-2"
                @click="onSubmitChangeRtmpInput"
              >
                RTMP配信に切り替え
              </v-btn>
              <v-btn
                block
                variant="outlined"
                color="secondary"
                class="mt-2"
                @click="onClickChangeMp4Input"
              >
                MP4配信に切り替え
              </v-btn>
            </v-list-item>

            <v-list-item class="px-0 mt-4">
              <v-list-item-subtitle> ふた絵の表示設定 </v-list-item-subtitle>
              <v-btn
                block
                variant="outlined"
                color="primary"
                class="mt-2"
                @click="onClickActivateStaticImage"
              >
                有効化
              </v-btn>
              <v-btn
                block
                variant="outlined"
                color="secondary"
                class="mt-2"
                @click="onClickDeactivateStaticImage"
              >
                無効化
              </v-btn>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>

      <v-card
        v-show="isVOD()"
        class="mt-4"
      >
        <v-card-text>
          <v-list>
            <v-list-item class="px-0">
              <v-list-item-subtitle>
                オンデマンド配信の設定
              </v-list-item-subtitle>
              <v-btn
                block
                variant="outlined"
                color="secondary"
                class="mt-2"
                @click="onClickUploadArchiveMp4"
              >
                アップロード
              </v-btn>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<style scoped>
#video {
  margin-bottom: 1.5em;
  width: 100%;
}
</style>

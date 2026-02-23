<script lang="ts" setup>
import {
  mdiPaperclip,
  mdiContentCopy,
  mdiPlay,
  mdiPause,
  mdiStop,
  mdiRefresh,
  mdiYoutube,
  mdiMonitor,
  mdiCog,
  mdiUpload,
  mdiClose,
  mdiCheck,
  mdiAlert,
  mdiVideo,
  mdiRadio,
} from '@mdi/js'
import {

  BroadcastStatus,

} from '~/types/api/v1'
import type { Broadcast, AuthYoutubeBroadcastRequest } from '~/types/api/v1'

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
      status: BroadcastStatus.BroadcastStatusUnknown,
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
  { title: 'リソース未作成', value: BroadcastStatus.BroadcastStatusDisabled },
  { title: 'リソース作成中', value: BroadcastStatus.BroadcastStatusWaiting },
  { title: '配信停止中', value: BroadcastStatus.BroadcastStatusIdle },
  { title: '配信中', value: BroadcastStatus.BroadcastStatusActive },
  { title: '不明', value: BroadcastStatus.BroadcastStatusUnknown },
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
  return props.broadcast?.status === BroadcastStatus.BroadcastStatusActive
}

const isVOD = (): boolean => {
  return props.broadcast?.status === BroadcastStatus.BroadcastStatusDisabled
}

const onClickVideo = async (): Promise<void> => {
  if (!videoRef.value || !props.broadcast) {
    return
  }
  if (!isLive()) {
    return
  }

  const video = videoRef.value
  const src = props.broadcast.outputUrl

  const { default: HlsLib } = await import('hls.js')
  if (HlsLib.isSupported()) {
    const hls = new HlsLib({ enableWorker: false })
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
  <!-- 停止確認ダイアログ -->
  <v-dialog
    v-model="pauseDialogValue"
    max-width="400"
  >
    <v-card class="dialog-card">
      <v-card-title class="d-flex align-center section-header pa-6">
        <v-icon
          :icon="mdiAlert"
          size="24"
          class="mr-3 text-warning"
        />
        <span class="text-h6 font-weight-medium">配信停止の確認</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <p class="text-body-1 mb-0">
          ライブ配信を停止しますか？
        </p>
        <p class="text-body-2 text-grey-darken-1 mt-2">
          停止後は再開に時間がかかる場合があります。
        </p>
      </v-card-text>
      <v-card-actions class="pa-6 pt-0">
        <v-spacer />
        <v-btn
          variant="text"
          @click="onClosePauseDialog"
        >
          <v-icon
            :icon="mdiClose"
            start
          />
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="error"
          variant="elevated"
          @click="onSubmitPause"
        >
          <v-icon
            :icon="mdiStop"
            start
          />
          停止
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- アーカイブアップロードダイアログ -->
  <v-dialog
    v-model="archiveMp4DialogValue"
    max-width="500"
  >
    <v-card class="dialog-card">
      <v-card-title class="d-flex align-center section-header pa-6">
        <v-icon
          :icon="mdiUpload"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">アーカイブ映像アップロード</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <v-file-input
          v-model="mp4FormDataValue"
          label="アーカイブ動画 *"
          :prepend-icon="mdiVideo"
          variant="outlined"
          density="comfortable"
          accept="video/mp4"
          :show-size="1000"
          counter
        >
          <template #selection="{ fileNames }">
            <template
              v-for="fileName in fileNames"
              :key="fileName"
            >
              <v-chip
                size="small"
                color="primary"
                variant="outlined"
                class="me-2"
              >
                {{ fileName }}
              </v-chip>
            </template>
          </template>
        </v-file-input>
      </v-card-text>
      <v-card-actions class="pa-6 pt-0">
        <v-spacer />
        <v-btn
          variant="text"
          @click="onCloseArchiveMp4Dialog"
        >
          <v-icon
            :icon="mdiClose"
            start
          />
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="elevated"
          @click="onSubmitUploadArchiveMp4"
        >
          <v-icon
            :icon="mdiUpload"
            start
          />
          アップロード
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- ライブ映像切り替えダイアログ -->
  <v-dialog
    v-model="liveMp4DialogValue"
    max-width="500"
  >
    <v-card class="dialog-card">
      <v-card-title class="d-flex align-center section-header pa-6">
        <v-icon
          :icon="mdiVideo"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">ライブ映像切り替え</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <v-file-input
          v-model="mp4FormDataValue"
          label="ライブ動画 *"
          :prepend-icon="mdiVideo"
          variant="outlined"
          density="comfortable"
          accept="video/mp4"
          :show-size="1000"
          counter
        >
          <template #selection="{ fileNames }">
            <template
              v-for="fileName in fileNames"
              :key="fileName"
            >
              <v-chip
                size="small"
                color="primary"
                variant="outlined"
                class="me-2"
              >
                {{ fileName }}
              </v-chip>
            </template>
          </template>
        </v-file-input>
      </v-card-text>
      <v-card-actions class="pa-6 pt-0">
        <v-spacer />
        <v-btn
          variant="text"
          @click="onCloseLiveMp4Dialog"
        >
          <v-icon
            :icon="mdiClose"
            start
          />
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="elevated"
          @click="onSubmitChangeMp4Input"
        >
          <v-icon
            :icon="mdiCheck"
            start
          />
          切り替え
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- メインコンテンツ -->
  <div class="streaming-container">
    <v-row>
      <v-col
        cols="12"
        lg="8"
      >
        <!-- 配信プレビューセクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiMonitor"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">配信プレビュー</span>
            <v-spacer />
            <v-btn
              v-show="isLive()"
              variant="outlined"
              size="small"
              @click="onClickVideo"
            >
              <v-icon
                :icon="mdiRefresh"
                start
              />
              更新
            </v-btn>
          </v-card-title>
          <v-card-text class="pa-6">
            <div class="video-container">
              <video
                v-show="isLive()"
                id="video"
                ref="videoRef"
                controls
                class="video-player"
              />
              <video
                v-show="isVOD()"
                id="video"
                controls
                class="video-player"
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
              <div
                v-if="!isLive() && !isVOD()"
                class="video-placeholder"
              >
                <v-icon
                  :icon="mdiMonitor"
                  size="64"
                  class="text-grey-lighten-1 mb-4"
                />
                <p class="text-body-1 text-grey-darken-1">
                  配信が開始されていません
                </p>
              </div>
            </div>
          </v-card-text>
        </v-card>

        <!-- YouTube連携セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiYoutube"
              size="24"
              class="mr-3 text-red"
            />
            <span class="text-h6 font-weight-medium">YouTube連携</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <div v-if="broadcastValue.youtubeAccount === ''">
              <template v-if="authYoutubeUrlValue === ''">
                <p class="text-body-2 text-grey-darken-1 mb-4">
                  YouTubeチャンネルと連携して配信を行います。ハンドル名を入力してください。
                </p>
                <v-text-field
                  v-model="authYoutubeFormDataValue.youtubeHandle"
                  label="YouTube ハンドル名 *"
                  placeholder="@から始まるハンドル名（例: @your-channel）"
                  variant="outlined"
                  density="comfortable"
                  class="mb-4"
                  :prepend-inner-icon="mdiYoutube"
                />
                <v-btn
                  color="red"
                  variant="elevated"
                  size="large"
                  @click="onClickLinkYouTube"
                >
                  <v-icon
                    :icon="mdiYoutube"
                    start
                  />
                  YouTube連携を開始
                </v-btn>
              </template>
              <template v-else>
                <v-alert
                  type="info"
                  variant="tonal"
                  class="mb-4"
                >
                  <v-icon
                    :icon="mdiYoutube"
                    start
                  />
                  配信者へ以下のURLを共有してください
                </v-alert>
                <v-text-field
                  v-model="authYoutubeUrlValue"
                  label="YouTube 連携用URL"
                  variant="outlined"
                  density="comfortable"
                  readonly
                  :append-inner-icon="mdiContentCopy"
                  @click:append-inner="onClickCopyUrl(authYoutubeUrlValue)"
                />
              </template>
            </div>
            <div v-else>
              <v-alert
                type="success"
                variant="tonal"
                class="mb-4"
              >
                <v-icon
                  :icon="mdiYoutube"
                  start
                />
                YouTube連携が完了しました
              </v-alert>
              <v-text-field
                v-model="broadcastValue.youtubeAccount"
                label="連携先ハンドル名"
                variant="outlined"
                density="comfortable"
                readonly
                class="mb-4"
                :prepend-inner-icon="mdiYoutube"
              />
              <v-text-field
                v-model="broadcastValue.youtubeViewerUrl"
                label="視聴画面URL"
                variant="outlined"
                density="comfortable"
                readonly
                class="mb-4"
                :append-inner-icon="mdiContentCopy"
                @click:append-inner="onClickCopyUrl(broadcast.youtubeViewerUrl)"
              />
              <v-text-field
                v-model="broadcastValue.youtubeAdminUrl"
                label="管理画面URL"
                variant="outlined"
                density="comfortable"
                readonly
                :append-inner-icon="mdiContentCopy"
                @click:append-inner="onClickCopyUrl(broadcast.youtubeAdminUrl)"
              />
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col
        cols="12"
        lg="4"
      >
        <!-- 配信状況セクション -->
        <v-card
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiRadio"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">配信状況</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <v-select
              v-model="broadcastValue.status"
              label="配信ステータス"
              :items="statuses"
              item-title="title"
              item-value="value"
              variant="outlined"
              density="comfortable"
              readonly
              class="mb-4"
            />
            <v-text-field
              v-show="isLive()"
              v-model="broadcastValue.inputUrl"
              label="入力エンドポイント"
              variant="outlined"
              density="comfortable"
              readonly
              class="mb-4"
            />
            <v-text-field
              v-show="isLive()"
              v-model="broadcastValue.outputUrl"
              label="出力エンドポイント"
              variant="outlined"
              density="comfortable"
              readonly
              class="mb-4"
            />
            <v-text-field
              v-show="isVOD()"
              v-model="broadcastValue.archiveUrl"
              label="オンデマンド配信URL"
              variant="outlined"
              density="comfortable"
              readonly
            />
          </v-card-text>
        </v-card>

        <!-- ライブ配信操作セクション -->
        <v-card
          v-show="isLive()"
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiCog"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">配信操作</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <div class="control-section">
              <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                ライブ配信制御
              </p>
              <div class="d-flex flex-column ga-3 mb-6">
                <v-btn
                  variant="outlined"
                  color="error"
                  size="large"
                  @click="onClickPause"
                >
                  <v-icon
                    :icon="mdiPause"
                    start
                  />
                  配信を停止
                </v-btn>
                <v-btn
                  variant="outlined"
                  color="success"
                  size="large"
                  @click="onSubmitUnpause"
                >
                  <v-icon
                    :icon="mdiPlay"
                    start
                  />
                  停止を解除
                </v-btn>
              </div>

              <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                入力チャンネル設定
              </p>
              <div class="d-flex flex-column ga-3 mb-6">
                <v-btn
                  variant="outlined"
                  color="primary"
                  size="large"
                  @click="onSubmitChangeRtmpInput"
                >
                  <v-icon
                    :icon="mdiRadio"
                    start
                  />
                  RTMP配信
                </v-btn>
                <v-btn
                  variant="outlined"
                  color="secondary"
                  size="large"
                  @click="onClickChangeMp4Input"
                >
                  <v-icon
                    :icon="mdiVideo"
                    start
                  />
                  MP4配信
                </v-btn>
              </div>

              <p class="text-subtitle-2 mb-3 text-grey-darken-1">
                ふた絵表示設定
              </p>
              <div class="d-flex flex-column ga-3">
                <v-btn
                  variant="outlined"
                  color="primary"
                  size="large"
                  @click="onClickActivateStaticImage"
                >
                  <v-icon
                    :icon="mdiCheck"
                    start
                  />
                  有効化
                </v-btn>
                <v-btn
                  variant="outlined"
                  color="secondary"
                  size="large"
                  @click="onClickDeactivateStaticImage"
                >
                  <v-icon
                    :icon="mdiClose"
                    start
                  />
                  無効化
                </v-btn>
              </div>
            </div>
          </v-card-text>
        </v-card>

        <!-- オンデマンド配信操作セクション -->
        <v-card
          v-show="isVOD()"
          class="form-section-card mb-6"
          elevation="2"
        >
          <v-card-title class="d-flex align-center section-header">
            <v-icon
              :icon="mdiUpload"
              size="24"
              class="mr-3 text-primary"
            />
            <span class="text-h6 font-weight-medium">オンデマンド配信</span>
          </v-card-title>
          <v-card-text class="pa-6">
            <p class="text-body-2 text-grey-darken-1 mb-4">
              アーカイブ動画をアップロードしてオンデマンド配信を開始できます。
            </p>
            <v-btn
              variant="elevated"
              color="primary"
              size="large"
              block
              @click="onClickUploadArchiveMp4"
            >
              <v-icon
                :icon="mdiUpload"
                start
              />
              アーカイブをアップロード
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>
.form-section-card {
  border-radius: 12px;
  max-width: none;
}

.section-header {
  background: linear-gradient(90deg, rgb(33 150 243 / 5%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(0 0 0 / 5%);
  padding: 20px 24px;
}

.dialog-card {
  border-radius: 12px;
}

.streaming-container {
  min-height: 400px;
}

.video-container {
  position: relative;
  width: 100%;
  min-height: 300px;
}

.video-player {
  width: 100%;
  max-height: 400px;
  border-radius: 8px;
  background: #000;
}

.video-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  background: rgb(245 245 245);
  border-radius: 8px;
  border: 2px dashed rgb(189 189 189);
}

.control-section {
  border-top: 1px solid rgb(0 0 0 / 10%);
  padding-top: 16px;
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }

  .section-header {
    padding: 16px 20px;
  }

  .dialog-card {
    border-radius: 8px;
    margin: 16px;
  }

  .video-placeholder {
    height: 200px;
  }
}
</style>

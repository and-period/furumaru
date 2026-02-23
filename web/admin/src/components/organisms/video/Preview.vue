<script lang="ts" setup>
import { dateTimeFormatter } from '~/lib/formatter'
import type { UpdateVideoRequest, VideoComment } from '~/types/api/v1'

const props = defineProps({
  formData: {
    type: Object as PropType<UpdateVideoRequest>,
    default: (): UpdateVideoRequest => ({
      videoUrl: '',
      coordinatorId: '',
      description: '',
      displayExperience: false,
      displayProduct: false,
      experienceIds: [],
      limited: false,
      productIds: [],
      _public: false,
      publishedAt: 0,
      thumbnailUrl: '',
      title: '',
    }),
  },
  comments: {
    type: Array<VideoComment>,
    default: () => [],
  },
})

const emit = defineEmits<{
  (e: 'click:ban-comment', commentId: string): void
}>()

const onClickBanComment = (commentId: string): void => {
  emit('click:ban-comment', commentId)
}
</script>

<template>
  <v-card-text>
    <v-container fluid>
      <v-row>
        <v-col
          cols="12"
          md="8"
        >
          <div
            v-if="formData.videoUrl"
            class="video-preview"
          >
            <video
              :src="formData.videoUrl"
              controls
              aria-label="動画プレビュー"
              style="width: 100%; height: auto;"
            />
          </div>
          <v-alert
            v-else
            type="info"
            variant="tonal"
          >
            動画がアップロードされていません。基本情報タブから動画をアップロードしてください。
          </v-alert>
        </v-col>
        <v-col
          cols="12"
          md="4"
        >
          <v-card>
            <v-card-title class="text-h6">
              コメント一覧
            </v-card-title>
            <v-divider />
            <v-card-text style="max-height: 500px; overflow-y: auto;">
              <v-list
                v-if="props.comments && props.comments.length > 0"
                lines="two"
              >
                <v-list-item
                  v-for="comment in props.comments"
                  :key="comment.id"
                  class="mb-2"
                >
                  <template #prepend>
                    <v-avatar
                      color="grey"
                      size="40"
                    >
                      <v-icon>mdi-account</v-icon>
                    </v-avatar>
                  </template>
                  <v-list-item-title class="font-weight-bold">
                    {{ comment.username }}
                  </v-list-item-title>
                  <v-list-item-subtitle>
                    {{ comment.comment }}
                  </v-list-item-subtitle>
                  <v-list-item-subtitle class="text-caption">
                    {{ dateTimeFormatter(comment.publishedAt) }}
                  </v-list-item-subtitle>
                  <template #append>
                    <v-btn
                      v-if="!comment.disabled"
                      icon="mdi-cancel"
                      size="small"
                      variant="text"
                      color="error"
                      @click="onClickBanComment(comment.id)"
                    >
                      <v-icon>mdi-cancel</v-icon>
                      <v-tooltip
                        activator="parent"
                        location="start"
                      >
                        コメントをBan
                      </v-tooltip>
                    </v-btn>
                    <v-chip
                      v-else
                      color="error"
                      size="small"
                      variant="tonal"
                    >
                      Ban済み
                    </v-chip>
                  </template>
                </v-list-item>
              </v-list>
              <v-alert
                v-else
                type="info"
                variant="tonal"
              >
                コメントがありません
              </v-alert>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-card-text>
</template>

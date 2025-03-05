import axios from 'axios'
import type { InstagramPost } from '~/types/api'

export const useInstagramStore = defineStore('instagram', {
  state: () => {
    return {
      instagramPostsPermalink: [] as string[],
    }
  },

  actions: {
    async listInstagramPostsPermalinkByHashTag(limit = 5) {
      const runtimeConfig = useRuntimeConfig()
      const url = `${runtimeConfig.public.INSTAGRAM_GRAPH_API_URL}/${runtimeConfig.public.INSTAGRAM_HASH_TAG_ID}/top_media`
      try {
        const response = await axios.get(url, {
          params: {
            fields: 'permalink',
            user_id: runtimeConfig.public.INSTAGRAM_USER_ID,
            access_token: runtimeConfig.public.INSTAGRAM_ACCESS_TOKEN,
            limit,
          },
        })
        this.instagramPostsPermalink = response.data.data.map((post: InstagramPost) => post.permalink)
      }
      catch (e) {
        console.error(e)
      }
      return
    },
  },

  getters: {
  },
})

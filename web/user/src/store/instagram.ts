import axios from 'axios'
import type { InstagramOEmbed } from '~/types/api'

export const useInstagramStore = defineStore('instagram', {
  state: () => {
    return {
      instagramOEmbed: {} as InstagramOEmbed,
      OEmbedHTML: '' as string,
    }
  },

  actions: {
    async getInstagramOEmbed() {
      console.log('getInstagramOEmbed')
      const url = 'https://graph.facebook.com/v22.0/instagram_oembed'
      const postUrl = 'https://www.instagram.com/p/DBwRdmPvSy7/?igsh=MWxsY291djNtcjQ4OA=='
      const runtimeConfig = useRuntimeConfig()
      try {
        const response = await axios.get(url, {
          params: {
            url: postUrl,
            access_token: runtimeConfig.public.INSTAGRAM_ACCESS_TOKEN,
            omitscript: true,
          },
        })
        return response.data
      }
      catch (e) {
        console.error(e)
      }
      return {} as InstagramOEmbed
    },
  },

  getters: {
  },
})

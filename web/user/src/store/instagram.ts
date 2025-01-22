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
      const response: InstagramOEmbed
        = await this.instagramApiClient().getInstagramOEmbedContent('https://www.instagram.com/p/DBwRdmPvSy7/?igsh=MWxsY291djNtcjQ4OA==')

      this.instagramOEmbed = response
      this.OEmbedHTML = response.html
    },
  },

  getters: {
  },
})

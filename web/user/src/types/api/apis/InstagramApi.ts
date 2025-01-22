import axios from "axios";
import type { InstagramOEmbed } from "../models/index";

export class InstagramApi {
  async getInstagramOEmbedContent(postUrl: string): Promise<InstagramOEmbed> {
    const runtimeConfig = useRuntimeConfig()
    const url = "https://graph.facebook.com/v22.0/instagram_oembed";
    try {
      const response = await axios.get(url, {
        params: {
          url: postUrl,
          access_token: runtimeConfig.public.INSTAGRAM_ACCESS_TOKEN,
          omitscript: true,
        },
      });
      return response.data;
    } catch (e) {
      console.error(e);
    }
    return {} as InstagramOEmbed;
  }
}

export interface VolunteerBlogContentInnerResponse {
  id: string
  createdAt: string
  updatedAt: string
  publishedAt: string
  revisedAt: string
  title: string
  content: string
}

export interface VolunteerBlogListResponse {
  contents: VolunteerBlogContentInnerResponse[]
  totalCount: number
  offset: number
  limit: number
}

export interface VolunteerBlogCategory {
  id: string
  name: string
}

export interface VolunteerBlogEyeCatch {
  url: string
}

export interface VolunteerBlogContentInnerResponse {
  id: string
  title: string
  name: string
  location: string
  content: string
  eyecatch: VolunteerBlogEyeCatch
  category: VolunteerBlogCategory[]
}

export interface VolunteerBlogListResponse {
  contents: VolunteerBlogContentInnerResponse[]
  totalCount: number
  offset: number
  limit: number
}

export interface VolunteerBlogItemResponse {
  title: string
  name: string
  location: string
  content: string
  eyecatch: VolunteerBlogEyeCatch
  category: VolunteerBlogCategory[]
}

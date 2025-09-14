export interface UploadStatus {
  isUploading: boolean
  hasError: boolean
  errorMessage: string
}

export interface VideoComment {
  id: string
  userId: string
  userName: string
  content: string
  createdAt: number
  isBanned?: boolean
}

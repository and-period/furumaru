import axios from 'axios'
import type { RawAxiosRequestHeaders } from 'axios'

import { apiClient } from '~/plugins/api-client'
import { UploadStatus } from '~/types/api'

/**
 * サーバーへファイルをアップロードする非同期関数
 * @param file アップロードするファイル
 * @param uploadUrl アップロード時に使用する署名付きURL
 * @returns 参照先URL
 */
export async function fileUpload(
  file: File,
  key: string,
  url: string,
): Promise<string> {
  // 署名付きURLを基にファイルをアップロード
  await upload(file, url)

  // アップロード後のバリデーション結果を取得
  return await getUploadResult(key)
}

/**
 * 処理を一定時間待機するための非同期関数
 * @param ms 待機する時間（ミリ秒）
 * @returns
 */
function sleep(ms: number): Promise<void> {
  return new Promise((resolve: any) => setTimeout(resolve, ms))
}

/**
 * サーバーへファイルをアップロードする非同期関数
 * @param file アップロードするファイル
 * @param uploadUrl アップロード時に使用する署名付きURL
 * @returns
 */
async function upload(file: File, uploadUrl: string): Promise<void> {
  const headers: RawAxiosRequestHeaders = {
    'Content-Type': file.type,
  }
  await axios.put(uploadUrl, file, { headers })
}

/**
 * ファイルアップロード後の実行結果を取得する非同期関数
 * @param uploadUrl アップロード先URL
 * @returns ファイルの参照先URL
 */
async function getUploadResult(uploadUrl: string): Promise<string> {
  while (true) {
    // アップロード処理が、サーバー側では非同期実行となるため
    await sleep(200)

    const event = await apiClient.otherApi().v1GetUploadState(uploadUrl)

    switch (event.data.status) {
      case UploadStatus.SUCEEDED:
        return event.data.url
      case UploadStatus.FAILED:
        throw new Error('ファイルのアップロードに失敗しました。')
      case UploadStatus.WAITING:
        continue // 再度状態を取得
    }
  }
}

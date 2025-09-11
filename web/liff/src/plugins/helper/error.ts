import { ResponseError, FetchError } from '@/types/api/facility/runtime';
import type { ErrorResponse } from '@/types/api/facility';

/**
 * APIエラーをUI表示向けの日本語メッセージに整形
 */
export async function buildApiErrorMessage(err: unknown): Promise<string> {
  if (err instanceof FetchError) {
    return 'ネットワークエラーが発生しました。通信状況をご確認ください。';
  }

  if (err instanceof ResponseError) {
    const status = err.response.status;
    let apiMsg: string | undefined;
    try {
      const data = (await err.response.clone().json()) as Partial<ErrorResponse> | undefined;
      apiMsg = data?.message || data?.detail || undefined;
    }
    catch {
      // 非JSONレスポンスは無視
    }

    // ステータスごとの既定文言（ステータス表記は最後にまとめて付与）
    let defaultMsg: string;
    switch (status) {
      case 400:
      case 422:
        defaultMsg = '入力内容に誤りがあります。ご確認ください。';
        break;
      case 401:
        defaultMsg = '認証エラーが発生しました。LINEで再ログインしてください。';
        break;
      case 403:
        defaultMsg = 'この操作の権限がありません。';
        break;
      case 404:
        defaultMsg = '施設が見つかりませんでした。';
        break;
      case 409:
        defaultMsg = '既に登録済みの可能性があります。';
        break;
      case 500:
      case 502:
      case 503:
      case 504:
        defaultMsg = 'サーバーエラーが発生しました。時間をおいて再度お試しください。';
        break;
      default:
        defaultMsg = 'エラーが発生しました。';
        break;
    }

    const detail = apiMsg ? `: ${apiMsg}` : '';
    return `${defaultMsg}${detail} (HTTP ${status})`;
  }

  if (err instanceof Error) {
    return err.message || 'エラーが発生しました。';
  }
  return 'エラーが発生しました。';
}

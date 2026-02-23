import { ref } from 'vue'
import type { CreditCardData } from '../types/index'

export function useKomojuTokenize(komojuHost: string, publishableKey: string) {
  const tokenizing = ref(false)
  const tokenError = ref<string | null>(null)

  async function tokenize(card: CreditCardData): Promise<string> {
    tokenizing.value = true
    tokenError.value = null
    try {
      const response = await fetch(`${komojuHost}/api/v1/tokens`, {
        method: 'POST',
        headers: {
          'Authorization': `Basic ${btoa(publishableKey + ':')}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          payment_details: {
            type: 'credit_card',
            number: card.number,
            month: String(card.month).padStart(2, '0'),
            year: String(card.year),
            verification_value: card.verificationValue,
            name: card.name,
          },
        }),
      })
      if (!response.ok) {
        throw new Error('カードのトークン化に失敗しました')
      }
      const data = await response.json()
      return data.id
    }
    catch (e) {
      const message = e instanceof Error ? e.message : 'カードのトークン化に失敗しました'
      tokenError.value = message
      throw e
    }
    finally {
      tokenizing.value = false
    }
  }

  return { tokenize, tokenizing, tokenError }
}

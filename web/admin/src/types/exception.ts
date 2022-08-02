class ApiBaseError<T> extends Error {
  constructor(public message: string, public detail: T) {
    super(message)
  }
}

export class ConnectionError<T> extends ApiBaseError<T> {
  constructor(public detail: T) {
    super(
      '現在、システムが停止中です。時間をおいてから再度アクセスしてください。',
      detail
    )
  }
}

export class InternalServerError<T> extends ApiBaseError<T> {
  constructor(public detail: T) {
    super('不明なエラーが発生しました。', detail)
  }
}

export class AuthError<T> extends ApiBaseError<T> {
  constructor(public status: number, public message: string, public detail: T) {
    super(message, detail)
  }
}

export class ValidationError<T> extends ApiBaseError<T> {
  constructor(public status: number, public message: string, public detail: T) {
    super(message, detail)
  }
}

export class NotFoundError<T> extends ApiBaseError<T> {
  constructor(public status: number, public message: string, public detail: T) {
    super(message, detail)
  }
}

export class ConflictError<T> extends ApiBaseError<T> {
  constructor(public status: number, public message: string, public detail: T) {
    super(message, detail)
  }
}

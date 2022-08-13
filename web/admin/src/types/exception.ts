class ApiBaseError<T> extends Error {
  constructor(public message: string, public cause?: T) {
    super(message)
    this.cause = cause
  }
}

export class ConnectionError<T> extends ApiBaseError<T> {
  constructor(public cause?: T) {
    super(
      '現在、システムが停止中です。時間をおいてから再度アクセスしてください。',
      cause
    )
  }
}

export class InternalServerError<T> extends ApiBaseError<T> {
  status = 500
  constructor(public cause?: T) {
    super('不明なエラーが発生しました。', cause)
  }
}

export class AuthError<T> extends ApiBaseError<T> {
  status = 401
  constructor(public message: string, public cause?: T) {
    super(message, cause)
  }
}

export class ValidationError<T> extends ApiBaseError<T> {
  status = 400
  constructor(public message: string, public cause?: T) {
    super(message, cause)
  }
}

export class NotFoundError<T> extends ApiBaseError<T> {
  status = 404
  constructor(public message: string, public cause?: T) {
    super(message, cause)
  }
}

export class ConflictError<T> extends ApiBaseError<T> {
  status = 409
  constructor(public message: string, public cause?: T) {
    super(message, cause)
  }
}

export class PerconditionError<T> extends ApiBaseError<T> {
  status = 412
  constructor(public message: string, public cause?: T) {
    super(message, cause)
  }
}

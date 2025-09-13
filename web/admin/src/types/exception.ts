export class ApiBaseError<T> extends Error {
  status: number
  constructor(public override message: string, public override cause?: T) {
    super(message)
    this.status = 0
    this.cause = cause
  }
}

export class ConnectionError<T> extends ApiBaseError<T> {
  constructor(public override cause?: T) {
    super(
      '現在、システムが停止中です。時間をおいてから再度アクセスしてください。',
      cause,
    )
  }
}

export class InternalServerError<T> extends ApiBaseError<T> {
  override status = 500
  constructor(public override cause?: T) {
    super('不明なエラーが発生しました。', cause)
  }
}

export class AuthError<T> extends ApiBaseError<T> {
  override status = 401
  constructor(public override message: string, public override cause?: T) {
    super(message, cause)
  }
}

export class ValidationError<T> extends ApiBaseError<T> {
  override status = 400
  constructor(public override message: string, public override cause?: T) {
    super(message, cause)
  }
}

export class PermissionError<T> extends ApiBaseError<T> {
  override status = 403
  constructor(public override message: string, public override cause?: T) {
    super(message, cause)
  }
}

export class NotFoundError<T> extends ApiBaseError<T> {
  override status = 404
  constructor(public override message: string, public override cause?: T) {
    super(message, cause)
  }
}

export class ConflictError<T> extends ApiBaseError<T> {
  override status = 409
  constructor(public override message: string, public override cause?: T) {
    super(message, cause)
  }
}

export class PreconditionError<T> extends ApiBaseError<T> {
  override status = 412
  constructor(public override message: string, public override cause?: T) {
    super(message, cause)
  }
}

export class TooManyRequestsError<T> extends ApiBaseError<T> {
  override status = 429
  constructor(public override message: string, public override cause?: T) {
    super(message, cause)
  }
}
export class CancelledError<T> extends ApiBaseError<T> {
  override status = 499
  constructor(public override message: string, public override cause?: T) {
    super(message, cause)
  }
}

export class NotImplementedError<T> extends ApiBaseError<T> {
  override status = 501
  constructor(public override cause?: T) {
    super('この操作はまだ利用できません。', cause)
  }
}

export class ServiceUnavailableError<T> extends ApiBaseError<T> {
  override status = 503
  constructor(public override cause?: T) {
    super(
      'サービスが一時的利用できません。時間をおいてから再度アクセスしてください。',
      cause,
    )
  }
}

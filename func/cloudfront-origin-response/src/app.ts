import {
  CloudFrontHeaders,
  CloudFrontRequest,
  CloudFrontResponse,
  CloudFrontResponseEvent,
  CloudFrontResponseResult,
} from 'aws-lambda';
import {
  GetObjectCommand,
  GetObjectCommandInput,
  PutObjectCommand,
  PutObjectCommandInput,
  S3Client,
} from '@aws-sdk/client-s3';
import * as querystring from 'querystring';
import sharp, { ResizeOptions } from 'sharp';

const s3Client = new S3Client({ region: process.env.AWS_REGION });
const cacheControl = 'max-age=0,s-maxage=2592000'; // 30 days

/**
 * Lambda@Edgeを利用して画像オブジェクトが存在するかを確認し、必要に応じて画像リサイズを実行する
 * @param {Object} event - CloudFront Request event Format
 * @returns {Object} object - CloudFront Result Format
 */
export const lambdaHandler = async (event: CloudFrontResponseEvent): Promise<CloudFrontResponseResult> => {
  const response: CloudFrontResponse = event.Records[0].cf.response;
  console.log('received event', JSON.stringify(event));

  // 画像変換が必要かの検証
  const request: CloudFrontRequest = event.Records[0].cf.request;
  if (request.querystring === '') {
    // キャッシュTTLの上書き
    response.headers = {
      ...response.headers,
      'cache-control': [{ key: 'Cache-Control', value: cacheControl }],
    };
    return response;
  }

  // オリジンサーバ情報の取得
  const bucketName: string = request.origin?.s3?.domainName.split('.')[0] || '';
  if (!bucketName) {
    console.log('failed to get bucket name');
    return response;
  }

  // 画像変換のパラメータ取得
  const params = querystring.parse(request.querystring);
  const details: FileDetails = getFileDetails(request.uri, params);
  if (!details.convertable) {
    console.log('this image is not convertable', { details });
    return response;
  }
  console.log('received to convertable image', { details });

  // S3から画像を取得して画像の加工処理
  let image: Buffer;
  try {
    const input: GetObjectCommandInput = {
      Bucket: bucketName,
      Key: details.objectKey,
    };
    const data = await s3Client.send(new GetObjectCommand(input));
    if (!data.Body) {
      throw new Error('this object is empty');
    }
    const object = await data.Body.transformToByteArray();
    image = await resizeImage(object, details.options);
  } catch (err) {
    console.log('failed to get object from S3 or convert image', err);
    return response;
  }
  console.log('finished to convert image', { details });

  // 加工後の画像をアップロード（検証用）
  try {
    const suffix = new Array([details.options.width, details.options.height, details.options.format]).join('_');
    const key = `${details.path.directory}/fixed/${details.path.filename}_${suffix}.${details.dstFormat}`;
    const input: PutObjectCommandInput = {
      Bucket: bucketName,
      Key: key,
      Body: image,
      ContentType: getMimeType(details.dstFormat),
    };
    console.log('put object to S3', { key, suffix, input });
    await s3Client.send(new PutObjectCommand(input));
  } catch (err) {
    // 画像のリサイズ処理は成功したがアップロードに失敗した状態であれば、エラーは返さずリサイズ後の画像を返す
    console.log('failed to put object to S3', err);
  }
  console.log('finished to put object to S3', { details });

  const headers: CloudFrontHeaders = {
    ...response.headers,
    'content-type': [{ key: 'Content-Type', value: getMimeType(details.dstFormat) }],
    'cache-control': [{ key: 'Cache-Control', value: cacheControl }],
  };
  const res: CloudFrontResponseResult = {
    status: '200',
    statusDescription: 'OK',
    headers: headers,
    body: image.toString('base64'),
    bodyEncoding: 'base64',
  };
  return res;
};

type FileDetails = {
  objectKey: string;
  convertable: boolean;
  path: FilePath;
  srcFormat: ImageFormat;
  dstFormat: ImageFormat;
  options: ConvertOptions;
};

type FilePath = {
  directory: string;
  filename: string;
  extension: string;
};

type ConvertOptions = {
  width?: number;
  height?: number;
  format?: ImageFormat;
  fit?: ImageFitType;
  blur?: number;
};

type ImageFormat = 'jpg' | 'jpeg' | 'png' | 'svg' | 'webp';

type ImageFitType = 'cover' | 'contain' | 'fill' | 'inside' | 'outside';

// リクエストされた画像のパスからS3オブジェクトのキーとリサイズオプションを取得
function getFileDetails(uri: string, params: querystring.ParsedUrlQuery): FileDetails {
  // S3オブジェクトのキー形式に変換
  // e.g. /images/image.jpg -> images/image.jpg
  const key: string = uri.substring(1);

  // S3オブジェクトのキーからディレクトリ、ファイル名、拡張子を取得
  const [, directory, filename, extension] = uri.match(/(.*)\/(.*)\.(.*)/) || [];

  // 画像のリサイズが可能かを判定
  const convertable: boolean = isConvertableExtension(extension);

  const details: FileDetails = {
    objectKey: key,
    srcFormat: extension as ImageFormat,
    dstFormat: extension as ImageFormat,
    convertable: false,
    path: { directory, filename, extension },
    options: {},
  };
  if (!convertable) {
    return details;
  }

  // リサイズが可能な場合、リサイズ設定を取得
  const options = getImageOptions(params);

  details.options = options;
  if (details.options.format) {
    details.dstFormat = details.options.format;
  }
  details.convertable = isConvertableOptions(details.options);
  return details;
}

// ファイル拡張子を基に画像のリサイズが必要かを判定
function isConvertableExtension(extension: string): boolean {
  console.log(`received extension='${extension}'`);
  return ['jpg', 'jpeg', 'png', 'svg', 'webp'].includes(extension);
}

// フィット種別が利用可能なものかを判定
function isConvertableFitType(fit: string): boolean {
  console.log(`received fit='${fit}'`);
  return ['cover', 'contain', 'fill', 'inside', 'outside'].includes(fit);
}

// クエリパラメータを基に画像のリサイズが必要かを判定
function isConvertableOptions(opts: ConvertOptions): boolean {
  if (opts.width || opts.height) {
    return true;
  }
  if (opts.format) {
    return true;
  }
  if (opts.blur) {
    return true;
  }
  return false;
}

// クエリパラメータから画像のリサイズオプションを取得
function getImageOptions(params: querystring.ParsedUrlQuery): ConvertOptions {
  const { width, height, format, fit, blur, dpr } = params;

  const options: ConvertOptions = {};
  if (width && Number(width) > 0) {
    options.width = Number(width);
  }
  if (height && Number(height) > 0) {
    options.height = Number(height);
  }
  if (dpr && Number(dpr) > 0) {
    options.width = options.width && options.width * Number(dpr);
    options.height = options.height && options.height * Number(dpr);
  }
  if (format && isConvertableExtension(String(format))) {
    options.format = format as ImageFormat;
  }
  if (fit && isConvertableFitType(String(fit))) {
    options.fit = fit as ImageFitType;
  }
  if (blur && Number(blur) >= 0.3 && Number(blur) <= 1000) {
    options.blur = Number(blur);
  }
  return options;
}

// 画像のMIMEタイプを取得
function getMimeType(extension: ImageFormat): string {
  switch (extension) {
    case 'jpg':
    case 'jpeg':
      return 'image/jpeg';
    case 'png':
      return 'image/png';
    case 'svg':
      return 'image/svg+xml';
    case 'webp':
      return 'image/webp';
    default:
      return 'image/jpeg';
  }
}

// 画像リサイズの実行 - 追加が必要な場合、以下ドキュメントを参照
// @see: https://sharp.pixelplumbing.com/api-operation
async function resizeImage(object: Uint8Array, options: ConvertOptions): Promise<Buffer> {
  let resize = sharp(object);
  if (options.height || options.height) {
    const opts: ResizeOptions = {
      width: options.width,
      height: options.height,
      fit: options.fit,
    };
    resize = resize.resize(opts);
  }
  if (options.blur) {
    resize = resize.blur(options.blur);
  }
  if (options.format) {
    resize = resize.toFormat(options.format);
  }
  return await resize.toBuffer();
}

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
import sharp, { ResizeOptions, RGBA, Color } from 'sharp';

const s3Client = new S3Client({ region: process.env.AWS_REGION });
const cacheControl = 'max-age=0,s-maxage=2592000'; // 30 days
const defaultBgColor: Color = { r: 0, g: 0, b: 0, alpha: 0 };

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
      Key: details.srcKey,
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
  const body: string = image.toString('base64');
  try {
    const input: PutObjectCommandInput = {
      Bucket: bucketName,
      Key: details.dstKey,
      Body: body,
      ContentType: getMimeType(details.dstFormat),
    };
    console.log('put object to S3', { input, details });
    s3Client.send(new PutObjectCommand(input));
  } catch (err) {
    // 画像のリサイズ処理は成功したがアップロードに失敗した状態であれば、エラーは返さずリサイズ後の画像を返す
    console.log('failed to put object to S3', err);
  }

  const headers: CloudFrontHeaders = {
    ...response.headers,
    'content-type': [{ key: 'Content-Type', value: getMimeType(details.dstFormat) }],
    'cache-control': [{ key: 'Cache-Control', value: cacheControl }],
  };
  const res: CloudFrontResponseResult = {
    status: '200',
    statusDescription: 'OK',
    headers: headers,
    body: body,
    bodyEncoding: 'base64',
  };
  return res;
};

type FileDetails = {
  srcKey: string;
  dstKey: string;
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
  bgColor?: Color;
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
    srcKey: key,
    dstKey: '',
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
  }
  details.convertable = isConvertableOptions(details.options);
  details.dstKey = `${directory.substring(1)}/${filename}_${getFileSuffix(details)}.${extension}`;
  return details;
}

function getFileSuffix(details: FileDetails): string {
  if (!details || !details.options) {
    return '';
  }
  const keys: string[] = [];
  if (details.options.width) {
    keys.push(`w${details.options.width}`);
  }
  if (details.options.height) {
    keys.push(`h${details.options.height}`);
  }
  if (details.options.format) {
    keys.push(`fmt${details.options.format}`);
  }
  if (details.options.fit) {
    keys.push(`fit${details.options.fit}`);
  }
  if (details.options.blur) {
    keys.push(`b${details.options.blur}`);
  }
  return keys.join('_');
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
  const { width, height, format, fit, blur, dpr, bg_color } = params;

  const options: ConvertOptions = {
    bgColor: defaultBgColor,
  };
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
  if (bg_color) {
    options.bgColor = getBackgroundColor(String(bg_color));
  }
  return options;
}

// RGBAの背景色を取得
function getBackgroundColor(color: string): Color | undefined {
  const strs: string[] = color.split(',');
  if (strs.length < 3) {
    return color;
  }

  // バリデーション検証
  if (strs.slice(0, 3).some((v): boolean => Number(v) < 0 || 255 < Number(v))) {
    console.log(`invalid rgba range. color='${color}'`);
    return defaultBgColor;
  }
  if (strs.length > 3 && (Number(strs[3]) < 0 || 1 < Number(strs[3]))) {
    console.log(`invalid alpha range. color='${color}'`);
    return defaultBgColor;
  }

  const rgba: RGBA = {
    r: Number(strs[0]),
    g: Number(strs[1]),
    b: Number(strs[2]),
    alpha: strs.length > 3 ? Number(strs[3]) : 1,
  };
  return rgba;
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
  // 横のみ指定の場合、4:3のアスペクト比になるように高さを計算
  if (options.width && !options.height) {
    options.height = Math.floor((options.width * 3) / 4);
  }
  // 縦のみ指定の場合、4:3のアスペクト比になるように幅を計算
  if (!options.width && options.height) {
    options.width = Math.floor((options.height * 4) / 3);
  }

  let resize = sharp(object);
  if (options.height || options.height || options.bgColor) {
    const opts: ResizeOptions = {
      width: options.width,
      height: options.height,
      fit: options.fit,
      background: options.bgColor,
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
